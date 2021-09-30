#!/bin/bash
ENVIRONMENT=$1
BASEPATH=$2
DOCKERREGISTRYURL=$3
ORGID=$4

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
CERT_DN="global.$ENVIRONMENT.io"

MSKPATH="."
mkdir -p $MSKPATH

# create kafka topics and grant acl for producer and consumer to use those topics
sed "s/{{ DOCKER_REGISTRY_URL }}/$DOCKERREGISTRYURL/g" $BASEPATH/create_topic_monitoring.template.yaml \
| sed "s/{{ ZOOKEEPER }}/$ZOOKEEPER/g" \
| sed "s/{{ DN }}/$CERT_DN/g" \
> $MSKPATH/create_topic_monitoring.yaml

kubectl label namespace kafka-topics istio-injection=disabled --overwrite

kubectl create -n kafka-topics -f $MSKPATH/create_topic_monitoring.yaml
kubectl wait -n kafka-topics --timeout=300s --for=condition=complete job/create-topic-monitoring