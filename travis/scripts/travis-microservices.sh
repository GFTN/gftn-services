#!/bin/bash

gererateMiddlewarePermissions()
{
  # # LONG FORM GENERATE: 
  # # install deps and generate the permisssions.go used by micro-service middleware
  # cd auth-service
  # sh -x install_dependencies.sh
  # bash ./deployment/build.sh --env dev --cloud gcloud

  # SHORT FORM GENERATE: 
  cd auth-service
  npm install --prefix deployment
  npm i -g typescript@3.4.5
  tsc deployment/src/permissions-only.ts
  node deployment/src/permissions-only.js
  ls authorization/middleware/permissions
  cd ..
}

beforeInstall()
{
  gererateMiddlewarePermissions
  chmod +x checkmarx/pullgit.sh
 # export GOBIN=${GOPATH}/bin
 # export GOPATH=${TRAVIS_BUILD_DIR}
  echo ${GOPATH}
  echo ${TRAVIS_BUILD_DIR}
 # export GOBIN="$EFFECTIVE_GOPATH/bin"
 # echo ${GOBIN}
 # export PATH="${PATH}:${GOBIN}"
  sudo apt-get install -y libxml2
  sudo apt-get install -y libxml2-dev
  sudo apt-get install -y pkg-config
  apt install librdkafka-dev
  # export PKG_CONFIG_PATH=${PKG_CONFIG_PATH}:/usr/lib/pkgconfig
  # git clone https://github.com/edenhill/librdkafka.git
  # cd librdkafka
  # ./configure --prefix /usr
  # make
  # sudo make install
  # sudo apt-get update && sudo apt-get install -y software-properties-common
  # sudo add-apt-repository ppa:0k53d-karl-f830m/openssl -y
  # sudo apt-get update
  sudo apt-get install openssl
  wget https://github.com/neo4j-drivers/seabolt/releases/download/v1.7.3/seabolt-1.7.3-Linux-ubuntu-16.04.deb
  sudo dpkg -i seabolt-1.7.3-Linux-ubuntu-16.04.deb
  sudo rm seabolt-1.7.3-Linux-ubuntu-16.04.deb
  # cd ${TRAVIS_BUILD_DIR}
  # Setup dependency management tool
  curl -L -s https://github.com/golang/dep/releases/download/v0.5.4/dep-linux-amd64 -o $GOPATH/bin/dep
  chmod +x $GOPATH/bin/dep
  dep ensure -vendor-only
}

run()
{

   gitBranch=`git branch | grep \* | cut -d " " -f2`
   # echo ${gitBranch}
   git status
   make swaggergen
   version=`cat VERSION`
   cd integration-tests/;
   while sleep 5m; do echo "=====[ $SECONDS seconds, docker images are still building... ]====="; done &
   make service-docker build=all || exit 1;
   kill %1;
   if [ ${TRAVIS_EVENT_TYPE} == 'cron' ] && [ ${TRAVIS_BRANCH} == "development" ]; then
     make push-dockers version="latest";
   elif [ ${TRAVIS_EVENT_TYPE} == 'cron' ] && [ ${TRAVIS_BRANCH} == "publish-docker-images" ]; then
     make push-dockers version=$version;
   fi

}

scanTwistlock()
{
  CONTAINER=$1
  sudo curl -k -ssl -u ${TL_USER}:${TL_PASS} --output /tmp/twistcli ${TL_CONSOLE_URL}/api/v1/util/twistcli && sudo chmod a+x /tmp/twistcli && sudo /tmp/twistcli images scan $CONTAINER --details -address ${TL_CONSOLE_URL} -u ${TL_USER} -p ${TL_PASS} || exit 1
}

runAdministrationService()
{

   gitBranch=`git branch | grep \* | cut -d " " -f2`
   # echo ${gitBranch}
   git status
   make swaggergen
   version=`cat VERSION`
   cd integration-tests/;
   while sleep 5m; do echo "=====[ $SECONDS seconds, docker images are still building... ]====="; done &
   make service-docker build=administration-service || exit 1;
   kill %1;
   scanTwistlock "gftn/administration-service:latest"
   if [ ${TRAVIS_EVENT_TYPE} == 'cron' ] && [ ${TRAVIS_BRANCH} == "development" ]; then
     make push-administration-service-dockers version="latest";
   elif [ ${TRAVIS_EVENT_TYPE} == 'cron' ] && [ ${TRAVIS_BRANCH} == "publish-docker-images" ]; then
     make push-administration-service-dockers version=$version;
   fi

}

