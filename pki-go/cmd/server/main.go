// pki-go/cmd/server/main.go
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

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
	certPath := cfg.Server.GetCertificatePath()

	// Verify certificate directory exists
	if _, err := os.Stat(certPath); os.IsNotExist(err) {
		log.Fatalf("Certificate directory not found: %s. Please run 'make cert-gen HOSTNAME=%s' first", certPath, cfg.Server.Hostname)
	}

	// Create server with configured values
	srv, err := server.NewServer(fmt.Sprintf(":%s", cfg.Server.Port), certPath)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Server created with address: %s and certificates from %s", srv.Addr, certPath)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// StartAsync() is useful only when you do NOT want to block. // Start server asynchronously
	log.Println("Starting server asynchronously...")
	errCh := srv.StartAsync()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-errCh:
		// log.Fatal logs the error and exits the program
		log.Fatalf("server error: %v", err)

	case sig := <-sigCh:
		log.Printf("received signal %s, shutting down...", sig)
		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("shutdown error: %v", err)
		} else {
			log.Println("Server shutdown gracefully.")
		}
	}
	log.Println("Exiting application.")
}
