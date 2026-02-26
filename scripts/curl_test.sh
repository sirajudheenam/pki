#!/bin/bash
HOSTNAME=${1:-"go-mtls-server-service"}  #localhost also possible
# curl -v --fail --cacert ../pki-go/certs/$HOSTNAME/client/inter-root-combined.cert.pem \
#      --cert ../pki-go/certs/$HOSTNAME/client/client.cert.pem \
#      --key ../pki-go/certs/$HOSTNAME/client/client.key.pem \
#      https://$HOSTNAME:8443/hello

curl -v -k --cacert ../pki-go/certs/$HOSTNAME/client/inter-root-combined.cert.pem \
     --cert ../pki-go/certs/$HOSTNAME/client/client.cert.pem \
     --key ../pki-go/certs/$HOSTNAME/client/client.key.pem \
     https://$HOSTNAME:8443/hello