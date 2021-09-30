#!/bin/bash

updateAWSSecret()
{
    VERSION=`cat ../../VERSION`
    echo $VERSION
    if [ $TRAVIS_BRANCH == 'dev-gftn' ] && [ $TRAVIS_PULL_REQUEST == 'false' ]; then
        ENVIRONMENT="eksdev"
        MSKNAME="ww-dev-msk"
        export REGION='us-west-2' && export CLUSTER_NAME='terraform-eks-ww-dev' && aws eks update-kubeconfig --name $CLUSTER_NAME --region $REGION
    elif [ $TRAVIS_BRANCH == 'qa-gftn' ] && [ $TRAVIS_PULL_REQUEST == 'false' ]; then
        ENVIRONMENT="eksqa"
        MSKNAME="ww-qa-msk"
        export REGION='us-west-2' && export CLUSTER_NAME='terraform-eks-ww-qa' && aws eks update-kubeconfig --name $CLUSTER_NAME --region $REGION
    elif [ $TRAVIS_BRANCH == 'staging-gftn' ] && [ $TRAVIS_PULL_REQUEST == 'false' ]; then
        ENVIRONMENT="st"
        MSKNAME="ww-st-msk"
        export REGION='us-east-2' && export CLUSTER_NAME='terraform-eks-ww-st' && aws eks update-kubeconfig --name $CLUSTER_NAME --region $REGION
    elif [ $TRAVIS_BRANCH == 'test-gftn' ] && [ $TRAVIS_PULL_REQUEST == 'false' ]; then
        ENVIRONMENT="pen"
        MSKNAME="ww-pen-msk"
        export REGION='us-east-2' && export CLUSTER_NAME='terraform-eks-ww-pen' && aws eks update-kubeconfig --name $CLUSTER_NAME --region $REGION
    elif [ $TRAVIS_BRANCH == 'live-gftn' ] && [ $TRAVIS_PULL_REQUEST == 'false' ]; then
        ENVIRONMENT="prod"
        MSKNAME="ww-prod-msk"
        export REGION='us-east-1' && export CLUSTER_NAME='terraform-eks-ww-prod' && aws eks update-kubeconfig --name $CLUSTER_NAME --region $REGION
    else
        echo " >>>>> INFO <<<<<: No conditions were matched for this deployment. This could be becuase this is not a branch that enables deployment or this is a pull-request(not a merge.)"
        exit 1
    fi

    # get Kafka brokers url
    CLUSTERARN=$(aws kafka list-clusters --cluster-name-filter $MSKNAME --region $REGION | jq '.ClusterInfoList[].ClusterArn' | sed -e 's/^"//' -e 's/"$//')
    BROKERSURL=$(aws kafka get-bootstrap-brokers --cluster-arn=$CLUSTERARN --region $REGION | jq '.BootstrapBrokerStringTls' | sed -e 's/^"//' -e 's/"$//')

    export AWS_DEFAULT_REGION=ap-southeast-1
    export SECRETNAME=/$ENVIRONMENT/aws/secret/$VERSION

    echo $SECRETNAME

    # get AWS secret template arn
    ARN=$(aws secretsmanager list-secrets | jq '.SecretList[] | select(.Name == env.SECRETNAME) | .ARN' | sed -e 's/^"//' -e 's/"$//')
    if [ -z $ARN ]; then
        echo " >>>>> INFO <<<<<: Can not found the template to create AWS secret"
        exit 1
    fi

    # get encoded AWS secret template and decode it to plaintext
    SECRETSTR=$(aws secretsmanager get-secret-value --secret-id $ARN | jq '.SecretString' | sed -e 's/^"//' -e 's/"$//' | sed 's/\\//g')
    ENCODEDSECRETTEMPLATE=$(echo $SECRETSTR | jq '.file' | sed -e 's/^"//' -e 's/"$//')
    DECODEDSECRETTEMPLATE=$(echo $ENCODEDSECRETTEMPLATE | base64 --decode)

    DEPLOYMENTLIST=$(kubectl get deployment -o go-template --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')

    for D in $DEPLOYMENTLIST
    do
        printf "Ready to delete AWS secret for %s on %s k8s cluster\n" "$D" "$TRAVIS_BRANCH"

        IFS='-' # hyphen (-) is set as delimiter
        read -ra NAME <<< "$D"
        IFS=''

        PARTICIPANTID='EMPTY'
        for N in ${NAME[@]}
        do
          if [ $N == 'api' ]; then
              SERVICENAME="api-service"
              break
          elif  [ $N == 'crypto'  ]; then
              SERVICENAME="crypto-service"
              break
          elif  [ $N == 'gateway'  ]; then
              SERVICENAME="ww-gateway"
              break
          elif  [ $N == 'listener'  ]; then
              SERVICENAME="payment-service"
              break
          elif  [ $N == 'send'  ]; then
              SERVICENAME="send-service"
              break
          elif  [ $N == 'administration'  ]; then
              SERVICENAME="admin-service"
              break
          elif  [ $N == 'anchor'  ]; then
              SERVICENAME="anchor-service"
              break
          elif  [ $N == 'fee'  ]; then
              SERVICENAME="fee-service"
              break
          elif  [ $N == 'gas'  ]; then
              SERVICENAME="gas-service"
              break
          elif  [ $N == 'payout'  ]; then
              SERVICENAME="payout-service"
              break
          elif  [ $N == 'pr'  ]; then
              SERVICENAME="pr-service"
              break
          elif  [ $N == 'quotes'  ]; then
              SERVICENAME="quotes-service"
              break
          elif [ $N == 'whitelist' ]; then
              SERVICENAME="whitelist-service"
              break
          else
            if [ $PARTICIPANTID == 'EMPTY' ]; then
              PARTICIPANTID=${N}
            else
              PARTICIPANTID="${PARTICIPANTID}-$N"
            fi
          fi
        done

        SECRETID="/$ENVIRONMENT/$PARTICIPANTID/$SERVICENAME/initialize"

        # Delete secret with id $SECRETID from AWS secret manager
        echo $SECRETID
        aws secretsmanager delete-secret --secret-id $SECRETID --force-delete-without-recovery
    done

    IFS=$'\n'
    read -rd '' -a DEPLOYMENTLIST2 <<< "$DEPLOYMENTLIST"
    IFS=''

    for D2 in ${DEPLOYMENTLIST2[@]}
    do
        printf "Ready to update AWS secret for %s on %s k8s cluster\n" "$D2" "$TRAVIS_BRANCH"

        IFS='-' # hyphen (-) is set as delimiter
        read -ra NAME <<< "$D2"
        IFS=''

        PARTICIPANTID='EMPTY'
        for N in ${NAME[@]}
        do
          if [ $N == 'api' ]; then
              SERVICENAME="api-service"
              break
          elif  [ $N == 'crypto'  ]; then
              SERVICENAME="crypto-service"
              break
          elif  [ $N == 'gateway'  ]; then
              SERVICENAME="ww-gateway"
              break
          elif  [ $N == 'listener'  ]; then
              SERVICENAME="payment-service"
              break
          elif  [ $N == 'send'  ]; then
              SERVICENAME="send-service"
              break
          elif  [ $N == 'administration'  ]; then
              SERVICENAME="admin-service"
              break
          elif  [ $N == 'anchor'  ]; then
              SERVICENAME="anchor-service"
              break
          elif  [ $N == 'fee'  ]; then
              SERVICENAME="fee-service"
              break
          elif  [ $N == 'gas'  ]; then
              SERVICENAME="gas-service"
              break
          elif  [ $N == 'payout'  ]; then
              SERVICENAME="payout-service"
              break
          elif  [ $N == 'pr'  ]; then
              SERVICENAME="pr-service"
              break
          elif  [ $N == 'quotes'  ]; then
              SERVICENAME="quotes-service"
              break
          elif [ $N == 'whitelist' ]; then
              SERVICENAME="whitelist-service"
              break
          else
            if [ $PARTICIPANTID == 'EMPTY' ]; then
              PARTICIPANTID=${N}
            else
              PARTICIPANTID="${PARTICIPANTID}-$N"
            fi
          fi
        done

        echo $SECRETID
        SECRETID="/$ENVIRONMENT/$PARTICIPANTID/$SERVICENAME/initialize"

        if [ $PARTICIPANTID == 'ww' ]; then
            SECRET=$(echo $DECODEDSECRETTEMPLATE | jq '.'$ENVIRONMENT'.ww["'$SERVICENAME'"]' | sed -e "s/kafka_broker_url/${BROKERSURL}/g")
        else
            SECRET=$(echo $DECODEDSECRETTEMPLATE | jq '.'$ENVIRONMENT'.participant["'$SERVICENAME'"]' | sed -e "s/kafka_broker_url/${BROKERSURL}/g")
        fi

        # Create secret with id $SECRETID with value $SECRET into AWS secret manager
        until aws secretsmanager create-secret --name $SECRETID --secret-string "$SECRET" --description "World Wire Service Secret"; do
          echo "Wait 5 seconds for the secret to be deleted"
          sleep 5
        done
    done
}

updateIAMPolicy()
{
    if [ $TRAVIS_BRANCH == 'dev-gftn' ] && [ $TRAVIS_PULL_REQUEST == 'false' ]; then
        ENVIRONMENT="eksdev"
        FOLDERNAME="eks-dev"
        export REGION='us-west-2' && export CLUSTER_NAME='terraform-eks-ww-dev' && aws eks update-kubeconfig --name $CLUSTER_NAME --region $REGION
    elif [ $TRAVIS_BRANCH == 'qa-gftn' ] && [ $TRAVIS_PULL_REQUEST == 'false' ]; then
        ENVIRONMENT="eksqa"
        FOLDERNAME="eks-qa"
        export REGION='us-west-2' && export CLUSTER_NAME='terraform-eks-ww-qa' && aws eks update-kubeconfig --name $CLUSTER_NAME --region $REGION
    elif [ $TRAVIS_BRANCH == 'staging-gftn' ] && [ $TRAVIS_PULL_REQUEST == 'false' ]; then
        ENVIRONMENT="st"
        FOLDERNAME="eks-st"
        export REGION='us-east-2' && export CLUSTER_NAME='terraform-eks-ww-st' && aws eks update-kubeconfig --name $CLUSTER_NAME --region $REGION
    elif [ $TRAVIS_BRANCH == 'test-gftn' ] && [ $TRAVIS_PULL_REQUEST == 'false' ]; then
        ENVIRONMENT="pen"
        FOLDERNAME="eks-pen"
        export REGION='us-east-2' && export CLUSTER_NAME='terraform-eks-ww-pen' && aws eks update-kubeconfig --name $CLUSTER_NAME --region $REGION
    elif [ $TRAVIS_BRANCH == 'live-gftn' ] && [ $TRAVIS_PULL_REQUEST == 'false' ]; then
        ENVIRONMENT="prod"
        FOLDERNAME="eks-prod"
        export REGION='us-east-1' && export CLUSTER_NAME='terraform-eks-ww-prod' && aws eks update-kubeconfig --name $CLUSTER_NAME --region $REGION
    else
        echo " >>>>> INFO <<<<<: No conditions were matched for this deployment. This could be becuase this is not a branch that enables deployment or this is a pull-request(not a merge.)"
        exit 1
    fi

    DBREGION="us-east-1"
    SECRETMANAGERREGION="ap-southeast-1"
    AWSORGID="000000000000"

    NEWFOLDER="../../automation-service/k8s/iam-policy/policies"
    TEMPLATEFOLDER="../../automation-service/k8s/iam-policy"

    # create a new folder to store the new policy
    mkdir -p $NEWFOLDER

    DEPLOYMENTLIST=$(kubectl get deployment -o go-template --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')

    for D in $DEPLOYMENTLIST
    do
        printf "Ready to update IAM policy for %s on %s k8s cluster\n" "$D" "$TRAVIS_BRANCH"

        IFS='-' # hyphen (-) is set as delimiter
        read -ra NAME <<< "$D"
        IFS=''

        PARTICIPANTID='EMPTY'
        for N in ${NAME[@]}
        do
          if [ $N == 'api' ]; then
              SERVICENAME="api-service"
              break
          elif  [ $N == 'crypto'  ]; then
              SERVICENAME="crypto-service"
              break
          elif  [ $N == 'gateway'  ]; then
              SERVICENAME="ww-gateway"
              break
          elif  [ $N == 'listener'  ]; then
              SERVICENAME="payment-service"
              break
          elif  [ $N == 'send'  ]; then
              SERVICENAME="send-service"
              break
          elif  [ $N == 'administration'  ]; then
              SERVICENAME="admin-service"
              break
          elif  [ $N == 'anchor'  ]; then
              SERVICENAME="anchor-service"
              break
          elif  [ $N == 'fee'  ]; then
              SERVICENAME="fee-service"
              break
          elif  [ $N == 'gas'  ]; then
              SERVICENAME="gas-service"
              break
          elif  [ $N == 'payout'  ]; then
              SERVICENAME="payout-service"
              break
          elif  [ $N == 'pr'  ]; then
              SERVICENAME="pr-service"
              break
          elif  [ $N == 'quotes'  ]; then
              SERVICENAME="quotes-service"
              break
          elif [ $N == 'whitelist' ]; then
              SERVICENAME="whitelist-service"
              break
          else
            if [ $PARTICIPANTID == 'EMPTY' ]; then
              PARTICIPANTID=${N}
            else
              PARTICIPANTID="${PARTICIPANTID}-$N"
            fi
          fi
        done

        DESC="policy for /$ENVIRONMENT/$PARTICIPANTID/$SERVICENAME"
        IAMUSERNAME="${ENVIRONMENT}_${PARTICIPANTID}_${SERVICENAME}"

        export POLICYNAME="${IAMUSERNAME}_policy"
        POLICYARN=$(aws iam list-policies | jq '.Policies[] | select(.PolicyName == env.POLICYNAME) | .Arn' | sed -e 's/^"//' -e 's/"$//')
        ACCESSKEYID=$(aws iam list-access-keys --user-name $IAMUSERNAME | jq '.AccessKeyMetadata[].AccessKeyId' | sed -e 's/^"//' -e 's/"$//')

        printf " >>>>> INFO <<<<<: Delete IAM policy %s for user %s\n" "$POLICYNAME" "$IAMUSERNAME"

        IFS=$'\n'
        read -rd '' -a KEYLIST <<< "$ACCESSKEYID"
        IFS=''

        # Delete all the user access keys
        for KEY in ${KEYLIST[@]}
        do
            echo $KEY
            aws iam delete-access-key --access-key-id $KEY --user-name $IAMUSERNAME
        done

        GLOBALGROUPNAME="${ENVIRONMENT}_global_service"
        LOCALGROUPNAME="${ENVIRONMENT}_local_service"
        COMMONGROUPNAME="${ENVIRONMENT}_service"

        # Detach user from the policy
        aws iam detach-user-policy --user-name $IAMUSERNAME --policy-arn $POLICYARN

        # If there are multiple policy version, first delete all the non-default versions
        POLICYVERSIONS=$(aws iam list-policy-versions --policy-arn $POLICYARN | jq '.Versions[] | select(.IsDefaultVersion == 'false') | .VersionId' | sed -e 's/^"//' -e 's/"$//')

        IFS=$'\n'
        read -rd '' -a VERSIONLIST <<< "$POLICYVERSIONS"
        IFS=''

        for V in ${VERSIONLIST[@]}
        do
            printf "None default policy version is: $V\n"
            aws iam delete-policy-version --policy-arn $POLICYARN --version-id $V
        done

        # Delete the policy
        aws iam delete-policy --policy-arn $POLICYARN

        # Remove user from the group
        if [ $PARTICIPANTID == 'ww' ]; then
            aws iam remove-user-from-group --group-name $GLOBALGROUPNAME --user-name $IAMUSERNAME
        else
            aws iam remove-user-from-group --group-name $LOCALGROUPNAME --user-name $IAMUSERNAME
        fi
        aws iam remove-user-from-group --group-name $COMMONGROUPNAME --user-name $IAMUSERNAME

        # Delete the user
        aws iam delete-user --user-name $IAMUSERNAME

        # Delete the credential from the k8s secret
        kubectl delete secret --namespace=default $PARTICIPANTID-$SERVICENAME-aws-iam-key

        # Define name of the resources
        POLICYBOUNDARY="\/$ENVIRONMENT\/$PARTICIPANTID\/$SERVICENAME\/*"
        PARTICIPANTBOUNDARY="\/$ENVIRONMENT\/$PARTICIPANTID\/participant\/*"
        ACCOUNTSECRET="\/$ENVIRONMENT\/$PARTICIPANTID\/account\/*"
        TOKENSECRET="\/$ENVIRONMENT\/ww\/account\/token*"
        ADMINKILLSWITCH="\/$ENVIRONMENT\/*\/killswitch\/accounts*"
        KILLSWITCH="\/$ENVIRONMENT\/$PARTICIPANTID\/killswitch\/accounts*"

        cp "$TEMPLATEFOLDER/$SERVICENAME.json" $NEWFOLDER

        cd $NEWFOLDER
        # replace
        find . -type f | xargs sed -i "s/<env>/$ENVIRONMENT/g"
        find . -type f | xargs sed -i "s/<participantId>/$PARTICIPANTID/g"
        find . -type f | xargs sed -i "s/<db_region>/$DBREGION/g"
        find . -type f | xargs sed -i "s/<region>/$SECRETMANAGERREGION/g"
        find . -type f | xargs sed -i "s/<accountId>/$AWSORGID/g"
        find . -type f | xargs sed -i "s/<service_secret>/$POLICYBOUNDARY/g"
        find . -type f | xargs sed -i "s/<participant_secret>/$PARTICIPANTBOUNDARY/g"
        find . -type f | xargs sed -i "s/<account_secret>/$ACCOUNTSECRET/g"
        find . -type f | xargs sed -i "s/<token_secret>/$TOKENSECRET/g"
        find . -type f | xargs sed -i "s/<killswitch_secret>/$KILLSWITCH/g"
        find . -type f | xargs sed -i "s/<admin_killswitch_secret>/$ADMINKILLSWITCH/g"

        printf " >>>>> INFO <<<<<: Create IAM policy %s for user %s\n" "$POLICYNAME" "$IAMUSERNAME"

        # Create the policy base on the policy file: ${SERVICENAME}.json
        POLICYARN=$(aws iam create-policy --description "$DESC" --policy-name $POLICYNAME --policy-document file://${SERVICENAME}.json | jq '.Policy.Arn' | sed -e 's/^"//' -e 's/"$//')
        if [ -z $POLICYARN ]; then
                echo "Empty POLICYARN because the policy was already there"
                POLICYARN=$(aws iam list-policies | jq '.Policies[] | select(.PolicyName == env.POLICYNAME) | .Arn' | sed -e 's/^"//' -e 's/"$//')
        fi

        # Create the user
        aws iam create-user --user-name $IAMUSERNAME

        # Add user to the group
        if [ $PARTICIPANTID == 'ww' ]; then
            aws iam add-user-to-group --group-name $GLOBALGROUPNAME --user-name $IAMUSERNAME
        else
            aws iam add-user-to-group --group-name $LOCALGROUPNAME --user-name $IAMUSERNAME
        fi
        aws iam add-user-to-group --group-name $COMMONGROUPNAME --user-name $IAMUSERNAME

        # Attach a policy to the user
        aws iam attach-user-policy --user-name $IAMUSERNAME --policy-arn $POLICYARN

        # Get the access credential and store into the k8s secret
        ACCESSKEY=$(aws iam create-access-key --user-name $IAMUSERNAME)
        AWSACCESSKEYID=$(echo $ACCESSKEY | jq '.AccessKey.AccessKeyId' | sed -e 's/^"//' -e 's/"$//')
        AWSSECRETACCESSKEY=$(echo $ACCESSKEY | jq '.AccessKey.SecretAccessKey' | sed -e 's/^"//' -e 's/"$//')

        # store iam ID and key into k8s secret
        kubectl create secret generic --namespace=default $PARTICIPANTID-$SERVICENAME-aws-iam-key --from-literal=aws-access-key-id=$AWSACCESSKEYID --from-literal=aws-secret-access-key=$AWSSECRETACCESSKEY

        NEWFOLDER="../../iam-policy/policies"
        TEMPLATEFOLDER="../../iam-policy"
    done
    kubectl delete pods --all -n default
}

updateAPIGateway()
{
    if [ $TRAVIS_BRANCH == 'dev-gftn' ] && [ $TRAVIS_PULL_REQUEST == 'false' ]; then
        N_ENVIRONMENT="-dev"
        ENV="dev"
        export AWS_DEFAULT_REGION="us-west-2"
        export REGION='us-west-2' && export CLUSTER_NAME='terraform-eks-ww-dev' && aws eks update-kubeconfig --name $CLUSTER_NAME --region $REGION
    elif [ $TRAVIS_BRANCH == 'qa-gftn' ] && [ $TRAVIS_PULL_REQUEST == 'false' ]; then
        N_ENVIRONMENT="-qa"
        ENV="qa"
        export AWS_DEFAULT_REGION="us-west-2"
        export REGION='us-west-2' && export CLUSTER_NAME='terraform-eks-ww-qa' && aws eks update-kubeconfig --name $CLUSTER_NAME --region $REGION
    elif [ $TRAVIS_BRANCH == 'staging-gftn' ] && [ $TRAVIS_PULL_REQUEST == 'false' ]; then
        N_ENVIRONMENT="-st"
        ENV="st"
        export AWS_DEFAULT_REGION="us-east-2"
        export REGION='us-east-2' && export CLUSTER_NAME='terraform-eks-ww-st' && aws eks update-kubeconfig --name $CLUSTER_NAME --region $REGION
    elif [ $TRAVIS_BRANCH == 'test-gftn' ] && [ $TRAVIS_PULL_REQUEST == 'false' ]; then
        N_ENVIRONMENT="-pen"
        ENV="pen"
        export AWS_DEFAULT_REGION="us-east-2"
        export REGION='us-east-2' && export CLUSTER_NAME='terraform-eks-ww-pen' && aws eks update-kubeconfig --name $CLUSTER_NAME --region $REGION
    elif [ $TRAVIS_BRANCH == 'live-gftn' ] && [ $TRAVIS_PULL_REQUEST == 'false' ]; then
        N_ENVIRONMENT=""
        ENV="prod"
        export AWS_DEFAULT_REGION="us-east-1"
        export REGION='us-east-1' && export CLUSTER_NAME='terraform-eks-ww-prod' && aws eks update-kubeconfig --name $CLUSTER_NAME --region $REGION
    else
        # Update the k8s configMap
        echo " >>>>> INFO <<<<<: No conditions were matched for this deployment. This could be becuase this is not a branch that enables deployment or this is a pull-request(not a merge.)"
        exit 1
    fi

    BASE_PATH="../../automation-service/k8s"
    OLDPARTICIPANT=participant_id_variable
    OLDENVIRONMENT=environment_variable
    OLDPARTICIPANTID="empty_name"

    echo $AWS_DEFAULT_REGION

    DEPLOYMENTLIST=$(kubectl get deployment -o go-template --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')

    COUNT=0

    for D in $DEPLOYMENTLIST
    do
        IFS='-' # hyphen (-) is set as delimiter
        read -ra NAME <<< "$D"
        IFS=''

        PARTICIPANTID='EMPTY'
        for N in ${NAME[@]}
        do
          if [ $N == 'api' ]; then
              ((COUNT++))
              SERVICENAME="api-service"
              break
          elif  [ $N == 'crypto'  ]; then
              ((COUNT++))
              SERVICENAME="crypto-service"
              break
          elif  [ $N == 'gateway'  ]; then
              ((COUNT++))
              SERVICENAME="ww-gateway"
              break
          elif  [ $N == 'listener'  ]; then
              ((COUNT++))
              SERVICENAME="payment-service"
              break
          elif  [ $N == 'send'  ]; then
              ((COUNT++))
              SERVICENAME="send-service"
              break
          elif  [ $N == 'administration'  ]; then
              SERVICENAME="admin-service"
              break
          elif  [ $N == 'anchor'  ]; then
              SERVICENAME="anchor-service"
              break
          elif  [ $N == 'fee'  ]; then
              SERVICENAME="fee-service"
              break
          elif  [ $N == 'gas'  ]; then
              SERVICENAME="gas-service"
              break
          elif  [ $N == 'payout'  ]; then
              SERVICENAME="payout-service"
              break
          elif  [ $N == 'pr'  ]; then
              SERVICENAME="pr-service"
              break
          elif  [ $N == 'quotes'  ]; then
              SERVICENAME="quotes-service"
              break
          elif [ $N == 'whitelist' ]; then
              SERVICENAME="whitelist-service"
              break
          else
            if [ $PARTICIPANTID == 'EMPTY' ]; then
              PARTICIPANTID=${N}
            else
              PARTICIPANTID="${PARTICIPANTID}-$N"
            fi
          fi
        done

        if [ $PARTICIPANTID == $OLDPARTICIPANTID ]; then
            continue
        elif [ $PARTICIPANTID == 'ww' ]; then
            COUNT=0
            continue
        elif [[ $COUNT -eq 5 ]]; then
            printf "Ready to update AWS API-Gateway for MM %s on %s k8s cluster\n" "$PARTICIPANTID" "$TRAVIS_BRANCH"

            cd $BASE_PATH
            OLDPARTICIPANTID=$PARTICIPANTID
            NEW_API_GATEWAY_PATH="./api-gateway/$PARTICIPANTID"
            mkdir -p $NEW_API_GATEWAY_PATH

            # copy files into new files
            cp -r "./api-gateway/template/" $NEW_API_GATEWAY_PATH/
            cd $NEW_API_GATEWAY_PATH/template

            # replace
            find . -type f | xargs sed -i "s/$OLDPARTICIPANT/$PARTICIPANTID/g"
            find . -type f | xargs sed -i "s/$OLDENVIRONMENT/$N_ENVIRONMENT/g"

            API_GATEWAY_FILE_PATH="file://./aws-api-gateway.yaml"

            export RESTAPINAME="$PARTICIPANTID"
            API_ID=$(aws apigateway get-rest-apis | jq '.items[] | select(.name == env.RESTAPINAME) | .id' | sed -e 's/^"//' -e 's/"$//')
            if [ -z $API_ID ]; then
                echo "Empty API ID"
                exit 1
            fi

            # Update API Gateway settings
            aws apigateway put-rest-api --rest-api-id $API_ID --body $API_GATEWAY_FILE_PATH --mode overwrite

            # Get vpc link id
            export VPC_LINK_NAME="istio-nlb-link$N_ENVIRONMENT"
            VPC_LINK_ID=$(aws apigateway get-vpc-links | jq '.items[] | select(.name == env.VPC_LINK_NAME) | select(.status == "AVAILABLE") | .id' | sed -e 's/^"//' -e 's/"$//')
            if [ -z $VPC_LINK_ID ]; then
                echo "Empty VPC link ID"
                exit 1
            fi

            # Deploy API
            aws apigateway create-deployment --rest-api-id=$API_ID --stage-name=$ENV --variables environment=$ENV,global='global',vpcLinkId=$VPC_LINK_ID,participant=$PARTICIPANTID
            BASE_PATH="../../.."
            COUNT=0
        elif [[ $COUNT -eq 1 ]] && [[ $SERVICENAME == 'ww-gateway' ]] ; then
            printf "Ready to update AWS API-Gateway for anchor %s on %s k8s cluster\n" "$PARTICIPANTID" "$TRAVIS_BRANCH"

            cd $BASE_PATH
            OLDPARTICIPANTID=$PARTICIPANTID
            NEW_API_GATEWAY_PATH="./global-only-api-gateway/$PARTICIPANTID"
            mkdir -p $NEW_API_GATEWAY_PATH

            # copy files into new files
            cp -r "./global-only-api-gateway/template/" $NEW_API_GATEWAY_PATH/
            cd $NEW_API_GATEWAY_PATH/template

            # replace
            find . -type f | xargs sed -i "s/$OLDPARTICIPANT/$PARTICIPANTID/g"
            find . -type f | xargs sed -i "s/$OLDENVIRONMENT/$N_ENVIRONMENT/g"

            API_GATEWAY_FILE_PATH="file://./aws-api-gateway.yaml"

            export RESTAPINAME="$PARTICIPANTID"
            API_ID=$(aws apigateway get-rest-apis | jq '.items[] | select(.name == env.RESTAPINAME) | .id' | sed -e 's/^"//' -e 's/"$//')
            if [ -z $API_ID ]; then
                echo "Empty API ID"
                exit 1
            fi

            # Update API Gateway settings
            aws apigateway put-rest-api --rest-api-id $API_ID --body $API_GATEWAY_FILE_PATH --mode overwrite

            # Get vpc link id
            export VPC_LINK_NAME="istio-nlb-link$N_ENVIRONMENT"
            VPC_LINK_ID=$(aws apigateway get-vpc-links | jq '.items[] | select(.name == env.VPC_LINK_NAME) | select(.status == "AVAILABLE") | .id' | sed -e 's/^"//' -e 's/"$//')
            if [ -z $VPC_LINK_ID ]; then
                echo "Empty VPC link ID"
                exit 1
            fi

            # Deploy API
            aws apigateway create-deployment --rest-api-id=$API_ID --stage-name=$ENV --variables environment=$ENV,global='global',vpcLinkId=$VPC_LINK_ID,participant=$PARTICIPANTID
            BASE_PATH="../../.."
            COUNT=0
        else
            continue
        fi

    done

}

updateServiceVersion()
{
    # Get the build version from VERSION file
    VERSION=`cat ../../VERSION`
    echo $VERSION
    if [ $TRAVIS_BRANCH == 'dev-gftn' ] && [ $TRAVIS_PULL_REQUEST == 'false' ]; then
        # Update the k8s configMap
        bash ../../automation-service/k8s/script/create-configs.sh dev
        export REGION='us-west-2' && export CLUSTER_NAME='terraform-eks-ww-dev' && aws eks update-kubeconfig --name $CLUSTER_NAME --region $REGION
    elif [ $TRAVIS_BRANCH == 'qa-gftn' ] && [ $TRAVIS_PULL_REQUEST == 'false' ]; then
        # Update the k8s configMap
        bash ../../automation-service/k8s/script/create-configs.sh qa
        export REGION='us-west-2' && export CLUSTER_NAME='terraform-eks-ww-qa' && aws eks update-kubeconfig --name $CLUSTER_NAME --region $REGION
    elif [ $TRAVIS_BRANCH == 'staging-gftn' ] && [ $TRAVIS_PULL_REQUEST == 'false' ]; then
        # Update the k8s configMap
        bash ../../automation-service/k8s/script/create-configs.sh st
        export REGION='us-east-2' && export CLUSTER_NAME='terraform-eks-ww-st' && aws eks update-kubeconfig --name $CLUSTER_NAME --region $REGION
    elif [ $TRAVIS_BRANCH == 'test-gftn' ] && [ $TRAVIS_PULL_REQUEST == 'false' ]; then
        # Update the k8s configMap
        bash ../../automation-service/k8s/script/create-configs.sh pen
        export REGION='us-east-2' && export CLUSTER_NAME='terraform-eks-ww-pen' && aws eks update-kubeconfig --name $CLUSTER_NAME --region $REGION
    elif [ $TRAVIS_BRANCH == 'live-gftn' ] && [ $TRAVIS_PULL_REQUEST == 'false' ]; then
        # Update the k8s configMap
        bash ../../automation-service/k8s/script/create-configs.sh prod
        export REGION='us-east-1' && export CLUSTER_NAME='terraform-eks-ww-prod' && aws eks update-kubeconfig --name $CLUSTER_NAME --region $REGION
    else
        # Update the k8s configMap
        echo " >>>>> INFO <<<<<: No conditions were matched for this deployment. This could be becuase this is not a branch that enables deployment or this is a pull-request(not a merge.)"
        exit 1
    fi

    # Change the environment variable `DOCKER_IMAGE_VERSION` to $VERSION for the deployment-service
    kubectl set env -n deployment deployment/deployment-service DOCKER_IMAGE_VERSION=$VERSION

    DEPLOYMENTLIST=$(kubectl get deployment -o go-template --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')

    for D in $DEPLOYMENTLIST
    do
        printf "Ready to update %s on %s k8s cluster to version %s\n" "$D" "$TRAVIS_BRANCH" "$VERSION"

        IFS='-' # hyphen (-) is set as delimiter
        read -ra NAME <<< "$D"
        IFS=''

        for N in ${NAME[@]}
        do
          if [ $N == 'api' ]; then
              SERVICENAME="api-service"
              break
          elif  [ $N == 'crypto'  ]; then
              if [ $TRAVIS_BRANCH == 'live-gftn' ]; then
                SERVICENAME="crypto-service-prod"
              else
                SERVICENAME="crypto-service"
              fi
              break
          elif  [ $N == 'gateway'  ]; then
              SERVICENAME="ww-gateway"
              break
          elif  [ $N == 'listener'  ]; then
              SERVICENAME="payment-listener"
              break
          elif  [ $N == 'send'  ]; then
              SERVICENAME="send-service"
              break
          elif  [ $N == 'administration'  ]; then
              SERVICENAME="administration-service"
              break
          elif  [ $N == 'anchor'  ]; then
              SERVICENAME="anchor-service"
              break
          elif  [ $N == 'fee'  ]; then
              SERVICENAME="fee-service"
              break
          elif  [ $N == 'gas'  ]; then
              SERVICENAME="gas-service"
              break
          elif  [ $N == 'payout'  ]; then
              SERVICENAME="payout-service"
              break
          elif  [ $N == 'pr'  ]; then
              SERVICENAME="participant-registry"
              break
          elif  [ $N == 'quotes'  ]; then
              SERVICENAME="quotes-service"
              break
          elif [ $N == 'whitelist' ]; then
              SERVICENAME="global-whitelist-service"
              break
          else
            continue
          fi
        done

        kubectl --record deployment.apps/$D set image deployment.v1.apps/$D $D=${DOCKER_REGISTRY}/gftn/$SERVICENAME:$VERSION -n default

    done

    kubectl --record deployment.apps/deployment-service set image deployment.v1.apps/deployment-service deployment-service=${DOCKER_REGISTRY}/gftn/automation-service:$VERSION -n deployment
    kubectl delete pods --all -n default

}

CMD=$1

echo "branch = $TRAVIS_BRANCH"
echo "Is pull request = $TRAVIS_PULL_REQUEST"

if [ $TRAVIS_PULL_REQUEST == 'false' ]; then
    if [ $CMD == "updateAWSSecret" ]; then
      updateAWSSecret
    elif [ $CMD == "updateAPIGateway" ]; then
      updateAPIGateway
    elif [ $CMD == "updateIAMPolicy" ]; then
      updateIAMPolicy
    elif [ $CMD == "updateServiceVersion" ]; then
      updateServiceVersion
    fi
else
    echo "Pull request number is = $TRAVIS_PULL_REQUEST"
fi