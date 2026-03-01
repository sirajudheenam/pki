package server

import (
	"context"
	"crypto/tls"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/sirajudheenam/pki/pki-go/internal/client"
	"github.com/sirajudheenam/pki/pki-go/internal/config"
)

const (
	testPort   = "8444"
	testHost   = "localhost" // Using existing certificates from certs/localhost
	urlToCheck = "https://" + testHost + ":" + testPort + "/hello"
)

func setupTestCertificates(t *testing.T) string {
	// Use existing certificates
	cfg := &config.ServerConfig{
		Hostname:    testHost,
		Port:        testPort,
		CertBaseDir: "../../certs",
		CertSubDir:  "server",
	}

	// Verify that the certificate directory exists
	certPath := cfg.GetCertificatePath()
	if _, err := os.Stat(certPath); os.IsNotExist(err) {
		t.Fatalf("Server certificate directory not found at %s", certPath)
	}

	return certPath
}

// No cleanup needed since we're using existing certificates

func TestServerStartStop(t *testing.T) {
	// Set up test certificates
	certDir := setupTestCertificates(t)

	// Create a new server
	srv, err := NewServer(":"+testPort, certDir)
	if err != nil {
		t.Fatalf("failed to create server: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	// Start server asynchronously
	errCh := srv.StartAsync()
	defer func() {
		if err := srv.Shutdown(ctx); err != nil {
			t.Errorf("failed to shutdown server: %v", err)
		}
	}()

	// Wait until server is ready
	waitForServer(t)

	// Case 1: client with valid certs should succeed
	clientCertPath := filepath.Join("../../certs", testHost, "client")
	c, err := client.NewClient(urlToCheck, clientCertPath)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Make request
	resp, err := c.DoRequest()
	if err != nil {
		t.Fatalf("expected success with valid certs, got error: %v", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			t.Errorf("failed to close response body: %v", err)
		}
	}()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200 OK, got %v", resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Unable to read Response Body")
	}

	expected := "Hello, client1!\n"
	if string(body) != expected {
		t.Errorf("expected %q, got %q", expected, string(body))
	}

	// Case 2: client without certs should fail
	insecureClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // ignore trust
		},
		Timeout: 2 * time.Second,
	}
	_, err = insecureClient.Get(urlToCheck)
	if err == nil {
		t.Errorf("expected handshake error without client cert, but got success")
	}

	// Ensure server exits cleanly
	select {
	case e := <-errCh:
		if e != nil {
			t.Logf("server exited with error: %v", e)
		}
	case <-time.After(1 * time.Second):
		t.Log("server did not exit (expected if still running)")
	}
}

// waitForServer retries until the server is reachable
func waitForServer(t *testing.T) {
	clientCertPath := filepath.Join("../../certs", testHost, "client")
	client, err := client.NewClient(urlToCheck, clientCertPath)
	if err != nil {
		t.Fatalf("Unable to create test client: %v", err)
	}
	deadline := time.Now().Add(5 * time.Second)

	for time.Now().Before(deadline) {
		resp, err := client.DoRequest()
		if err == nil {
			if err := resp.Body.Close(); err != nil {
				t.Errorf("failed to close response body: %v", err)
			}
			return
		}
		time.Sleep(100 * time.Millisecond)
	}
	t.Fatal("server did not start in time")
}
