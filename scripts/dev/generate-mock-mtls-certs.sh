#!/usr/bin/env bash
set -euo pipefail

OUT_DIR="${1:-test-certs}"
ORG_NR="${2:-5560001111}"

mkdir -p "$OUT_DIR"

openssl req -x509 -newkey rsa:2048 -days 3650 -nodes \
  -keyout "$OUT_DIR/ca-key.pem" \
  -out "$OUT_DIR/ca.pem" \
  -subj "/CN=ExpiTrust Test CA v8/O=Expisoft AB/C=SE"

cat > "$OUT_DIR/server.cnf" <<EOF
[ req ]
distinguished_name = dn
prompt = no
req_extensions = req_ext

[ dn ]
CN = localhost
O = Redofri
C = SE

[ req_ext ]
subjectAltName = @alt_names
extendedKeyUsage = serverAuth

[ alt_names ]
DNS.1 = localhost
IP.1 = 127.0.0.1
EOF

openssl req -newkey rsa:2048 -nodes \
  -keyout "$OUT_DIR/server-key.pem" \
  -out "$OUT_DIR/server.csr" \
  -config "$OUT_DIR/server.cnf"

openssl x509 -req -days 3650 \
  -in "$OUT_DIR/server.csr" \
  -CA "$OUT_DIR/ca.pem" \
  -CAkey "$OUT_DIR/ca-key.pem" \
  -CAcreateserial \
  -out "$OUT_DIR/server.pem" \
  -extensions req_ext \
  -extfile "$OUT_DIR/server.cnf"

cat > "$OUT_DIR/client.cnf" <<EOF
[ req ]
distinguished_name = dn
prompt = no
req_extensions = req_ext

[ dn ]
CN = Redofri Mock Client
O = Example Client AB
C = SE
serialNumber = 16${ORG_NR}

[ req_ext ]
extendedKeyUsage = clientAuth
EOF

openssl req -newkey rsa:2048 -nodes \
  -keyout "$OUT_DIR/client-key.pem" \
  -out "$OUT_DIR/client.csr" \
  -config "$OUT_DIR/client.cnf"

openssl x509 -req -days 3650 \
  -in "$OUT_DIR/client.csr" \
  -CA "$OUT_DIR/ca.pem" \
  -CAkey "$OUT_DIR/ca-key.pem" \
  -CAcreateserial \
  -out "$OUT_DIR/client.pem" \
  -extensions req_ext \
  -extfile "$OUT_DIR/client.cnf"

rm -f "$OUT_DIR"/*.csr "$OUT_DIR"/*.cnf "$OUT_DIR"/*.srl

cat <<EOF
Generated mock mTLS certificates in $OUT_DIR

Server cert: $OUT_DIR/server.pem
Server key:  $OUT_DIR/server-key.pem
Client cert: $OUT_DIR/client.pem
Client key:  $OUT_DIR/client-key.pem
CA cert:     $OUT_DIR/ca.pem

Client certificate serialNumber is set to 16${ORG_NR}
EOF
