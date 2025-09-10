#!/bin/bash
# create-certs.sh

set -e

# Demo PKI folder
PKI_DIR=demo-pki

# Remove if folder exists
rm -rf $PKI_DIR 

mkdir -p $PKI_DIR/{root,intermediate,server,crl,private,newcerts,client}
chmod 700 $PKI_DIR/private
touch $PKI_DIR/index.txt
echo 1000 > $PKI_DIR/serial

# Create OpenSSL config with proper extensions
cat > $PKI_DIR/openssl.cnf <<EOF
[ ca ]
default_ca = CA_default

[ CA_default ]
dir               = .
database          = \$dir/index.txt
new_certs_dir     = \$dir/newcerts
serial            = \$dir/serial
RANDFILE          = \$dir/private/.rand

private_key       = \$dir/root/root.key.pem
certificate       = \$dir/root/root.cert.pem

[ req ]
string_mask        = utf8only
default_bits        = 2048
prompt              = no
default_md          = sha256
req_extensions      = req_ext
prompt              = no
distinguished_name  = dn

[ dn ]
countryName                 = Country Name (2 letter code)
countryName_default         = DE
stateOrProvinceName         = State or Province Name
stateOrProvinceName_default = Berlin
localityName                = Locality Name
localityName_default        = Berlin
organizationName            = Organization Name
organizationName_default    = DemoCA
commonName                  = Common Name
CN                          = localhost

[ req_ext ]
subjectAltName = @alt_names

[ alt_names ]
DNS.1 = localhost

[ v3_ca ]
basicConstraints = critical,CA:TRUE
keyUsage = critical, digitalSignature, cRLSign, keyCertSign

[ v3_intermediate_ca ]
basicConstraints = critical,CA:TRUE, pathlen:0
keyUsage = critical, digitalSignature, cRLSign, keyCertSign
authorityKeyIdentifier = keyid:always,issuer

[ v3_server_cert ]
basicConstraints = CA:FALSE
nsCertType = server
nsComment = "Demo Server Certificate"
subjectKeyIdentifier = hash
authorityKeyIdentifier = keyid,issuer
keyUsage = critical, digitalSignature, keyEncipherment
extendedKeyUsage = serverAuth

[ v3_client_cert ]
basicConstraints = CA:FALSE
keyUsage = critical, digitalSignature, keyEncipherment
extendedKeyUsage = clientAuth
subjectKeyIdentifier = hash
authorityKeyIdentifier = keyid,issuer

EOF

echo "=== 1. Create Root CA ==="
openssl genrsa -out $PKI_DIR/root/root.key.pem 4096
chmod 400 $PKI_DIR/root/root.key.pem

openssl req -x509 -new -nodes -key $PKI_DIR/root/root.key.pem \
  -sha256 -days 3650 -out $PKI_DIR/root/root.cert.pem \
  -subj "/C=DE/ST=Berlin/L=Berlin/O=DemoCA/OU=Root/CN=DemoRootCA" \
  -config $PKI_DIR/openssl.cnf -extensions v3_ca

echo "Root CA created:"
openssl x509 -noout -text -in $PKI_DIR/root/root.cert.pem

echo "******************"

echo "=== 2. Create Intermediate CA ==="
openssl genrsa -out $PKI_DIR/intermediate/intermediate.key.pem 4096
chmod 400 $PKI_DIR/intermediate/intermediate.key.pem

openssl req -new -sha256 -key $PKI_DIR/intermediate/intermediate.key.pem \
  -out $PKI_DIR/intermediate/intermediate.csr.pem \
  -subj "/C=DE/ST=Berlin/L=Berlin/O=DemoCA/OU=Intermediate/CN=DemoIntermediateCA"

# Sign Intermediate with Root (CA extensions)
openssl x509 -req -in $PKI_DIR/intermediate/intermediate.csr.pem \
  -CA $PKI_DIR/root/root.cert.pem -CAkey $PKI_DIR/root/root.key.pem \
  -CAcreateserial -out $PKI_DIR/intermediate/intermediate.cert.pem \
  -days 1825 -sha256 -extfile $PKI_DIR/openssl.cnf -extensions v3_intermediate_ca

echo "Intermediate CA created:"
openssl x509 -noout -text -in $PKI_DIR/intermediate/intermediate.cert.pem

# Create CA bundle (root + intermediate)
echo "Creating CA bundle (root + intermediate)"
cat $PKI_DIR/root/root.cert.pem $PKI_DIR/intermediate/intermediate.cert.pem > $PKI_DIR/root/ca-bundle.pem

