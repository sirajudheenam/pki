package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	// Allow server URL override via env var
	serverURL := os.Getenv("SERVER_URL")
	if serverURL == "" {
		serverURL = "https://localhost:8443/hello"
	}

	// Load client cert + key
	cert, err := tls.LoadX509KeyPair("certs/client/client.cert.pem", "certs/client/client.key.pem")
	if err != nil {
		log.Fatal("Failed loading client cert/key:", err)
	}

	// Load CA bundle (root + intermediate)
	caCert, err := os.ReadFile("certs/client/inter-root-combined.cert.pem")
	if err != nil {
		log.Fatal(err)
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// TLS config
	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		RootCAs:            caCertPool,
		InsecureSkipVerify: false, // enforce validation
		MinVersion:         tls.VersionTLS12,
	}

	// HTTPS client
	client := &http.Client{
		Transport: &http.Transport{TLSClientConfig: tlsConfig},
	}

	// Send request
	resp, err := client.Get(serverURL)
	if err != nil {
		log.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("Server response: %s\n", body)
}
