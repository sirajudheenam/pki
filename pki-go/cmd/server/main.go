package main

import (
	"fmt"
	"log"
	"os"

	"github.com/sirajudheenam/pki/pki-go/internal/config"
	"github.com/sirajudheenam/pki/pki-go/internal/server"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Printf("Warning: Error loading config: %v", err)
		log.Println("Using default configuration...")
	}

	// Get certificate path
	certPath := cfg.GetCertificatePath()

	// Verify certificate directory exists
	if _, err := os.Stat(certPath); os.IsNotExist(err) {
		log.Fatalf("Certificate directory not found: %s. Please run 'make cert-gen HOSTNAME=%s' first", certPath, cfg.Hostname)
	}

	// Create server with configured values
	srv, err := server.NewServer(fmt.Sprintf(":%s", cfg.Port), certPath)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Starting server on port %s using certificates from %s", cfg.Port, certPath)
	log.Fatal(srv.Start())
}
