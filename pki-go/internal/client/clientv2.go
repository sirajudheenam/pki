package client

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/sirajudheenam/pki/pki-go/internal/config"
)

type ClientV2 struct {
	Addr string
	http *http.Client
}

type Config *config.ClientConfig

func NewClientV2(cfg Config) (*ClientV2, error) {

	serverURL := strings.Join([]string{
		"https://",
		cfg.Hostname,
		":",
		cfg.Port,
		"/",
		cfg.ServerRootPath,
	}, "")

	// -----------------------------
	// Parse and validate server URL
	// -----------------------------
	parsedURL, err := url.Parse(serverURL)
	if err != nil {
		log.Fatalf("invalid server URL: %v", err)
	}

	if parsedURL.Scheme != "https" {
		log.Fatal("only https scheme is allowed")
	}

	host := parsedURL.Hostname()
	if host == "" {
		log.Fatal("server URL must include hostname")
	}

	// Allowed hosts (prevent SSRF)
	allowedHosts := map[string]bool{
		"go-mtls-server-service": true,
		"localhost":              true,
	}

	if !allowedHosts[host] {
		log.Fatalf("host not allowed: %s", host)
	}

	certPath := strings.Join([]string{cfg.CertBaseDir, "/", cfg.Hostname, "/", cfg.CertSubDir}, "")
	certFile := filepath.Join(certPath, "client.cert.pem")
	keyFile := filepath.Join(certPath, "client.key.pem")
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		log.Fatalf("Failed loading client cert/key: %v", err)
	}

	root, err := os.OpenRoot(certPath)
	if err != nil {
		log.Fatalf("Unable Load *certPath")
	}

	certDirFile, err := root.Open("inter-root-combined.cert.pem")
	if err != nil {
		log.Fatalf("unable to load inter-root-combined.cert.pem")
	}

	caCert, err := io.ReadAll(certDirFile)
	if err != nil {
		log.Printf("Unable to load caCert")
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
		MinVersion:   tls.VersionTLS12,
		MaxVersion:   tls.VersionTLS13,
	}

	// -----------------------------
	// Secure Transport
	// -----------------------------
	dialer := &net.Dialer{
		Timeout:   5 * time.Second,
		KeepAlive: 30 * time.Second,
	}
	// Verify Server Availability: (Additional)
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(cfg.Hostname, cfg.Port), 5*time.Second)
	if err != nil {
		log.Printf("Server not available: %v\n", err)
	}

	if err := conn.Close(); err != nil {
		fmt.Printf("Unbale to close connection: %v", err)
	}

	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			h, _, err := net.SplitHostPort(addr)
			if err != nil {
				return nil, err
			}
			if !allowedHosts[h] {
				return nil, fmt.Errorf("blocked host: %s", h)
			}
			return dialer.DialContext(ctx, network, addr)
		},
		TLSHandshakeTimeout:   5 * time.Second,
		ResponseHeaderTimeout: 5 * time.Second,
		IdleConnTimeout:       90 * time.Second,
	}

	httpClient := &http.Client{
		Transport: transport,
		Timeout:   10 * time.Second,
	}

	client := &ClientV2{
		Addr: serverURL,
		http: httpClient,
	}

	return client, nil

}

// DoRequest calls /hello endpoint
func (c *ClientV2) DoRequest() (*http.Response, error) {

	// -----------------------------
	// Build request safely
	// -----------------------------
	req, err := http.NewRequest(http.MethodGet, c.Addr, nil)
	if err != nil {
		log.Fatalf("failed to create HTTP request: %v", err)
	}

	// Request  Dump
	dump, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		log.Fatalf("Failed to dump HTTP request: %v", err)
	}
	log.Printf("HTTP Request (DUMP): %s\n", dump)

	resp, err := c.http.Do(req)
	if err != nil {
		if tlsErr, ok := err.(x509.UnknownAuthorityError); ok {
			log.Printf("httpRequest v2 - Unknown authority: %s", tlsErr.Error())
		} else if tlsErr, ok := err.(x509.HostnameError); ok {
			log.Printf("httpRequest v2 - Hostname error: %s", tlsErr.Error())
		} else {
			log.Printf("httpRequest v2- General error: %v", err)
		}
		return nil, err
	}

	return resp, nil
}
