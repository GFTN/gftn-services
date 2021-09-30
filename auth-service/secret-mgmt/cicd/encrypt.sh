#! /bin/bash

inputDir=$1 
outputName=$2
newVersion=$3

tar -zcf ${outputName}.tgz $inputDir

# create pepper and iterate and output to file
randKey=$(openssl rand -hex 32)
iv=$(openssl dgst -r -md5 ${outputName}.tgz | cut -d' ' -f 1)
randIter=41 # default val, set random using - $(shuf -i 0-100 -n 1)
# echo iter $randIter
salt='c971dc7fda3e400a36bdfcd3036d8111' # randomly hardcoded default val - create rand val with $ openssl rand -hex 16
# echo salt $salt

# # ubuntu bionic or higher:
# openssl enc -aes-256-cbc -base64 -pbkdf2 -e \
#     -salt -S $salt \
#     -iv $iv \
#     -in ${outputName}.tgz -out ${outputName}.tgz.enc \
#     -iter $randIter \
#     -K $randKey

# ubuntu xenial:
openssl enc -aes-256-cbc -base64 -e \
    -salt -S $salt \
    -iv $iv \
    -in ${outputName}.tgz -out ${outputName}.tgz.enc \
    -K $randKey

rm ${outputName}.tgz

echo "CICD_CRED_KEY_V${newVersion}=${randKey}"
echo "CICD_CRED_IV_V${newVersion}=${iv}"
echo ""
echo "# Run decryption using: \nsh ./secret-mgmt/cicd/decrypt.sh ./${outputName}.tgz.enc ${randKey} ${iv} # end"