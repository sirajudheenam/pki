# Go mTLS Helm Chart

This Helm chart deploys the Go mTLS server and client applications in Kubernetes.

## Prerequisites

- Kubernetes 1.19+
- Helm 3.0+
- A Kubernetes Ingress Controller (for server ingress)
- cert-manager (optional, for automatic certificate management)

## Installing the Chart

```bash
# Add local repository
helm package .
helm install go-mtls ./go-mtls-0.1.1.tgz

# With custom values
helm install go-mtls ./go-mtls-0.1.1.tgz -f my-values.yaml
```

## Configuration

### Global Parameters

| Parameter | Description | Default |
|-----------|-------------|---------|
| `nameOverride` | Override the name of the chart | `""` |
| `fullnameOverride` | Override the full name of the resources | `""` |

### Server Parameters

| Parameter | Description | Default |
|-----------|-------------|---------|
| `server.image` | Server image repository | `sirajudheenam/go-mtls-server` |
| `server.tag` | Server image tag | `latest` |
| `server.replicas` | Number of server replicas | `1` |
| `server.resources` | Server pod resources | `{}` |
| `server.certs.server_cert_pem` | Server certificate in PEM format | `""` |
| `server.certs.server_key_pem` | Server private key in PEM format | `""` |

### Client Parameters

| Parameter | Description | Default |
|-----------|-------------|---------|
| `client.image` | Client image repository | `sirajudheenam/go-mtls-client` |
| `client.tag` | Client image tag | `latest` |
| `client.replicas` | Number of client replicas | `1` |
| `client.resources` | Client pod resources | `{}` |
| `client.certs.client_cert_pem` | Client certificate in PEM format | `""` |
| `client.certs.client_key_pem` | Client private key in PEM format | `""` |

## Certificate Management

### Manual Certificate Configuration

```yaml
server:
  certs:
    server_cert_pem: |
      -----BEGIN CERTIFICATE-----
      ...
      -----END CERTIFICATE-----
    server_key_pem: |
      -----BEGIN PRIVATE KEY-----
      ...
      -----END PRIVATE KEY-----

client:
  certs:
    client_cert_pem: |
      -----BEGIN CERTIFICATE-----
      ...
      -----END CERTIFICATE-----
    client_key_pem: |
      -----BEGIN PRIVATE KEY-----
      ...
      -----END PRIVATE KEY-----
```

### Using cert-manager

```yaml
server:
  certManager:
    enabled: true
    issuerName: my-issuer
    issuerKind: ClusterIssuer
```

## Examples

### Basic Installation
```bash
helm install go-mtls . --set server.replicas=2
```

### Custom Configuration
```bash
helm install go-mtls . -f values.production.yaml
```

### Testing

```bash
# Run helm tests
helm test go-mtls

# Manual testing
kubectl port-forward svc/go-mtls-server 8443:8443
curl --cacert ca.crt --cert client.crt --key client.key https://localhost:8443/hello
```

## Monitoring

The server exposes Prometheus metrics at `/metrics`. Configure your ServiceMonitor:

```yaml
server:
  metrics:
    enabled: true
    serviceMonitor:
      enabled: true
```

## Troubleshooting

1. Check certificate secrets:
   ```bash
   kubectl get secrets -l app.kubernetes.io/name=go-mtls
   ```

2. Verify client-server communication:
   ```bash
   kubectl logs -l app.kubernetes.io/name=go-mtls-client
   ```

3. Test server connectivity:
   ```bash
   kubectl exec -it $(kubectl get pod -l app.kubernetes.io/name=go-mtls-server -o name | head -1) -- curl -k https://localhost:8443/health
   ```