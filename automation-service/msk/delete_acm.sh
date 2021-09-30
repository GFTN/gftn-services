#!/usr/bin/env bash

PARTICIPANTID=$1
ENVIRONMENT=$2

if  [ $ENVIRONMENT == eksdev  ]; then
  ENV='-dev'
elif  [ $ENVIRONMENT == eksqa  ]; then
  ENV='-qa'
elif  [ $ENVIRONMENT == st  ]; then
  ENV='-st'
elif [ $ENVIRONMENT == prod ]; then
  ENV=''
fi

DN="$PARTICIPANTID.worldwire$ENV.io"

CERTNAME=$(aws acm list-certificates | jq '.CertificateSummaryList[].DomainName' | sed -e 's/^"//' -e 's/"$//')
CERTARN=$(aws acm list-certificates | jq '.CertificateSummaryList[].CertificateArn' | sed -e 's/^"//' -e 's/"$//')

IDX=0
for N in $CERTNAME:
do
    if [ $N == $DN ]; then
        I=0
        for A in $CERTARN:
        do
            if [ $I == $IDX ]; then
                aws acm delete-certificate --certificate-arn $A
            fi
            ((I++))
        done
    fi
    ((IDX++))
done
