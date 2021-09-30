#!/bin/sh
# sudo yum check-update
sudo yum update kernel -y
echo "kernel updated"
sudo yum update -y && sudo yum upgrade -y
sudo package-cleanup --oldkernels --count=1 -y
# sudo yum install -y yum-utils device-mapper-persistent-data lvm2
# sudo yum-config-manager --add-repo https://download.docker.com/linux/centos/docker-ce.repo
sudo yum install -y docker
echo "installed docker"
sleep 10
sudo systemctl start docker
echo "started docker"
sleep 10
# echo "starting twistlock"
# curl -sSL -k --header "authorization: Bearer $TL_HEADER" https://ec2-18-139-84-127.ap-southeast-1.compute.amazonaws.com:8083/api/v1/scripts/defender.sh | sudo bash -s -- -c "ec2-18-139-84-127.ap-southeast-1.compute.amazonaws.com" -d "none"  --install-host
# echo "ending twistlock"
# sleep 30