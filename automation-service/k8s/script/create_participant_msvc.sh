#!/bin/bash
# param 1 participant
# param 2 environment
# param 3 docker-tag
# param 4 remote_participant

# decalarations

PARTICIPANT=$1
ENVIRONMENT=$2
DOCKERTAG=$3
REPLICAS=$4
DOCKERREGISTRYURL=$5
ORGID=$6

if [ $ENVIRONMENT == eksdev ]; then
  aws eks update-kubeconfig --name terraform-eks-ww-dev --region us-west-2
  kubectl config use-context arn:aws:eks:us-west-2:$ORGID:cluster/terraform-eks-ww-dev
  ENVIRONMENT="dev"
  HOST="ww-dev"
  CRYPTOSERVICENAME="crypto-service"
elif  [ $ENVIRONMENT == eksqa  ]; then
  aws eks update-kubeconfig --name terraform-eks-ww-qa --region us-west-2
  kubectl config use-context arn:aws:eks:us-west-2:$ORGID:cluster/terraform-eks-ww-qa
  ENVIRONMENT="qa"
  HOST="ww-qa"
  CRYPTOSERVICENAME="crypto-service"
elif [ $ENVIRONMENT == st ]; then
  aws eks update-kubeconfig --name terraform-eks-ww-st --region us-east-2
  kubectl config use-context arn:aws:eks:us-east-2:$ORGID:cluster/terraform-eks-ww-st
  ENVIRONMENT="st"
  HOST="ww-st"
  CRYPTOSERVICENAME="crypto-service"
elif [ $ENVIRONMENT == pen ]; then
  aws eks update-kubeconfig --name terraform-eks-ww-pen --region us-east-2
  kubectl config use-context arn:aws:eks:us-east-2:$ORGID:cluster/terraform-eks-ww-pen
  ENVIRONMENT="pen"
  HOST="ww-pen"
  CRYPTOSERVICENAME="crypto-service"
elif [ $ENVIRONMENT == prod ]; then
  aws eks update-kubeconfig --name terraform-eks-ww-prod --region us-east-1
  kubectl config use-context arn:aws:eks:us-east-1:$ORGID:cluster/terraform-eks-ww-prod
  ENVIRONMENT="prod"
  HOST="ww-prod"
  CRYPTOSERVICENAME="crypto-service-prod"
fi

BASE_PATH="/var/k8s"
PARTICIPANTS_PATH="$BASE_PATH/participant/$ENVIRONMENT"

OLDPARTICIPANT=participant_id_variable
OLDDOCKERTAG=docker_tag_variable
OLDHOST=host_variable
OLDREPLICAS=replica_variable
OLDDOCKERREGISTRYURL=docker_registry_url
OLDCRYPTOSERVICENAME=crypto_service_name_variable

mkdir -p $PARTICIPANTS_PATH/$PARTICIPANT

NEW_PARTICIPANT_PATH="$PARTICIPANTS_PATH/$PARTICIPANT"

# copy files into new files
cp -r $BASE_PATH/participant/template/* $NEW_PARTICIPANT_PATH/

cd $NEW_PARTICIPANT_PATH

# replace
find . -type f | xargs sed -i "s/$OLDPARTICIPANT/$PARTICIPANT/g"
find . -type f | xargs sed -i "s/$OLDDOCKERTAG/$DOCKERTAG/g"
find . -type f | xargs sed -i "s/$OLDHOST/$HOST/g"
find . -type f | xargs sed -i "s/$OLDREPLICAS/$REPLICAS/g"
find . -type f | xargs sed -i "s/$OLDDOCKERREGISTRYURL/$DOCKERREGISTRYURL/g"
find . -type f | xargs sed -i "s/$OLDCRYPTOSERVICENAME/$CRYPTOSERVICENAME/g"

mkdir -p "/var/files/participant"

cp -r $NEW_PARTICIPANT_PATH "/var/files/participant"

printf "Ready to deploy %s on %s k8s cluster using %s image\n" "$PARTICIPANT" "$ENVIRONMENT" "$DOCKERTAG"

istioctl kube-inject -f $NEW_PARTICIPANT_PATH/api-deployment.yaml | kubectl apply -n default -f  -
istioctl kube-inject -f $NEW_PARTICIPANT_PATH/crypto-deployment.yaml | kubectl apply -n default -f  -
istioctl kube-inject -f $NEW_PARTICIPANT_PATH/send-deployment.yaml | kubectl apply -n default -f  -
istioctl kube-inject -f $NEW_PARTICIPANT_PATH/listener-deployment.yaml | kubectl apply -n default -f  -
istioctl kube-inject -f $NEW_PARTICIPANT_PATH/wwgateway-deployment.yaml | kubectl apply -n default -f  -
kubectl apply -f $NEW_PARTICIPANT_PATH/participant-service-account.yaml -n default
kubectl apply -f $NEW_PARTICIPANT_PATH/participant-svc.yaml -n default
kubectl apply -f $NEW_PARTICIPANT_PATH/participant-vservice.yaml -n default