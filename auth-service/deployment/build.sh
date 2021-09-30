# USAGE: bash ./deployment/build.sh -h
#  local: $ bash ./deployment/build.sh --env dev --debug --cloud gcloud

# defaults
env=''
cloud='gcloud'
secrets="use_travis"
npmInstall="no-install"

show_help()
{ 
cat << EOF

Build authentication-service in nodejs:

--cloud <gcloud | docker> (default is "gcloud") REQUIRED
  Target build for specific cloud. eg: gcloud does not need docker 
  image, nor does docker container need does not need app.yaml or .gcloudgitignore

--env (alias: -e) <local|dev|qa|st|prod> OPTIONAL (REQUIRED when '--cloud gcloud' present)
  Target auth-service environment to build deploymet files  
  NOTE: Not required when build is docker becuase app.yml is 
  not required for docker. 

--debug (alias: -d) OPTIONAL
  Used to run local and developement environment applications locally.
  By default, travis stored encryption env vars are used
  NOTE: production credentials (eg: st, qa, prod) are not available is local development 

--npm-install (alias: -i) OPTIONAL
  Install dependencies along with build, which might be important for running build locally.
  By default, node_modules are NOT installed

EOF

}

# flags usage - https://archive.is/5jGpl#selection-709.0-715.1 or https://archive.is/TRzn4
while :; do
    case $1 in
          -h|-\?|--help) # Call a "show_help" function to display a synopsis, then exit.
              show_help
              exit
              ;;
          -d|--debug) 
              secrets="use_debug"
              ;;
          -i|--npm-install) 
              npmInstall="install"
              ;;              
          -e|--env) # Takes an option argument, ensuring it has been specified.
              if [ -n "$2" ]; then
                  env=$2
                  shift
              else
                  printf "$(tput setaf 1)ERROR: \"-e|--env\" requires a non-empty option argument.\n" >&2
                  exit 1
              fi
              ;;
          -c|--cloud) # Takes an option argument, ensuring it has been specified.
              if [ -n "$2" ]; then
                  cloud=$2
                  if [[ $cloud = 'docker' ]]
                    then
                        # env can be anything when building for docker
                        # this is becuase env is only important to create
                        # the deployment app.yml for gcloud
                        env='dev'
                    fi
                  shift
              else
                  printf '$(tput setaf 1)ERROR: "-t|--cloud" requires a non-empty option argument.\n' >&2
                  exit 1
              fi
              ;;              
          -v|--verbose)
              verbose=$((verbose + 1)) # Each -v argument adds 1 to verbosity.
              ;;
          --) # End of all options.
              shift
              break
              ;;
          -?*)
              printf 'WARN: Unknown option (ignored): %s\n' "$1" >&2
              ;;
          *)  # Default case: If no more options then break out of the loop.
              break
      esac
  
      shift
  done

# validate that required fields were provided, otherwise error:
if [ -z "$env" ] && [ "${env+xxx}" = "xxx" ]
then 
    echo $(tput setaf 1)"======= errors: "$errors" ======="
    echo $(tput setaf 1)"ERROR: input missing - need defined 'env' - ie: dev | qa | st | prod, see build.sh"
    echo $(tput setaf 1)"For help Run $ bash ./deployment/build.sh -h"$(tput sgr0)
    exit 1
fi

checkerrors(){

    errors=$?

    # helps to alert failure if code does not build properly
    # and helps vscode debuger exit on typscript/lint error
    if [ $errors -ne 0 ]
    then
        echo $(tput setaf 1)"======= errors: "$errors" ======="
        echo $(tput setaf 1)"ERROR: build process error, see build.sh"
        echo $(tput setaf 1)"For help Run $ bash ./deployment/build.sh -h"$(tput sgr0)
        exit 1
    fi  
    
}

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
        bash -x ./secret-mgmt/cicd/decrypt.sh ./cicd-cred-v29.tgz.enc ${key} ${iv}
        checkerrors
    fi
fi

# start with clean build dir by removing all files 
rm -Rf authentication/build || true
mkdir authentication/build || true

# ========= start - used by /deployment to create /build ===========: 

# copy over files after transpiling ts to js
bash ./authentication/prelaunch.sh --prod ; checkerrors
bash ./deployment/deployment_prelaunch.sh ; checkerrors

# ========= end - used by /deployment to create /build ===========: 

# used pass.txt only used during build to decrypt
# this is passed in dynamically to the node app.js in prod
# via a secret manager like kubernetes secret manager or in app.yaml for gae

env=$env \
cred_version=29 \
node deployment/lib/index.js

checkerrors

if [[ $cloud = 'gcloud' ]]
then
    # remove unecessary files for cloud build

    # Don't need Dockerfile for kubernetes deployment 
    # to store with source and for potentially 
    # app-engine in the future if container 
    # deployment is preferable on gcloud
    rm -f authentication/build/docker-entrypoint.sh
    rm -f authentication/build/Dockerfile
    rm -Rf authentication/build/.certs

    # self-signed-certs are needed for local debug of IBMiD
    # rm -Rf authentication/build/.self-signed-certs 

    checkerrors
fi

if [[ $cloud = 'docker' ]]
then
    # remove unecessary files for cloud build
    rm -f authentication/build/app.yaml
    rm -f authentication/build/.gcloudignore
    rm -f authentication/build/Dockerfile
    rm -Rf authentication/build/.certs
    
    # self-signed-certs are needed for local debug of IBMiD
    # rm -Rf authentication/build/.self-signed-certs 

    checkerrors
fi

# install dependencies (include npm_modules in ./build)
if [[ $npmInstall = 'install' ]]
then
    cd authentication/build 
    npm i
    cd ../..
    checkerrors
fi
printf "$(tput setaf 3)\nNOTE: node_modules not installed with build. Run \"$ sh ./deployment/build.sh -i\" if node_modules install desired.\n\n"$(tput sgr0)