echo "******************"
# Verify intermediate is signed by root
echo "Verifying Intermediate CA:"
openssl verify -CAfile $PKI_DIR/root/root.cert.pem $PKI_DIR/intermediate/intermediate.cert.pem

echo "=== 3. Create Server Certificate ==="
openssl genrsa -out $PKI_DIR/server/server.key.pem 2048
chmod 400 $PKI_DIR/server/server.key.pem


openssl req -new -key $PKI_DIR/server/server.key.pem \
  -out $PKI_DIR/server/server.csr.pem \
  -subj "/C=DE/ST=Berlin/L=Berlin/O=DemoServer/OU=IT/CN=localhost" \
  -config $PKI_DIR/openssl.cnf -extensions req_ext

# Sign server certificate with intermediate CA (server extensions)
# openssl x509 -req -in $PKI_DIR/server/server.csr.pem \
#   -CA $PKI_DIR/intermediate/intermediate.cert.pem \
#   -CAkey $PKI_DIR/intermediate/intermediate.key.pem \
#   -CAcreateserial -out $PKI_DIR/server/server.cert.pem \
#   -days 825 -sha256 -extfile $PKI_DIR/openssl.cnf -extensions v3_server_cert

openssl x509 -req -in $PKI_DIR/server/server.csr.pem \
  -CA $PKI_DIR/intermediate/intermediate.cert.pem \
  -CAkey $PKI_DIR/intermediate/intermediate.key.pem \
  -CAcreateserial -out $PKI_DIR/server/server.cert.pem \
  -days 825 -sha256 -extfile $PKI_DIR/openssl.cnf -extensions req_ext

echo "Server certificate created:"
openssl x509 -noout -text -in $PKI_DIR/server/server.cert.pem
# echo "Server certificate: $PKI_DIR/server/server.cert.pem"

# Verify server certificate against root and intermediate
echo "Verifying Server Certificate:"
openssl verify -CAfile $PKI_DIR/root/root.cert.pem \
  -untrusted $PKI_DIR/intermediate/intermediate.cert.pem \
  $PKI_DIR/server/server.cert.pem

echo "Verify Subject Alternative Name (SAN)"
openssl x509 -in $PKI_DIR/server/server.cert.pem -text -noout | grep -A1 "Subject Alternative Name"

echo "******************************************************"

echo "=== 4. Create full chain for server ==="
cat $PKI_DIR/server/server.cert.pem $PKI_DIR/intermediate/intermediate.cert.pem > $PKI_DIR/server/server.chain.pem

echo "Full chain created: $PKI_DIR/server/server.chain.pem"

# echo "Verifying full chain:"
# openssl verify -CAfile $PKI_DIR/root/root.cert.pem $PKI_DIR/server/server.chain.pem

echo "=== DONE: PKI setup verified successfully! ==="

echo "******************************************************"

echo "Create Client Certificates ==="

echo "Remove Client certs first"
rm -rf $PKI_DIR/client/client.key.pem $PKI_DIR/client.csr.pem $PKI_DIR/client/client.cert.pem

echo "1. Generating a client private key"
openssl genrsa -out $PKI_DIR/client/client.key.pem 2048
chmod 400 $PKI_DIR/client/client.key.pem

echo "2. Generating a CSR (Certificate Signing Request) for the client certificate:"
openssl req -new -key $PKI_DIR/client/client.key.pem \
  -subj "/C=DE/ST=Berlin/L=Berlin/O=DemoClient/OU=IT/CN=client1" \
  -out $PKI_DIR/client/client.csr.pem

# Here /CN=client1 is the client’s identity.
# You can also add OU, O, emailAddress, etc.

echo "3. Sign the CSR with your CA and create client.cert.pem"
openssl x509 -req -in $PKI_DIR/client/client.csr.pem \
  -CA $PKI_DIR/intermediate/intermediate.cert.pem \
  -CAkey $PKI_DIR/intermediate/intermediate.key.pem \
  -CAcreateserial \
  -out $PKI_DIR/client/client.cert.pem \
  -days 365 -sha256 \
  -extfile $PKI_DIR/openssl.cnf -extensions v3_client_cert

echo "4. Verify the Client certificate"
echo "=====================!!====================="
openssl x509 -in $PKI_DIR/client/client.cert.pem -text -noout
openssl verify -CAfile $PKI_DIR/root/root.cert.pem -untrusted $PKI_DIR/intermediate/intermediate.cert.pem $PKI_DIR/client/client.cert.pem
echo "=====================!!====================="
# Look for:
# Subject: CN=client1
# Issuer: CN=Demo Root CA

