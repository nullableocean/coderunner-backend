#!/bin/bash
# Generate TLS-certs for Docker Daemon

set -e  

CERTS_DIR="${1:-../docker/certs}"

CA_KEY="$CERTS_DIR/key.pem"         # private CA
CA_CERT="$CERTS_DIR/ca.pem"         # CA
CA_SERIAL="$CERTS_DIR/ca.srl"       # serial CA
CERT_PEM="$CERTS_DIR/cert.pem"      # client CA
SERVER_KEY="$CERTS_DIR/server-key.pem"  # private server key
SERVER_CSR="$CERTS_DIR/server.csr"  # server req signature
SERVER_CERT="$CERTS_DIR/server-cert.pem" # server cert


if [ ! -d "$CERTS_DIR" ]; then
    echo "create certs dir: $CERTS_DIR"
    mkdir -p "$CERTS_DIR"
fi

echo "generate Root CA..."
openssl genrsa -out "$CA_KEY" 4096
openssl req -new -x509 -days 365 -key "$CA_KEY" -subj "/CN=Docker Root CA" -out "$CA_CERT"


echo "generate server certs..."
openssl genrsa -out "$SERVER_KEY" 4096
openssl req -new -key "$SERVER_KEY" -subj "/CN=docdemon" -out "$SERVER_CSR"

# SAN config
cat > "$CERTS_DIR/san.cnf" <<EOF
[req]
req_extensions = v3_req
distinguished_name = req_distinguished_name

[req_distinguished_name]
CN = docdemon

[v3_req]
basicConstraints = CA:FALSE
keyUsage = digitalSignature, keyEncipherment
extendedKeyUsage = serverAuth
subjectAltName = @alt_names

[alt_names]
DNS.1 = docdemon
DNS.2 = localhost
IP.1 = 127.0.0.1
EOF

# server cert with SAN
openssl x509 -req -days 365 -in "$SERVER_CSR" -CA "$CA_CERT" -CAkey "$CA_KEY" \
    -CAcreateserial -out "$SERVER_CERT" -extfile "$CERTS_DIR/san.cnf" -extensions v3_req

# client cert (copy)
cp "$CA_CERT" "$CERT_PEM" 


chmod 0400 "$CA_KEY" "$SERVER_KEY"
chmod 0444 "$CA_CERT" "$SERVER_CERT" "$CERT_PEM"
# chown -R 1001:1001 "$CERTS_DIR"  # Устанавливаем владельца (user из Dockerfile)

echo "certificates created"