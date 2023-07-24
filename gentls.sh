#!/usr/bin/env bash

set -euo pipefail

rm -rf ./tls
mkdir -p ./tls

openssl ecparam -genkey -name secp521r1 -out tls/ca.key

openssl req -x509 -new -sha512 -nodes -key tls/ca.key -days 3650 -out tls/ca.crt -addext "basicConstraints=critical,CA:TRUE" -addext "keyUsage=cRLSign,keyCertSign" -subj "/CN=Vault-AU Root CA 01"

openssl ecparam -genkey -name secp521r1 -out tls/tls.key

cat <<-EOF > tls/tls.conf
[req]
default_md = sha512
basicConstraints = CA:FALSE
keyUsage = nonRepudiation, digitalSignature, keyEncipherment
[req]
distinguished_name = req_distinguished_name
req_extensions = req_ext
prompt = no
[req_distinguished_name]
CN  = localhost
[req_ext]
subjectAltName = @alt_names
[alt_names]
DNS.1 = localhost
IP.1 = 127.0.0.1
EOF

openssl req -new -sha512 -nodes -key tls/tls.key -out tls/tls.csr -config tls/tls.conf

cat <<-EOF > tls/tls-ext.conf
basicConstraints = CA:FALSE
nsCertType = server
nsComment = "VAU Certificate"
subjectKeyIdentifier = hash
authorityKeyIdentifier = keyid,issuer:always
keyUsage = nonRepudiation, digitalSignature, keyEncipherment
extendedKeyUsage = serverAuth
subjectAltName = @alt_names
[alt_names]
DNS.1 = localhost
EOF

openssl x509 -req -sha512 -days 365 -in tls/tls.csr -CA tls/ca.crt -CAkey tls/ca.key -CAcreateserial -out tls/tls.crt -extfile tls/tls-ext.conf

