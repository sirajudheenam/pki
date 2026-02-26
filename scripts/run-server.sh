#/bin/bash
HOSTNAME=${2:-localhost}
PORT=${3:-4433}
# Get the directory where the script itself is located
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# Go one directory up
PARENT_DIR="$(dirname "$SCRIPT_DIR")"
PKI_DIR="$PARENT_DIR/demo-pki"
echo "Starting demo TLS server on port 4433..."
echo "Press Ctrl+C to stop the server."
openssl s_server -accept $PORT -www \
  -key $PKI_DIR/$HOSTNAME/server/server.key.pem \
  -cert $PKI_DIR/$HOSTNAME/server/server.chain.pem
