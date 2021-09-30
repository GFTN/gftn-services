#!/bin/bash

# Example for the Docker Hub V2 API
# Returns all images and tags associated with a Docker Hub organization account.
# Requires 'jq': https://stedolan.github.io/jq/

# set username, password, and organization
UNAME=$1
UPASS=$2
ORG="gftn"
SVC=$3
# -------

set -e
#echo

# get token
#echo "Retrieving token ..."
TOKEN=$(curl -s -H "Content-Type: application/json" -X POST -d '{"username": "'${UNAME}'", "password": "'${UPASS}'"}' https://hub.docker.com/v2/users/login/ | jq -r .token)

# get list of repositories
#echo "Retrieving repository list ..."
REPO_LIST=$(curl -s -H "Authorization: JWT ${TOKEN}" https://hub.docker.com/v2/repositories/${ORG}/?page_size=100 | jq -r '.results|.[]|.name')

# output images & tags
#echo
#echo "Images and tags for organization: ${ORG}"
#echo
#for i in ${REPO_LIST}
#do
#  echo "${i}:"
#  # tags
#
#done
IMAGE_TAGS=$(curl -s -H "Authorization: JWT ${TOKEN}" https://hub.docker.com/v2/repositories/${ORG}/${SVC}/tags/?page_size=100 | jq -r '.results|.[]|.name')
for j in ${IMAGE_TAGS}
do
   echo "${j}"
done
echo