#!/bin/bash
export VER="1.4.2"
wget https://releases.hashicorp.com/packer/${VER}/packer_${VER}_linux_amd64.zip
unzip packer_${VER}_linux_amd64.zip
sudo mv packer /usr/local/bin
packer --version
echo 'Starting Packer build'
packer build -machine-readable packer-alpine-linux.json
docker login -u ${DOCKER_USER} -p ${DOCKER_PASSWORD} ${DOCKER_REGISTRY}
sudo curl -k -ssl -u ${TL_USER}:${TL_PASS} --output /tmp/twistcli ${TL_CONSOLE_URL}/api/v1/util/twistcli
sudo chmod a+x /tmp/twistcli
sudo /tmp/twistcli images scan ip-team-worldwire-docker-local.artifactory.swg-devops.com/gftn/alpine:latest --details -address ${TL_CONSOLE_URL} -u ${TL_USER} -p ${TL_PASS} || exit 1
docker push ip-team-worldwire-docker-local.artifactory.swg-devops.com/gftn/alpine:latest
