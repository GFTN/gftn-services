#!/bin/sh

ENVIRONMENT=$1
#34.80.67.4

if [ $ENVIRONMENT == 'dev' ]; then
  ENV="eksdev"
  HOST="ww-dev"
  HORIZONURL="http:\/\/35.197.35.7:1234"
  STELLARNETWORK="Standalone World Wire Network ; Mar 2019"
  FIREBASEURL="https:\/\/dev-2-c8774.firebaseio.com"
  ENABLEJWT="true"
  export REGION='us-west-2' && export CLUSTER_NAME='terraform-eks-ww-dev' && aws eks update-kubeconfig --name $CLUSTER_NAME --region $REGION
elif  [ $ENVIRONMENT == 'qa'  ]; then
  ENV="eksqa"
  HOST="ww-qa"
  HORIZONURL="http:\/\/35.197.35.7:1234"
  STELLARNETWORK="Standalone World Wire Network ; Mar 2019"
  FIREBASEURL="https:\/\/qatwo-251007.firebaseio.com"
  ENABLEJWT="true"
  export REGION='us-west-2' && export CLUSTER_NAME='terraform-eks-ww-qa' && aws eks update-kubeconfig --name $CLUSTER_NAME --region $REGION
elif [ $ENVIRONMENT == 'st' ]; then
  ENV="st"
  HOST="ww-st"
  HORIZONURL="http:\/\/35.197.35.7:1234"
  STELLARNETWORK="Standalone World Wire Network ; Mar 2019"
  FIREBASEURL="https:\/\/st21-251107.firebaseio.com"
  ENABLEJWT="true"
  export REGION='us-east-2' && export CLUSTER_NAME='terraform-eks-ww-st' && aws eks update-kubeconfig --name $CLUSTER_NAME --region $REGION
elif [ $ENVIRONMENT == 'pen' ]; then
  ENV="pen"
  HOST="ww-pen"
  HORIZONURL="http:\/\/35.197.35.7:1234"
  STELLARNETWORK="Standalone World Wire Network ; Mar 2019"
  FIREBASEURL="https:\/\/pen1-260919.firebaseio.com"
  ENABLEJWT="true"
  export REGION='us-east-2' && export CLUSTER_NAME='terraform-eks-ww-pen' && aws eks update-kubeconfig --name $CLUSTER_NAME --region $REGION
elif [ $ENVIRONMENT == 'prod' ]; then
  ENV="prod"
  HOST="ww"
  HORIZONURL="https:\/\/horizon.worldwire.io"
  STELLARNETWORK="Public Global Stellar Network ; September 2015"
  FIREBASEURL="https:\/\/prod-251807.firebaseio.com"
  ENABLEJWT="true"
  export REGION='us-east-1' && export CLUSTER_NAME='terraform-eks-ww-prod' && aws eks update-kubeconfig --name $CLUSTER_NAME --region $REGION
fi

OLDENVIRONMENT=environment_variable
OLDHOST=host_variable
OLDHORIZONURL=horizon_url_variable
OLDSTELLARNETWORK=stellar_network_variable
OLDENABLEJWT=enable_jwt_variable
OLDFIREBASEURL=firebase_url_variable

cd "../../automation-service/k8s" || exit 1
dir=$(pwd)
echo $dir

NEWPATH="$dir/files/$ENVIRONMENT/configs"

mkdir -p $NEWPATH

# copy files into new files
cp -r "$dir/configs/template/" $NEWPATH/
cd "$NEWPATH/template"
pwd

# replace
if [[ "$OSTYPE" == "linux-gnu" ]]; then
  find . -type f | xargs sed -i "s/$OLDENVIRONMENT/$ENV/g"
  find . -type f | xargs sed -i "s/$OLDHOST/$HOST/g"
  find . -type f | xargs sed -i "s/$OLDHORIZONURL/$HORIZONURL/g"
  find . -type f | xargs sed -i "s/$OLDSTELLARNETWORK/$STELLARNETWORK/g"
  find . -type f | xargs sed -i "s/$OLDENABLEJWT/$ENABLEJWT/g"
  find . -type f | xargs sed -i "s/$OLDFIREBASEURL/$FIREBASEURL/g"
elif [[ "$OSTYPE" == "darwin"* ]]; then
  find . -type f | xargs sed -i '' -e "s/$OLDENVIRONMENT/$ENV/g"
  find . -type f | xargs sed -i '' -e "s/$OLDHOST/$HOST/g"
  find . -type f | xargs sed -i '' -e "s/$OLDHORIZONURL/$HORIZONURL/g"
  find . -type f | xargs sed -i '' -e "s/$OLDSTELLARNETWORK/$STELLARNETWORK/g"
  find . -type f | xargs sed -i '' -e "s/$OLDENABLEJWT/$ENABLEJWT/g"
  find . -type f | xargs sed -i '' -e "s/$OLDFIREBASEURL/$FIREBASEURL/g"
fi

printf "Ready to deploy global-configs on %s k8s cluster\n" "$ENVIRONMENT"

kubectl apply -f "./env-configmap.yaml"