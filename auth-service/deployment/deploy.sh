# USAGE: 
#     local: $ bash ./deployment/deploy.sh --debug --env dev

# PURPOSE:
# The purpose of this script is to deploy the auth-service and 
# firebase portal, database secruity rules and trigger functions.
# This can be used by travis and locally to deploy. Travis requires
# some initial configuration for setting environment variable and 
# and deployment branches, see setTravisDeployEnvVars() below.

# TRAVIS NOTE - Encrypted Env Var as file: 
# Travis env vars are uploaded to travis as described
# here: https://docs.travis-ci.com/user/encrypting-files/.
# To make the process even simplier, if you already have extracted the 
# [credentials.zip](https://ibm.box.com/s/7mgtvv5yrvjs1k95pm74ib2jnpmstt4n) file to 
# gftn-service/auth-service/.credentials
# you can update these env vars by opening up .vscode for the auth-service 
# and running the process in the launch.json for "deploy travis env".

# # mock the followng for testing travis dev env deployment
# TRAVIS_BRANCH="dev-gftn"
# TRAVIS_PULL_REQUEST="false"
# export CICD_CRED_KEY_V11="REPLACE_HERE" ; source $CICD_CRED_KEY_V11
# export CICD_CRED_IV_V11="REPLACE_HERE" ; source $CICD_CRED_IV_V11

# init default env vars
gcloud_project=""
serviceAccount=""
firebase_ci_token=""
domain_https=""
domain_short=""
key=""
iv=""
secrets="use_travis"
cloud='gcloud'
debugPath=""
verbose=0

show_help()
{ 
cat << EOF

Deploy auth-service in nodejs:

--env (alias: -e) <local|dev|qa|st|prod> OPTIONAL
  Target auth-service environment to build deploymet files
  If omitted, env will be set based on the deployment branch  
  NOTE: Automated CI/CD deployment occurs only when running build 
  in travis from proper deployment braches
  ie: gftn-dev, gftn-qa, gftn-st, gftn-prod

--debug (alias: -d) OPTIONAL
  Used to run local and developement environment applications locally.
  By default, travis stored encryption env vars are used
  NOTE: production credentials (eg: st, qa, prod) are not available is local development 

EOF

}

# flags usage - https://archive.is/5jGpl#selection-709.0-715.1 or https://archive.is/TRzn4
while :; do
    case $1 in
          -h|-\?|--help)   # Call a "show_help" function to display a synopsis, then exit.
              show_help
              exit
              ;;
          --debug) 
              secrets="use_debug"
              debugPath="-debug" # appended to debug path for development/local creds
              ;;
          -e|--env) # Takes an option argument, ensuring it has been specified.
              if [ -n "$2" ]; then
                  env=$2
                  echo "env" $2
                  shift
              else
                  printf "$(tput setaf 1)ERROR: \"-e|--env\" requires a non-empty option argument.\n" >&2
                  exit 1
              fi
              ;;
          -v|--verbose)
              verbose=$((verbose + 1)) # Each -v argument adds 1 to verbosity.
              ;;
          --)              # End of all options.
              shift
              break
              ;;
          -?*)
              printf 'WARN: Unknown option (ignored): %s\n' "$1" >&2
              ;;
          *)               # Default case: If no more options then break out of the loop.
              break
      esac
  
      shift
  done

checkerrors(){

    errors=$?

    # helps to alert failure if code does not build properly
    # and helps vscode debuger exit on typscript/lint error
    if [ $errors -ne 0 ]
    then
        echo $(tput setaf 1)"======= errors: "$errors" ======="
        echo $(tput setaf 1)"ERROR: build process error, see deploy.sh"
        echo $(tput setaf 1)"For help Run $ bash ./deployment/deploy.sh -h"$(tput sgr0)
        exit 1
    fi  
    
}

# build dependent files for portal & auth-service
build(){

    echo "building files for deployment..."

    if [ $secrets = 'use_debug' ]
    then
        # debug creds
        bash ./deployment/build.sh --env $env --debug --cloud gcloud 
    else 
        # travis creds
        bash ./deployment/build.sh --env $env --cloud gcloud 
    fi 
}

