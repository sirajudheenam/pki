// pki-go/internal/server/server.go
package server

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// Server holds the HTTP server
type Server struct {
	Addr string
	http *http.Server
}

// New returns a configured server
func NewServer(addr, certDir string) (*Server, error) {

	log.Println("internal/server/server.go - called")

	certPath := filepath.Join(certDir, "server.chain.pem")
	keyPath := filepath.Join(certDir, "server.key.pem")

	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return nil, fmt.Errorf("failed loading server cert/key: %w", err)
	}

	// log.Println("LoadX509KeyPair - called")
	root, err := os.OpenRoot(certDir)
	if err != nil {
		return nil, err
	}

	caRootFile, err := root.Open("root.cert.pem")
	if err != nil {
		return nil, err
	}
	caRootCertContent, err := io.ReadAll(caRootFile)
	if err != nil {
		return nil, err
	}
	// log.Println("caRootFile: root.cert.pem is read")

	caInterCertFile, err := root.Open("intermediate.cert.pem")
	if err != nil {
		log.Fatalf("Unable to load caInterCertFile ")
	}
	caInterCertContent, err := io.ReadAll(caInterCertFile)
	if err != nil {
		log.Fatalf("Unable to read caInterCertFile ")
		return nil, err
	}
	// log.Println("certDirFile: intermediate.cert.pem is read")

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caRootCertContent)
	caCertPool.AppendCertsFromPEM(caInterCertContent)

	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		ClientCAs:          caCertPool,
		ClientAuth:         tls.RequireAndVerifyClientCert,
		MinVersion:         tls.VersionTLS11,
		MaxVersion:         tls.VersionTLS12,
		InsecureSkipVerify: false, // keep this false for production
	}

	// log.Printf("TLS configuration loaded with min version: %x", tlsConfig.MinVersion)

	mux := http.NewServeMux()
	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		if len(r.TLS.PeerCertificates) > 0 {
			clientCert := r.TLS.PeerCertificates[0]
			name := html.EscapeString(clientCert.Subject.CommonName)
			if _, err := fmt.Fprintf(w, "Hello, %s!\n", name); err != nil {
				log.Printf("Error writing response: %s", err)
			}
		} else {
			if _, err := fmt.Fprintf(w, "Hello, unknown client!\n"); err != nil {
				log.Printf("Error writing response: %v", err)
			}
		}
	})

	httpServer := &http.Server{
		Addr:              addr,
		Handler:           mux,
		TLSConfig:         tlsConfig,
		ReadHeaderTimeout: 2 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       90 * time.Second,
	}

	return &Server{
		Addr: addr,
		http: httpServer,
	}, nil
}

// Start starts the server (blocking)
// This is standard production behavior for:
// CLI servers
// Containers
// Systemd services
// Kubernetes pods
func (s *Server) Start() error {
	log.Printf("Server listening on %s...", s.Addr)
	if err := s.http.ListenAndServeTLS("", ""); err != nil {
		return fmt.Errorf("server failed: %w", err)
	}
	return nil
}

// StartAsync starts the server in a goroutine and returns an error channel
// StartAsync() is useful only when you do NOT want to block.
// Typical use cases:
// Integration tests
// Embedded servers
// Running multiple servers in one process
// Graceful shutdown handling
// Starting the server + doing other work
func (s *Server) StartAsync() <-chan error {
	errCh := make(chan error, 1)
	go func() {
		errCh <- s.Start()
	}()
	return errCh
}

// Shutdown gracefully shuts down the server with Context
func (s *Server) Shutdown(ctx context.Context) error {
	log.Println("Server is shutting down...")
	return s.http.Shutdown(ctx)
}