runAuthService()
{

   gitBranch=`git branch | grep \* | cut -d " " -f2`
   # echo ${gitBranch}
   git status
  #  make swaggergen
   version=`cat VERSION`
   cd integration-tests/;
   while sleep 5m; do echo "=====[ $SECONDS seconds, docker images are still building... ]====="; done &
   make service-docker build=auth-service || exit 1;
   kill %1;
   scanTwistlock "gftn/auth-service:latest"
   if [ ${TRAVIS_EVENT_TYPE} == 'cron' ] && [ ${TRAVIS_BRANCH} == "development" ]; then
     make push-auth-service-dockers version="latest";
   elif [ ${TRAVIS_EVENT_TYPE} == 'cron' ] && [ ${TRAVIS_BRANCH} == "publish-docker-images" ]; then
     make push-auth-service-dockers version=$version;
   fi

}

runAnchorService()
{

   gitBranch=`git branch | grep \* | cut -d " " -f2`
   # echo ${gitBranch}
   git status
   make swaggergen
   version=`cat VERSION`
   cd integration-tests/;
   while sleep 5m; do echo "=====[ $SECONDS seconds, docker images are still building... ]====="; done &
   make service-docker build=anchor-service || exit 1;
   kill %1;
   scanTwistlock "gftn/anchor-service:latest"
   if [ ${TRAVIS_EVENT_TYPE} == 'cron' ] && [ ${TRAVIS_BRANCH} == "development" ]; then
     make push-anchor-service-dockers version="latest";
   elif [ ${TRAVIS_EVENT_TYPE} == 'cron' ] && [ ${TRAVIS_BRANCH} == "publish-docker-images" ]; then
     make push-anchor-service-dockers version=$version;
   fi

}

runApiService()
{

   gitBranch=`git branch | grep \* | cut -d " " -f2`
   # echo ${gitBranch}
   git status
   make swaggergen
   version=`cat VERSION`
   cd integration-tests/;
   while sleep 5m; do echo "=====[ $SECONDS seconds, docker images are still building... ]====="; done &
   make service-docker build=api-service || exit 1;
   kill %1;
   scanTwistlock "gftn/api-service:latest"
   if [ ${TRAVIS_EVENT_TYPE} == 'cron' ] && [ ${TRAVIS_BRANCH} == "development" ]; then
     make push-api-service-dockers version="latest";
   elif [ ${TRAVIS_EVENT_TYPE} == 'cron' ] && [ ${TRAVIS_BRANCH} == "publish-docker-images" ]; then
     make push-api-service-dockers version=$version;
   fi

}

runCryptoService()
{

   gitBranch=`git branch | grep \* | cut -d " " -f2`
   # echo ${gitBranch}
   git status
   make swaggergen
   version=`cat VERSION`
   cd integration-tests/;
   while sleep 5m; do echo "=====[ $SECONDS seconds, docker images are still building... ]====="; done &
   make service-docker build=crypto-service || exit 1;
   kill %1;
   scanTwistlock "gftn/crypto-service:latest"
   scanTwistlock "gftn/crypto-service-prod:latest"
   if [ ${TRAVIS_EVENT_TYPE} == 'cron' ] && [ ${TRAVIS_BRANCH} == "development" ]; then
     make push-crypto-service-dockers version="latest";
   elif [ ${TRAVIS_EVENT_TYPE} == 'cron' ] && [ ${TRAVIS_BRANCH} == "publish-docker-images" ]; then
     make push-crypto-service-dockers version=$version;
   fi

}

runFeeService()
{

   gitBranch=`git branch | grep \* | cut -d " " -f2`
   # echo ${gitBranch}
   git status
   make swaggergen
   version=`cat VERSION`
   cd integration-tests/;
   while sleep 5m; do echo "=====[ $SECONDS seconds, docker images are still building... ]====="; done &
   make service-docker build=fee-service || exit 1;
   kill %1;
   scanTwistlock "gftn/fee-service:latest"
   if [ ${TRAVIS_EVENT_TYPE} == 'cron' ] && [ ${TRAVIS_BRANCH} == "development" ]; then
     make push-fee-service-dockers version="latest";
   elif [ ${TRAVIS_EVENT_TYPE} == 'cron' ] && [ ${TRAVIS_BRANCH} == "publish-docker-images" ]; then
     make push-fee-service-dockers version=$version;
   fi

}

