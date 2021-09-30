#!/bin/bash

set -o nounset \
    -o errexit \
    -o verbose \
    -o xtrace
    
SECRET_STRING=${SECRET_STRING:-}
SECRET_NAME=${SECRET_NAME:-}
DESCRIPTION=${DESCRIPTION:-}
AWS_REGION=${AWS_REGION:-}

aws secretsmanager delete-secret --secret-id $SECRET_NAME --force-delete-without-recovery --region $AWS_REGION || true
if [ $? -eq 0 ]; then
    echo "secret deleted"
    sleep 5
else
    echo "secret not deleted or exist"
fi

aws secretsmanager create-secret --name $SECRET_NAME --description $DESCRIPTION  --secret-string $SECRET_STRING --region $AWS_REGION