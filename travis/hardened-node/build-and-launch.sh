#!/bin/bash
export VER="1.4.2"
wget https://releases.hashicorp.com/packer/${VER}/packer_${VER}_linux_amd64.zip
unzip packer_${VER}_linux_amd64.zip
sudo mv packer /usr/local/bin
packer --version
docker login -u ${DOCKER_USER} -p ${DOCKER_PASSWORD} ${DOCKER_REGISTRY}
echo 'Starting Packer build'
packer build -machine-readable packer-node-alpine.json
docker push ip-team-worldwire-docker-local.artifactory.swg-devops.com/gftn/node-alpine:latest