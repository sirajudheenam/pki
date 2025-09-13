package client

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"net/http"
	"os"
)

// Client wraps an HTTP client
type Client struct {
	Addr string
	http *http.Client
}

// NewClient returns a configured HTTPS client
func NewClient(addr, certDir string) (*Client, error) {
	cert, err := tls.LoadX509KeyPair(certDir+"/client.cert.pem", certDir+"/client.key.pem")
	if err != nil {
		return nil, fmt.Errorf("failed loading client cert/key: %w", err)
	}

	caCert, err := os.ReadFile(certDir + "/inter-root-combined.cert.pem")
	if err != nil {
		return nil, err
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		RootCAs:            caCertPool,
		InsecureSkipVerify: false,
	}

	httpClient := &http.Client{
		Transport: &http.Transport{TLSClientConfig: tlsConfig},
	}

	return &Client{
		Addr: addr,
		http: httpClient,
	}, nil
}

// DoRequest calls /hello endpoint
func (c *Client) DoRequest() (string, error) {
	// resp, err := c.http.Get(c.Addr + "/hello")
	resp, err := c.http.Get(c.Addr)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	return string(body), err
}