runGasService()
{

   gitBranch=`git branch | grep \* | cut -d " " -f2`
   # echo ${gitBranch}
   git status
   version=`cat VERSION`
   cd integration-tests/;
   while sleep 5m; do echo "=====[ $SECONDS seconds, docker images are still building... ]====="; done &
   make service-docker build=gas-service || exit 1;
   kill %1;
   scanTwistlock "gftn/gas-service:latest"
   if [ ${TRAVIS_EVENT_TYPE} == 'cron' ] && [ ${TRAVIS_BRANCH} == "development" ]; then
     make push-gas-service-dockers version="latest";
   elif [ ${TRAVIS_EVENT_TYPE} == 'cron' ] && [ ${TRAVIS_BRANCH} == "publish-docker-images" ]; then
     make push-gas-service-dockers version=$version;
   fi

}

runGlobalWhitelistService()
{

   gitBranch=`git branch | grep \* | cut -d " " -f2`
   # echo ${gitBranch}
   git status
   make swaggergen
   version=`cat VERSION`
   cd integration-tests/;
   while sleep 5m; do echo "=====[ $SECONDS seconds, docker images are still building... ]====="; done &
   make service-docker build=global-whitelist-service || exit 1;
   kill %1;
   scanTwistlock "gftn/global-whitelist-service:latest"
   if [ ${TRAVIS_EVENT_TYPE} == 'cron' ] && [ ${TRAVIS_BRANCH} == "development" ]; then
     make push-global-whitelist-service-dockers version="latest";
   elif [ ${TRAVIS_EVENT_TYPE} == 'cron' ] && [ ${TRAVIS_BRANCH} == "publish-docker-images" ]; then
     make push-global-whitelist-service-dockers version=$version;
   fi

}

runParticipantRegistry()
{

   gitBranch=`git branch | grep \* | cut -d " " -f2`
   # echo ${gitBranch}
   git status
   make swaggergen
   version=`cat VERSION`
   cd integration-tests/;
   while sleep 5m; do echo "=====[ $SECONDS seconds, docker images are still building... ]====="; done &
   make service-docker build=participant-registry || exit 1;
   kill %1;
   scanTwistlock "gftn/participant-registry:latest"
   if [ ${TRAVIS_EVENT_TYPE} == 'cron' ] && [ ${TRAVIS_BRANCH} == "development" ]; then
     make push-participant-registry-dockers version="latest";
   elif [ ${TRAVIS_EVENT_TYPE} == 'cron' ] && [ ${TRAVIS_BRANCH} == "publish-docker-images" ]; then
     make push-participant-registry-dockers version=$version;
   fi

}

runPaymentListener()
{

   gitBranch=`git branch | grep \* | cut -d " " -f2`
   # echo ${gitBranch}
   git status
   make swaggergen
   version=`cat VERSION`
   cd integration-tests/;
   while sleep 5m; do echo "=====[ $SECONDS seconds, docker images are still building... ]====="; done &
   make service-docker build=payment-listener || exit 1;
   kill %1;
   scanTwistlock "gftn/payment-listener:latest"
   if [ ${TRAVIS_EVENT_TYPE} == 'cron' ] && [ ${TRAVIS_BRANCH} == "development" ]; then
     make push-payment-listener-dockers version="latest";
   elif [ ${TRAVIS_EVENT_TYPE} == 'cron' ] && [ ${TRAVIS_BRANCH} == "publish-docker-images" ]; then
     make push-payment-listener-dockers version=$version;
   fi

}

