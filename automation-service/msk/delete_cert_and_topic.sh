#!/bin/bash
PARTICIPANT_ID=$1
ENVIRONMENT=$2
BASEPATH=$3
DOCKERREGISTRYURL=$4
ORGID=$5

if [ $ENVIRONMENT == eksdev ]; then
  aws eks update-kubeconfig --name terraform-eks-ww-dev --region us-west-2
  kubectl config use-context arn:aws:eks:us-west-2:$ORGID:cluster/terraform-eks-ww-dev
  export AWS_DEFAULT_REGION=us-west-2
  MSKNAME="ww-dev-msk"
  ENVIRONMENT="worldwire-dev"
elif  [ $ENVIRONMENT == eksqa  ]; then
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
  MSKNAME="ww-prod-msk"
  ENVIRONMENT="worldwire"
fi

# get zookeeper connect
ZOOKEEPER=$(aws kafka list-clusters --cluster-name-filter $MSKNAME | jq '.ClusterInfoList[].ZookeeperConnectString' | sed -e 's/^"//' -e 's/"$//')

MSKPATH="/var/files/msk/$PARTICIPANT_ID"

# delete the private certificate
export CERT_DN="$PARTICIPANT_ID.$ENVIRONMENT.io"
CERT_ARN=$(aws acm list-certificates | jq '.CertificateSummaryList[] | select(.DomainName == env.CERT_DN) | .CertificateArn' | sed -e 's/^"//' -e 's/"$//' -e 's/\\n/\n/g' )
aws acm delete-certificate --certificate-arn $CERT_ARN

# delete the cert and private key from the k8s secret
kubectl delete secret -n default kafka-secret-$PARTICIPANT_ID

# delete kafka topics
sed "s/{{ PARTICIPANT_ID }}/$PARTICIPANT_ID/g" $BASEPATH/delete_topic_participant.template.yaml \
| sed "s/{{ DOCKER_REGISTRY_URL }}/$DOCKERREGISTRYURL/g" \
| sed "s/{{ ZOOKEEPER }}/$ZOOKEEPER/g" \
| sed "s/{{ DN }}/$CERT_DN/g" \
> $MSKPATH/delete_topic_participant.$PARTICIPANT_ID.yaml

kubectl label namespace kafka-topics istio-injection=disabled --overwrite

kubectl create -n kafka-topics -f $MSKPATH/delete_topic_participant.$PARTICIPANT_ID.yaml
kubectl wait -n kafka-topics --timeout=300s --for=condition=complete job/$PARTICIPANT_ID-delete-topic-participant