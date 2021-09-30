#!/bin/bash

set -o errexit \
    -o verbose \
    -o nounset \
    -o xtrace

ENVIRONMENT=${ENVIRONMENT:-test}
DOMAIN=${DOMAIN:-ww}

CERT_NAME=${CA_NAME:-ibmca}
CA_PASSWORD=${CA_PASSWORD:-$(openssl rand -base64 16)}
CN=${CN:-WorldWire}
OU=${OU:-IBMBlockchain}
O=${O:-IBM}
L=${L:-SG}
C=${C:-SG}

DIR="$(cd "$(dirname "$0")" && pwd)"

source $DIR/generate_cert_ca.sh
SECRET_STRING=$(base64 $CERT_NAME.crt | tr -d "\n\r")
SECRET_NAME=/$ENVIRONMENT/$DOMAIN/kafka_ca_cert
DESCRIPTION=Kafka_ca_cert
source $DIR/secret_push.sh
SECRET_STRING=$(base64 $CERT_NAME.key | tr -d "\n\r")
SECRET_NAME=/$ENVIRONMENT/$DOMAIN/kafka_ca_key
DESCRIPTION=Kafka_ca_key
source $DIR/secret_push.sh
SECRET_STRING=$(echo $CA_PASSWORD | base64 | tr -d "\n\r")
SECRET_NAME=/$ENVIRONMENT/$DOMAIN/kafka_ca_password
DESCRIPTION=Kafka_ca_password
source $DIR/secret_push.sh