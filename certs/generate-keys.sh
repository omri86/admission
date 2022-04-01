#!/usr/bin/env bash

# generate-keys.sh
#
# Generate a (self-signed) CA certificate and a certificate and private key to be used by the webhook demo server.
#
# NOTE: THIS SCRIPT EXISTS FOR DEMO PURPOSES ONLY. DO NOT USE IT FOR YOUR PRODUCTION WORKLOADS.

# Generate the CA cert and private key
openssl req -nodes -new -x509 -keyout ca.key -out ca.crt -subj "/CN=Admission Controller Webhook Demo CA"

# Generate a Certificate Signing Request (CSR) for the private key, and sign it with the private key of the CA.
openssl req -new -out webhook-server-tls.csr -newkey rsa:2048 -nodes -sha256 -keyout webhook-server-tls.key -config openssl-san.cnf
openssl x509 -req -CAcreateserial -in webhook-server-tls.csr -days 365 -CA ca.crt -CAkey ca.key -out webhook-server-tls.crt -extensions v3_req -extfile openssl-san.cnf

# Inject CA in the webhook template file
export CA_BUNDLE=$(cat ca.crt | base64 | tr -d '\n')
cat webhook-template.yaml | envsubst > webhook.yaml