package main

import (
	"flag"
	"log"
	"os"

	"github.com/sirajudheenam/pki/pki-go/internal/client"
)

func main() {
	// Default values
	defaultServerName := "go-mtls-server-service"
	defaultPort := "8443"
	defaultServerRootPath := "/hello"
	defaultCertPath := "./certs/client"

	// Override defaults with environment variables if set
	if envServerName := os.Getenv("SERVER_NAME"); envServerName != "" {
		log.Println("Using env: SERVER_NAME: ", envServerName)
		defaultServerName = envServerName
	}
	if envServerPort := os.Getenv("SERVER_PORT"); envServerPort != "" {
		log.Println("Using env: SERVER_PORT: ", envServerPort)
		defaultPort = envServerPort
	}
	if envServerRootPath := os.Getenv("SERVER_ROOT_PATH"); envServerRootPath != "" {
		log.Println("Using env: SERVER_ROOT_PATH: ", envServerRootPath)
		defaultServerRootPath = envServerRootPath
	}
	if envCert := os.Getenv("CERT_PATH"); envCert != "" {
		defaultCertPath = envCert
	}
	// Command-line flags (take precedence over env vars)

	serverName := flag.String("server-name", defaultServerName, "Server name")
	serverPort := flag.String("server-port", defaultPort, "Server port")
	serverRootPath := flag.String("server-root-path", defaultServerRootPath, "Server root path")
	certPath := flag.String("cert-path", defaultCertPath, "Path to client certificates")
	flag.Parse()

	serverURL := "https://" + *serverName + ":" + *serverPort + *serverRootPath
	log.Println("Connecting to server at:", serverURL)

	log.Println("Certificate Path: ", *certPath)

	// Create client
	c, err := client.NewClient(serverURL, *certPath)
	if err != nil {
		log.Fatal("\n Failed to create client:", err)
	}

	// Make request
	resp, err := c.DoRequest()
	if err != nil {
		log.Fatal("\n Request failed:", err)
	}

	log.Println("Server Responded with : ", resp)
}
