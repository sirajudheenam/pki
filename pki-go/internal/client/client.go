// pki-go/internal/client/client.go
package client

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// Client wraps an HTTP client
type Client struct {
	Addr string
	http *http.Client
}

// NewClient returns a configured HTTPS client
func NewClient(addr, certDir string) (*Client, error) {

	certFile := filepath.Join(certDir, "client.cert.pem")
	keyFile := filepath.Join(certDir, "client.key.pem")

	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, fmt.Errorf("failed loading client cert/key: %w", err)
	}

	root, err := os.OpenRoot(certDir)
	if err != nil {
		return nil, err
	}

	certDirFile, err := root.Open("inter-root-combined.cert.pem")
	if err != nil {
		log.Fatalf("unable to load inter-root-combined.cert.pem")
	}

	// /* Ensure a proper close and error check */
	defer func() {
		err := certDirFile.Close()
		if err != nil {
			log.Fatalf("Unable to close file: %v", err)
		}
	}()

	caCert, err := io.ReadAll(certDirFile)
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

	client := &Client{
		Addr: addr,
		http: httpClient,
	}

	return client, nil
}

// DoRequest calls /hello endpoint
func (c *Client) DoRequest() (*http.Response, error) {
	resp, err := c.http.Get(c.Addr)
	if err != nil {
		if tlsErr, ok := err.(x509.UnknownAuthorityError); ok {
			log.Printf("Unknown authority: %s", tlsErr.Error())
		} else if tlsErr, ok := err.(x509.HostnameError); ok {
			log.Printf("Hostname error: %s", tlsErr.Error())
		} else {
			log.Printf("General error: %v", err)
		}
		return nil, err
	}

	return resp, nil
}
