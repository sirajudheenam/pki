# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go-based PKI (Public Key Infrastructure) learning project that demonstrates mutual TLS (mTLS) authentication between a client and server. The repository includes certificate generation scripts, a Go client-server application, Docker containerization, Kubernetes deployment manifests, and Helm charts.

## Repository Structure

- `pki-go/` - Main Go application with mTLS client-server implementation
  - `cmd/client/` - Client application entry point
  - `cmd/server/` - Server application entry point
  - `internal/client/` - Client implementation with TLS configuration
  - `internal/server/` - Server implementation with TLS configuration
  - `internal/config/` - Configuration loader supporting YAML and environment variables
  - `simple/` - Simplified examples without the internal package structure
- `scripts/` - Certificate generation and helper scripts
  - `create-certs.sh` - Main certificate generation script (hostname-based)
  - `create-certs-k8s.sh` - K8s-specific certificate generation
  - `deploy-k8s.sh`, `cleanup-k8s.sh`, `helm-operations.sh` - K8s helpers
- `k8s/` - Kubernetes deployment manifests
- `helm-chart/` - Helm chart for deploying the application
- `demo-pki/` - Generated certificates organized by hostname (gitignored)

## Certificate Architecture

Certificates are organized by hostname in `demo-pki/{hostname}/`:
- `root/` - Root CA certificates
- `intermediate/` - Intermediate CA certificates
- `server/` - Server certificates (server.chain.pem, server.key.pem)
- `client/` - Client certificates (client.cert.pem, client.key.pem)

The Go application uses hostname-based certificate paths: `certs/{hostname}/server/` or `certs/{hostname}/client/`.

## Development Commands

All commands should be run from the `pki-go/` directory unless otherwise specified.

### Certificate Generation
```bash
# From repo root - generates certs for localhost in demo-pki/localhost/
bash scripts/create-certs.sh all localhost

# For Kubernetes (generates for go-mtls-server-service)
bash scripts/create-certs-k8s.sh

# For custom hostname
bash scripts/create-certs.sh all my-server.com
```

### Building and Testing
```bash
cd pki-go

# Run tests with coverage
make test
# Or run tests without regenerating certs
make test-only

# Run linters
make lint

# Format code
go fmt ./...
```

### Running Locally
```bash
cd pki-go

# Terminal 1 - Start server (uses config.yaml or environment variables)
make run-server
# Or with specific hostname
make run-server HOSTNAME=localhost

# Terminal 2 - Run client
make run-client
# Or with specific hostname
make run-client HOSTNAME=localhost

# Direct go run commands (from pki-go/)
go run cmd/server/main.go
SERVER_NAME=localhost SERVER_PORT=8443 SERVER_ROOT_PATH=/hello CLIENT_CERTS=./certs/localhost/client go run cmd/client/main.go
```

### Docker Operations
```bash
cd pki-go

# Build images
make build VERSION=1.0.0

# Build, tag, and push
make release VERSION=1.0.0

# Run with docker-compose
make up
make logs
make down

# Clean up images
make clean VERSION=1.0.0
```

### Kubernetes Deployment
```bash
# Generate K8s certificates first
bash scripts/create-certs-k8s.sh

# Create secrets
kubectl create secret generic go-mtls-server-certs \
  --from-file=server.chain.pem=./certs/server/server.chain.pem \
  --from-file=server.key.pem=./certs/server/server.key.pem \
  --from-file=root.cert.pem=./certs/server/root.cert.pem \
  --from-file=intermediate.cert.pem=./certs/server/intermediate.cert.pem

kubectl create secret generic go-mtls-client-certs \
  --from-file=client.cert.pem=./certs/client/client.cert.pem \
  --from-file=client.key.pem=./certs/client/client.key.pem \
  --from-file=inter-root-combined.cert.pem=./certs/client/inter-root-combined.cert.pem \
  --from-file=root.cert.pem=./certs/server/root.cert.pem

# Deploy to K8s
cd k8s
kubectl create -f .

# Port forward for local testing (minikube)
kubectl port-forward svc/go-mtls-server-service 8443:8443

# Test with curl (requires /etc/hosts entry: 127.0.0.1 go-mtls-server-service)
curl -vk https://go-mtls-server-service:8443/hello \
  --cert ./certs/client/client.cert.pem \
  --key ./certs/client/client.key.pem \
  --cacert ./certs/client/root.cert.pem
```

### Helm Operations
```bash
cd helm-chart

# Install
helm install go-mtls ./go-mtls

# Uninstall
helm uninstall go-mtls

# Test
helm test go-mtls
```

## Configuration System

The application uses a layered configuration approach:

1. **Defaults** (in code): `localhost`, port `8443`, cert dir `certs`
2. **config.yaml** (optional): Override defaults with YAML
3. **Environment Variables** (highest priority):
   - Server: `HOSTNAME`, `PORT`, `SERVER_CERT_BASE_DIR`, `SERVER_CERT_SUB_DIR`
   - Client: `SERVER_NAME`, `SERVER_PORT`, `SERVER_ROOT_PATH`, `CLIENT_CERTS`

Certificate paths are constructed as: `{CertBaseDir}/{Hostname}/{CertSubDir}/`

## Key Implementation Details

### Server (`internal/server/server.go`)
- Loads server certificate chain (`server.chain.pem`) and private key (`server.key.pem`)
- Loads Root CA (`root.cert.pem`) to verify client certificates
- Configures TLS with `ClientAuthType: tls.RequireAndVerifyClientCert`
- Serves HTTPS on port 8443 with mTLS required

### Client (`internal/client/client.go`)
- Loads client certificate (`client.cert.pem`) and private key (`client.key.pem`)
- Loads combined intermediate and root CA chain (`inter-root-combined.cert.pem`)
- Configures HTTP client with custom TLS transport
- Connects to server with mutual TLS authentication

### Certificate Chain
- Root CA signs Intermediate CA
- Intermediate CA signs Server and Client certificates
- Server uses `server.chain.pem` (server cert + intermediate)
- Client uses separate cert file and verifies against combined root chain

## CI/CD Pipeline

GitHub Actions workflow (`.github/workflows/ci.yml`) runs on push/PR:
1. Generates test certificates for `go-mtls-server-service`
2. Runs tests with coverage (`go test -v -race -coverprofile=coverage.out`)
3. Uploads coverage to Codecov
4. Builds Docker images (tagged with commit SHA)
5. Pushes images to Docker Hub on main branch

## Testing Strategy

Tests require certificates to be generated first. The Makefile `test` target automatically generates certs before running tests. Test files are located alongside implementation:
- `internal/server/server_test.go` - Server tests
- `internal/client/client_test.go` - Client tests
- `simple/server_client_test.go` - Simplified integration tests

## Common Patterns

- All certificate operations use `os.OpenRoot()` for secure directory access
- TLS configuration uses strict cipher suites and TLS 1.2 minimum
- Error handling wraps errors with context using `fmt.Errorf("%w", err)`
- Logging uses standard `log` package with descriptive messages
- File paths are constructed with `filepath.Join()` for cross-platform compatibility

## Troubleshooting

- **Certificate errors**: Regenerate certs for the correct hostname matching your server/client configuration
- **Port in use**: Change `SERVER_PORT` environment variable or kill existing process
- **Docker network issues on macOS**: Add `/etc/hosts` entry for `host.docker.internal` or service hostname
- **K8s port forwarding**: Required on minikube with Docker driver to access services from host
