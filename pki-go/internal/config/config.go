package config

import (
	"fmt"
	"os"
	"path/filepath"

	yaml "gopkg.in/yaml.v3"
)

type ServerConfig struct {
	Hostname       string `yaml:"hostname"`
	Port           string `yaml:"port"`
	CertBaseDir    string `yaml:"certBaseDir"`
	CertSubDir     string `yaml:"certSubDir"`
	ServerRootPath string `yaml:"serverRootPath"`
}

type ClientConfig struct {
	Hostname       string `yaml:"hostname"`
	Port           string `yaml:"port"`
	CertBaseDir    string `yaml:"certBaseDir"`
	CertSubDir     string `yaml:"certSubDir"`
	ServerRootPath string `yaml:"serverRootPath"`
}

// Config represents the server configuration
type Config struct {
	Server ServerConfig `yaml:"server"`
	Client ClientConfig `yaml:"client"`
}

// LoadServerConfig loads server configuration from environment variables or config file
func LoadConfig() (*Config, error) {
	serverConfig := &ServerConfig{
		Hostname:       "localhost", // default hostname
		Port:           "8443",      // default port
		CertBaseDir:    "certs",     // default base directory for certificates
		CertSubDir:     "server",    // default subdirectory for server certificates
		ServerRootPath: "/hello",    // default handler path
	}
	clientConfig := &ClientConfig{
		Hostname:       "localhost", // default hostname
		Port:           "8443",      // default port
		CertBaseDir:    "certs",     // default base directory for certificates
		CertSubDir:     "server",    // default subdirectory for server certificates
		ServerRootPath: "/hello",    // default handler path
	}
	config := &Config{
		Server: *serverConfig,
		Client: *clientConfig,
	}

	fmt.Printf("config (before loading config.yaml) : %+v \n", config)
	// Check environment variables first
	if hostname := os.Getenv("HOSTNAME"); hostname != "" {
		config.Server.Hostname = hostname
	}
	if port := os.Getenv("PORT"); port != "" {
		config.Server.Port = port
	}
	if certDir := os.Getenv("CERT_BASE_DIR"); certDir != "" {
		config.Server.CertBaseDir = certDir
	}
	if certSubDir := os.Getenv("SERVER_CERT_SUB_DIR"); certSubDir != "" {
		config.Server.CertSubDir = certSubDir
	}

	// Try to load config.yaml if it exists
	if configFile, err := os.Open("config.yaml"); err == nil {
		defer func() {
			if err := configFile.Close(); err != nil {
				fmt.Printf("Error closing config file: %v\n", err)
			}
		}()
		decoder := yaml.NewDecoder(configFile)
		if err := decoder.Decode(config); err != nil {
			return nil, err
		}
	}

	fmt.Printf("Config: %+v\n", config)
	fmt.Printf("Config Type: %T\n", config)
	return config, nil
}

// GetCertificatePath returns the full path to the certificate directory
func (c *ServerConfig) GetCertificatePath() string {
	return filepath.Join(c.CertBaseDir, c.Hostname, c.CertSubDir)
}

// GetCertificatePath returns the full path to the certificate directory
func (c *ClientConfig) GetCertificatePath() string {
	return filepath.Join(c.CertBaseDir, c.Hostname, c.CertSubDir)
}
