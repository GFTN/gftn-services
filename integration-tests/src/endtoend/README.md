# wwAutomationTest

    ├── failCase
    │   ├── 1_RegisterParticipant
    │   │   ├── 0_failCase_RegisterParticipant_WrongPayLoad.feature
    │   │   └── 1_failCase_with_wrong_payload.feature
    │   ├── 2_create_accounts
    │   │   ├── 0_failCase_create_issuing_account_duplicated.feature
    │   │   └── 1_failCase_create_operating_account_duplicated.feature
    │   ├── 3_issue_assets
    │   │   └── 1_failCase_issue_asset.feature
    │   ├── 4_whitelist_trustline
    │   │   └── 0_failCase_trust_asset_without_whitelist.feature
    │   ├── 5_payment
    │   │   ├── 0_failCase_payment_without_OFI_whitelist_RFI.feature
    │   │   ├── 1_failCase_payment_exceed_balance.feature
    │   │   ├── 2_failCase_payment_without_RFI_whitelist_OFI.feature
    │   │   ├── 3_RFI_receive_endToendID_not_exitst_fail.feature
    │   │   ├── 4_failCase_OFI_send_account_not_exitst.feature
    │   │   └── 5_failCase_OFI_receive_account_not_exitst.feature
    │   └── 6_exchange
    │       ├── 0_failCase_exchange_without_whitelist.feature
    │       └── 1_failCase_exchange_with_expire_quotes.feature
    └── goodCase
        ├── 0_onboarding
        │   ├── 1_onboard_participant
        │   │   ├── 1_create_participants.feature
        │   │   ├── 2_create_issuing_accounts.feature
        │   │   └── 3_create_operating_accounts.feature
        │   ├── 2_onboard_anchor
        │   │   ├── 1_create_anchor.feature
        │   │   └── 2_registe_anchor.feature
        │   └── 3_onboard_sweep_accounts
        │       └── 0_create_sweep_accounts.feature
        ├── 1_issue_assets
        │   ├── 1_participant_DO
        │   │   └── 1_issue_asset.feature
        │   └── 2_anchor_DA
        │       └── 1_anchor_onboard_asset.feature
        ├── 2_manage_whitelist
        │   └── 1_add_whitelist.feature
        ├── 3_trust_asset
        │   ├── 1_request_trust_asset_DO.feature
        │   ├── 2_request_trust_asset_DA.feature
        │   └── 3_sweep_accounts_request_trust_asset_DA.feature
        ├── 4_payment
        │   ├── DA
        │   │   └── 1_payment_comple.feature
        │   └── DO
        │       ├── 1_payment_comple_settle.feature
        │       ├── 2_payment_cancelation_agree.feature
        │       └── 3_payment_cancelation_reject.feature
        ├── 5_exchange
        │   └── 1_exchange_asset_between_participants.feature
        ├── 6_funding
        │   ├── 1_anchor_funding_participant.feature
        │   └── 2_anchor_funding_sweep_accounts.feature
        ├── 7_healthy_check
        │   ├── 0_healthy_check_pr.feature
        │   ├── 1_healthy_check_api.feature
        │   ├── 2_healthy_check_whitelist.feature
        │   ├── 3_healthy_check_send_fee.feature
        │   └── 4_healthy_check_quote_crypto.feature
        └── 8_sweep
            ├── 1_sweep_single_asset.feature
            └── 2_sweep_multi_accounts
                ├── 0_apply_multi_accounts_balance.feature
                ├── 1_sweep_multi_accounts.feature
                └── 3_check_after_sweep_balance.feature
    
    




# Install

## Node

```brew install node ```

## Install package

`$ npm install `

# Running test

## Running test from local

0. Setting testing environment in local

Make docker images in local or pull from Jfrog

- make images from local
    ```
    cd worldwire-service && make dep
    cd gftn-models/ && make
    cd $SERVICES && make docker
    ```    
- pull from Jfrog and tag as latest or replace docker images version to which version you want to test in cicd/file/worldwireServices/docker-compose.yaml
    ```
    docker login $DOCKER_REGISTRY 
    LOGIN BY YOUR IBM ACCOUNT 

    docker pull $DOCKER_REGISTRY/gftn/$SERVICES:$version
    docker tag $DOCKER_REGISTRY/gftn/$SERVICES:$version gftn/$SERVICES:latest
    ```    

