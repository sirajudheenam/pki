#!/usr/bin/env bash
set -euo pipefail

sed -i -e "s/<BASE64_ENCODED_CLIENT_CERT_PEM>/$(base64 -w0 -i pki-go/certs/client/client.cert.pem)/" helm-chart/go-mtls/values2.yaml
sed -i -e "s/<BASE64_ENCODED_CLIENT_KEY_PEM>/$(base64 -w0 -i pki-go/certs/client/client.key.pem)/" helm-chart/go-mtls/values2.yaml
sed -i -e "s/<BASE64_ENCODED_INTER_ROOT_CERT_PEM>/$(base64 -w0 -i pki-go/certs/client/inter-root-combined.cert.pem)/" helm-chart/go-mtls/values2.yaml

sed -i -e "s/<BASE64_ENCODED_CHAIN_PEM>/$(base64 -w0 -i pki-go/certs/server/server.chain.pem)/" helm-chart/go-mtls/values2.yaml
sed -i -e "s/<BASE64_ENCODED_KEY_PEM>/$(base64 -w0 -i pki-go/certs/server/server.key.pem)/" helm-chart/go-mtls/values2.yaml
sed -i -e "s/<BASE64_ENCODED_ROOT_PEM>/$(base64 -w0 -i pki-go/certs/server/root.cert.pem)/" helm-chart/go-mtls/values2.yaml
sed -i -e "s/<BASE64_ENCODED_INTER_PEM>/$(base64 -w0 -i pki-go/certs/server/intermediate.cert.pem)/" helm-chart/go-mtls/values2.yaml

rm helm-chart/go-mtls/values2.yaml-e
# cat helm-chart/go-mtls/values2.yaml

# # Original values2.yaml file
# client:
#   image: sirajudheenam/go-mtls-client:latest
#   replicas: 1
#   certs:
#     client.cert.pem: "<BASE64_ENCODED_CLIENT_CRT>"
#     client.key.pem: "<BASE64_ENCODED_CLIENT_KEY>"
#     inter-root-combined.cert.pem: "<BASE64_ENCODED_CA_CRT>"
# server:
#   image: sirajudheenam/go-mtls-server:1.0.0
#   replicas: 2
#   servicePort: 8443
#   nodePort: 30443
#   ingressHost: go-mtls-server.local
#   certs:
#     server.chain.pem: "<BASE64_ENCODED_CHAIN_PEM>"
#     server.key.pem: "<BASE64_ENCODED_KEY_PEM>"
#     root.cert.pem: "<BASE64_ENCODED_ROOT_PEM>"
#     intermediate.cert.pem: "<BASE64_ENCODED_INTER_PEM>"