echo "5. Verify trust"
openssl verify -CAfile $PKI_DIR/root/root.cert.pem \
  -untrusted $PKI_DIR/intermediate/intermediate.cert.pem \
  $PKI_DIR/client/client.cert.pem

# should return
# $PKI_DIR/client/client.cert.pem: OK

# Copy Certs to our go server and client

# TEST SERVER folder
TEST_SERVER_DIR=pki-server/certs
TEST_CLIENT_DIR=pki-client/certs
TEST_CA_SERVER_DIR=pki-ca-server/certs

# Remove if folder exists
rm -rf $TEST_SERVER_DIR $TEST_CLIENT_DIR $TEST_CA_SERVER_DIR
mkdir -p $TEST_SERVER_DIR $TEST_CLIENT_DIR $TEST_CA_SERVER_DIR

cp $PKI_DIR/root/root.cert.pem $TEST_SERVER_DIR
echo "Root CA is copied to Go server:"
cp $PKI_DIR/intermediate/intermediate.cert.pem $TEST_SERVER_DIR
echo "Intermediate CA is copied to Go server:"
cp $PKI_DIR/server/server.key.pem $TEST_SERVER_DIR
echo "Server Key is copied to Go server:"
cp $PKI_DIR/server/server.cert.pem $TEST_SERVER_DIR
echo "Server cert is copied to Go server:"

cp $PKI_DIR/server/server.chain.pem $TEST_SERVER_DIR
echo "Server chain is copied to Go server:"

ls -l $TEST_SERVER_DIR

echo "Copying client certs to Go client:"
cp $PKI_DIR/client/client.cert.pem $TEST_CLIENT_DIR
echo "Client cert is copied to Go client:"
cp $PKI_DIR/client/client.key.pem $TEST_CLIENT_DIR
echo "Client key is copied to Go client:"
cp $PKI_DIR/root/root.cert.pem $TEST_CLIENT_DIR
echo "Root CA is copied to Go client:"
ls -l $TEST_CLIENT_DIR

echo "=== 5. Test HTTPS server with OpenSSL s_server ==="
echo "You can run:"
echo "openssl s_server -accept 4433 -www -cert $PKI_DIR/server/server.chain.pem -key $PKI_DIR/server/server.key.pem"
echo "Then test with:"
echo "openssl s_client -connect localhost:4433 -CAfile $PKI_DIR/root/root.cert.pem"

# echo "All certificates created and verified successfully!"
echo "******************"


# If any previous openssl s_server is running, kill it
# PID=$(pgrep -f "openssl s_server -accept 4433")
# if [ -n "$PID" ]; then
#   echo "killing previous openssl s_server with PID $PID"
#   kill -9 -f "$PID"
#   echo "Killed openssl s_server (PID $PID)"
# fi

echo "=== 6. Run server & client test ==="
# Start server in background
openssl s_server -accept 4433 -www \
  -key $PKI_DIR/server/server.key.pem \
  -cert $PKI_DIR/server/server.chain.pem \
  -CAfile $PKI_DIR/root/root.cert.pem > $PKI_DIR/server/server.log 2>&1 &

SERVER_PID=$!
echo "******************"

echo "Starting server..."
echo "Server PID: $SERVER_PID"
sleep 1  # give server a moment to start

# Run client test
# Combine both root and intermediate for client verification
mkdir -p $PKI_DIR/test
cat $PKI_DIR/root/root.cert.pem $PKI_DIR/intermediate/intermediate.cert.pem > $PKI_DIR/client/inter-root-combined.cert.pem
openssl verify -CAfile $PKI_DIR/client/inter-root-combined.cert.pem $PKI_DIR/server/server.chain.pem
cp $PKI_DIR/client/inter-root-combined.cert.pem $TEST_CLIENT_DIR
echo "Combined Cert is copied to Go client:"
echo "Running client test..."
# openssl s_client -connect localhost:4433 -CAfile $PKI_DIR/test/inter-root.cert.pem </dev/null
# echo "Q" | openssl s_client -connect localhost:4433 -CAfile $PKI_DIR/root/root.cert.pem 
# gives Verify return code: 21 (unable to verify the first certificate)
# Now we are making the client to trust both root and intermediate by combining them.
# Run client test (send a dummy "Q" so server stays alive long enough)
echo "Q" | openssl s_client -connect localhost:4433 -CAfile $PKI_DIR/client/inter-root-combined.cert.pem
# Look for "Verify return code: 0 (ok)" in output

# Kill server after test
kill $SERVER_PID
echo "openssl Server with PID $SERVER_PID stopped."

# echo "Server log:"
# cat $PKI_DIR/server/server.log

echo "******************"
echo "Client test completed."
echo "******************"