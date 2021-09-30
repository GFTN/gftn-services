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
ENVIRONMENT=${ENVIRONMENT:-test}
DOMAIN=${DOMAIN:-ww}
DIR="$(cd "$(dirname "$0")" && pwd)"

echo $DIR
#pull ca cert, key and password
SECRET_NAME=/$ENVIRONMENT/$DOMAIN/kafka_ca_cert
source $DIR/secret_pull.sh | grep SecretString |awk '{print $2}' | tr -d '",' > temp
python -m base64 -d  temp > ./$CA_CERT

SECRET_NAME=/$ENVIRONMENT/$DOMAIN/kafka_ca_key
source $DIR/secret_pull.sh | grep SecretString |awk '{print $2}' | tr -d '",' > temp
python -m base64 -d  temp >  ./$CA_KEY

SECRET_NAME=/$ENVIRONMENT/$DOMAIN/kafka_ca_password
source $DIR/secret_pull.sh | grep SecretString |awk '{print $2}' | tr -d '",' > temp
python -m base64 -d  temp >  passtemp
CA_PASSWORD=$(cat passtemp)

source $DIR/generate_sign_store_kafka.sh

if [ $ENVIRONMENT == dev ]; then
  export KOPS_CLUSTER_NAME=development.worldwire-dev.io.k8s.local
  export KOPS_STATE_STORE=s3://ww-kube-dev
  ENVIRONMENT="dev"
elif  [ $ENVIRONMENT == qa  ]; then
  export KOPS_CLUSTER_NAME=qa.worldwire-qa.io.k8s.local
  export KOPS_STATE_STORE=s3://ww-kube-qa
  ENVIRONMENT="qa"
elif [ $ENVIRONMENT == st ]; then
  export KOPS_CLUSTER_NAME=staging.worldwire-st.io.k8s.local
  export KOPS_STATE_STORE=s3://kubernetes-staging
  ENVIRONMENT="st"
elif [ $ENVIRONMENT == tn ]; then
  export KOPS_CLUSTER_NAME=tn.worldwire-tn.io.k8s.local
  export KOPS_STATE_STORE=s3://ww-kube-tn
  ENVIRONMENT="tn"
elif [ $ENVIRONMENT == prod ]; then
  export KOPS_CLUSTER_NAME=prod.worldwire.io.k8s.local
  export KOPS_STATE_STORE=s3://kubernetes-worldwire
  ENVIRONMENT=""
elif [ $ENVIRONMENT == pen ]; then
  export KOPS_CLUSTER_NAME=demo.worldwire.io.k8s.local
  export KOPS_STATE_STORE=s3://kubernetes-worldwire
  ENVIRONMENT="pen"
fi

kops export kubecfg

for i in $LIST
do
	echo "pushing" $i "stores"
	SECRET_STRING=$(base64 $STORE_LOCATION/$i/kafka.$i.keystore.jks | tr -d "\n\r")
	SECRET_NAME=/$ENVIRONMENT/$DOMAIN/$i/kafka_keystore
	DESCRIPTION="$i.kafka_keystore"
	source $DIR/secret_push.sh

	SECRET_STRING=$(base64 $STORE_LOCATION/$i/kafka.$i.truststore.jks | tr -d "\n\r")
	SECRET_NAME=/$ENVIRONMENT/$DOMAIN/$i/kafka_truststore
	DESCRIPTION="$i.kafka_truststore"
	source $DIR/secret_push.sh

	SECRET_STRING=$(echo $STORE_PASS | base64 | tr -d "\n\r")
	SECRET_NAME=/$ENVIRONMENT/$DOMAIN/$i/store_pass
	DESCRIPTION="$i.store_pass"
	source $DIR/secret_push.sh

	kubectl create secret generic kafka-store-generation-secret-$i --from-file=$STORE_LOCATION/$i/kafka.$i.keystore.jks --from-file=$STORE_LOCATION/$i/kafka.$i.truststore.jks --from-literal=store_pass=$STORE_PASS
done

rm ./$CA_CERT ./$CA_KEY