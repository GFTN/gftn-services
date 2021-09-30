#!/bin/bash

set -o nounset \
    -o errexit \
    -o verbose \
    -o xtrace
    
SECRET_STRING=${SECRET_STRING:-}
SECRET_NAME=${SECRET_NAME:-}
DESCRIPTION=${DESCRIPTION:-}
AWS_REGION=${AWS_REGION:-}

aws secretsmanager get-secret-value --secret-id $SECRET_NAME --region $AWS_REGION