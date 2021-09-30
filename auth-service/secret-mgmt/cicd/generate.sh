#! /bin/bash

# USAGE:
# $ bash secret-mgmt/cicd/generate.sh NEW_VERSION_NUMBER_HERE

# set version
old_version=$1
next_version=$((old_version + 1))

# copy .credentials dir (so that when it gets 
# extracted the root dirname is .credentials)
# echo $old_version

# encrypt production creds for cicd
printf "$(tput setaf 4)$(tput bold)\nProduction cicd decryption keys: \n$(tput sgr0)"
deployKeys=$(sh secret-mgmt/cicd/encrypt.sh .credentials-v${next_version} cicd-cred-v${next_version} ${next_version})

# create local debug creds by deleting out deployment creds
cp -Rf .credentials-v${next_version} .credentials-debug-v${next_version} 
rm -Rf .credentials-debug-v${next_version}/enc/prod .credentials-debug-v${next_version}/enc/qa .credentials-debug-v${next_version}/enc/st
rm -Rf .credentials-debug-v${next_version}/k8-secrets/prod .credentials-debug-v${next_version}/k8-secrets/qa .credentials-debug-v${next_version}/k8-secrets/st
rm -Rf .credentials-debug-v${next_version}/raw/prod .credentials-debug-v${next_version}/raw/qa .credentials-debug-v${next_version}/raw/st 

# encrypt development creds for local development
printf "$(tput setaf 4)$(tput bold)Local decryption keys generated (see output added to ./extract-debug-creds.sh) \n\n$(tput sgr0)"
debugKeys=$(sh secret-mgmt/cicd/encrypt.sh .credentials-debug-v${next_version} cicd-cred-debug-v${next_version} ${next_version})

# update values in extract-debug-creds.sh
read -d '' sql << EOF
# These are development ONLY credentials
# NOTE: these decryption keys are visible to
# aid developers with **development** environment only. 
# these are not sensitive becuase these can only be used 
# for accessing development
# other environment credentials are encrypted via cicd-cred-vN.tgz.enc
# and these decryption values are only known to the "buildEnv" such as travis
# as they are input manually via the secret manager for the cicd tool 
# (ie: travis.com env var secrets console)
${debugKeys}
EOF

# copy contents to file extract-debug-creds.sh
echo "$sql" > extract-debug-creds.sh

# clean-up temporarily used folder
rm -Rf .credentials-debug-v${next_version}

# remove previous out-dated cicd credentials this is important 
# becuase it forces the ci/cd (ie: travis) to error out if 
# new decryption keys are not added cicd secret manager 
# (ie: travis secret env var via web console at travis.com)
rm ./cicd-cred-v${old_version}.tgz.enc
rm ./cicd-cred-debug-v${old_version}.tgz.enc

# update dependent file paths (eg: launch.json, .travis.yml)
printf "$(tput setaf 4)$(tput bold)Updating dependent files for rotated credentials...$(tput sgr0)"
sed -i "s/-v${old_version}/-v${next_version}/g" .vscode/launch.json .vscode/tasks.json ./deployment/deploy.sh ./deployment/build.sh ./authentication/src/config.ts
sed -i "s/\"cred_version\": \"${old_version}\"/\"cred_version\": \"${next_version}\"/g" .vscode/launch.json
sed -i "s/\${CICD_CRED_KEY_V${old_version}}/\${CICD_CRED_KEY_V${next_version}}/g" ../.travis.yml ./deployment/build.sh ./deployment/deploy.sh ./secret-mgmt/cicd/generate.sh
sed -i "s/\${CICD_CRED_IV_V${old_version}}/\${CICD_CRED_IV_V${next_version}}/g" ../.travis.yml ./deployment/build.sh ./deployment/deploy.sh ./secret-mgmt/cicd/generate.sh
sed -i "s/cicd-cred-v${old_version}.tgz.enc/cicd-cred-v${next_version}.tgz.enc/g" ../.travis.yml ../.travis.yml ./deployment/build.sh ./deployment/deploy.sh
sed -i "s/cred_version=${old_version}/cred_version=${next_version}/g" ./deployment/build.sh
printf "done\n"

printf "$(tput setaf 3)Go to $(tput bold)https://travis.ibm.com/gftn/gftn-services/settings$(tput sgr0)$(tput setaf 3) and the following decryption secrets: $(tput sgr0)\n"
printf "${deployKeys}"
