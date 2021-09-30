#!/usr/bin/env bash

PARTICIPANT=$1
ENVIRONMENT=$2
RESOURCES=$3

if [ $ENVIRONMENT == eksdev ]; then
  N_ENVIRONMENT="-dev"
  export AWS_DEFAULT_REGION="us-west-2"
  HOSTED_ZONE_ID="ZVWJ2MV8IP8AK"
  ALIASID="Z2OJLYMUO9EFXC"
  ENV="dev"
elif  [ $ENVIRONMENT == eksqa  ]; then
  N_ENVIRONMENT="-qa"
  export AWS_DEFAULT_REGION="us-west-2"
  HOSTED_ZONE_ID="Z28ZSNR7XW0ION"
  ALIASID="Z2OJLYMUO9EFXC"
  ENV="qa"
elif [ $ENVIRONMENT == st ]; then
  N_ENVIRONMENT="-st"
  export AWS_DEFAULT_REGION="us-east-2"
  HOSTED_ZONE_ID="Z2U1DCOOQTE6AH"
  ALIASID="ZOJJZC49E0EPZ"
#  ZOJJZC49E0EPZ us-east-2
  ENV="st"
elif [ $ENVIRONMENT == pen ]; then
  N_ENVIRONMENT="-pen"
  export AWS_DEFAULT_REGION="us-east-2"
  HOSTED_ZONE_ID="ZLGHTTJ2R5OBV"
  ALIASID="ZOJJZC49E0EPZ"
  ENV="pen"
elif [ $ENVIRONMENT == prod ]; then
  N_ENVIRONMENT=""
  export AWS_DEFAULT_REGION="us-east-1"
  HOSTED_ZONE_ID="Z3V7YQCN2R884C"
  ALIASID="Z1UJRXOUMOOFQ8"
  ENV="prod"
fi

if [ $RESOURCES == dynamodb ]; then
    # remove dynamodb
    TABLENAME="${ENVIRONMENT}_${PARTICIPANT}_cursor"
    aws dynamodb delete-table --table-name $TABLENAME --region=us-east-1
fi

if [ $RESOURCES == customdomainname ]; then
    # remove custom domain name
    DOMAINNAME="$PARTICIPANT.worldwire$N_ENVIRONMENT.io"
    aws apigateway delete-domain-name --domain-name=$DOMAINNAME
fi

if [ $RESOURCES == apigateway ]; then
#    # remove custom domain name
#    DOMAINNAME="$PARTICIPANT.worldwire$N_ENVIRONMENT.io"
#    aws apigateway delete-domain-name --domain-name=$DOMAINNAME
#
#    sleep 1m
    # remove api gateway rest api
    ID="_API_ID"
    FILENAME="$PARTICIPANT$ID"
    API_ID=$(cat /var/files/apigateway/$FILENAME.txt)

    aws apigateway delete-rest-api --rest-api-id $API_ID

    rm /var/files/apigateway/$FILENAME.txt
fi

if [ $RESOURCES == route53domain ]; then
    ID="_API_ID"
    FILENAME="$PARTICIPANT$ID"
    API_ID=$(cat /var/files/apigateway/$FILENAME.txt)

    # remove route53 domain
    BASE_PATH="/var/k8s"

    ROUTE53_PATH="$BASE_PATH/route53"
    mkdir -p $ROUTE53_PATH/$PARTICIPANT
    NEW_ROUTE53_PATH="$ROUTE53_PATH/$PARTICIPANT"

    cp -r $BASE_PATH/route53/delete/* $NEW_ROUTE53_PATH/
    cd $NEW_ROUTE53_PATH

    OLDPARTICIPANT=participant_id_variable
    OLDENVIRONMENT=environment_variable
    OLDAPIID=api_id_variable
    OLDAWSREGION=aws_region_variable
    OLDALIASID=alias_hosted_zone_id_variable

    # replace
    find . -type f | xargs sed -i "s/$OLDPARTICIPANT/$PARTICIPANT/g"
    find . -type f | xargs sed -i "s/$OLDENVIRONMENT/$N_ENVIRONMENT/g"
    find . -type f | xargs sed -i "s/$OLDAPIID/$API_ID/g"
    find . -type f | xargs sed -i "s/$OLDAWSREGION/$AWS_DEFAULT_REGION/g"
    find . -type f | xargs sed -i "s/$OLDALIASID/$ALIASID/g"

    ROUTE53_FILE_PATH="file://$NEW_ROUTE53_PATH/delete.json"
    aws route53 change-resource-record-sets --hosted-zone-id $HOSTED_ZONE_ID --change-batch $ROUTE53_FILE_PATH
fi