- Create network in local, order to make connection between kafka and services

    `docker network create --driver bridge  wwcicdnet`

- Running Kafka in local

    `cd cicd/file/kafka-cluster && docker-compose up -d`

- Running Services in local

    `cd cicd/file/worldwireServices && docker-compose up -d`

- Monitoring logs

    `docker logs -f $container_name`  : Print and follow single service log.
    ex:   `docker logs -f travis1-api`

    `docker-compose logs -f` : Print and follow the services spin up from docker-compose.

1. Setting test case config 

    `source environment/.travis.env`
2. Running services in docker environment

    `cd cicd/file/kafka-cluster/ && docker-compose up -d`
    
    `cd ../worldwireServices/ && docker-compose up -d` 

3. Running test

    `npm run test $TEST_FEATURE`

You can change the scenarion to which you want to test 

```
ex:

         npm run test features/failCase/ \
         npm run test features/goodCase/1_issue_assets/
         npm run test features/goodCase/2_manage_whitelist/
         npm run test features/goodCase/3_trust_asset/1_request_trust_asset_DO.feature
         npm run test features/goodCase/3_trust_asset/2_request_trust_asset_DA.feature
         npm run test features/goodCase/6_funding/1_anchor_funding_participant.feature
         npm run test features/goodCase/4_payment/DO/
         npm run test features/goodCase/4_payment/DA/
         npm run test features/goodCase/5_exchange/
         npm run test features/goodCase/6_funding/2_anchor_funding_sweep_accounts.feature 
         npm run test features/goodCase/8_sweep/
```


## Automation test by BASH FILE

In order to run automated test need to 

1. Setting RDO Client Database connection and docker refistry information.
    ```
    export RDO_DB_HOST_2= 
    export RDO_DB_PASSWORD_2= 
    export RDO_DB_PASSWORD_1= 
    export RDO_DB_HOST_1= 
    export RDO_DB_PORT_1= 
    export RDO_DB_PORT_2= 
    export RDO_DB_USER_1= 
    export RDO_DB_USER_2= 
    export RDO_DB_NAME_1= 
    export RDO_DB_NAME_2= 
    export DOCKER_REGISTRY=
    ```

2. Modefy cicd.sh file

- Change file path

    `E2E_PATH=~/build/gftn/gftn-services/integration-tests/src/endtoend` 

    change to

    `E2E_PATH=$GOPATH/src/github.com/GFTN/gftn-services/integration-tests/src/endtoend`

- Login to docker 

    `docker login -u $DOCKER_USER -p $DOCKER_PASSWORD $DOCKER_REGISTRY`
    
    
    change to

    `docker login $DOCKER_REGISTRY`
    

3. Execute bash file

- ` cd gftn-services/`

- `bash integration-tests/src/endtoend/cicd.sh `

- Login docker by IBM account (If ask)

# What cicd.sh doing

## Initialize testing environment, clean up the package

1. Clean/Remove node_modules

    ```rm -rf node_modules/```

2. Creating wwcicdnet network

    ```docker network create wwcicdnet```

3. Clean/Remove cicd/file/worldwireServices/

    ```rm cicd/file/worldwireServices/docker-compose.yml```

4. Clean/Remove cicd/file/worldwireServices/configMap.env

    ```rm cicd/file/worldwireServices/configMap.env```

5. Clean/Remove Docker environment

    ``` docker stop $(docker ps -qa) && docker rm $(docker ps -qa)```

## Install dependencies

```npm install```


## Create docker-compose 

1. Create docker-compose file

    ```touch cicd/file/worldwireServices/docker-compose.yml```

2.  Generate docker-compose file

    ```node cicd/createDockerCompose.js```

3. Running Kafka service in local

    ```cd cicd/file/kafka-cluster/ && docker-compose up -d```

3. Running WW service in local

    ```cd ../worldwireServices/ && docker-compose up -d```

