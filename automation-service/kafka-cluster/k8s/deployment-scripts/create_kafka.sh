#!/bin/bash
# param 1 environment
# param 2 aws region

ENVIRONMENT=$1
AWS_REGION=$2
BASEPATH=$3
DOCKERREGISTRYURL=$4
STORE_PASS=${STORE_PASS:-Worldwire-teststore}
LIST=${LIST:-kafka-1 kafka-2 kafka-3}
DOMAIN=${DOMAIN:-ww}
#ZONE=a
AWS_AZ_ZONE=$AWS_REGION

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

OLDDOCKERREGISTRYURL=docker_registry_url

sed "s/{{ AWS_AZ_ZONE }}/$AWS_REGION/g" $BASEPATH/generate_storage_class.template.yaml \
> generate_storage_class.yaml

sed "s/{{ STORE_PASS }}/$STORE_PASS/g" $BASEPATH/generate_push_sign_store_kafka.template.yaml \
| sed "s/{{ LIST }}/$LIST/g" \
| sed "s/{{ ENVIRONMENT }}/$ENVIRONMENT/g" \
| sed "s/{{ DOMAIN }}/$DOMAIN/g" \
| sed "s/{{ AWS_REGION }}/$AWS_REGION/g" \
| sed "s/{{ AWS_AZ_ZONE }}/$AWS_REGION/g" \
| sed "s/$OLDDOCKERREGISTRYURL/$DOCKERREGISTRYURL/g" \
> generate_push_sign_store_kafka.yaml

#sed "s/{{ STORE_PASS }}/$STORE_PASS/g" $BASEPATH/generate_push_sign_store_kafka-2.template.yaml \
#| sed "s/{{ LIST }}/$LIST/g" \
#| sed "s/{{ ENVIRONMENT }}/$ENVIRONMENT/g" \
#| sed "s/{{ DOMAIN }}/$DOMAIN/g" \
#| sed "s/{{ AWS_REGION }}/$AWS_REGION/g" \
#| sed "s/{{ AWS_AZ_ZONE }}/$AWS_REGION/g" \
#| sed "s/$OLDDOCKERREGISTRYURL/$DOCKERREGISTRYURL/g" \
#> generate_push_sign_store_kafka-2.yaml
#
#sed "s/{{ STORE_PASS }}/$STORE_PASS/g" $BASEPATH/generate_push_sign_store_kafka-3.template.yaml \
#| sed "s/{{ LIST }}/$LIST/g" \
#| sed "s/{{ ENVIRONMENT }}/$ENVIRONMENT/g" \
#| sed "s/{{ DOMAIN }}/$DOMAIN/g" \
#| sed "s/{{ AWS_REGION }}/$AWS_REGION/g" \
#| sed "s/{{ AWS_AZ_ZONE }}/$AWS_REGION/g" \
#| sed "s/$OLDDOCKERREGISTRYURL/$DOCKERREGISTRYURL/g" \
#> generate_push_sign_store_kafka-3.yaml

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

kubectl label namespace default istio-injection=disabled --overwrite
kubectl create -f generate_storage_class.yaml
kubectl create -f $BASEPATH/kafka-broker-volume/
kubectl create -f $BASEPATH/zookeeper-volume/

kubectl create -f generate_push_sign_store_kafka.yaml
kubectl wait --timeout=120s --for=condition=complete job/kafka-store-generation

#kubectl create -f generate_push_sign_store_kafka-2.yaml
#kubectl wait --timeout=90s --for=condition=complete job/kafka-store-generation-2
#
#kubectl create -f generate_push_sign_store_kafka-3.yaml
#kubectl wait --timeout=90s --for=condition=complete job/kafka-store-generation-3

#kubectl create -f $BASEPATH/zookeeper-service-entry
kubectl label namespace default istio-injection=enabled --overwrite
kubectl create -f $BASEPATH/zookeeper/
kubectl create -f $BASEPATH/kafka-broker/
#sleep 2m
#kubectl delete -f $BASEPATH/zookeeper-service-entry
