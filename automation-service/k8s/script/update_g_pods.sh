#!/bin/bash

ENVIRONMENT=$1

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

# update services
PR_POD_NAME=$(kubectl get pod -l app=ww-pr -o jsonpath='{.items[0].metadata.name}')
echo $PR_POD_NAME
kubectl delete pod $PR_POD_NAME

ADMIN_POD_NAME=$(kubectl get pod -l app=ww-administration -o jsonpath='{.items[0].metadata.name}')
echo $ADMIN_POD_NAME
kubectl delete pod $ADMIN_POD_NAME

PAYOUT_POD_NAME=$(kubectl get pod -l app=ww-payout -o jsonpath='{.items[0].metadata.name}')
echo $PAYOUT_POD_NAME
kubectl delete pod $PAYOUT_POD_NAME

QUOTE_POD_NAME=$(kubectl get pod -l app=ww-quotes -o jsonpath='{.items[0].metadata.name}')
echo $QUOTE_POD_NAME
kubectl delete pod $QUOTE_POD_NAME

WL_POD_NAME=$(kubectl get pod -l app=ww-whitelist -o jsonpath='{.items[0].metadata.name}')
echo $WL_POD_NAME
kubectl delete pod $WL_POD_NAME

GAS_POD_NAME=$(kubectl get pod -l app=ww-gas -o jsonpath='{.items[0].metadata.name}')
echo $GAS_POD_NAME
kubectl delete pod $GAS_POD_NAME

FEE_POD_NAME=$(kubectl get pod -l app=ww-fee -o jsonpath='{.items[0].metadata.name}')
echo $FEE_POD_NAME
kubectl delete pod $FEE_POD_NAME

ANCHOR_POD_NAME=$(kubectl get pod -l app=ww-anchor -o jsonpath='{.items[0].metadata.name}')
echo $ANCHOR_POD_NAME
kubectl delete pod $ANCHOR_POD_NAME