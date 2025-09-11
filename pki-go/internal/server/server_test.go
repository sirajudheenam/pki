package server

import (
	"net/http"
	"testing"
	"time"
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
	defer srv.Shutdown()

	// Give server some time to start
	time.Sleep(500 * time.Millisecond)

	// Try a simple GET request using default client with InsecureSkipVerify
	resp, err := http.Get("https://localhost:8444/hello")
	if err == nil {
		defer resp.Body.Close()
		t.Errorf("expected TLS handshake error without client cert, got %v", resp.Status)
	}

	// Stop server
	if err := srv.Shutdown(); err != nil {
		t.Fatalf("failed to shutdown server: %v", err)
	}

	select {
	case e := <-errCh:
		// server exited, should be nil or TLS handshake errors
		t.Logf("server exited: %v", e)
	default:
	}
}
