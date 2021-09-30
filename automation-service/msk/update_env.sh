#!/usr/bin/env bash

ENVIRONMENT=$1
ORGID=$2

if [ $ENVIRONMENT == eksdev ]; then
  aws eks update-kubeconfig --name terraform-eks-ww-dev --region us-west-2
  kubectl config use-context arn:aws:eks:us-west-2:$ORGID:cluster/terraform-eks-ww-dev
  export AWS_DEFAULT_REGION=us-west-2
  MSKNAME="ww-dev-msk"
elif  [ $ENVIRONMENT == eksqa  ]; then
  aws eks update-kubeconfig --name terraform-eks-ww-qa --region us-west-2
  kubectl config use-context arn:aws:eks:us-west-2:$ORGID:cluster/terraform-eks-ww-qa
  export AWS_DEFAULT_REGION=us-west-2
  MSKNAME="ww-qa-msk"
elif [ $ENVIRONMENT == st ]; then
  aws eks update-kubeconfig --name terraform-eks-ww-st --region us-east-2
  kubectl config use-context arn:aws:eks:us-east-2:$ORGID:cluster/terraform-eks-ww-st
  export AWS_DEFAULT_REGION=us-east-2
  MSKNAME="ww-st-msk"
elif [ $ENVIRONMENT == prod ]; then
  aws eks update-kubeconfig --name terraform-eks-ww-prod --region us-east-1
  kubectl config use-context arn:aws:eks:us-east-1:$ORGID:cluster/terraform-eks-ww-prod
  export AWS_DEFAULT_REGION=us-east-1
  MSKNAME="ww-prod-msk"
fi

CLUSTERARN=$(aws kafka list-clusters --cluster-name-filter $MSKNAME | jq '.ClusterInfoList[].ClusterArn' | sed -e 's/^"//' -e 's/"$//')
BROKERSURL=$(aws kafka get-bootstrap-brokers --cluster-arn=$CLUSTERARN | jq '.BootstrapBrokerStringTls' | sed -e 's/^"//' -e 's/"$//')

kubectl set env -n deployment deployment/deployment-service TEST="TEST"
kubectl set env -n deployment deployment/deployment-service BROKERS_URL=$BROKERSURL