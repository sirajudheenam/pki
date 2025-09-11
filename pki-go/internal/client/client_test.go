package client

import (
	"testing"
	"time"

	"github.com/sirajudheenam/pki/pki-go/internal/server"
)

func TestClientRequest(t *testing.T) {
	// Use the local certs path
	certDir := "../../certs/client"

	// Start the server first
	srv, err := server.NewServer(":8445", "../../certs/server")
	if err != nil {
		t.Fatalf("failed to start server: %v", err)
	}
	defer srv.Shutdown()
	go srv.Start()

	// Give server some time to start
	time.Sleep(500 * time.Millisecond)

	// Create a client
	c, err := NewClient("https://localhost:8445", certDir)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Make request
	body, err := c.DoRequest()
	if err != nil {
		t.Fatalf("client request failed: %v", err)
	}

	if body == "" {
		t.Errorf("expected non-empty response from server")
	} else {
		t.Logf("Server response: %s", body)
	}
}