# deploy portal, triggers, and firebase database rules
deployPortal() {

    echo "deploying portal, triggers, and database rules to firebase..."

    # IMPORTANT: run after deployAuthServiceGae() - becuase the directory location after deployment is important

    # count errors up to this point
    errors=$?

    if [ $errors -ne 0 ]
    then
        echo $(tput setaf 1)"======= errors: "$errors" ======="
        echo $(tput setaf 1)"ERROR: Erorr: not deploying portal, see deploy.sh"
        exit 1
    fi    

    # update firebase.json for dynamic content header
    # get https domain ie: https://xxxx.com
    domain_https=$(file=".credentials$debugPath-v29/raw/$env/img/env.json" key="site_root"  node ./project-provisioning/src/read_val_from_json.js)
    # get short domain ie: xxxx.com
    domain_short=$(str_original="$domain_https" str_substr="https://" str_replace="" node ./project-provisioning/src/replace.js)
    # update world wire web CSP Headers definition template
    # ../gftn-web/firebase-template.json is used to generate
    # ../gftn-web/firebase.json which is used to apply
    # CSP headers in firebase, see docs:
    #  https://developer.mozilla.org/en-US/docs/Web/HTTP/CSP
    #  https://firebase.google.com/docs/hosting/full-config#headers
    template_file="../gftn-web/firebase-template.json" replace_file="../gftn-web/firebase.json" template_txt="{{{PORTAL_DOMAIN}}}" replace_txt="$domain_short"  node ./project-provisioning/src/replace_from_template.js

    # navigate to portal
    cd ../gftn-web

    npm i ; ng build -c="$env" # IMPORTANT: env name must be the same as the build in angular.json ie: st | dev | qa | prod 

    # deploy firebase database rules
    firebase-bolt database.rules.bolt || exit 1 
    firebase deploy --only database --project="$gcloud_project" --token=$firebase_ci_token

    # deploy firebase portal (hosting)
    firebase deploy --only hosting --project="$gcloud_project" --token=$firebase_ci_token
    
    # deploy firebase functions
    cd functions ; npm i ; cd ..
    firebase deploy --only functions --project="$gcloud_project" --token=$firebase_ci_token

}

# deploy auth-service to google app-engine (see `make docker` for dockerized auth-service)
# IMPORTANT: run this before deployPortal() - becuase the directory location after deployment is important
deployAuthServiceGae(){

    echo "deploying authentication-service to gcloud..."

    if [ $cloud = 'gcloud' ]
    then

        gcloud auth activate-service-account --key-file=$serviceAccount || exit 1
        gcloud config set project $gcloud_project || exit 1      

        # cd into file with app.yaml to deploy to GAE
        cd authentication/build

        checkerrors
        
        if [ $secrets = "use_travis" ]
        then
            # no user prompt
            echo "no gcloud user prompt"
            gcloud app deploy --quiet
        else
            # with user prompt
            gcloud app deploy
        fi

        # back to ./auth-service
        cd ../..
    fi

}

# update micro-serivces, via effectively pulling the updated version and restaring when the pod gets deleted out
deployMicroservices(){

    echo "restarting pods for k8, effectivly updating deployment..."

    # IMPORTANT: run this before deployPortal() - becuase the directory location after deployment is important

    if [ $TRAVIS_BRANCH = 'dev-gftn' ] && [ $TRAVIS_PULL_REQUEST = 'false' ]
    then
        env="dev" # name of the folder in ./.credentials/{env} to use
        export REGION='us-west-2' && export CLUSTER_NAME='terraform-eks-ww-dev' && aws eks update-kubeconfig --name $CLUSTER_NAME --region $REGION && kubectl -n default delete --all pods
    elif [ $TRAVIS_BRANCH = 'qa-gftn' ] && [ $TRAVIS_PULL_REQUEST = 'false' ]
    then
        env="qa" # name of the folder in ./.credentials/{env} to use
        export REGION='us-west-2' && export CLUSTER_NAME='terraform-eks-ww-qa' && aws eks update-kubeconfig --name $CLUSTER_NAME --region $REGION && kubectl -n default delete --all pods
    elif [ $TRAVIS_BRANCH = 'staging-gftn' ] && [ $TRAVIS_PULL_REQUEST = 'false' ]
    then
        env="st" # name of the folder in ./.credentials/{env} to use
        export REGION='us-east-2' && export CLUSTER_NAME='terraform-eks-ww-st' && aws eks update-kubeconfig --name $CLUSTER_NAME --region $REGION && kubectl -n default delete --all pods
    elif [ $TRAVIS_BRANCH = 'test-gftn' ] && [ $TRAVIS_PULL_REQUEST = 'false' ]
    then
        env="tn" # name of the folder in ./.credentials/{env} to use
        export REGION='us-east-2' && export CLUSTER_NAME='terraform-eks-ww-st' && aws eks update-kubeconfig --name $CLUSTER_NAME --region $REGION && kubectl -n default delete --all pods
    elif [ $TRAVIS_BRANCH = 'live-gftn' ] && [ $TRAVIS_PULL_REQUEST = 'false' ]
    then
        env="prod"
        export REGION='us-east-1' && export CLUSTER_NAME='terraform-eks-ww-prod' && aws eks update-kubeconfig --name $CLUSTER_NAME --region $REGION && kubectl -n default delete --all pods
        else
        echo "\n\n >>>>> INFO <<<<<: No conditions were matched for this deployment. This could be becuase this is not a branch that enables deployment or this is a pull-request(not a merge.)\n\n"
        exit 0
    fi

}

