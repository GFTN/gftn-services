#!/bin/bash
# param 1 environment
# decalarations

ENVIRONMENT=$1

if [ $ENVIRONMENT == dev ]; then
  export KOPS_CLUSTER_NAME=development.worldwire-dev.io.k8s.local
  export KOPS_STATE_STORE=s3://ww-kube-dev
  ENVIRONMENT="dev"
  HOST="worldwire-dev"
  HORIZONURL="https:\/\/horizon-testnet.stellar.org"
  STELLARNETWORK="Test SDF Network ; September 2015"
  FIREBASEURL="https:\/\/next-gftn.firebaseio.com"
  ENABLEJWT="false"
elif  [ $ENVIRONMENT == qa  ]; then
  export KOPS_CLUSTER_NAME=qa.worldwire-qa.io.k8s.local
  export KOPS_STATE_STORE=s3://ww-kube-qa
  ENVIRONMENT="qa"
  HOST="worldwire-qa"
  HORIZONURL="http:\/\/34.80.67.4:1234"
  STELLARNETWORK="Standalone World Wire Network ; Mar 2019"
  FIREBASEURL="https:\/\/next-gftn.firebaseio.com"
  ENABLEJWT="true"
elif [ $ENVIRONMENT == st ]; then
  export KOPS_CLUSTER_NAME=staging.worldwire-st.io.k8s.local
  export KOPS_STATE_STORE=s3://kubernetes-staging
  ENVIRONMENT="st"
  HOST="worldwire-st"
  HORIZONURL="http:\/\/34.80.67.4:1234"
  STELLARNETWORK="Standalone World Wire Network ; Mar 2019"
  FIREBASEURL="https:\/\/next-gftn.firebaseio.com"
  ENABLEJWT="false"
elif [ $ENVIRONMENT == tn ]; then
  export KOPS_CLUSTER_NAME=tn.worldwire-tn.io.k8s.local
  export KOPS_STATE_STORE=s3://ww-kube-tn
  ENVIRONMENT="tn"
  HOST="worldwire-tn"
  HORIZONURL="http:\/\/34.80.67.4:1234"
  STELLARNETWORK="Standalone World Wire Network ; Mar 2019"
  FIREBASEURL="https:\/\/next-gftn.firebaseio.com"
  ENABLEJWT="true"
elif [ $ENVIRONMENT == prod ]; then
  export KOPS_CLUSTER_NAME=prod.worldwire.io.k8s.local
  export KOPS_STATE_STORE=s3://kubernetes-worldwire
  ENVIRONMENT=""
  HOST="worldwire"
  HORIZONURL="https:\/\/horizon.stellar.org"
  STELLARNETWORK="Public Global Stellar Network ; September 2015"
  FIREBASEURL="https:\/\/gftn.firebaseio.com/"
  ENABLEJWT="true"
elif [ $ENVIRONMENT == pen ]; then
  export KOPS_CLUSTER_NAME=demo.worldwire.io.k8s.local
  export KOPS_STATE_STORE=s3://kubernetes-worldwire
  ENVIRONMENT="pen"
  HOST="worldwire"
  HORIZONURL="https:\/\/horizon.stellar.org"
  STELLARNETWORK="Public Global Stellar Network ; September 2015"
  FIREBASEURL="https:\/\/gftn.firebaseio.com/"
  ENABLEJWT="true"
fi

kops export kubecfg

BASE_PATH="/root/files/global/$ENVIRONMENT"

printf "Ready to tear down global-service on %s k8s cluster\n" "$ENVIRONMENT"

kubectl delete -f $BASE_PATH
