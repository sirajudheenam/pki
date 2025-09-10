#/bin/bash
PKI_DIR=demo-pki
echo "Running client test..."
echo "Q" | openssl s_client -connect localhost:4433 -CAfile $PKI_DIR/client/inter-root-combined.cert.pem
