package server

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/sirajudheenam/pki/pki-go/internal/client"
)

const (
	urlToCheck = "https://go-mtls-server-service:8444/hello"
)

func TestServerStartStop(t *testing.T) {
	// Use the local certs path
	certDir := "../../certs/server"

	// Create a new server
	srv, err := NewServer(":8444", certDir) // use a test port
	if err != nil {
		t.Fatalf("failed to create server: %v", err)
	}

	// Start server asynchronously
	errCh := srv.StartAsync()
	defer func() {
		if err := srv.Shutdown(); err != nil {
			t.Errorf("failed to shutdown server: %v", err)
		}
	}()

	// Wait until server is ready instead of fixed sleep
	waitForServer(t, urlToCheck)

	// Case 1: client with valid certs should succeed
	c, err := client.NewClient(urlToCheck, "../../certs/client")

	if err != nil {
		fmt.Println(" Failed to create client:", err)
	}

	// Make request
	resp, err := c.DoRequest()
	if err != nil {
		t.Fatalf("expected success with valid certs, got error: %v", err)
	}
	defer resp.Body.Close()
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
func waitForServer(t *testing.T, url string) {
	client, err := client.NewClient(url, "../../certs/client")

	if err != nil {
		t.Fatal("Unable to create a client to test ")
	}
	deadline := time.Now().Add(2 * time.Second)

	for time.Now().Before(deadline) {
		resp, err := client.DoRequest()
		if err == nil {
			resp.Body.Close()
			return
		}
		time.Sleep(100 * time.Millisecond)
	}
	t.Fatal("server did not start in time")
}
