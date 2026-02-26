package config

import (
	"fmt"
	"os"
	"path/filepath"

	yaml "gopkg.in/yaml.v3"
)

// Config represents the server configuration
type Config struct {
	Hostname    string `yaml:"hostname"`
	Port        string `yaml:"port"`
	CertBaseDir string `yaml:"cert_base_dir"`
	CertSubDir  string `yaml:"cert_sub_dir"`
}

// LoadConfig loads configuration from environment variables or config file
func LoadConfig() (*Config, error) {
	config := &Config{
		Hostname:    "localhost", // default hostname
		Port:        "8443",      // default port
		CertBaseDir: "certs",     // default base directory for certificates
		CertSubDir:  "server",    // default subdirectory for server certificates
	}

	// Check environment variables first
	if hostname := os.Getenv("HOSTNAME"); hostname != "" {
		config.Hostname = hostname
	}
	if port := os.Getenv("PORT"); port != "" {
		config.Port = port
	}
	if certDir := os.Getenv("CERT_BASE_DIR"); certDir != "" {
		config.CertBaseDir = certDir
	}
	if certSubDir := os.Getenv("CERT_SUB_DIR"); certSubDir != "" {
		config.CertSubDir = certSubDir
	}

	// Try to load config.yaml if it exists
	if configFile, err := os.Open("config.yaml"); err == nil {
		defer configFile.Close()
		decoder := yaml.NewDecoder(configFile)
		if err := decoder.Decode(config); err != nil {
			return nil, err
		}
	}

	fmt.Println("Config: ", config)
	return config, nil
}

// GetCertificatePath returns the full path to the certificate directory
func (c *Config) GetCertificatePath() string {
	fmt.Printf("c.CertBaseDir: %s\n", c.CertBaseDir)
	fmt.Printf("c.Hostname: %s\n", c.Hostname)
	fmt.Printf("c.CertSubDir: %s\n", c.CertSubDir)
	return filepath.Join(c.CertBaseDir, c.Hostname, c.CertSubDir)
}
