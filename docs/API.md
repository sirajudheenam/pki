# PKI Server API Documentation

## Overview

This document describes the API endpoints provided by the PKI server for managing certificates and performing TLS operations.

## Authentication

All endpoints require mutual TLS (mTLS) authentication. Clients must present a valid client certificate signed by the trusted CA.

## Endpoints

### Health Check
```
GET /health
```
Returns the health status of the server.

### Certificate Status
```
GET /api/v1/cert/status
```
Returns the status of the server's certificate including expiry information.

### Certificate Renewal
```
POST /api/v1/cert/renew
```
Triggers certificate renewal process.

### Metrics
```
GET /metrics
```
Prometheus metrics endpoint for monitoring server health and certificate status.

## Error Responses

All errors follow this format:
```json
{
  "error": {
    "code": "string",
    "message": "string",
    "details": {}
  }
}
```

## Certificate Operations

### Certificate Lifecycle
1. Initial Certificate Generation
2. Certificate Validation
3. Regular Rotation
4. Expiry Monitoring
5. Revocation (if needed)

## Monitoring

The server exposes the following metrics:
- Certificate expiry time
- TLS handshake success/failure rates
- Connection duration
- Error counts

## Security Considerations

1. Keep private keys secure
2. Monitor certificate expiry
3. Implement certificate rotation
4. Use secure cipher suites
5. Regular security audits