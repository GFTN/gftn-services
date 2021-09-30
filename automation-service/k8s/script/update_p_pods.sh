#!/bin/bash

ENVIRONMENT=$1
PARTICIPANT=$2

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
API_POD_NAME=$(kubectl get pod -l app=$PARTICIPANT-api -o jsonpath='{.items[0].metadata.name}')
echo $API_POD_NAME
kubectl delete pod $API_POD_NAME

CRYPTO_POD_NAME=$(kubectl get pod -l app=$PARTICIPANT-crypto -o jsonpath='{.items[0].metadata.name}')
echo $CRYPTO_POD_NAME
kubectl delete pod $CRYPTO_POD_NAME

LISTENER_POD_NAME=$(kubectl get pod -l app=$PARTICIPANT-listener -o jsonpath='{.items[0].metadata.name}')
echo $LISTENER_POD_NAME
kubectl delete pod $LISTENER_POD_NAME

RDO_POD_NAME=$(kubectl get pod -l app=$PARTICIPANT-rdo -o jsonpath='{.items[0].metadata.name}')
echo $RDO_POD_NAME
kubectl delete pod $RDO_POD_NAME

SEND_POD_NAME=$(kubectl get pod -l app=$PARTICIPANT-send -o jsonpath='{.items[0].metadata.name}')
echo $SEND_POD_NAME
kubectl delete pod $SEND_POD_NAME

CALLBACK_POD_NAME=$(kubectl get pod -l app=$PARTICIPANT-callback -o jsonpath='{.items[0].metadata.name}')
echo $CALLBACK_POD_NAME
kubectl delete pod $CALLBACK_POD_NAME

RDOCLIENT_POD_NAME=$(kubectl get pod -l app=$PARTICIPANT-rdo-client -o jsonpath='{.items[0].metadata.name}')
echo $RDOCLIENT_POD_NAME
kubectl delete pod $RDOCLIENT_POD_NAME