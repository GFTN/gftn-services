# MSK topic configurations
Some scripts for setting up the Kafka topic and ACLs

###Create topics and grant ACLs for participant

Before create topics, please make sure this is a newly onboarding participant. Some topic can not recreate in the
some MSK cluster, unless the MSK cluster is a brand new one.

#### For newly onboarding participant
Normally the Kafka topic creation will be handled during the onboarding process. 
If the process failed to create the topics, use the below script.
````
bash create_cert_and_topic.sh $PARTICIPANT_ID $KEY_PASSWORD $ENVIRONMENT . ip-team-worldwire-docker-local.artifactory.swg-devops.com
````
Arguments
1. PARTICIPANT_ID : participant's id
2. KEY_PASSWORD : password for the private key
3. ENVIRONMENT : cluster environment [eksdev, eksqa, st, prod]
4. DIRECTORY : default directory is current folder
5. DOCKERREGISTRYURL : URL for docker registry

#### For existing participant
After the MSK cluster is reconfigure, we have to create topics and grant ACLs for all the existing participants again. 
Please replace the `{participant_id}` in the below scripts with the actual id of the participant and `{environment}`
with the cluster environment that you are working on.

Delete the existing k8s secret which store the outdated Kafka certificate and private key. 
````
kubectl delete secret kafka-secret-{participant_id}
````
Delete the existing k8s jobs in namespace `kafka-topics`.
````
kubectl delete job -n kafka-topics {participant_id}-create-topic-participant
````
Delete the previous certificate in the AWS certificate manager
````
bash delete_acm.sh {participant_id} {envrionment}
````
Update environment variable for deployment-service
````
bash update_env.sh {environment}
````
Run the script to create topics
````
bash create_cert_and_topic.sh $PARTICIPANT_ID $KEY_PASSWORD $ENVIRONMENT . ip-team-worldwire-docker-local.artifactory.swg-devops.com
````