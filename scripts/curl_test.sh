#!/bin/bash
set -e
HOSTNAME=${1:-"localhost"}  #localhost also possible

cd ..
BASE_DIR="demo-pki"         # NEW: base PKI folder
PKI_DIR="$(pwd)/$BASE_DIR/$HOSTNAME"  # NEW: certificates per-host
TEST_CLIENT_DIR="pki-go/certs/$HOSTNAME/client"

if [ -d $PKI_DIR ]; then
     echo "PKI_DIR $PKI_DIR exists"
fi

if [ -d $TEST_CLIENT_DIR ]; then
     echo "TEST_CLIENT_DIR $TEST_CLIENT_DIR exists"
fi


# curl -v --fail --cacert ../pki-go/certs/$HOSTNAME/client/inter-root-combined.cert.pem \
#      --cert ../pki-go/certs/$HOSTNAME/client/client.cert.pem \
#      --key ../pki-go/certs/$HOSTNAME/client/client.key.pem \
#      https://$HOSTNAME:8443/hello

# curl -v -k --cacert ../pki-go/certs/$HOSTNAME/client/inter-root-combined.cert.pem \
#      --cert ../pki-go/certs/$HOSTNAME/client/client.cert.pem \
#      --key ../pki-go/certs/$HOSTNAME/client/client.key.pem \
#      https://$HOSTNAME:8443/hello

# curl -v --cacert $TEST_CLIENT_DIR/inter-root-combined.cert.pem \
#      --cert $TEST_CLIENT_DIR/client.cert.pem \
#      --key $TEST_CLIENT_DIR/client.key.pem \
#      https://$HOSTNAME:8443/hello

curl -v --ipv4 --cacert $TEST_CLIENT_DIR/inter-root-combined.cert.pem \
     --cert $TEST_CLIENT_DIR/client.cert.pem \
     --key $TEST_CLIENT_DIR/client.key.pem \
     https://$HOSTNAME:8443/hello
