package server

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net/http"
	"os"
)

// Server holds the HTTP server
type Server struct {
	Addr string
	http *http.Server
}

// New returns a configured server
func NewServer(addr, certDir string) (*Server, error) {
	cert, err := tls.LoadX509KeyPair(certDir+"/server.chain.pem", certDir+"/server.key.pem")
	if err != nil {
		return nil, fmt.Errorf("failed loading server cert/key: %w", err)
	}

	caCert, err := os.ReadFile(certDir + "/root.cert.pem")
	if err != nil {
		return nil, err
	}
	interCert, err := os.ReadFile(certDir + "/intermediate.cert.pem")
	if err != nil {
		return nil, err
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	caCertPool.AppendCertsFromPEM(interCert)

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientCAs:    caCertPool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
		MinVersion:   tls.VersionTLS12,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		if len(r.TLS.PeerCertificates) > 0 {
			clientCert := r.TLS.PeerCertificates[0]
			fmt.Fprintf(w, "Hello, %s!\n", clientCert.Subject.CommonName)
		} else {
			fmt.Fprintf(w, "Hello, unknown client!\n")
		}
	})

	httpServer := &http.Server{
		Addr:      addr,
		Handler:   mux,
		TLSConfig: tlsConfig,
	}

	return &Server{
		Addr: addr,
		http: httpServer,
	}, nil
}

// Start starts the server (blocking)
func (s *Server) Start() error {
	log.Println("Server listening on", s.Addr)
	return s.http.ListenAndServeTLS("", "")
}

// StartAsync starts the server in a goroutine and returns an error channel
func (s *Server) StartAsync() <-chan error {
	errCh := make(chan error, 1)
	go func() {
		errCh <- s.Start()
	}()
	return errCh
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown() error {
	return s.http.Close()
}
