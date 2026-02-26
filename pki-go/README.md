# PKI Go Client-Server Application

A Go-based implementation of mutual TLS (mTLS) authentication using PKI certificates. This application demonstrates secure client-server communication with certificate-based authentication.

## 🚀 Quick Start

```bash
# Generate certificates for localhost
make cert-gen

# Or generate for a specific hostname (creates certs in certs/my-server.com/)
make cert-gen HOSTNAME=my-server.com

# Run the server (uses certificates from certs/localhost/ by default)
make run-server

# In another terminal, run the client (automatically uses matching certificates)
make run-client

# Or run with specific hostname
make run-server HOSTNAME=my-server.com
make run-client HOSTNAME=my-server.com
```

## 🛠 Development Setup

### Prerequisites

- Go 1.19 or higher
- Docker (for container builds)
- make
- OpenSSL (for certificate generation)

### Available Make Commands

Use `make help` to see all available commands. Here are the key ones:

```bash
# Generate certificates for a specific hostname
make cert-gen HOSTNAME=my-server.com

# Clean up generated certificates
make cert-clean

# Run tests with coverage
make test

# Format code and run linters
make lint

# Build Docker images
make build

# Complete release process (lint, test, build, tag, push)
make release VERSION=1.1.0
```

### Environment Variables

The following environment variables can be configured:

| Variable | Description | Default |
|----------|-------------|---------|
| HOSTNAME | Target hostname for certificates | localhost |
| VERSION | Image version tag | 1.0.2 |
| SERVER_PORT | Server port | 8443 |

## 🏗️ Project Structure

```
.
├── cmd/
│   ├── client/        # Client application entry point
│   └── server/        # Server application entry point
├── internal/
│   ├── client/        # Client implementation
│   └── server/        # Server implementation
├── certs/                         # Generated certificates (gitignored)
│   ├── localhost/                  # Certificates for localhost
│   │   ├── ca/                    # Certificate Authority
│   │   ├── server/                # Server certificates
│   │   └── client/                # Client certificates
│   └── my-server.com/             # Certificates for other hosts
│       ├── ca/                    # CA for my-server.com
│       ├── server/                # Server certs for my-server.com
│       └── client/                # Client certs for my-server.com
└── Makefile          # Build and development commands
```

## 🐳 Docker Support

### Building Images

```bash
# Build both client and server images
make build

# Tag images as latest
make tag

# Push to registry
make push
```

### Running with Docker Compose

```bash
# Start services
make up

# View logs
make logs

# Stop services
make down
```

## 🔒 Certificate Management

Certificates are automatically managed through make targets and organized by hostname:

```bash
# Generate new certificates (creates in certs/my-server.com/)
make cert-gen HOSTNAME=my-server.com

# Generate for multiple hosts (creates separate directories)
make cert-gen HOSTNAME=dev.example.com
make cert-gen HOSTNAME=prod.example.com

# Clean up all certificates
make cert-clean

# Clean specific hostname certificates
rm -rf certs/my-server.com/

# Certificates are automatically generated when needed:
make run-server HOSTNAME=dev.example.com  # Uses/generates certs/dev.example.com/
make run-client HOSTNAME=dev.example.com  # Uses matching certificates
```

### Certificate Structure

Certificates are organized by hostname in separate directories:

```
certs/
└── {hostname}/           # Directory named after the target hostname
    ├── ca/
    │   ├── ca.crt       # CA certificate
    │   └── ca.key       # CA private key
    ├── server/
    │   ├── server.crt   # Server certificate
    │   └── server.key   # Server private key
    └── client/
        ├── client.crt   # Client certificate
        └── client.key   # Client private key

# Example for multiple hosts:
certs/
├── localhost/           # Certificates for localhost
│   ├── ca/
│   ├── server/
│   └── client/
└── prod.example.com/   # Certificates for production
    ├── ca/
    ├── server/
    └── client/
```

## 🧪 Testing

```bash
# Run all tests with coverage
make test

# Run linting checks
make lint
```

## 📦 Release Process

To create a new release:

1. Update version in Makefile if needed
2. Run tests and create release:
   ```bash
   make release VERSION=1.1.0
   ```
   This will:
   - Run linting checks
   - Run all tests
   - Build Docker images
   - Tag images
   - Push to registry

## 🐛 Troubleshooting

### Common Issues

1. **Certificate Issues**
   ```bash
   # Clean specific hostname certificates
   rm -rf certs/your-hostname/
   
   # Or clean all certificates
   make cert-clean
   
   # Regenerate certificates (creates in certs/your-hostname/)
   make cert-gen HOSTNAME=your-hostname
   
   # Verify certificate structure
   tree certs/your-hostname/
   ```

2. **Port Already in Use**
   ```bash
   # Use a different port
   make run-server SERVER_PORT=9443
   ```

3. **Docker Build Fails**
   ```bash
   # Clean up and rebuild
   make clean
   make build
   ```

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch
3. Run tests and linting:
   ```bash
   make lint
   make test
   ```
4. Submit a pull request

## 📜 License

This project is licensed under the MIT License - see the LICENSE file for details.