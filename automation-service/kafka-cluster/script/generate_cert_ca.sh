#!/bin/bash

set -o errexit \
    -o verbose \
    -o nounset \
    -o xtrace

CERT_NAME=${CERT_NAME:-ibmca}
CA_PASSWORD=${CA_PASSWORD:-WorldWire-test}
CN=${CN:-WorldWire}
OU=${OU:-IBMBlockchain}
O=${O:-IBM}
L=${L:-SG}
C=${C:-SG}

# Generate CA key
openssl req -new -x509 -keyout $CERT_NAME.key -out $CERT_NAME.crt -days 365 -subj "/CN=$CN/OU=$OU/O=$O/L=$L/C=$C" -passin pass:$CA_PASSWORD -passout pass:$CA_PASSWORD
