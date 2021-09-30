#!/bin/bash

PARTICIPANTID=$1
SERVICENAME=$2
AWSACCESSKEYID=$3
AWSSECRETACCESSKEY=$4
ENVIRONMENT=$5
ORGID=$6

if  [ $ENVIRONMENT == eksdev  ]; then
  aws eks update-kubeconfig --name terraform-eks-ww-dev --region us-west-2
  kubectl config use-context arn:aws:eks:us-west-2:$ORGID:cluster/terraform-eks-ww-dev
elif  [ $ENVIRONMENT == eksqa  ]; then
  aws eks update-kubeconfig --name terraform-eks-ww-qa --region us-west-2
  kubectl config use-context arn:aws:eks:us-west-2:$ORGID:cluster/terraform-eks-ww-qa
elif  [ $ENVIRONMENT == st  ]; then
  aws eks update-kubeconfig --name terraform-eks-ww-st --region us-east-2
  kubectl config use-context arn:aws:eks:us-east-2:$ORGID:cluster/terraform-eks-ww-st
elif  [ $ENVIRONMENT == pen  ]; then
  aws eks update-kubeconfig --name terraform-eks-ww-pen --region us-east-2
  kubectl config use-context arn:aws:eks:us-east-2:$ORGID:cluster/terraform-eks-ww-pen
elif [ $ENVIRONMENT == prod ]; then
  aws eks update-kubeconfig --name terraform-eks-ww-prod --region us-east-1
  kubectl config use-context arn:aws:eks:us-east-1:$ORGID:cluster/terraform-eks-ww-prod
fi

#kops export kubecfg

kubectl create secret generic --namespace=default $PARTICIPANTID-$SERVICENAME-aws-iam-key --from-literal=aws-access-key-id=$AWSACCESSKEYID --from-literal=aws-secret-access-key=$AWSSECRETACCESSKEY
