#/bin/bash
# Get the directory where the script itself is located
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# Go one directory up
PARENT_DIR="$(dirname "$SCRIPT_DIR")"
PKI_DIR="$PARENT_DIR/demo-pki"
echo "Running client test..."
echo "Q" | openssl s_client -connect localhost:4433 -CAfile $PKI_DIR/client/inter-root-combined.cert.pem