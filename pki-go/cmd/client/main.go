package main

import (
	"flag"
	"log"
	"os"

	"github.com/sirajudheenam/pki/pki-go/internal/client"
)

func main() {
	// Default values
	defaultURL := "https://localhost:8443"
	defaultCertPath := "./certs/client"

	// Override defaults with environment variables if set
	if envURL := os.Getenv("SERVER_URL"); envURL != "" {
		defaultURL = envURL
	}
	if envCert := os.Getenv("CLIENT_CERTS"); envCert != "" {
		defaultCertPath = envCert
	}

	// Command-line flags (take precedence over env vars)
	serverURL := flag.String("url", defaultURL, "Server base URL")
	certPath := flag.String("certs", defaultCertPath, "Path to client certificates")
	flag.Parse()

	// Create client
	c, err := client.NewClient(*serverURL, *certPath)
	if err != nil {
		log.Fatal("Failed to create client:", err)
	}

	// Make request
	resp, err := c.DoRequest()
	if err != nil {
		log.Fatal("Request failed:", err)
	}

	log.Println("Server response:", resp)
}
