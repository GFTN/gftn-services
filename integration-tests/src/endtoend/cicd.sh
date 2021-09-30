
E2E_PATH=$GOPATH/src/github.com/GFTN/gftn-services/integration-tests/src/endtoend

# # FOR GENERATE THE DOCKER-COMTAINER, IN CASE IF WE NEED TO CHANGE THE ENV,PARTICIPANT ID, URL REPLACE RULE
# #########################################################################
# echo -e '\x1B[32m\xE2\x9C\x94 Initialised MicroService docker-compose.yml \x1B[0m'

# cd $E2E_PATH
# # Create directory
# if [ ! -d "cicd/file/worldwireServices" ]; then
#     echo -e '\x1B[33m  Creating directory cicd/file/worldwireServices ... \x1B[0m'
#     mkdir -p $E2E_PATH/cicd/file/worldwireServices
# fi

# cd $E2E_PATH
# # if file exist remove file
# if [ -f "cicd/file/worldwireServices/configMap.env" ]; then
#     rm $E2E_PATH/cicd/file/worldwireServices/configMap.env
#     echo -e '\x1B[31m\xE2\x9C\x94 Clean/Remove cicd/file/worldwireServices/configMap.env \x1B[0m'
# fi

# cd $E2E_PATH
# # if file exist remove file
# if [ -f "cicd/file/worldwireServices/docker-compose.yml" ]; then
#     rm cicd/file/worldwireServices/docker-compose.yml
#     echo -e '\x1B[31m\xE2\x9C\x94 Clean/Remove cicd/file/worldwireServices/docker-compose.yml \x1B[0m'
# fi

# # Create file
# touch $E2E_PATH/cicd/file/worldwireServices/docker-compose.yml

# # Generate docker-compose file
# echo -e '\x1B[33m  Generate MicroService docker-compose.yml \x1B[0m'
# export CICD_PATH=$E2E_PATH/cicd
# node $E2E_PATH/cicd/createDockerCompose.js
# echo -e '\x1B[33m  docker-compose.yml Generation complete\x1B[0m'

# #########################################################################

version=`cat VERSION`


echo -e '\x1B[32m\xE2\x9C\x94 Login Docker \x1B[0m'
docker login -u $DOCKER_USER -p $DOCKER_PASSWORD $DOCKER_REGISTRY

echo -e '\x1B[32m\xE2\x9C\x94 pull Docker images \x1B[0m'
docker pull $DOCKER_REGISTRY/gftn/crypto-service:$version
docker pull $DOCKER_REGISTRY/gftn/quotes-service:$version
docker pull $DOCKER_REGISTRY/gftn/api-service:$version
docker pull $DOCKER_REGISTRY/gftn/gas-service:$version
docker pull $DOCKER_REGISTRY/gftn/send-service:$version
docker pull $DOCKER_REGISTRY/gftn/payout-service:$version
docker pull $DOCKER_REGISTRY/gftn/participant-registry:$version
docker pull $DOCKER_REGISTRY/gftn/fee-service:$version
docker pull $DOCKER_REGISTRY/gftn/anchor-service:$version
docker pull $DOCKER_REGISTRY/gftn/payment-listener:$version
docker pull $DOCKER_REGISTRY/gftn/administration-service:$version
docker pull $DOCKER_REGISTRY/gftn/global-whitelist-service:$version
docker pull $DOCKER_REGISTRY/gftn/ww-gateway:$version


# docker tag $DOCKER_REGISTRY/gftn/crypto-service:$version gftn/crypto-service:latest
# docker tag $DOCKER_REGISTRY/gftn/quotes-service:$version gftn/quotes-service:latest
# docker tag $DOCKER_REGISTRY/gftn/api-service:$version gftn/api-service:latest
# docker tag $DOCKER_REGISTRY/gftn/gas-service:$version gftn/gas-service:latest
# docker tag $DOCKER_REGISTRY/gftn/send-service:$version gftn/send-service:latest
# docker tag $DOCKER_REGISTRY/gftn/payout-service:$version gftn/payout-service:latest
# docker tag $DOCKER_REGISTRY/gftn/participant-registry:$version gftn/participant-registry:latest
# docker tag $DOCKER_REGISTRY/gftn/fee-service:$version gftn/fee-service:latest
# docker tag $DOCKER_REGISTRY/gftn/anchor-service:$version gftn/anchor-service:latest
# docker tag $DOCKER_REGISTRY/gftn/payment-listener:$version gftn/payment-listener:latest
# docker tag $DOCKER_REGISTRY/gftn/administration-service:$version gftn/administration-service:latest
# docker tag $DOCKER_REGISTRY/gftn/global-whitelist-service:$version gftn/global-whitelist-service:latest
# docker tag $DOCKER_REGISTRY/gftn/ww-gateway:$version gftn/ww-gateway:latest



