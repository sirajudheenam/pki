package client

import (
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/sirajudheenam/pki/pki-go/internal/server"
)

const (
	urlToCheck = "https://go-mtls-server-service:8444/hello"
)

func TestClientRequest(t *testing.T) {
	// Use the local certs path
	serverCertDir := "../../certs/server"
	clientCertDir := "../../certs/client"

	// Create a new server
	srv, err := server.NewServer(":8444", serverCertDir) // use a test port
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
	defer resp.Body.Close()
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
	client, err := NewClient(url, "../../certs/client")

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
