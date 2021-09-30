#!/bin/bash
export VER="1.4.2"
wget https://releases.hashicorp.com/packer/${VER}/packer_${VER}_linux_amd64.zip
unzip packer_${VER}_linux_amd64.zip
sudo chmod a+x packer
sudo mv packer /usr/local/bin
packer --version
echo 'Starting Packer build'
packer build -machine-readable packer-amazon-linux.json