#!/bin/bash
PARTICIPANT_ID=$1
ENVIRONMENT=$2
DOCKERREGISTRYURL=${DOCKERREGISTRYURL-ip-team-worldwire-docker-local.artifactory.swg-devops.com}
ORGID=$3
VERSION=$4

if [ $ENVIRONMENT == dev ]; then
  aws eks update-kubeconfig --name terraform-eks-ww-qa --region us-west-2
  kubectl config use-context arn:aws:eks:us-west-2:$ORGID:cluster/terraform-eks-ww-qa
  export AWS_DEFAULT_REGION=us-west-2
  MSKNAME="ww-dev-msk"
  ENVIRONMENT="worldwire-dev"
elif  [ $ENVIRONMENT == qa  ]; then
  aws eks update-kubeconfig --name terraform-eks-ww-qa --region us-west-2
  kubectl config use-context arn:aws:eks:us-west-2:$ORGID:cluster/terraform-eks-ww-qa
  export AWS_DEFAULT_REGION=us-west-2
  MSKNAME="ww-qa-msk"
  ENVIRONMENT="worldwire-qa"
elif [ $ENVIRONMENT == st ]; then
  aws eks update-kubeconfig --name terraform-eks-ww-st --region us-east-2
  kubectl config use-context arn:aws:eks:us-east-2:$ORGID:cluster/terraform-eks-ww-st
  export AWS_DEFAULT_REGION=us-east-2
  MSKNAME="ww-st-msk"
  ENVIRONMENT="worldwire-st"
elif [ $ENVIRONMENT == prod ]; then
  aws eks update-kubeconfig --name terraform-eks-ww-prod --region us-east-1
  kubectl config use-context arn:aws:eks:us-east-1:$ORGID:cluster/terraform-eks-ww-prod
  export AWS_DEFAULT_REGION=us-east-1
  ENVIRONMENT="worldwire"
fi

# get zookeeper connect
ZOOKEEPER=$(aws kafka list-clusters --cluster-name-filter $MSKNAME | jq '.ClusterInfoList[].ZookeeperConnectString' | sed -e 's/^"//' -e 's/"$//')

# create participate certs, push and save in k8s secret
#CA_ARN=$(aws acm-pca list-certificate-authorities | jq '.CertificateAuthorities[].Arn' | sed -e 's/^"//' -e 's/"$//')
#CERT_DN="$PARTICIPANT_ID.$ENVIRONMENT.io"
#CERT_ARN=$(aws acm request-certificate --domain-name $CERT_DN --certificate-authority-arn $CA_ARN | jq '.CertificateArn' | sed -e 's/^"//' -e 's/"$//' )

# get the private certificate
#aws acm get-certificate --certificate-arn $CERT_ARN | jq '.Certificate' | sed -e 's/^"//' -e 's/"$//' -e 's/\\n/\n/g' >> ./kafka_cert.crt

# key password from secret manager
#KAFKA_KEY_PASSWORD="$PARTICIPANT_ID-$KEY_PASSWORD"
# get the private key
#aws acm export-certificate --certificate-arn $CERT_ARN --passphrase $KAFKA_KEY_PASSWORD | jq '.PrivateKey' | sed -e 's/^"//' -e 's/"$//' -e 's/\\n/\n/g' >> "./kafka_key.key"

# store the cert and private key into the k8s secret
#kubectl create secret generic kafka-secret-$PARTICIPANT_ID --from-file=./kafka_cert.crt --from-file=./kafka_key.key --from-literal=kafka_key_password=$KAFKA_KEY_PASSWORD

#rm ./kafka_cert.crt ./kafka_key.key

#DN="CN="
#DN="*"
CERT_DN="$PARTICIPANT_ID.$ENVIRONMENT.io"

# create kafka topics and grant acl for producer and consumer to use those topics
sed "s/{{ PARTICIPANT_ID }}/$PARTICIPANT_ID/g" ./add_participant_acl_template.yaml \
| sed "s/{{ DOCKER_REGISTRY_URL }}/$DOCKERREGISTRYURL/g" \
| sed "s/{{ VERSION }}/$VERSION/g" \
| sed "s/{{ ZOOKEEPER }}/$ZOOKEEPER/g" \
| sed "s/{{ DN }}/$CERT_DN/g" \
> add_participant_acl_$PARTICIPANT_ID.yaml

kubectl label namespace kafka-topics istio-injection=disabled --overwrite

kubectl create -n kafka-topics -f ./add_participant_acl_$PARTICIPANT_ID.yaml
kubectl wait -n kafka-topics --timeout=300s --for=condition=complete job/$PARTICIPANT_ID-add-acl-participant