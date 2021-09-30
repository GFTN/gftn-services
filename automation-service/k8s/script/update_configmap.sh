#!/bin/bash

ENVIRONMENT=$1
FILEPATH=$2
DATA=$3

if [ $ENVIRONMENT == dev ]; then
  export KOPS_CLUSTER_NAME=development.worldwire-dev.io.k8s.local
  export KOPS_STATE_STORE=s3://ww-kube-dev
  ENVIRONMENT="dev"
  HORIZONURL="http:\/\/34.80.67.4:1234"
  STELLARNETWORK="Standalone World Wire Network ; Mar 2019"
  FIREBASEURL="https:\/\/dev-gftn-addcf.firebaseio.com"
  ENABLEJWT="false"
elif  [ $ENVIRONMENT == qa  ]; then
  export KOPS_CLUSTER_NAME=qa.worldwire-qa.io.k8s.local
  export KOPS_STATE_STORE=s3://ww-kube-qa
  ENVIRONMENT="qa"
  HORIZONURL="http:\/\/34.80.67.4:1234"
  STELLARNETWORK="Standalone World Wire Network ; Mar 2019"
  FIREBASEURL="https:\/\/qa-gftn-b61ef.firebaseio.com"
  ENABLEJWT="true"
elif [ $ENVIRONMENT == st ]; then
  export KOPS_CLUSTER_NAME=staging.worldwire-st.io.k8s.local
  export KOPS_STATE_STORE=s3://kubernetes-staging
  ENVIRONMENT="st"
  HORIZONURL="http:\/\/34.80.67.4:1234"
  STELLARNETWORK="Standalone World Wire Network ; Mar 2019"
  FIREBASEURL="https:\/\/staging-gftn-efac5.firebaseio.com"
  ENABLEJWT="false"
elif [ $ENVIRONMENT == tn ]; then
  export KOPS_CLUSTER_NAME=tn.worldwire-tn.io.k8s.local
  export KOPS_STATE_STORE=s3://ww-kube-tn
  ENVIRONMENT="tn"
  HORIZONURL="https:\/\/horizon-testnet.stellar.org"
  STELLARNETWORK="Test SDF Network ; September 2015"
  FIREBASEURL="https:\/\/test-gftn-22fc3.firebaseio.com"
  ENABLEJWT="true"
elif [ $ENVIRONMENT == prod ]; then
  export KOPS_CLUSTER_NAME=prod.worldwire.io.k8s.local
  export KOPS_STATE_STORE=s3://kubernetes-worldwire
  ENVIRONMENT=""
  HORIZONURL="https:\/\/horizon.stellar.org"
  STELLARNETWORK="Public Global Stellar Network ; September 2015"
  FIREBASEURL="https:\/\/live-gftn-97e03.firebaseio.com"
  ENABLEJWT="true"
elif [ $ENVIRONMENT == pen ]; then
  export KOPS_CLUSTER_NAME=demo.worldwire.io.k8s.local
  export KOPS_STATE_STORE=s3://kubernetes-worldwire
  ENVIRONMENT="pen"
  HORIZONURL="https:\/\/horizon.stellar.org"
  STELLARNETWORK="Public Global Stellar Network ; September 2015"
  FIREBASEURL="https:\/\/pen-gftn.firebaseio.com"
  ENABLEJWT="true"
fi

kops export kubecfg

mkdir -p $FILEPATH

cd $FILEPATH

echo $DATA | base64 -d | sed -e 's/\\n/\n/g' >> ./env-configmap.yaml

# update config map
OLDENVIRONMENT=environment_variable
OLDHORIZONURL=horizon_url_variable
OLDSTELLARNETWORK=stellar_network_variable
OLDENABLEJWT=enable_jwt_variable
OLDFIREBASEURL=firebase_url_variable

# replace
find . -type f | xargs sed -i "s/$OLDENVIRONMENT/$ENVIRONMENT/g"
find . -type f | xargs sed -i "s/$OLDHORIZONURL/$HORIZONURL/g"
find . -type f | xargs sed -i "s/$OLDSTELLARNETWORK/$STELLARNETWORK/g"
find . -type f | xargs sed -i "s/$OLDENABLEJWT/$ENABLEJWT/g"
find . -type f | xargs sed -i "s/$OLDFIREBASEURL/$FIREBASEURL/g"

kubectl delete -f $FILEPATH/env-configmap.yaml
kubectl apply -f $FILEPATH/env-configmap.yaml