runPayoutService()
{

   gitBranch=`git branch | grep \* | cut -d " " -f2`
   # echo ${gitBranch}
   git status
   make swaggergen
   version=`cat VERSION`
   cd integration-tests/;
   while sleep 5m; do echo "=====[ $SECONDS seconds, docker images are still building... ]====="; done &
   make service-docker build=payout-service || exit 1;
   kill %1;
   scanTwistlock "gftn/payout-service:latest"
   if [ ${TRAVIS_EVENT_TYPE} == 'cron' ] && [ ${TRAVIS_BRANCH} == "development" ]; then
     make push-payout-service-dockers version="latest";
   elif [ ${TRAVIS_EVENT_TYPE} == 'cron' ] && [ ${TRAVIS_BRANCH} == "publish-docker-images" ]; then
     make push-payout-service-dockers version=$version;
   fi

}

runQuotesService()
{

   gitBranch=`git branch | grep \* | cut -d " " -f2`
   # echo ${gitBranch}
   git status
   make swaggergen
   version=`cat VERSION`
   cd integration-tests/;
   while sleep 5m; do echo "=====[ $SECONDS seconds, docker images are still building... ]====="; done &
   make service-docker build=quotes-service || exit 1;
   kill %1;
   scanTwistlock "gftn/quotes-service:latest"
   if [ ${TRAVIS_EVENT_TYPE} == 'cron' ] && [ ${TRAVIS_BRANCH} == "development" ]; then
     make push-quotes-service-dockers version="latest";
   elif [ ${TRAVIS_EVENT_TYPE} == 'cron' ] && [ ${TRAVIS_BRANCH} == "publish-docker-images" ]; then
     make push-quotes-service-dockers version=$version;
   fi

}

runSendService()
{

   gitBranch=`git branch | grep \* | cut -d " " -f2`
   # echo ${gitBranch}
   git status
   make swaggergen
   version=`cat VERSION`
   cd integration-tests/;
   while sleep 5m; do echo "=====[ $SECONDS seconds, docker images are still building... ]====="; done &
   make service-docker build=send-service || exit 1;
   kill %1;
   scanTwistlock "gftn/send-service:latest"
   if [ ${TRAVIS_EVENT_TYPE} == 'cron' ] && [ ${TRAVIS_BRANCH} == "development" ]; then
     make push-send-service-dockers version="latest";
   elif [ ${TRAVIS_EVENT_TYPE} == 'cron' ] && [ ${TRAVIS_BRANCH} == "publish-docker-images" ]; then
     make push-send-service-dockers version=$version;
   fi

}


runWWGateway()
{

   gitBranch=`git branch | grep \* | cut -d " " -f2`
   # echo ${gitBranch}
   git status
   make swaggergen
   version=`cat VERSION`
   cd integration-tests/;
   while sleep 5m; do echo "=====[ $SECONDS seconds, docker images are still building... ]====="; done &
   make service-docker build=ww-gateway || exit 1;
   kill %1;
   scanTwistlock "gftn/ww-gateway:latest"
   if [ ${TRAVIS_EVENT_TYPE} == 'cron' ] && [ ${TRAVIS_BRANCH} == "development" ]; then
     make push-ww-gateway-dockers version="latest";
   elif [ ${TRAVIS_EVENT_TYPE} == 'cron' ] && [ ${TRAVIS_BRANCH} == "publish-docker-images" ]; then
     make push-ww-gateway-dockers version=$version;
   fi

}

runAutomationService()
{

   gitBranch=`git branch | grep \* | cut -d " " -f2`
   # echo ${gitBranch}
   git status
   make swaggergen
   version=`cat VERSION`
   cd integration-tests/;
   while sleep 5m; do echo "=====[ $SECONDS seconds, docker images are still building... ]====="; done &
   make service-docker build=automation-service || exit 1;
   kill %1;
   scanTwistlock "gftn/automation-service:latest"
   if [ ${TRAVIS_EVENT_TYPE} == 'cron' ] && [ ${TRAVIS_BRANCH} == "development" ]; then
     make push-automation-service-dockers version="latest";
   elif [ ${TRAVIS_EVENT_TYPE} == 'cron' ] && [ ${TRAVIS_BRANCH} == "publish-docker-images" ]; then
     make push-automation-service-dockers version=$version;
   fi

}