echo -e '\x1B[32m\xE2\x9C\x94 Clean up the exist docker container \x1B[0m'
# if container exist remove container
if [ "$(docker ps -q -f name=ww-pr)" ]; then
    docker stop $(docker ps -qa)
    docker rm $(docker ps -qa)
    echo -e '\x1B[31m\xE2\x9C\x94 Clean/Remove Docker environment \x1B[0m'
fi

# if kafka folder exist remove container
if [ ! -d "cicd/file/worldwireServices" ]; then
    echo -e '\x1B[33m  Creating directory cicd/file/worldwireServices ... \x1B[0m'
    mkdir -p $E2E_PATH/cicd/file/worldwireServices
fi

# Create docker network
if [ ! "$(docker network ls | grep wwcicdnet)" ]; then
    echo -e '\x1B[32m\xE2\x9C\x94 Creating wwcicdnet network ... \x1B[0m'
    docker network create --driver bridge --subnet=172.19.0.0/24 --gateway 172.19.0.1 \
    --opt com.docker.network.bridge.enable_icc=true \
    --opt com.docker.network.bridge.enable_ip_masquerade=true \
    --opt com.docker.network.driver.mtu=1350 wwcicdnet
    
else
    echo -e '\x1B[32m\xE2\x9C\x94 wwcicdnet network exists. \x1B[0m'
fi



echo -e '\x1B[32m\xE2\x9C\x94 Clean up the exist docker container \x1B[0m'
# if container exist remove container
if [ "$(docker ps -q -f name=ww-pr)" ]; then
    docker stop $(docker ps -qa)
    docker rm $(docker ps -qa)
    echo -e '\x1B[31m\xE2\x9C\x94 Clean/Remove Docker environment \x1B[0m'
fi


# Running WW service in local
echo -e '\x1B[32m\xE2\x9C\x94 Running Zoomkeeper, Kafka in docker environment \x1B[0m'
rm -rf $E2E_PATH/cicd/file/kafka-cluster/kafka
rm -rf $E2E_PATH/cicd/file/kafka-cluster/zk
# mkdir -p $E2E_PATH/cicd/file/kafka-cluster/kafka
cd $E2E_PATH/cicd/file/kafka-cluster/ && docker-compose up -d

sleep 90

echo -e '\x1B[32m\xE2\x9C\x94 Running MicroService in docker environment \x1B[0m'
cd $E2E_PATH/cicd/file/worldwireServices/  && docker-compose up -d --no-recreate

sleep 90

# Setting test target environment
echo -e '\x1B[33m  Setting E2E configration {$CICD_PATH} \x1B[0m'
cd $E2E_PATH
source ./environment/.travis.env

# Initialize testing environment, clean up the package
if [ -d "node_modules" ]; then
    rm -rf $E2E_PATH/node_modules/
    echo -e '\x1B[31m\xE2\x9C\x94 Clean/Remove node_modules \x1B[0m'
fi

# Install dependencies
echo -e '\x1B[32m\xE2\x9C\x94 Install E2E,CICD package\x1B[0m'
npm install

# List all the scenatio are going to test
# echo -e '\x1B[32m\xE2\x9C\x94 Running E2E test \x1B[0m'
# echo -e '\x1B[33m Scenario: \x1B[0m'
# echo -e '\x1B[94m   failCase \x1B[0m'
# echo -e '\x1B[94m   1_issue_assets \x1B[0m'
# echo -e '\x1B[94m   2_manage_whitelist \x1B[0m'
# echo -e '\x1B[94m   3_trust_asset/1_request_trust_asset_DO.feature  \x1B[0m'
# echo -e '\x1B[94m   3_trust_asset/2_request_trust_asset_DA.feature  \x1B[0m'
# echo -e '\x1B[94m   6_funding/1_anchor_funding_participant.feature  \x1B[0m'
# echo -e '\x1B[94m   4_payment/ \x1B[0m'
# echo -e '\x1B[94m   5_exchange/ \x1B[0m'

# Start E2E funtional test
npm run test features/goodCase
# features/goodCase/2_manage_whitelist/ \
# features/goodCase/3_trust_asset/1_request_trust_asset_DO.feature \
# features/goodCase/3_trust_asset/2_request_trust_asset_DA.feature \
# features/goodCase/6_funding/1_anchor_funding_participant.feature 
# features/goodCase/4_payment/ \
# features/goodCase/5_exchange/