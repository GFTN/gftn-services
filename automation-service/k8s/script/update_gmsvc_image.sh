#!/bin/bash

ENVIRONMENT=$1
DOCKERTAG=$2

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

LIST='ww-pr ww-administration ww-anchor ww-payout ww-gas ww-quotes ww-whitelist ww-fee'
for j in $LIST
do

if [ $j == ww-pr ]; then
  SERVICENAME="participant-registry"
elif  [ $j == ww-administration  ]; then
  SERVICENAME="administration-service"
elif [ $j == ww-anchor ]; then
  SERVICENAME="anchor-service"
elif [ $j == ww-payout ]; then
  SERVICENAME="payout-service"
elif [ $j == ww-gas ]; then
  SERVICENAME="gas-service"
elif [ $j == ww-quotes ]; then
  SERVICENAME="quotes-service"
elif [ $j == ww-whitelist ]; then
  SERVICENAME="global-whitelist-service"
elif [ $j == ww-fee ]; then
  SERVICENAME="fee-service"
fi

printf "Ready to update %s on %s k8s cluster to version %s\n" "$SERVICENAME" "$ENVIRONMENT" "$DOCKERTAG"
kubectl --record deployment.apps/$j set image deployment.v1.apps/$j $j=ip-team-worldwire-docker-local.artifactory.swg-devops.com/gftn/$SERVICENAME:$DOCKERTAG

done
