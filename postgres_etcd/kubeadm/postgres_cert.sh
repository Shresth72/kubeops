#!/bin/bash

# Install and start a postgres service
apt -y install postgresql postgresql-contrib
systemctl start postgresql.service

# Generate self signed root CA cert
openssl req -addext "subjectAltName = DNS:localhost" -nodes \
    -x509 -newkey rsa:2048 -keyout ca.key -out ca.crt -subj "/CN=localhost"

# Generate server cert to be signed
openssl req -addext "subjectAltName = DNS:localhost" -nodes \
    -newkey rsa:2048 -keyout server.key -out server.csr -subj "/CN=localhost"

# Sign the server cert
openssl x509 -extfile <(printf "subjectAltName=DNS:localhost") -req \
    -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out server.crt

chmod og-rwx ca.key
chmod og-rwx server.key

cp {server.crt,server.key,ca.crt} /var/lib/postgresql/
chown postgres.postgres /var/lib/postgresql/server.key

# Copy both the certificates and key to /var/lib/postgresql/ and make postgres user an owner of the server key
# And modify the postgresql.conf to tell PostgreSQL to use our new SSL certs and key
sed -i -e "s|ssl_cert_file.*|ssl_cert_file = '/var/lib/postgresql/server.crt'|g" /etc/postgresql/14/main/postgresql.conf
sed -i -e "s|ssl_key_file.*|ssl_key_file = '/var/lib/postgresql/server.key'|g" /etc/postgresql/14/main/postgresql.conf
sed -i -e "s|#ssl_ca_file.*|ssl_ca_file = '/var/lib/postgresql/ca.crt'|g" /etc/postgresql/14/main/postgresql.conf

systemctl restart postgresql.service