# set environment vars from travis via './.credentials'
setTravisDeployEnvVars()
{
    echo "setting environment variables for deployment based on travis branch..."

    echo "branch = $TRAVIS_BRANCH"
    echo "pull request (no. or false) = $TRAVIS_PULL_REQUEST"
    echo "\n"

    if [ $TRAVIS_BRANCH = 'dev-gftn' ] && [ $TRAVIS_PULL_REQUEST = 'false' ]
    then
        env="dev" # name of the folder in ./.credentials/{env} to use 
    
    elif [ $TRAVIS_BRANCH = 'qa-gftn' ] && [ $TRAVIS_PULL_REQUEST = 'false' ]
    then
        env="qa" # name of the folder in ./.credentials/{env} to use
    
    elif [ $TRAVIS_BRANCH = 'staging-gftn' ] && [ $TRAVIS_PULL_REQUEST = 'false' ]
    then
        env="st" # name of the folder in ./.credentials/{env} to use
    
    elif [ $TRAVIS_BRANCH = 'test-gftn' ] && [ $TRAVIS_PULL_REQUEST = 'false' ]
    then
        env="tn" # name of the folder in ./.credentials/{env} to use
    
    elif [ $TRAVIS_BRANCH = 'pen1-gftn' ] && [ $TRAVIS_PULL_REQUEST = 'false' ]
    then
        env="pen1" # name of the folder in ./.credentials/{env} to use
    
    elif [ $TRAVIS_BRANCH = 'live-gftn' ] && [ $TRAVIS_PULL_REQUEST = 'false' ]
    then
        env="prod"
    
    else
        echo "\n\n >>>>> INFO <<<<<: No conditions were matched for this deployment. This could be becuase this is not a branch that enables deployment or this is a pull-request(not a merge.)\n\n"
        exit 0
    fi

}

# set environment vars via './.credentials'
setEnvs()
{
    echo "setting environment variables for deployment..."

    # transpile ts to js helpers
    tsc "./project-provisioning/src/replace.ts"
    tsc "./project-provisioning/src/read_val_from_json.ts" 
    tsc "./project-provisioning/src/replace_from_template.ts"

    # firebase deployment ci token
    firebase_ci_token=$(file=".credentials$debugPath-v29/fb_deploy.json" key="ci_token"  node ./project-provisioning/src/read_val_from_json.js)

    # gcloud projectId
    gcloud_project=$(file=".credentials$debugPath-v29/raw/$env/img/env.json" key="gae_service"  node ./project-provisioning/src/read_val_from_json.js)
    echo 'Deploying using gcloud project: '$gcloud_project
    # gcloud service account credential .json file
    serviceAccount=".credentials$debugPath-v29/raw/$env/deploy/deploy.json"
}

extractCICDCreds(){

    echo "decrypting application layer secrets for deployment..."
    
    if [ $secrets = "use_debug" ]
    then
        # if directory does not already exist
        if [ ! -d ".credentials-debug-v29" ] 
        then
            # local debug 
            bash extract-debug-creds.sh
            checkerrors
        fi
    else
        # if directory does not already exist
        if [ ! -d ".credentials-v29" ] 
        then    
            # travis creds
            key=${CICD_CRED_KEY_V29} # set in travis as secret env var 
            iv=${CICD_CRED_IV_V29} # set in travis as secret env var
            bash ./secret-mgmt/cicd/decrypt.sh ./cicd-cred-v29.tgz.enc ${key} ${iv}
            checkerrors
        fi
    fi

}

# if no arguments use travis
if [ $secrets = "use_travis" ]
then
    
    # travis deployment 
    # ================
    
    printf "\n================== Using travis for deployment ==================\n"

    sh -x install_dependencies.sh # install gcloud from travis yml
    # bash run-tests.sh
    setTravisDeployEnvVars
    extractCICDCreds
    setEnvs
    build
    deployAuthServiceGae 
    deployPortal
    # deployMicroservices
else
    
    # local deployment 
    # ================
    
    printf "\n================== NOT using travis for deployment ==================\n"
    
    # sh install_dependencies.sh
    # bash run-tests.sh
    extractCICDCreds
    setEnvs
    build
    deployAuthServiceGae
    deployPortal
    deployMicroservices
fi