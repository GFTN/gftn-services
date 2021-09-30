
#!/bin/bash

set -o nounset \
    -o errexit \
    -o verbose \
    -o xtrace

CA_CERT=${CA_CERT:-ibmca.crt}
CA_KEY=${CA_KEY:-ibmca.key}
CA_PASSWORD=${CA_PASSWORD:-WorldWire-test}
KAFKA_KEY_PASSWORD=${KAFKA_KEY_PASSWORD:-participant-WorldWire-test}
PARTICIPANT_ID=${PARTICIPANT_ID:-participant_id_test}
STORE_LOCATION=${STORE_LOCATION:-$PARTICIPANT_ID}
CERT_NAME=${CERT_NAME:-participant}
CN=${CN:-$PARTICIPANT_ID}
mkdir -p $STORE_LOCATION

# Create Cert
openssl req -new -sha256 -keyout $STORE_LOCATION/$CERT_NAME.key -out $STORE_LOCATION/$CERT_NAME.csr -subj "/CN=$CN" -passin pass:$KAFKA_KEY_PASSWORD -passout pass:$KAFKA_KEY_PASSWORD

openssl x509 -req -CA $CA_CERT -CAkey $CA_KEY -in $STORE_LOCATION/$CERT_NAME.csr -out $STORE_LOCATION/$CERT_NAME.crt -days 9999 -CAcreateserial -passin pass:$CA_PASSWORD
