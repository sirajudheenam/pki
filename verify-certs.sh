#/bin/bash
PKI_DIR=demo-pki

openssl verify -CAfile $PKI_DIR/root/root.cert.pem $PKI_DIR/root/root.cert.pem
openssl verify -CAfile $PKI_DIR/root/root.cert.pem $PKI_DIR/intermediate/intermediate.cert.pem
openssl verify -CAfile $PKI_DIR/root/root.cert.pem -untrusted $PKI_DIR/intermediate/intermediate.cert.pem $PKI_DIR/client/client.cert.pem
openssl verify -CAfile $PKI_DIR/root/root.cert.pem -untrusted $PKI_DIR/intermediate/intermediate.cert.pem $PKI_DIR/server/server.cert.pem
openssl verify -CAfile $PKI_DIR/root/root.cert.pem -untrusted $PKI_DIR/intermediate/intermediate.cert.pem $PKI_DIR/client/inter-root-combined.cert.pem
echo "Verifying against combined root+intermediate cert:"
openssl verify -CAfile $PKI_DIR/root/ca-bundle.pem $PKI_DIR/server/server.chain.pem
openssl verify -CAfile $PKI_DIR/root/ca-bundle.pem $PKI_DIR/server/server.cert.pem
openssl verify -CAfile $PKI_DIR/root/ca-bundle.pem $PKI_DIR/client/client.cert.pem
openssl verify -CAfile $PKI_DIR/root/ca-bundle.pem $PKI_DIR/client/inter-root-combined.cert.pem