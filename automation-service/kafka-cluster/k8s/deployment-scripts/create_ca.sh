#!/bin/bash
# param 1 environment
# param 2 aws region

ENVIRONMENT=$1
BASEPATH=$2
DOCKERREGISTRYURL=$3
DOMAIN=${DOMAIN:-ww}

if [ $ENVIRONMENT == dev ]; then
  export KOPS_CLUSTER_NAME=development.worldwire-dev.io.k8s.local
  export KOPS_STATE_STORE=s3://ww-kube-dev
  ENVIRONMENT="dev"
elif [ $ENVIRONMENT == dev2 ]; then
  export KOPS_CLUSTER_NAME=dev2.worldwire-dev.io.k8s.local
  export KOPS_STATE_STORE=s3://kops-peter
  ENVIRONMENT="dev2"
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

sed "s/{{ ENVIRONMENT }}/$ENVIRONMENT/g" $BASEPATH/generate_push_cert_ca.template.yaml \
| sed "s/{{ DOMAIN }}/$DOMAIN/g" \
> generate_push_cert_ca.yaml

kubectl label namespace default istio-injection=disabled --overwrite
kubectl create -f generate_push_cert_ca.yaml
kubectl wait --timeout=90s --for=condition=complete job/ca-cert-generation