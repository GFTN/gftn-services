#!/usr/bin/env bash

PARTICIPANTID=$1
ENVIRONMENT=$2
ORGID=$3

if  [ $ENVIRONMENT == eksdev  ]; then
  aws eks update-kubeconfig --name terraform-eks-ww-dev --region us-west-2
  kubectl config use-context arn:aws:eks:us-west-2:$ORGID:cluster/terraform-eks-ww-dev
elif  [ $ENVIRONMENT == eksqa  ]; then
  aws eks update-kubeconfig --name terraform-eks-ww-qa --region us-west-2
  kubectl config use-context arn:aws:eks:us-west-2:$ORGID:cluster/terraform-eks-ww-qa
elif  [ $ENVIRONMENT == st  ]; then
  aws eks update-kubeconfig --name terraform-eks-ww-st --region us-east-2
  kubectl config use-context arn:aws:eks:us-east-2:$ORGID:cluster/terraform-eks-ww-st
elif  [ $ENVIRONMENT == pen  ]; then
  aws eks update-kubeconfig --name terraform-eks-ww-pen --region us-east-2
  kubectl config use-context arn:aws:eks:us-east-2:$ORGID:cluster/terraform-eks-ww-pen
elif [ $ENVIRONMENT == prod ]; then
  aws eks update-kubeconfig --name terraform-eks-ww-prod --region us-east-1
  kubectl config use-context arn:aws:eks:us-east-1:$ORGID:cluster/terraform-eks-ww-prod
fi

SVCS="api-service crypto-service send-service payment-service ww-gateway"

for SERVICENAME in $SVCS
do
export POLICYNAME="${ENVIRONMENT}_${PARTICIPANTID}_${SERVICENAME}_policy"
POLICYARN=$(aws iam list-policies | jq '.Policies[] | select(.PolicyName == env.POLICYNAME) | .Arn' | sed -e 's/^"//' -e 's/"$//')
IAMUSERNAME="${ENVIRONMENT}_${PARTICIPANTID}_${SERVICENAME}"
ACCESSKEYID=$(aws iam list-access-keys --user-name $IAMUSERNAME | jq '.AccessKeyMetadata[].AccessKeyId' | sed -e 's/^"//' -e 's/"$//')

for KEY in $ACCESSKEYID
do
aws iam delete-access-key --access-key-id $KEY --user-name $IAMUSERNAME
done

LOCALGROUPNAME="${ENVIRONMENT}_local_service"
COMMONGROUPNAME="${ENVIRONMENT}_service"

aws iam detach-user-policy --user-name $IAMUSERNAME --policy-arn $POLICYARN
aws iam delete-policy --policy-arn $POLICYARN
aws iam remove-user-from-group --group-name $LOCALGROUPNAME --user-name $IAMUSERNAME
aws iam remove-user-from-group --group-name $COMMONGROUPNAME --user-name $IAMUSERNAME
aws iam delete-user --user-name $IAMUSERNAME
kubectl delete secret --namespace=default $PARTICIPANTID-$SERVICENAME-aws-iam-key
done

