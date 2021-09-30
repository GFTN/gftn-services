#!/bin/bash
# param 1 participant
# param 2 environment
# param 3 docker-tag
# param 4 remote_participant

# decalarations

PARTICIPANT=$1
ENVIRONMENT=$2
DOCKERTAG=$3

if [ $ENVIRONMENT == dev ]; then
  export KOPS_CLUSTER_NAME=development.worldwire-dev.io.k8s.local
  export KOPS_STATE_STORE=s3://ww-kube-dev
  ENVIRONMENT="dev"
elif  [ $ENVIRONMENT == qa  ]; then
  export KOPS_CLUSTER_NAME=qa.worldwire-qa.io.k8s.local
  export KOPS_STATE_STORE=s3://ww-kube-qa
  ENVIRONMENT="qa"
elif [ $ENVIRONMENT == st ]; then
  export KOPS_CLUSTER_NAME=staging.worldwire-st.io.k8s.local
  export KOPS_STATE_STORE=s3://kubernetes-staging
  ENVIRONMENT="st"
  CLIENT="true"
elif [ $ENVIRONMENT == tn ]; then
  export KOPS_CLUSTER_NAME=tn.worldwire-tn.io.k8s.local
  export KOPS_STATE_STORE=s3://ww-kube-tn
  ENVIRONMENT="tn"
  CLIENT="true"
elif [ $ENVIRONMENT == prod ]; then
  export KOPS_CLUSTER_NAME=prod.worldwire.io.k8s.local
  export KOPS_STATE_STORE=s3://kubernetes-worldwire
  ENVIRONMENT=""
elif [ $ENVIRONMENT == pen ]; then
  export KOPS_CLUSTER_NAME=demo.worldwire.io.k8s.local
  export KOPS_STATE_STORE=s3://kubernetes-worldwire
  ENVIRONMENT="pen"
fi

kops export kubecfg

LIST='api crypto listener rdo send'
for j in $LIST
do

if [ $j == api ]; then
  SERVICENAME="api-service"
elif  [ $j == crypto ]; then
  SERVICENAME="crypto-service"
elif [ $j ==  listener ]; then
  SERVICENAME="payment-listener"
elif [ $j == send ]; then
  SERVICENAME="send-service"
elif [ $j == rdo ]; then
  SERVICENAME="rdo-service"
fi

printf "Ready to update %s on %s k8s cluster to version %s\n" "$SERVICENAME" "$ENVIRONMENT" "$DOCKERTAG"
kubectl --record deployment.apps/$PARTICIPANT-$j set image deployment.v1.apps/$PARTICIPANT-$j $PARTICIPANT-$j=ip-team-worldwire-docker-local.artifactory.swg-devops.com/gftn/$SERVICENAME:$DOCKERTAG

done

CLIENT_LIST="callback rdo-client"

if [ "$CLIENT" == true ]; then
    for j in $CLIENT_LIST
    do

    if [ $j == callback ]; then
        SERVICENAME="callback-service"
    elif  [ $j == rdo-client ]; then
        SERVICENAME="rdo-client"
    fi

    printf "Ready to update %s on %s k8s cluster to version %s\n" "$SERVICENAME" "$ENVIRONMENT" "$DOCKERTAG"
    kubectl --record deployment.apps/$PARTICIPANT-$j set image deployment.v1.apps/$PARTICIPANT-$j $PARTICIPANT-$j=ip-team-worldwire-docker-local.artifactory.swg-devops.com/gftn/$SERVICENAME:$DOCKERTAG

    done
fi