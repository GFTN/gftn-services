#!/bin/bash

PARTICIPANT=$1
ENVIRONMENT=$2
RESOURCES=$3
ORGID=$4

if [ $ENVIRONMENT == eksdev ]; then
  aws eks update-kubeconfig --name terraform-eks-ww-dev --region us-west-2
  kubectl config use-context arn:aws:eks:us-west-2:$ORGID:cluster/terraform-eks-ww-dev
  N_ENVIRONMENT="-dev"
  export AWS_DEFAULT_REGION="us-west-2"
  HOSTED_ZONE_ID="ZVWJ2MV8IP8AK"
  ALIASID="Z2OJLYMUO9EFXC"
  ENV="dev"
elif  [ $ENVIRONMENT == eksqa  ]; then
  aws eks update-kubeconfig --name terraform-eks-ww-qa --region us-west-2
  kubectl config use-context arn:aws:eks:us-west-2:$ORGID:cluster/terraform-eks-ww-qa
  N_ENVIRONMENT="-qa"
  export AWS_DEFAULT_REGION="us-west-2"
  HOSTED_ZONE_ID="Z28ZSNR7XW0ION"
  ALIASID="Z2OJLYMUO9EFXC"
  ENV="qa"
elif [ $ENVIRONMENT == st ]; then
  aws eks update-kubeconfig --name terraform-eks-ww-st --region us-east-2
  kubectl config use-context arn:aws:eks:us-east-2:$ORGID:cluster/terraform-eks-ww-st
  N_ENVIRONMENT="-st"
  export AWS_DEFAULT_REGION="us-east-2"
  HOSTED_ZONE_ID="Z2U1DCOOQTE6AH"
  ALIASID="ZOJJZC49E0EPZ"
#  ZOJJZC49E0EPZ us-east-2
  ENV="st"
elif [ $ENVIRONMENT == pen ]; then
  aws eks update-kubeconfig --name terraform-eks-ww-pen --region us-east-2
  kubectl config use-context arn:aws:eks:us-east-2:$ORGID:cluster/terraform-eks-ww-pen
  N_ENVIRONMENT="-pen"
  export AWS_DEFAULT_REGION="us-east-2"
  HOSTED_ZONE_ID="ZLGHTTJ2R5OBV"
  ALIASID="ZOJJZC49E0EPZ"
  ENV="pen"
elif [ $ENVIRONMENT == prod ]; then
  aws eks update-kubeconfig --name terraform-eks-ww-prod --region us-east-1
  kubectl config use-context arn:aws:eks:us-east-1:$ORGID:cluster/terraform-eks-ww-prod
  N_ENVIRONMENT=""
  export AWS_DEFAULT_REGION="us-east-1"
  HOSTED_ZONE_ID="Z3V7YQCN2R884C"
  ALIASID="Z1UJRXOUMOOFQ8"
  ENV="prod"
fi

BASE_PATH="/var/k8s"
OLDPARTICIPANT=participant_id_variable
OLDENVIRONMENT=environment_variable

echo $AWS_DEFAULT_REGION