runMSKCLI()
{

   gitBranch=`git branch | grep \* | cut -d " " -f2`
   # echo ${gitBranch}
   git status
   make swaggergen
   version=`cat VERSION`
   cd integration-tests/;
   while sleep 5m; do echo "=====[ $SECONDS seconds, docker images are still building... ]====="; done &
   make service-docker build=msk-cli || exit 1;
   kill %1;
   scanTwistlock "gftn/msk-cli:latest"
   if [ ${TRAVIS_EVENT_TYPE} == 'cron' ] && [ ${TRAVIS_BRANCH} == "development" ]; then
     make push-msk-cli-dockers version="latest";
   elif [ ${TRAVIS_EVENT_TYPE} == 'cron' ] && [ ${TRAVIS_BRANCH} == "publish-docker-images" ]; then
     make push-msk-cli-dockers version=$version;
   fi

}

runMSKMonitoringUI()
{

   gitBranch=`git branch | grep \* | cut -d " " -f2`
   # echo ${gitBranch}
   git status
   make swaggergen
   version=`cat VERSION`
   cd integration-tests/;
   while sleep 5m; do echo "=====[ $SECONDS seconds, docker images are still building... ]====="; done &
   make service-docker build=msk-ui || exit 1;
   kill %1;
   scanTwistlock "gftn/msk-burrow-ui:latest"
   if [ ${TRAVIS_EVENT_TYPE} == 'cron' ] && [ ${TRAVIS_BRANCH} == "development" ]; then
     make push-msk-monitoring-ui-dockers version="latest";
   elif [ ${TRAVIS_EVENT_TYPE} == 'cron' ] && [ ${TRAVIS_BRANCH} == "publish-docker-images" ]; then
     make push-msk-monitoring-ui-dockers version=$version;
   fi

}

runMSKMonitoringServer()
{

   gitBranch=`git branch | grep \* | cut -d " " -f2`
   # echo ${gitBranch}
   git status
   make swaggergen
   version=`cat VERSION`
   cd integration-tests/;
   while sleep 5m; do echo "=====[ $SECONDS seconds, docker images are still building... ]====="; done &
   make service-docker build=msk-server || exit 1;
   kill %1;
   scanTwistlock "gftn/msk-burrow:latest"
   if [ ${TRAVIS_EVENT_TYPE} == 'cron' ] && [ ${TRAVIS_BRANCH} == "development" ]; then
     make push-msk-monitoring-server-dockers version="latest";
   elif [ ${TRAVIS_EVENT_TYPE} == 'cron' ] && [ ${TRAVIS_BRANCH} == "publish-docker-images" ]; then
     make push-msk-monitoring-server-dockers version=$version;
   fi

}

afterFailure()
{

  # dump the last 2000 lines of our build to show error
  tail -n 2000 build.log

}

afterSuccess()
{

  # Log that the build was a success
  echo "build succeeded.."

}

CMD=$1

if [ $CMD = "beforeInstall" ]; then
  beforeInstall
elif [ $CMD = "run" ]; then
  run
elif [ $CMD = "runAdministrationService" ]; then
  runAdministrationService
elif [ $CMD = "runAnchorService" ]; then
  runAnchorService
elif [ $CMD = "runApiService" ]; then
  runApiService
elif [ $CMD = "runAuthService" ]; then
  runAuthService
elif [ $CMD = "runAutomationService" ]; then
  runAutomationService
elif [ $CMD = "runCryptoService" ]; then
  runCryptoService
elif [ $CMD = "runFeeService" ]; then
  runFeeService
elif [ $CMD = "runGasService" ]; then
  runGasService
elif [ $CMD = "runGlobalWhitelistService" ]; then
  runGlobalWhitelistService
elif [ $CMD = "runParticipantRegistry" ]; then
  runParticipantRegistry
elif [ $CMD = "runPaymentListener" ]; then
  runPaymentListener
elif [ $CMD = "runPayoutService" ]; then
  runPayoutService
elif [ $CMD = "runQuotesService" ]; then
  runQuotesService
elif [ $CMD = "runSendService" ]; then
  runSendService
elif [ $CMD = "runWWGateway" ]; then
  runWWGateway
elif [ $CMD = "runMSKCLI" ]; then
  runMSKCLI
elif [ $CMD = "runMSKMonitoringServer" ]; then
  runMSKMonitoringServer
elif [ $CMD = "runMSKMonitoringUI" ]; then
  runMSKMonitoringUI
elif [ $CMD = "afterFailure" ]; then
 afterFailure
elif [ $CMD = "afterSuccess" ]; then
 afterSuccess
fi