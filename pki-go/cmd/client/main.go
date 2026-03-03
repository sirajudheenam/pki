// pki-go/cmd/client/main.go
package main

import (
	"flag"
	"io"
	"log"
	"os"
	"strings"

	internalClient "github.com/sirajudheenam/pki/pki-go/internal/client"
	"github.com/sirajudheenam/pki/pki-go/internal/config"
)

func main() {

	// read config.yaml file and assign
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Printf("error loading config from file: %v", err)
	}
	defaultServerName := cfg.Client.Hostname
	defaultPort := cfg.Client.Port
	defaultServerRootPath := cfg.Client.ServerRootPath
	defaultClientCertPath := cfg.Client.GetCertificatePath()

	cwd, err := os.Getwd()

	if err != nil {
		log.Fatal(err)
	}

	// Override config defaults with environment variables if set
	envServerName := os.Getenv("SERVER_NAME")
	if envServerName != "" {
		safe := strings.ReplaceAll(envServerName, "\n", "")
		log.Printf("Using env: SERVER_NAME: %q", safe)
		defaultServerName = safe
	}
	envServerPort := os.Getenv("SERVER_PORT")
	if envServerPort != "" {
		safe := strings.ReplaceAll(envServerPort, "\n", "")
		log.Printf("Using env: SERVER_PORT: %q", safe)
		defaultPort = safe
	}
	envServerRootPath := os.Getenv("SERVER_ROOT_PATH")
	if envServerRootPath != "" {
		safe := strings.ReplaceAll(envServerRootPath, "\n", "")
		log.Printf("Using env: SERVER_ROOT_PATH: %q", safe)
		defaultServerRootPath = safe
	}
	envCertPath := os.Getenv("CERT_BASE_DIR")
	if envCertPath != "" {
		safe := strings.ReplaceAll(envCertPath, "\n", "")
		log.Printf("Using env: CERT_BASE_DIR: %q", safe)
	}
	envSubPath := os.Getenv("CLIENT_CERT_SUB_DIR")
	if envSubPath != "" {
		safe := strings.ReplaceAll(envSubPath, "\n", "")
		log.Printf("Using env: CLIENT_CERT_SUB_DIR: %q", safe)
	}
	if envCertPath != "" && envServerName != "" && envSubPath != "" {
		defaultClientCertPath = strings.Join([]string{cwd, "/", envCertPath, "/", envServerName, "/", envSubPath}, "")
	}

	// Default values for prod
	if defaultServerName == "" {
		defaultServerName = "go-mtls-server-service"
	}

	if defaultPort == "" {
		defaultPort = "8443"
	}

	if defaultServerRootPath == "" {
		defaultServerRootPath = "/hello"
	}

	if defaultClientCertPath == "" {
		defaultClientCertPath = strings.Join([]string{cwd, "/", "certs", "/", defaultServerName, "/", "client"}, "")
	}

	// Command-line flags (take precedence over env vars)

	serverName := flag.String("server-name", defaultServerName, "Server name")
	serverPort := flag.String("server-port", defaultPort, "Server port")
	serverRootPath := flag.String("server-root-path", defaultServerRootPath, "Server root path")
	certPath := flag.String("cert-path", defaultClientCertPath, "Path to client certificates")
	flag.Parse()

	serverURL := strings.Join([]string{"https://", *serverName, ":", *serverPort, *serverRootPath}, "")

	log.Printf("Connecting to server at: %q", serverURL)

	// Version 1
	// Create client
	c, err := internalClient.NewClient(serverURL, *certPath)
	if err != nil {
		log.Fatalf("\n Failed to create client: %q", err)
	}

	// Make request
	resp, err := c.DoRequest()
	if err != nil {
		log.Fatal("\n Request failed:", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("\n 1 Failed to read response body:", err)
	}
	if err := resp.Body.Close(); err != nil {
		log.Fatal("\n Failed to close response body:", err)
	}

	log.Printf("1. Server responded with: %q", string(body))

	// Version 2
	clientConfig := &config.ClientConfig{
		Hostname:       defaultServerName,
		Port:           defaultPort,
		CertBaseDir:    "certs",
		CertSubDir:     "client",
		ServerRootPath: "/hello",
	}

	_ = certPath
	c2, err := internalClient.NewClientV2(clientConfig)
	if err != nil {
		log.Fatalf("unable to create a client : %v", err)
	}

	resp2, err := c2.DoRequest()
	if err != nil {
		log.Fatal("\n 2 Failed to read response body:", err)
	}

	body2, err := io.ReadAll(resp2.Body)
	if err != nil {
		log.Fatal("\n Failed to read response body:", err)
	}
	log.Printf("2. Server responded with: %q", string(body2))
	defer func() {
		if err := resp2.Body.Close(); err != nil {
			log.Fatal("\n Failed to close response body:", err)
		}
	}()
}
