# gas-service
IBM controlled lumen service to fund transactions in WW

## Script

### Delete table
* 1. IBM_TOKEN_ACCOUNTS: accounts info.

* 2. IBM_CONTACTS: accounts contact info.

* 3. IBM_GROUPS: groups topic ingfo.

Data Table name can be modefied by environment variables .


Environment variables : 
* 1. DYNAMODB_ACCOUNTS_TABLE_NAME
* 2. DYNAMODB_CONTACTS_TABLE_NAME
* 3. DYNAMODB_GROUPS_TABLE_NAME
    
          npm run deleteTB $tablename

### Create table
    npm run createAccountsTB
    npm run createContactsTB
    npm run createGroupsTB


### Make Docekr image 
    npm run makedocker


### Run Docekr container (Running in port:18080 ,You can change to whatever port)
    docker run -p 18080:8080 -it gftn/gas-service

### Start Server (port :8080) 
    npm run start


## API

### GET /lockaccount  
* payload (empty)

* Succesful (HTTP status :200 )


        {
        "pkey": "publick address",
        "sequenceNumber": "number"
        }

* Fail (HTTP status :500 )

        {
        "failure_reason": "no avaible account now"
        }



### POST /signXDRAndExecuteXDR
* payload (Object )

        {
        "oneSignedXDR": "XDR('base64')"
        } 

* Account expire response: (HTTP status:400)


        {   
            "title": "Source Account Expire",
            "failure_reason": "something"
        }



* Tx success response: (HTTP status:200)

        {
            "title": "Transaction successful",
            "hash": "something",
            "ledger": something
        }


* Tx fail response: (HTTP status:403)


        {
            "title": "Transaction Failed",
            "failure_reason": {
                "envelope_xdr": "something",
                "result_codes": {
                    "transaction": "something"
                },
                "result_xdr": "something"
            }
        }   


### POST /createAccounts
* payload (Object Array)



        [
            {
            "key": {

                "Object": "Vault tag"
            },
            "seed": {
                "Object": "Vault tag"
            },
            "accountStatus": true/false,
            "groupName": "groupName"

            }
        ]


### POST /createGroups
* payload (Object Array)

        [

            {
                "TopicName":"name",
                "DisplayName":"name"
            }
        ]


### POST /createContacts
* payload (Object Array)

        [

            {
                "groupName": "Group1",
                "email":"your.user@your.domain",
                "phoneNumber":"+9999999999"
            }
        ]



## Lock accounts
When client send a request to get an account .

 * 0.check whether $unlockAccounts still have account can use.   If not ,response no avaible account now
 * 1.pop from unlockAccounts ($account)
 * 2.get timestamp
 * 3.update to dynamoDB ,
 * [ update account from table where pkey = account.pkey and status = unlock (status,timestamp) ]
 * 3.add timestamp to $account 
 * 4.push $account to lockArray (lockAccounts.push($account))
 * 5.get sequence number $account.seqNum
 * 5.return $account.seqNum,$account.pkey , $lockAccounts , $unlockAccounts

 
## Sign and Execute 

### Logic 

#### sign
 * 1. decode signedXDRin , get pkey 
 * 2. use pkey to get secret
 * 3. signed by the secret 
 * 4. return txeB64 , using $account


#### execute 
 * 1. decode signedXDRin , get source account 
 * 2. check whether account is in $lockAccounts , if account is not in $lockAccounts , return can not execute 
 * 3. unlock account
 * 4. execute signedXDR
 * 5. return $result ($result=stallar response detail)


## DynamoDB

### Tables : 


* DYNAMODB_ACCOUNTS_TABLE_NAME(using environment variables) : saving accounts Data

    SCHEMA

            {
            "pkey": "string,
            "accountStatus": true/false,
            "secret": "string",
            "groupName": "string"

            }



* DYNAMODB_CONTACTS_TABLE_NAME(using environment variables) : saving contacts Data

    SCHEMA

            {
            "groupName": "string",
            "email": "string",
            "phoneNumber": "string"
            
            }


* DYNAMODB_GROUPS_TABLE_NAME(using environment variables) : saving Groups Data


    SCHEMA

            {
            "TopicName": "string",
            "TopicArn": "string",
            "displayName": "string"

            }
            
## Monitor

There are 3 monitor in gas service 
* monirtor for lock accounts
* monitor for high threshold accounts
* monitor for low threshold accounts


## 1. monitor Lock accounts
$account.lockTimestamp = the account been lock time

$expireTime = the time range client can use this account

Every monitoringTime , monitor will check the lock array , whether there has account been lock .

If yes, then run foreach loop ($lockAccounts).

    If (account.lockTimestamp + expireTime < now) unlock automatically

## 2. monitor high threshold accounts

* 0. read all IBM accounts from dynamoDB
* 1. get the balance from stellar
* 2. check each accounts whether balance is ($balance<MONITOR_HIGH_THRESHOLD_BALANCE)
* 3. if ($balance<MONITOR_HIGH_THRESHOLD_BALANCE) and ($lowThresholdAccounts.length<1) { create thread (call monitorBalanceL)}
* 4. else ($balance<MONITOR_HIGH_THRESHOLD_BALANCE) { push accounts to lowThresholdAccounts }
     

## 3. monitor low threshold accounts

* 0. if there has no accounts in the low threshold list stop
* 1. read account balance from stellar
* 2. check whetheer the account balancs is lower the NOTIFY_BALANCE
* 3. if $balance>MONITOR_LOW_THRESHOLD_BALANCE move to another list (high threshold list)
* 4. if lower than the NOTIFY_BALANCE , send SMS,EMAIL,SLACK
 



## set up env variables

      export GAS_SERVICE_URL=http://gas-service-io
      export SERVICE_NAME=gas-service
      export GAS_SERVICE_VERSION=v1
      export GAS_SERVICE_PORT=8099
      export GAS_SERVICE_MONITOR_LOCKACCOUNT_FEQ=5000
      export GAS_SERVICE_EXPIRE_TIME=30
      export GAS_SERVICE_EMAIL_NOTIFICATION=false
      export GAS_SERVICE_SMS_NOTIFICATION=false
      export HIGH_THRESHOLD_BALANCE=9996
      export HIGH_THRESHOLD_TIMEOUT=10000
      export LOW_THRESHOLD_BALANCE=9996
      export LOW_THRESHOLD_TIMEOUT=5000
      export DYNAMODB_ENDPOINT_URL=http://dynamodb.us-east-1.amazonaws.com
      export DYNAMODB_ACCESSKEYID=
      export DYNAMODB_SECRECT_ACCESS_KEY=
      export DYNAMODB_REGION=us-east-1
      export DYNAMODB_ACCOUNTS_TABLE_NAME=IBM_TOKEN_ACCOUNTS
      export DYNAMODB_CONTACTS_TABLE_NAME=IBM_CONTACTS
      export DYNAMODB_GROUPS_TABLE_NAME=IBM_GROUPS
      export VAULT_APPID=SSLcert
      export VAULT_SAFE=IBM
      export VAULT_FOLDER=Root
      export VAULT_URL=https://3.0.15.221
      export VAULT_CERT_PATH=./config/certificate.crt
      export VAULT_KEY_PATH=./config/privateKey.key
      export SENDGRID_API_KEY=
      export SENDGRIID_ENDPINT=https://api.sendgrid.com/v3/mail/send
      export EMAIL_SENDER=
      export SERVICE_LOG_FILE=$GOPATH/src/github.com/GFTN/gftn-services/gas-service
