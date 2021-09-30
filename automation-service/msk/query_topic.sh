#!/bin/bash
PARTICIPANT_ID=$1
ENVIRONMENT=$2
BASEPATH=$3
DOCKERREGISTRYURL=$4
ORGID=$5

if [ $ENVIRONMENT == dev ]; then
  aws eks update-kubeconfig --name terraform-eks-ww-dev --region us-west-2
  kubectl config use-context arn:aws:eks:us-west-2:$ORGID:cluster/terraform-eks-ww-dev
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
  MSKNAME="ww-prod-msk"
  ENVIRONMENT="worldwire"
fi

#kops export kubecfg

# get zookeeper connect
CLUSTERARN=$(aws kafka list-clusters --cluster-name-filter $MSKNAME | jq '.ClusterInfoList[].ClusterArn' | sed -e 's/^"//' -e 's/"$//')
BROKERSURL=$(aws kafka get-bootstrap-brokers --cluster-arn $CLUSTERARN | jq '.BootstrapBrokerStringTls' | sed -e 's/^"//' -e 's/"$//')

# create kafka topics and grant acl for producer and consumer to use those topics
sed "s/{{ PARTICIPANT_ID }}/$PARTICIPANT_ID/g" $BASEPATH/query_topic.template.yaml \
| sed "s/{{ DOCKER_REGISTRY_URL }}/$DOCKERREGISTRYURL/g" \
| sed "s/{{ BROKER }}/$BROKERSURL/g" \
> query_topic.$PARTICIPANT_ID.yaml

kubectl label namespace kafka-topics istio-injection=disabled --overwrite

kubectl create -n kafka-topics -f ./query_topic.$PARTICIPANT_ID.yaml
kubectl wait -n kafka-topics --timeout=60s --for=condition=complete job/$PARTICIPANT_ID-query-topic