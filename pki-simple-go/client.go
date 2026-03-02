package main

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

func RunClient() {
	// -----------------------------
	// Configuration (env + defaults)
	// -----------------------------
	serverURL := getenv("SERVER_URL", "https://localhost:8443/hello")

	clientCertPath := getenv("CLIENT_CERT_PATH", filepath.Join("certs", "client.cert.pem"))
	clientKeyPath := getenv("CLIENT_KEY_PATH", filepath.Join("certs", "client.key.pem"))
	caCertPath := getenv("CA_CERT_PATH", filepath.Join("certs", "client", "inter-root-combined.cert.pem"))

	// -----------------------------
	// Parse & validate URL
	// -----------------------------
	u, err := url.Parse(serverURL)
	if err != nil {
		log.Fatalf("invalid SERVER_URL: %v", err)
	}

	if u.Scheme != "https" {
		log.Fatal("only https scheme is allowed")
	}

	host := u.Hostname()
	if host == "" {
		log.Fatal("URL must include hostname")
	}

	// Allowed hosts (SSRF protection)
	allowedHosts := map[string]bool{
		"localhost":   true,
		"example.com": true,
	}

	if !allowedHosts[host] {
		log.Fatalf("host not allowed: %s", host)
	}

	// -----------------------------
	// Load client certificate
	// -----------------------------
	cert, err := tls.LoadX509KeyPair(clientCertPath, clientKeyPath)
	if err != nil {
		log.Fatalf("failed loading client cert/key: %v", err)
	}

	// -----------------------------
	// Load CA certificates
	// -----------------------------
	caCertPEM, err := os.ReadFile(caCertPath)
	if err != nil {
		log.Fatalf("failed reading CA bundle: %v", err)
	}

	caPool := x509.NewCertPool()
	if !caPool.AppendCertsFromPEM(caCertPEM) {
		log.Fatal("failed to append CA certificates")
	}

	// -----------------------------
	// TLS configuration (tight)
	// -----------------------------
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caPool,
		ServerName:   host, // IMPORTANT: hostname verification
		MinVersion:   tls.VersionTLS12,
	}

	// -----------------------------
	// HTTP transport (production)
	// -----------------------------
	dialer := &net.Dialer{
		Timeout:   5 * time.Second,
		KeepAlive: 30 * time.Second,
	}

	transport := &http.Transport{
		TLSClientConfig:       tlsConfig,
		DialContext:           dialer.DialContext,
		TLSHandshakeTimeout:   5 * time.Second,
		ResponseHeaderTimeout: 5 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		IdleConnTimeout:       90 * time.Second,
		MaxIdleConns:          10,
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   10 * time.Second, // total request timeout
	}

	// -----------------------------
	// Execute request with retries
	// -----------------------------
	resp, err := doGetWithRetry(client, u.String(), 3)
	if err != nil {
		log.Fatalf("request failed: %v", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Fatalf("Unable to Close the Body")
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("failed reading response body: %v", err)
	}

	fmt.Printf("Server response: %s\n", body)
}

// -----------------------------
// Helpers
// -----------------------------

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func doGetWithRetry(client *http.Client, url string, attempts int) (*http.Response, error) {
	var lastErr error

	for i := 1; i <= attempts; i++ {
		resp, err := client.Get(url)
		if err == nil {
			return resp, nil
		}

		lastErr = err
		time.Sleep(time.Duration(i) * 500 * time.Millisecond)
	}

	return nil, errors.New("all retry attempts failed: " + lastErr.Error())
}