4. Check if all the servicea ready.

    ```docker ps -a```
    ```docker logs $CONTAINER_NAME```

5. Setting test target environment

    ```source ../../../environment/.travis.env```

## Start E2E funtional test

### Teting Negative cases

1. Register Participant Flow
    ```
    npm run test features/failCase/1_RegisterParticipant/0_failCase_RegisterParticipant_WrongPayLoad.feature 
    npm run test features/failCase/1_RegisterParticipant/1_failCase_with_wrong_payload.feature
    ```

2. Create Accounts Flow
    ```
    npm run test features/failCase/2_create_accounts/0_failCase_create_issuing_account_duplicated.feature 
    npm run test features/failCase/2_create_accounts/1_failCase_create_operating_account_duplicated.feature 
    ```
    
3. Issue Asset Flow
    ```
    npm run test features/failCase/3_issue_assets/1_failCase_issue_asset.feature
    ```
    
4. Whitelist And Trustline Flow
    ```
    npm run test features/failCase/4_whitelist_trustline/0_failCase_trust_asset_without_whitelist.feature 
    ```
    
5. Payment Flow
    ```
    npm run test features/failCase/5_payment/0_failCase_payment_without_OFI_whitelist_RFI.feature 
    npm run test features/failCase/5_payment/1_failCase_payment_exceed_balance/
    npm run test features/failCase/5_payment/2_failCase_payment_without_RFI_whitelist_OFI.feature 
    npm run test features/failCase/5_payment/3_RFI_receive_endToendID_not_exitst_fail.feature 
    npm run test features/failCase/5_payment/4_failCase_OFI_send_account_not_exitst/
    npm run test features/failCase/5_payment/5_failCase_OFI_receive_account_not_exitst/
    ```
    
6. Exchange Flow
    ```
    npm run test features/failCase/6_exchange/0_failCase_exchange_without_whitelist.feature
    npm run test features/failCase/6_exchange/1_failCase_exchange_with_expire_quotes.feature
    ```
    
### Testing Positive Cases

1. Onboarding Flow
    ```
    npm run test features/goodCase/0_onboarding/1_onboard_participant/1_create_participants.feature
    npm run test features/goodCase/0_onboarding/1_onboard_participant/2_create_issuing_accounts.feature
    npm run test features/goodCase/0_onboarding/1_onboard_participant/3_create_operating_accounts.feature
    npm run test features/goodCase/0_onboarding/2_onboard_anchor/1_create_anchor.feature
    npm run test features/goodCase/0_onboarding/2_onboard_anchor/2_registe_anchor.feature
    npm run test features/goodCase/0_onboarding/3_onboard_sweep_accounts/0_create_sweep_accounts.feature
    ```
    
2. Issue Asset Flow
    ```
    npm run test features/goodCase/1_issue_assets/1_participant_DO/1_issue_asset.feature
    npm run test features/goodCase/1_issue_assets/2_anchor_DA/1_anchor_onboard_asset.feature
    ```
    
3. Manage Whitelist Flow
    ```
    npm run test features/goodCase/2_manage_whitelist/1_add_whitelist.feature
    ```
    
4. Trustline Flow
    ```
    npm run test features/goodCase/3_trust_asset/1_request_trust_asset_DO.feature
    npm run test features/goodCase/3_trust_asset/2_request_trust_asset_DA.feature
    npm run test features/goodCase/3_trust_asset/3_sweep_accounts_request_trust_asset_DA.feature
    ```
    
5. Anchor Flow
    ```
    npm run test features/goodCase/6_funding/1_anchor_funding_participant.feature
    npm run test features/goodCase/6_funding/2_issue_asset.js
    npm run test features/goodCase/6_funding/3_funding_participant.js
    ```
    
6. Payment Flow
    ```
    npm run test features/goodCase/4_payment/DA/
    npm run test features/goodCase/4_payment/DO/1_payment_comple_settle/
    npm run test features/goodCase/4_payment/DO/2_payment_cancelation_agree/
    npm run test features/goodCase/4_payment/DO/3_payment_cancelation_reject
    ```
    
7. Exchange Flow
    ```
    npm run test features/goodCase/5_exchange/
    ```

