#!/bin/bash

PARTICIPANTID=$1
ENVIRONMENT=$2

DOMAIN='api-service crypto-service send-service payment-service ww-gateway'

for i in $DOMAIN
do
aws secretsmanager delete-secret --secret-id /$ENVIRONMENT/$PARTICIPANTID/$i/initialize --force-delete-without-recovery
done
