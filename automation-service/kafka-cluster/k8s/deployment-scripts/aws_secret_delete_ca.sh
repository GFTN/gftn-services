#!/bin/bash
ENVIRONEMENT=$1
DOMAIN=${DOMAIN:-ww}
LIST='kafka_ca_cert kafka_ca_key kafka_ca_password'

for i in $LIST
do
aws secretsmanager delete-secret --secret-id /$ENVIRONEMENT/$DOMAIN/$i \
   --force-delete-without-recovery
done
