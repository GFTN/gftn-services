#!/bin/bash

PARTICIPANT=$1
ENVIRONMENT=$2

if [ $ENVIRONMENT == dev ]; then
  export KOPS_CLUSTER_NAME=development.worldwire-dev.io.k8s.local
  export KOPS_STATE_STORE=s3://ww-kube-dev
  N_ENVIRONMENT="-dev"
  export AWS_DEFAULT_REGION="us-west-2"
  HOSTED_ZONE_ID="ZVWJ2MV8IP8AK"
elif  [ $ENVIRONMENT == qa  ]; then
  export KOPS_CLUSTER_NAME=qa.worldwire-qa.io.k8s.local
  export KOPS_STATE_STORE=s3://ww-kube-qa
  N_ENVIRONMENT="-qa"
  export AWS_DEFAULT_REGION="ap-southeast-1"
  HOSTED_ZONE_ID="Z13NQBOATBTV50"
elif [ $ENVIRONMENT == st ]; then
  export KOPS_CLUSTER_NAME=staging.worldwire-st.io.k8s.local
  export KOPS_STATE_STORE=s3://kubernetes-staging
  N_ENVIRONMENT="-st"
  export AWS_DEFAULT_REGION="us-west-2"
  HOSTED_ZONE_ID="Z2U1DCOOQTE6AH"
elif [ $ENVIRONMENT == tn ]; then
  export KOPS_CLUSTER_NAME=tn.worldwire-tn.io.k8s.local
  export KOPS_STATE_STORE=s3://ww-kube-tn
  N_ENVIRONMENT="-tn"
  export AWS_DEFAULT_REGION="ap-southeast-1"
  HOSTED_ZONE_ID="Z1H1FK58PAH4L9"
elif [ $ENVIRONMENT == prod ]; then
  export KOPS_CLUSTER_NAME=prod.worldwire.io.k8s.local
  export KOPS_STATE_STORE=s3://kubernetes-worldwire
  N_ENVIRONMENT=""
  export AWS_DEFAULT_REGION="ap-southeast-1"
  HOSTED_ZONE_ID="Z3V7YQCN2R884C"
elif [ $ENVIRONMENT == pen ]; then
  export KOPS_CLUSTER_NAME=demo.worldwire.io.k8s.local
  export KOPS_STATE_STORE=s3://kubernetes-worldwire
  N_ENVIRONMENT="-pen"
  export AWS_DEFAULT_REGION="us-east-2"
  HOSTED_ZONE_ID="Z2U1DCOOQTE6AH"
fi

kops export kubecfg

BASE_PATH="/root/k8s"

# Delete Route53 Domain
ROUTE53_PATH="$BASE_PATH/route53"
NEW_ROUTE53_PATH="$ROUTE53_PATH/$PARTICIPANT"

cd $NEW_ROUTE53_PATH

OLDACTION=CREATE
NEWACTION=DELETE

# replace
find . -type f | xargs sed -i "s/$OLDACTION/$NEWACTION/g"

# Delete Route 53 record
ROUTE53_FILE_PATH="file://$NEW_ROUTE53_PATH/config.json"
aws route53 change-resource-record-sets --hosted-zone-id $HOSTED_ZONE_ID --change-batch $ROUTE53_FILE_PATH

# Delete Custom Domain Name
CUSTOM_DOMAIN_NAME="$PARTICIPANT.worldwire$N_ENVIRONMENT.io"
DOMAIN_NAME=$(aws apigateway delete-domain-name --domain-name=$CUSTOM_DOMAIN_NAME)

# Delete API Gateway settings
API_GATEWAY_PATH="$BASE_PATH/api-gateway/$PARTICIPANT"

ID="_API_ID"
FILENAME="$PARTICIPANT$ID"
APIID=`cat $FILENAME.txt`
echo APIID

aws apigateway delete-rest-api --rest-api-id=$APIID

# Delete AWS secret
LIST='api-service crypto-service send-service rdo-service payment-service participant'
LIST_KAFKA='kafka_key kafka_cert kafka_password'

for i in $LIST
do
aws secretsmanager delete-secret --secret-id /$ENVIRONMENT/$PARTICIPANT/$i/initialize \
   --force-delete-without-recovery
done

for i in $LIST_KAFKA
do
aws secretsmanager delete-secret --secret-id /$ENVIRONMENT/$PARTICIPANT/$i \
   --force-delete-without-recovery
done