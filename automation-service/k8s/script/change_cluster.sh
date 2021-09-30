#!/bin/bash
# param 1 environment

ENVIRONMENT=$1

if [ $ENVIRONMENT == eksdev ]; then
  aws eks update-kubeconfig --name terraform-eks-ww-dev --region us-west-2
elif  [ $ENVIRONMENT == eksqa  ]; then
  aws eks update-kubeconfig --name terraform-eks-ww-qa --region us-west-2
elif [ $ENVIRONMENT == st ]; then
  aws eks update-kubeconfig --name terraform-eks-ww-st --region us-east-2
elif [ $ENVIRONMENT == pen ]; then
  aws eks update-kubeconfig --name terraform-eks-ww-pen --region us-east-2
elif [ $ENVIRONMENT == prod ]; then
  aws eks update-kubeconfig --name terraform-eks-ww-prod --region us-east-1
fi

istioctl version