#! /bin/bash

# USAGE: 
#     $ sh secret-mgmt/cicd/decrypt.sh cicd-cred-v${VERSION_NUMBER_HERE}.tgz.enc someKey someIv


checkerrors(){

    errors=$?

    # helps to alert failure if code does not build properly
    # and helps vscode debuger exit on typscript/lint error
    if [ $errors -ne 0 ]
    then
        echo $(tput setaf 1)"======= errors: "$errors" ======="
        echo $(tput setaf 1)"ERROR: build process error, see decrypt.sh"$(tput sgr0)
        exit 1
    fi  
    
}

inputFile=${1%.tgz.enc} # used /someDir to create

# create pepper and iterate and output to file
key=$2
# echo key $key
iv=$3
# echo iv  $iv
iter=41 # default val
# echo iter $iter
salt='c971dc7fda3e400a36bdfcd3036d8111' # randomly hardcoded default val - create rand val with $ openssl rand -hex 16
# echo salt $salt

# # ubuntu bionic or higher:
# openssl enc -aes-256-cbc -base64 -pbkdf2 -d \
#     -salt -S $salt \
#     -iv $iv \
#     -in ${inputFile}.tgz.enc -out ${inputFile}-out.tgz \
#     -iter $iter \
#     -K $key

# ubuntu xenial:
openssl enc -aes-256-cbc -base64 -d \
    -salt -S $salt \
    -iv $iv \
    -in ${inputFile}.tgz.enc -out ${inputFile}-out.tgz \
    -K $key

checkerrors

resultingHash=$(openssl dgst -r -md5 ${inputFile}-out.tgz | cut -d' ' -f 1)

checkerrors

# # output results of vals
# echo $resultingHash $iv

# check if the output file's hash matches the iv
if [ $resultingHash = $iv ]
then

    # extract compressed tar
    tar -zxf ${inputFile}-out.tgz

    # clean-up zip file
    rm ${inputFile}-out.tgz

    checkerrors

    echo $(tput setaf 2)'encryption integrity maintained for' ${inputFile}'.tgz.enc'$(tput sgr0)

else
    checkerrors
    echo 'IMPORTANT: decyption integrity compromised rotate travis creds immediately'
fi


