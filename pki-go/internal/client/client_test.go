package client

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/sirajudheenam/pki/pki-go/internal/server"
)

const (
	urlToCheck = "https://localhost:8444/hello"
)

const (
	testHost = "localhost" // Using existing certificates from certs/localhost
	testPort = "8444"
)

func setupTestCertificates(t *testing.T) (string, string) {
	// Use existing certificates
	serverCertDir := filepath.Join("../../certs", testHost, "server")
	clientCertDir := filepath.Join("../../certs", testHost, "client")

	// Verify that the certificate directories exist
	if _, err := os.Stat(serverCertDir); os.IsNotExist(err) {
		t.Fatalf("Server certificate directory not found at %s", serverCertDir)
	}
	if _, err := os.Stat(clientCertDir); os.IsNotExist(err) {
		t.Fatalf("Client certificate directory not found at %s", clientCertDir)
	}
	return serverCertDir, clientCertDir
}

// No cleanup needed since we're using existing certificates

func TestClientRequest(t *testing.T) {
	// Use existing certificates
	serverCertDir, clientCertDir := setupTestCertificates(t)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	// Create a new server
	srv, err := server.NewServer(":"+testPort, serverCertDir)
	if err != nil {
		t.Fatalf("failed to create server: %v", err)
	}

	// Start server asynchronously
	errCh := srv.StartAsync()
	defer func() {
		if err := srv.Shutdown(ctx); err != nil {
			t.Errorf("failed to shutdown server: %v", err)
		}
	}()

	// Wait until server is ready instead of fixed sleep
	waitForServer(t, urlToCheck)

	// Create a client
	c, err := NewClient(urlToCheck, clientCertDir)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
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
		t.Errorf("Unable to read Response body %s", err.Error())
	}

	if len(body) == 0 {
		t.Errorf("expected non-empty response from server")
	}
	expected := "Hello, client1!\n"
	if string(body) != expected {
		t.Errorf("expected %q, got %q", expected, string(body))
	} else {
		t.Logf("Server response: %s", string(body))
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
func waitForServer(t *testing.T, url string) {
	fmt.Printf("Waiting for server at %s...\n", url)
	// extract hostname from url
	hostname := strings.Split(url, ":")[1]
	client, err := NewClient(url, "../../certs/"+hostname+"/client")

	if err != nil {
		t.Fatal("Unable to create a client to test ")
	}
	deadline := time.Now().Add(2 * time.Second)

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
