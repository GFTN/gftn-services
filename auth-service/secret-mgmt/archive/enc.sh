# USAGE: 
#     encrypt file - sh enc.sh e somepass input.txt input.txt.enc pepper_iter.txt
#     decrypt file - sh enc.sh d somepass input.txt.enc input.txt pepper_iter.txt

# NOTE: for debugging environment.ts you need to create 
# a .credentials.tar.enc in /auth-service. Steps involved 
# # 1. remove existing encrypted tar
# # 2. create new temp-unencrypted tar
# # 3. update the encrypted tar file and pepper
# # 4. clean-up the temp un-encrypted tar
# Run the following to execute steps above:
# $ rm .credentials/local/.credentials.tar.enc ; tar zcf .credentials/local/.credentials.tar .credentials/dev ; sh ./src/encrypt/enc.sh e somepassphase .credentials/local/.credentials.tar .credentials/local/.credentials.tar.enc .credentials/local/pepper_iter.txt ; rm .credentials/local/.credentials.tar 

# e=encrypt | d=decrypt
action=$1 
# stored in kubernetes secrets
passpharse=$2 # to be in the kubernetes secrets
# for encryption - path to tar file to be encrypted ie: credentials.tar
# for decryption - path to encrypted file to be decypted ie: credentials.tar.enc
inputFile=$3 # used /deployment to create
# for encryption - path to where the resulting encrypted file should reside ie: credentials.tar.enc
# for decryption - path to where the resulting decrypted file should reside ie: credentials.tar
outputFile=$4 # used /deployment to create
# for encryption - path to where the randomly generated iterate and pepper should be written ie: pepper_iter.txt
# for decryption - path to where the file where the pepper and iterate should reside ie: pepper_iter.txt
pepperIterFile=$5 # to be stored in the docker image

if [ $action = 'e' ] 
then

    # clear out existing pepper file 
    rm $pepperIterFile

    # create pepper and iterate and output to file
    pepperSecret=$(openssl rand -hex 16)
    echo $pepperSecret >> $pepperIterFile
    iter=$(shuf -i 0-100 -n 1)
    echo $iter >> $pepperIterFile

    encrypt the using passphrase, pepper, and iterate 
    openssl enc -aes-256-cbc -base64 -nosalt -e \
            -in $inputFile -out $outputFile \
            -iter $iter \
            -pass pass:$passpharse$pepperSecret

    # openssl enc -aes-256-cbc -base64 -pbkdf2 -e \
    #         -salt -S $pepperSecret
    #         -
    #         -in $inputFile -out $outputFile \
    #         -iter $iter \
    #         -pass pass:$passpharse$pepperSecret
fi

if [ $action = 'd' ] 
then

    # get the pepper and iterate from local file
    pepperSecret=$(sed -n '1p' $pepperIterFile)
    iter=$(sed -n '2p' $pepperIterFile)

    openssl enc -aes-256-cbc -base64 -nosalt -d \
            -in $inputFile -out $outputFile \
            -iter $iter \
            -pass pass:$passpharse$pepperSecret
fi