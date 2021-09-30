#!/bin/bash

ORGID=$1

aws eks update-kubeconfig --name terraform-eks-ww-dev --region us-west-2
kubectl config use-context arn:aws:eks:us-west-2:$ORGID:cluster/terraform-eks-ww-dev
export AWS_DEFAULT_REGION=us-west-2

CA_ARN=$(aws acm-pca list-certificate-authorities | jq '.CertificateAuthorities[].Arn' | sed -e 's/^"//' -e 's/"$//')

aws acm-pca get-certificate-authority-certificate --certificate-authority-arn $CA_ARN | jq '.Certificate' | sed -e 's/^"//' -e 's/"$//' -e 's/\\n/\n/g' >> ./ca_cert.crt

export CERT_DN="msk.worldwire-dev.io"
CERT_ARN=$(aws acm request-certificate --domain-name $CERT_DN --certificate-authority-arn $CA_ARN | jq '.CertificateArn' | sed -e 's/^"//' -e 's/"$//' )

#sleep 10s
#
echo $CERT_ARN
## get the private certificate and store in a file
#aws acm get-certificate --certificate-arn $CERT_ARN | jq '.Certificate' | sed -e 's/^"//' -e 's/"$//' -e 's/\\n/\n/g' >> ./kafka_cert.crt
#
## key password for private key
#KAFKA_KEY_PASSWORD="msk-worldwire-dev"
#
## get the private key and store in a file
#aws acm export-certificate --certificate-arn $CERT_ARN --passphrase $KAFKA_KEY_PASSWORD | jq '.PrivateKey' | sed -e 's/^"//' -e 's/"$//' -e 's/\\n/\n/g' >> ./kafka_key.key
#
#openssl rsa -in ./kafka_key.key -out ./decrypted_key.key

## store the cert and private key into the k8s secret
kubectl create secret generic -n kafka-topics kafka-secret-global --from-file=./ca_cert.crt --from-file=./kafka_cert.crt --from-file=./decrypted_key.key
#
## remove certificate and private key
#rm ./ca_cert.crt ./kafka_cert.crt ./kafka_key.key ./decrypted_key.key
#
#kubectl apply -f burrow-svc.yaml
#$HOME/Downloads/istio-1.2.4/bin/istioctl kube-inject -f ./burrow-deployment.yaml | kubectl apply -n kafka-topics -f  -