package main

import (
	"crypto/tls"
	"crypto/x509"
	"io"
	"net/http"
	"os"
	"testing"
	"time"
)

// helper to start server in a goroutine
func startTestServer(t *testing.T) *http.Server {
	cert, err := tls.LoadX509KeyPair("certs/server/server.chain.pem", "certs/server/server.key.pem")
	if err != nil {
		t.Fatalf("Failed loading server cert/key: %v", err)
	}

	caCert, err := os.ReadFile("certs/server/root.cert.pem")
	if err != nil {
		t.Fatalf("Failed reading root cert: %v", err)
	}

	interCert, err := os.ReadFile("certs/server/intermediate.cert.pem")
	if err != nil {
		t.Fatalf("Failed reading intermediate cert: %v", err)
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

	server := &http.Server{
		Addr:      ":8443",
		TLSConfig: tlsConfig,
	}

	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		if len(r.TLS.PeerCertificates) > 0 {
			clientCert := r.TLS.PeerCertificates[0]
			_, _ = w.Write([]byte("Hello, " + clientCert.Subject.CommonName + "!\n"))
		} else {
			_, _ = w.Write([]byte("Hello, unknown client!\n"))
		}
	})

	go func() {
		if err := server.ListenAndServeTLS("", ""); err != nil && err != http.ErrServerClosed {
			t.Fatalf("Server failed: %v", err)
		}
	}()

	// Give the server a moment to start
	time.Sleep(500 * time.Millisecond)
	return server
}

func TestClientServer(t *testing.T) {
	server := startTestServer(t)
	defer server.Close()

	// Load client cert
	clientCert, err := tls.LoadX509KeyPair("certs/client/client.cert.pem", "certs/client/client.key.pem")
	if err != nil {
		t.Fatalf("Failed loading client cert/key: %v", err)
	}

	caCert, err := os.ReadFile("certs/client/inter-root-combined.cert.pem")
	if err != nil {
		t.Fatalf("Failed reading client CA bundle: %v", err)
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{clientCert},
		RootCAs:            caCertPool,
		InsecureSkipVerify: false,
	}

	client := &http.Client{
		Transport: &http.Transport{TLSClientConfig: tlsConfig},
		Timeout:   5 * time.Second,
	}

	resp, err := client.Get("https://localhost:8443/hello")
	if err != nil {
		t.Fatalf("Client request failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed reading response body: %v", err)
	}

	expected := "Hello, client1!\n"
	if string(body) != expected {
		t.Errorf("Unexpected response: got %q, want %q", string(body), expected)
	}
}
