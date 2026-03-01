package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"path/filepath"
	"strings"
	"time"

	internalClient "github.com/sirajudheenam/pki/pki-go/internal/client"
)

func main() {
	// Default values
	defaultServerName := "go-mtls-server-service"
	defaultPort := "8443"
	defaultServerRootPath := "/hello"

	cwd, err := os.Getwd()

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(cwd)
	defaultCertPath := strings.Join([]string{cwd + "/certs/" + defaultServerName + "/client/"}, "")
	// log.Printf("1. defaultCertPath %s", defaultCertPath)
	// Override defaults with environment variables if set
	if envServerName := os.Getenv("SERVER_NAME"); envServerName != "" {
		safe := strings.ReplaceAll(envServerName, "\n", "")
		// log.Printf("Using env: SERVER_NAME: %q", safe)
		defaultServerName = safe
		defaultCertPath = strings.Join([]string{cwd + "/certs/" + defaultServerName + "/client/"}, "")
	}
	if envServerPort := os.Getenv("SERVER_PORT"); envServerPort != "" {
		safe := strings.ReplaceAll(envServerPort, "\n", "")
		// log.Printf("Using env: SERVER_PORT: %q", safe)
		defaultPort = safe
	}
	if envServerRootPath := os.Getenv("SERVER_ROOT_PATH"); envServerRootPath != "" {
		safe := strings.ReplaceAll(envServerRootPath, "\n", "")
		// log.Printf("Using env: SERVER_ROOT_PATH: %q", safe)
		defaultServerRootPath = safe
	}
	if envCertPath := os.Getenv("CLIENT_CERT_BASE_DIR"); envCertPath != "" {
		safe := strings.ReplaceAll(envCertPath, "\n", "")
		// log.Printf("Using env: CLIENT_CERT_BASE_DIR: %q", safe)
		defaultCertPath = strings.Join([]string{cwd, "/", safe}, "")
	}
	// log.Printf("2. defaultCertPath %s", defaultCertPath)
	// Command-line flags (take precedence over env vars)

	serverName := flag.String("server-name", defaultServerName, "Server name")
	serverPort := flag.String("server-port", defaultPort, "Server port")
	serverRootPath := flag.String("server-root-path", defaultServerRootPath, "Server root path")
	certPath := flag.String("cert-path", defaultCertPath, "Path to client certificates")
	flag.Parse()

	serverURL := strings.Join([]string{"https://", *serverName, ":", *serverPort, *serverRootPath}, "")
	// log.Printf("Connecting to server at: %q", serverURL)
	// log.Printf("Certificate Path: %q", *certPath)

	// Additional starts here. // Can be safely REMOVED

	// log.Println("Additional starts here.")
	// fmt.Println("Additional starts here.")

	certFile := filepath.Join(*certPath, "client.cert.pem")
	keyFile := filepath.Join(*certPath, "client.key.pem")
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		log.Fatalf("Failed loading client cert/key: %v", err)
	}
	// log.Println("Successfully loaded client certificate and key.")

	root, err := os.OpenRoot(*certPath)
	if err != nil {
		log.Fatalf("Unable Load *certPath")
	} else {
		log.Printf("OpenRoot: %v", &root)
	}

	certDirFile, err := root.Open("inter-root-combined.cert.pem")
	if err != nil {
		log.Fatalf("unable to load inter-root-combined.cert.pem")
	} else {
		log.Println("loaded inter-root-combined.cert.pem")
	}

	caCert, err := io.ReadAll(certDirFile)
	if err != nil {
		log.Printf("Unable to load caCert")
	} else {
		log.Println("caCert loaded")
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// log.Printf("caCertPool: %v", &caCertPool)
	// log.Println("Successfully appended CA certificate to pool.")

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
		MinVersion:   tls.VersionTLS12,
		MaxVersion:   tls.VersionTLS13,
	}
	// log.Println("TLS versions set between 1.2 and 1.3.")

	// Verify Server Availability: (Additional)
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(*serverName, *serverPort), 5*time.Second)
	if err != nil {
		log.Printf("Server not available: %v\n", err)
		return
	}
	log.Println("Server is available.")
	conn.Close()

	// Request  Dump
	req, err := http.NewRequest("GET", serverURL, nil)
	if err != nil {
		log.Fatalf("Failed to create HTTP request: %v", err)
	}
	dump, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		log.Fatalf("Failed to dump HTTP request: %v", err)
	}
	log.Printf("HTTP Request (DUMP): %s\n", dump)

	client := &http.Client{Transport: &http.Transport{TLSClientConfig: tlsConfig}}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("HTTP request failed: %v", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("Unable to Close the response body: %v", err)
		}
	}()

	log.Printf("Response Status: %s", resp.Status)
	// If `c.DoRequest` internally constructs its own request,
	// move the dumping logic there.
	// log.Println("Additional ends here.")
	// fmt.Println("Additional ends here.")

	// Additional ends here. // Can be safely REMOVED

	// Create client
	c, err := internalClient.NewClient(serverURL, *certPath)
	if err != nil {
		log.Fatalf("\n Failed to create client: %q", err)
	}

	// Make request
	respo, err := c.DoRequest()
	if err != nil {
		log.Fatal("\n Request failed:", err)
	}

	body, err := io.ReadAll(respo.Body)
	if err != nil {
		log.Fatal("\n Failed to read response body:", err)
	}
	if err := respo.Body.Close(); err != nil {
		log.Fatal("\n Failed to close response body:", err)
	}

	log.Printf("Server responded with: %q", string(body))
}
