#!/bin/env
curl --fail --cacert ../pki-go/certs/client/inter-root-combined.cert.pem \
     --cert ../pki-go/certs/client/client.cert.pem \
     --key ../pki-go/certs/client/client.key.pem \
     https://go-mtls-server-service:8443/hello