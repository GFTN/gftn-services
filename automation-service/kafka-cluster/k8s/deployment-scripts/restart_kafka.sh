#!/bin/bash

ENVIRONMENT=$1
AWS_REGION=$2
BASEPATH=$3
DOCKERREGISTRYURL=$4
ZONE=a
#AWS_AZ_ZONE="$AWS_REGION$ZONE"

if [ $ENVIRONMENT == dev ]; then
  export KOPS_CLUSTER_NAME=development.worldwire-dev.io.k8s.local
  export KOPS_STATE_STORE=s3://ww-kube-dev
elif  [ $ENVIRONMENT == dev2  ]; then
  export KOPS_CLUSTER_NAME=dev2.worldwire-dev.io.k8s.local
  export KOPS_STATE_STORE=s3://kops-peter
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

OLDDOCKERREGISTRYURL=docker_registry_url

LIST="kafka-1 kafka-2 kafka-3"
for i in $LIST
do
kubectl delete deployment -l app=$i
kubectl delete service -l app=$i
kubectl delete service -l app=$i-headless
done

# tear down zookeeper
LIST="zookeeper-1 zookeeper-2 zookeeper-3"
for i in $LIST
do
kubectl delete deployment -l app=$i
kubectl delete service -l app=$i
kubectl delete service -l app=$i-headless
done

# replace AWS_AZ_ZONE in kafka deployment
LIST="kafka-1 kafka-2 kafka-3"
for i in $LIST
do
sed "s/{{ AWS_AZ_ZONE }}/$AWS_REGION/g" $BASEPATH/kafka-broker-template/$i-deployment.yaml \
| sed "s/$OLDDOCKERREGISTRYURL/$DOCKERREGISTRYURL/g" \
> $BASEPATH/kafka-broker/$i-deployment.yaml
done

# replace AWS_AZ_ZONE in zookeeper deployment
LIST="zookeeper-1 zookeeper-2 zookeeper-3"
for i in $LIST
do
sed "s/{{ AWS_AZ_ZONE }}/$AWS_REGION/g" $BASEPATH/zookeeper-template/$i-deployment.yaml \
| sed "s/$OLDDOCKERREGISTRYURL/$DOCKERREGISTRYURL/g" \
> $BASEPATH/zookeeper/$i-deployment.yaml
done

#kubectl create -f $BASEPATH/zookeeper-service-entry
kubectl label namespace default istio-injection=enabled --overwrite
kubectl create -f $BASEPATH/zookeeper/
kubectl create -f $BASEPATH/kafka-broker/
#sleep 2m
#kubectl delete -f $BASEPATH/zookeeper-service-entry
