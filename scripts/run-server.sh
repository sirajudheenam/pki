#/bin/bash
# Get the directory where the script itself is located
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# Go one directory up
PARENT_DIR="$(dirname "$SCRIPT_DIR")"
PKI_DIR="$PARENT_DIR/demo-pki"
echo "Starting demo TLS server on port 4433..."
echo "Press Ctrl+C to stop the server."
openssl s_server -accept 4433 -www \
  -key $PKI_DIR/server/server.key.pem \
  -cert $PKI_DIR/server/server.chain.pem