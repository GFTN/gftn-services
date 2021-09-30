# Automate deployment service

Service for deploy participant services to Kubernetes cluster.

```
https://deployment.worldwire{-env}.io/global/automation/v1
```
    
## APIs for WorldWire admin
User authentication needs to include two tokens in the request header
1. X-Fid (FireBase ID)
2. X-Permission

Deploy participant services
> POST: **/deploy/participant**

````
{
    "institutionId": "institution ID",
    "participantId": "participant ID",
    "env": "dev",
    "version": "latest",
    "callbackUrl": "callback URL",
    "status": "deploying status",
    "countryCode": "country code",
    "replica": "1",
    "bic": "BIC code",
    "role": "MM or IS"
}
````

model:[automation.yaml](model/models/automation.yaml)

Update Docker image version and other cluster configurations
> POST: **/update/image**

````
{
    "env": "st",
    "version": "2.9.3.11_RC1",
    "participants": [
        {
            "institutionID": "institution ID",
            "id": "participant ID",
            "bic": "BIC code",
            "callbackUrl": "callback URL",
            "countryCode": "country code",
            "role": "MM or IS"
        },
        ...
    ],
    "awsSecret": "encoded aws secret file",
    "configMap": "encoded k8s config map",
    "k8sSecret": "encoded k8s secret"
}
````

model:[update.yaml](model/models/update.yaml)

model:[participant.yaml](model/models/participant.yaml)