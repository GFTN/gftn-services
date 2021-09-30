#!/bin/bash
set -o nounset \
    -o errexit \
    -o verbose \
    -o xtrace

DIR="$(cd "$(dirname "$0")" && pwd)"
#BASEPATH="/root/kafka-cluster"
ENVIRONMENT=$1
DOMAIN=${DOMAIN:-ww}

#pull ca cert, key and password
SECRET_NAME=/$ENVIRONMENT/$DOMAIN/kafka_ca_cert
source $DIR/secret_pull.sh | grep SecretString |awk '{print $2}' | tr -d '",' > temp
python -m base64 -d  temp > ./ibmca.crt

SECRET_NAME=/$ENVIRONMENT/$DOMAIN/kafka_ca_key
source $DIR/secret_pull.sh | grep SecretString |awk '{print $2}' | tr -d '",' > temp
python -m base64 -d  temp >  ./ibmca.key

SECRET_NAME=/$ENVIRONMENT/$DOMAIN/kafka_ca_password
source $DIR/secret_pull.sh | grep SecretString |awk '{print $2}' | tr -d '",' > temp
python -m base64 -d  temp >  passtemp
CA_PASSWORD=$(cat passtemp)

CA_CERT=./ibmca.crt
CA_KEY=./ibmca.key
CA_PASSWORD=${CA_PASSWORD:-WorldWire-test}
KAFKA_KEY_PASSWORD=${KAFKA_KEY_PASSWORD:-participant-WorldWire-test}
PARTICIPANT_ID=$2
STORE_LOCATION=${STORE_LOCATION:-./$PARTICIPANT_ID}
CN=${CN:-$PARTICIPANT_ID}
mkdir -p $STORE_LOCATION

# Create Cert
CERT_NAME=${CERT_NAME:-participant}
source $DIR/generate_sign_cert_participant.sh

# Push Cert
SECRET_STRING=$(base64 $STORE_LOCATION/$CERT_NAME.crt | tr -d "\n\r")
SECRET_NAME=/$ENVIRONMENT/$PARTICIPANT_ID/kafka_cert
DESCRIPTION=Kafka_cert
source $DIR/secret_push.sh
SECRET_STRING=$(base64 $STORE_LOCATION/$CERT_NAME.key | tr -d "\n\r")
SECRET_NAME=/$ENVIRONMENT/$PARTICIPANT_ID/kafka_key
DESCRIPTION=Kafka_key
source $DIR/secret_push.sh
SECRET_STRING=$(echo $KAFKA_KEY_PASSWORD | base64 | tr -d "\n\r")
SECRET_NAME=/$ENVIRONMENT/$PARTICIPANT_ID/kafka_password
DESCRIPTION=Kafka_password
source $DIR/secret_push.sh

#move ibm.crt to store

mv ./ibmca.crt $STORE_LOCATION
#clean the dir after pushing

rm *.key *.srl *temp