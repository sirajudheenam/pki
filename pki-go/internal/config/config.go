package config

import (
	"fmt"
	"os"
	"path/filepath"

	yaml "gopkg.in/yaml.v3"
)

type ServerConfig struct {
	Hostname    string `yaml:"hostname"`
	Port        string `yaml:"port"`
	CertBaseDir string `yaml:"certBaseDir"`
	CertSubDir  string `yaml:"certSubDir"`
}

type ClientConfig struct {
	Hostname    string `yaml:"hostname"`
	Port        string `yaml:"port"`
	CertBaseDir string `yaml:"certBaseDir"`
	CertSubDir  string `yaml:"certSubDir"`
}

// Config represents the server configuration
type Config struct {
	Server ServerConfig `yaml:"server"`
	Client ClientConfig `yaml:"client"`
}

// LoadConfig loads configuration from environment variables or config file
func LoadConfig() (*ServerConfig, error) {
	config := &ServerConfig{
		Hostname:    "localhost", // default hostname
		Port:        "8443",      // default port
		CertBaseDir: "certs",     // default base directory for certificates
		CertSubDir:  "server",    // default subdirectory for server certificates
	}

	fmt.Printf("config (before loading config.yaml) : %+v \n", config)
	// Check environment variables first
	if hostname := os.Getenv("HOSTNAME"); hostname != "" {
		config.Hostname = hostname
	}
	if port := os.Getenv("PORT"); port != "" {
		config.Port = port
	}
	if certDir := os.Getenv("SERVER_CERT_BASE_DIR"); certDir != "" {
		config.CertBaseDir = certDir
	}
	if certSubDir := os.Getenv("SERVER_CERT_SUB_DIR"); certSubDir != "" {
		config.CertSubDir = certSubDir
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
	fmt.Printf("c.CertBaseDir: %s\n", c.CertBaseDir)
	fmt.Printf("c.Hostname: %s\n", c.Hostname)
	fmt.Printf("c.CertSubDir: %s\n", c.CertSubDir)
	return filepath.Join(c.CertBaseDir, c.Hostname, c.CertSubDir)
}
