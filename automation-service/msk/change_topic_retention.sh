#!/bin/bash
PARTICIPANT_ID=$1
KEY_PASSWORD=$2
ENVIRONMENT=$3
BASEPATH=$4
DOCKERREGISTRYURL=$5
TOPIC=$6
ORGID=$7

if [ $ENVIRONMENT == eksdev ]; then
  aws eks update-kubeconfig --name terraform-eks-ww-dev --region us-west-2
  kubectl config use-context arn:aws:eks:us-west-2:$ORGID:cluster/terraform-eks-ww-dev
  export AWS_DEFAULT_REGION=us-west-2
  MSKNAME="ww-dev-msk"
elif  [ $ENVIRONMENT == eksqa  ]; then
  aws eks update-kubeconfig --name terraform-eks-ww-qa --region us-west-2
  kubectl config use-context arn:aws:eks:us-west-2:$ORGID:cluster/terraform-eks-ww-qa
  export AWS_DEFAULT_REGION=us-west-2
  MSKNAME="ww-qa-msk"
elif [ $ENVIRONMENT == st ]; then
  aws eks update-kubeconfig --name terraform-eks-ww-st --region us-east-2
  kubectl config use-context arn:aws:eks:us-east-2:$ORGID:cluster/terraform-eks-ww-st
  export AWS_DEFAULT_REGION=us-east-2
  MSKNAME="ww-st-msk"
elif [ $ENVIRONMENT == prod ]; then
  aws eks update-kubeconfig --name terraform-eks-ww-prod --region us-east-1
  kubectl config use-context arn:aws:eks:us-east-1:$ORGID:cluster/terraform-eks-ww-prod
  export AWS_DEFAULT_REGION=us-east-1
  MSKNAME="ww-prod-msk"
fi

# get zookeeper connect
ZOOKEEPER=$(aws kafka list-clusters --cluster-name-filter $MSKNAME | jq '.ClusterInfoList[].ZookeeperConnectString' | sed -e 's/^"//' -e 's/"$//')
TIME=604800000

# create kafka topics and grant acl for producer and consumer to use those topics
sed "s/{{ PARTICIPANT_ID }}/$PARTICIPANT_ID/g" $BASEPATH/change_topic_retention.template.yaml \
| sed "s/{{ DOCKER_REGISTRY_URL }}/$DOCKERREGISTRYURL/g" \
| sed "s/{{ ZOOKEEPER }}/$ZOOKEEPER/g" \
| sed "s/{{ TOPIC }}/$TOPIC/g" \
| sed "s/{{ TIME }}/$TIME/g" \
> change_topic_retention.$PARTICIPANT_ID.yaml

kubectl label namespace kafka-topics istio-injection=disabled --overwrite

kubectl create -n kafka-topics -f ./change_topic_retention.$PARTICIPANT_ID.yaml
kubectl wait -n kafka-topics --timeout=90s --for=condition=complete job/$PARTICIPANT_ID-change-topic-retention