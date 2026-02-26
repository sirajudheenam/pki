#/bin/bash
HOSTNAME=${2:-localhost}
PORT=${3:-4433}
# Get the directory where the script itself is located
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# Go one directory up
PARENT_DIR="$(dirname "$SCRIPT_DIR")"
PKI_DIR="$PARENT_DIR/demo-pki"
echo "Running client test..."
echo "Q" | openssl s_client -connect $HOSTNAME:$PORT -CAfile $PKI_DIR/$HOSTNAME/client/inter-root-combined.cert.pem