if [ $RESOURCES == apigateway ]; then
    API_GATEWAY_PATH="$BASE_PATH/api-gateway"
    mkdir -p $API_GATEWAY_PATH/$PARTICIPANT
    NEW_API_GATEWAY_PATH="$API_GATEWAY_PATH/$PARTICIPANT"

    # copy files into new files
    cp -r $BASE_PATH/api-gateway/template/* $NEW_API_GATEWAY_PATH/
    cd $NEW_API_GATEWAY_PATH

    # replace
    find . -type f | xargs sed -i "s/$OLDPARTICIPANT/$PARTICIPANT/g"
    find . -type f | xargs sed -i "s/$OLDENVIRONMENT/$N_ENVIRONMENT/g"

    API_GATEWAY_FILE_PATH="file://$NEW_API_GATEWAY_PATH/aws-api-gateway.yaml"

    # Import API Gateway settings
    API_ID=$(aws apigateway import-rest-api --body $API_GATEWAY_FILE_PATH | jq '.id' | sed -e 's/^"//' -e 's/"$//')
    if [ -z $API_ID ]; then
        exit 1
    fi

    ID="_API_ID"
    FILENAME="$PARTICIPANT$ID"
    mkdir -p "/var/files/apigateway"
    echo $API_ID >> "/var/files/apigateway/$FILENAME.txt"

    # Get NLB ARN
    NLB_ID=$(kubectl get svc -n istio-system istio-ingressgateway -o json | jq '.status.loadBalancer.ingress[].hostname' | sed -e 's/^"//' -e 's/"$//' | cut -d'-' -f1)
    NLB_ARN=$(aws elbv2 describe-load-balancers --names=$NLB_ID | jq '.LoadBalancers[].LoadBalancerArn' | sed -e 's/^"//' -e 's/"$//')

    # If VPC-link not exist, create a new VPC-link using NLB ARN
    export VPC_LINK_NAME="istio-nlb-link$N_ENVIRONMENT"
    #export VPC_LINK_NAME="$ENVIRONMENT-link"
    VPC_LINK_ID=$(aws apigateway get-vpc-links | jq '.items[] | select(.name == env.VPC_LINK_NAME) | select(.status == "AVAILABLE") | .id' | sed -e 's/^"//' -e 's/"$//')
    if [ -z $VPC_LINK_ID ]; then
        VPC_LINK_ID=$(aws apigateway create-vpc-link --name=$VPC_LINK_NAME --target-arns=$NLB_ARN | jq '.id' | sed -e 's/^"//' -e 's/"$//')
    fi

    # Deploy API
    aws apigateway create-deployment --rest-api-id=$API_ID --stage-name=$ENV --variables environment=$ENV,global='global',vpcLinkId=$VPC_LINK_ID,participant=$PARTICIPANT | jq '.id' | sed -e 's/^"//' -e 's/"$//'
fi

if [ $RESOURCES == customdomainname ]; then
    # Get ACM Certificate
    ORIGINAL_REGION=$AWS_DEFAULT_REGION
    export AWS_DEFAULT_REGION="us-east-1"
    CUSTOM_DOMAIN_NAME="$PARTICIPANT.worldwire$N_ENVIRONMENT.io"
    export CERT_NAME="*.worldwire$N_ENVIRONMENT.io"
    CERT_ARN=$(aws acm list-certificates | jq '.CertificateSummaryList[] | select(.DomainName == env.CERT_NAME) | .CertificateArn' | sed -e 's/^"//' -e 's/"$//')

    export AWS_DEFAULT_REGION=$ORIGINAL_REGION

    TLSVERSION="TLS_1_2"
    # Create Custom Domain Name
    DOMAIN_NAME=$(aws apigateway create-domain-name --domain-name=$CUSTOM_DOMAIN_NAME --certificate-arn=$CERT_ARN --security-policy=$TLSVERSION --endpoint-configuration types=EDGE | jq '.domainName' | sed -e 's/^"//' -e 's/"$//')

    ID="_API_ID"
    FILENAME="$PARTICIPANT$ID"
    API_ID=$(cat /var/files/apigateway/$FILENAME.txt)
    # Create base-path mapping
    aws apigateway create-base-path-mapping --domain-name=$DOMAIN_NAME --rest-api-id=$API_ID --stage=$ENV
fi

if [ $RESOURCES == route53domain ]; then
    # Create Route53 Domain
    ROUTE53_PATH="$BASE_PATH/route53"
    mkdir -p $ROUTE53_PATH/$PARTICIPANT
    NEW_ROUTE53_PATH="$ROUTE53_PATH/$PARTICIPANT"

    # copy files into new files
    cp -r $BASE_PATH/route53/template/* $NEW_ROUTE53_PATH/
    cd $NEW_ROUTE53_PATH

    OLDAPIID=api_id_variable
    OLDAWSREGION=aws_region_variable
    OLDALIASID=alias_hosted_zone_id_variable

    ID="_API_ID"
    FILENAME="$PARTICIPANT$ID"
    API_ID=$(cat /var/files/apigateway/$FILENAME.txt)

    # replace
    find . -type f | xargs sed -i "s/$OLDPARTICIPANT/$PARTICIPANT/g"
    find . -type f | xargs sed -i "s/$OLDENVIRONMENT/$N_ENVIRONMENT/g"
    find . -type f | xargs sed -i "s/$OLDAPIID/$API_ID/g"
    find . -type f | xargs sed -i "s/$OLDAWSREGION/$AWS_DEFAULT_REGION/g"
    find . -type f | xargs sed -i "s/$OLDALIASID/$ALIASID/g"

    ROUTE53_FILE_PATH="file://$NEW_ROUTE53_PATH/config.json"
    aws route53 change-resource-record-sets --hosted-zone-id $HOSTED_ZONE_ID --change-batch $ROUTE53_FILE_PATH
fi

if [ $RESOURCES == dynamodb ]; then
    #create DynamoDB table for payment-listener
    export AWS_DEFAULT_REGION="us-east-1"
    TABLENAME="${ENVIRONMENT}_${PARTICIPANT}_cursor"
    cd $BASE_PATH/script
    dir=$(pwd)
    aws dynamodb create-table --attribute-definitions=file://$dir/dynamodb/attribute.json --table-name=$TABLENAME --key-schema=file://$dir/dynamodb/schema.json --provisioned-throughput=file://$dir/dynamodb/throughput.json
fi