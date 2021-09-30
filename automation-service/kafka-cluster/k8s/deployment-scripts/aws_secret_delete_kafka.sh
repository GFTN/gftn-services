#!/bin/bash
ENVIRONMENT=$1

if [ $ENVIRONMENT == dev ]; then
  export KOPS_CLUSTER_NAME=development.worldwire-dev.io.k8s.local
  export KOPS_STATE_STORE=s3://ww-kube-dev
elif  [ $ENVIRONMENT == qa  ]; then
  export KOPS_CLUSTER_NAME=qa.worldwire-qa.io.k8s.local
  export KOPS_STATE_STORE=s3://ww-kube-qa
elif [ $ENVIRONMENT == st ]; then
  export KOPS_CLUSTER_NAME=staging.worldwire-st.io.k8s.local
  export KOPS_STATE_STORE=s3://kubernetes-staging
elif [ $ENVIRONMENT == tn ]; then
  export KOPS_CLUSTER_NAME=tn.worldwire-tn.io.k8s.local
  export KOPS_STATE_STORE=s3://ww-kube-tn
elif [ $ENVIRONMENT == prod ]; then
  export KOPS_CLUSTER_NAME=prod.worldwire.io.k8s.local
  export KOPS_STATE_STORE=s3://kubernetes-worldwire
elif [ $ENVIRONMENT == pen ]; then
  export KOPS_CLUSTER_NAME=demo.worldwire.io.k8s.local
  export KOPS_STATE_STORE=s3://kubernetes-worldwire
fi

kops export kubecfg

DOMAIN=${DOMAIN:-ww}
BROKER_LIST="kafka-1 kafka-2 kafka-3"
LIST='kafka_keystore kafka_truststore store_pass'
for j in $BROKER_LIST
do
for i in $LIST
do
aws secretsmanager delete-secret --secret-id /$ENVIRONMENT/$DOMAIN/$j/$i \
   --force-delete-without-recovery
done
done