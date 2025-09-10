// client.go
package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	// Load client cert + key
	cert, err := tls.LoadX509KeyPair("certs/client.cert.pem", "certs/client.key.pem")
	if err != nil {
		log.Fatal("Failed loading client cert/key:", err)
	}

	// Load Root CA (to trust server’s cert)
	caCert, err := os.ReadFile("certs/inter-root-combined.cert.pem")
	if err != nil {
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// TLS config
	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		RootCAs:            caCertPool,
		InsecureSkipVerify: false,
	}

	client := &http.Client{
		Transport: &http.Transport{TLSClientConfig: tlsConfig},
	}

	// Call server
	resp, err := client.Get("https://localhost:8443/hello")
	if err != nil {
		log.Fatal("Request failed:", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Println("Server response:", string(body))
}
