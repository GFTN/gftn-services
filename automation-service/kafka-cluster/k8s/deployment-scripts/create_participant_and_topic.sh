#!/bin/bash
PARTICIPANT_ID=$1
STORE_LOCATION=${STORE_LOCATION:-'/store'}
KAFKA_KEY_PASSWORD=${KAFKA_KEY_PASSWORD:-Worldwire-test}
ENVIRONMENT=$2
AWS_REGION="ap-southeast-1"
BASEPATH=$3
DOMAIN=${DOMAIN:-ww}
DOCKERREGISTRYURL=$4

if [ $ENVIRONMENT == dev ]; then
  export KOPS_CLUSTER_NAME=development.worldwire-dev.io.k8s.local
  export KOPS_STATE_STORE=s3://ww-kube-dev
  ENVIRONMENT="dev"
elif  [ $ENVIRONMENT == dev2  ]; then
  export KOPS_CLUSTER_NAME=dev2.worldwire-dev.io.k8s.local
  export KOPS_STATE_STORE=s3://kops-peter
  ENVIRONMENT="dev2"
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

kubectl label namespace default istio-injection=disabled --overwrite

# create participate certs, push and save in pvc
#if [[ $STORE_LOCATION == /* ]]; then STORE_LOCATION='\'"$STORE_LOCATION"; fi
STORE_LOCATION='\'"$STORE_LOCATION"

sed "s/{{ PARTICIPANT_ID }}/$PARTICIPANT_ID/g" $BASEPATH/k8s/generate_push_sign_cert_participant.template.yaml \
| sed "s/{{ STORE_LOCATION }}/$STORE_LOCATION/g" \
| sed "s/{{ KAFKA_KEY_PASSWORD }}/$KAFKA_KEY_PASSWORD/g" \
| sed "s/{{ ENVIRONMENT }}/$ENVIRONMENT/g" \
| sed "s/{{ DOCKER_REGISTRY_URL }}/$DOCKERREGISTRYURL/g" \
> generate_push_sign_cert_participant.$PARTICIPANT_ID.yaml

kubectl create -f ./generate_push_sign_cert_participant.$PARTICIPANT_ID.yaml
kubectl wait --timeout=120s --for=condition=complete job/$PARTICIPANT_ID-onboard-cert-generation

SECRET_NAME=/$ENVIRONMENT/$DOMAIN/kafka_ca_cert
source $BASEPATH/script/secret_pull.sh | grep SecretString |awk '{print $2}' | tr -d '",' > temp
python -m base64 -d  temp > ./ibmca.crt

SECRET_NAME=/$ENVIRONMENT/$PARTICIPANT_ID/kafka_cert
source $BASEPATH/script/secret_pull.sh | grep SecretString |awk '{print $2}' | tr -d '",' > temp
python -m base64 -d  temp > ./kafka_cert.crt

SECRET_NAME=/$ENVIRONMENT/$PARTICIPANT_ID/kafka_key
source $BASEPATH/script/secret_pull.sh | grep SecretString |awk '{print $2}' | tr -d '",' > temp
python -m base64 -d  temp > ./kafka_key.key

SECRET_NAME=/$ENVIRONMENT/$PARTICIPANT_ID/kafka_password
source $BASEPATH/script/secret_pull.sh | grep SecretString |awk '{print $2}' | tr -d '",' > KAFKA_KEY_PASSWORD

kubectl create secret generic kafka-secret-$PARTICIPANT_ID --from-file=./ibmca.crt --from-file=./kafka_cert.crt --from-file=./kafka_key.key --from-literal=kafka_key_password=$KAFKA_KEY_PASSWORD

rm ./ibmca.crt ./kafka_cert.crt ./kafka_key.key ./temp ./KAFKA_KEY_PASSWORD
# create kafka topic
PARTICIPANT_ID=$1

sed "s/{{ PARTICIPANT_ID }}/$PARTICIPANT_ID/g" $BASEPATH/k8s/create_topic_participant.template.yaml \
| sed "s/{{ DOCKER_REGISTRY_URL }}/$DOCKERREGISTRYURL/g" \
> create_topic_participant.$PARTICIPANT_ID.yaml

kubectl create -f ./create_topic_participant.$PARTICIPANT_ID.yaml
kubectl wait --timeout=120s --for=condition=complete job/$PARTICIPANT_ID-create-topic-participant
