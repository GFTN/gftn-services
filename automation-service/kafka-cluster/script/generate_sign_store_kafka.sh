#!/bin/bash

set -o nounset \
    -o errexit \
    -o verbose \
    -o xtrace

CA_CERT=${CA_CERT:-ibmca.crt}
CA_KEY=${CA_KEY:-ibmca.key}
CA_PASSWORD=${CA_PASSWORD:-WorldWire-test}
STORE_PASS=${STORE_PASS:-Worldwire-teststore}
STORE_LOCATION=${STORE_LOCATION:-"."}
LIST=${LIST:-kafka-1 kafka-2 kafka-3}
OU=${OU:-IBMBlockchain}
O=${O:-IBM}
L=${L:-SG}
C=${C:-SG}
for i in $LIST
do
	echo $i
	mkdir -p $STORE_LOCATION/$i/

	# Create keystores
	keytool -genkey -noprompt \
				 -alias $i \
				 -dname "CN=$i, OU=$OU, O=$O, L=$L, C=$C" \
				 -keystore $STORE_LOCATION/$i/kafka.$i.keystore.jks \
				 -keyalg RSA \
				 -storepass $STORE_PASS \
				 -keypass $STORE_PASS

	# Create CSR, sign the cert and import back into keystore
	
	keytool -keystore $STORE_LOCATION/$i/kafka.$i.keystore.jks -alias $i -certreq -file $i.csr -storepass $STORE_PASS 

	openssl x509 -req -CA $CA_CERT -CAkey $CA_KEY -in $i.csr -out $i-ca-signed.crt -days 9999 -CAcreateserial -passin pass:$CA_PASSWORD

	keytool -keystore $STORE_LOCATION/$i/kafka.$i.keystore.jks -alias CARoot -import -file $CA_CERT -storepass $STORE_PASS -noprompt

	keytool -keystore $STORE_LOCATION/$i/kafka.$i.keystore.jks -alias $i -import -file $i-ca-signed.crt -storepass $STORE_PASS -noprompt

	# Create truststore and import the CA cert.
	keytool -keystore $STORE_LOCATION/$i/kafka.$i.truststore.jks -alias CARoot -import -file $CA_CERT -storepass $STORE_PASS -noprompt

done
