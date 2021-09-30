# if your using Ubuntu Desktop - installs kops, kubectl, aws cli and sets aws credentials

checkAndInstallKops()
{
    dpkg -s kops
    if [ $? -ne 0 ]; then
        curl -Lo kops https://github.com/kubernetes/kops/releases/download/$(curl -s https://api.github.com/repos/kubernetes/kops/releases/latest | grep tag_name | cut -d '"' -f 4)/kops-linux-amd64
        chmod +x ./kops
        mv ./kops /usr/local/bin/
    fi

}

checkAndInstallKubectl()
{
    dpkg -s kubectl
    if [ $? -ne 0 ]; then
        curl -Lo kubectl https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/linux/amd64/kubectl
        chmod +x ./kubectl
        mv ./kubectl /usr/local/bin/kubectl
    fi

}

# needed to install awscli
checkAndInstallPip3()
{
    dpkg -s pip3
    if [ $? -ne 0 ]; then
        apt-get -y install python3-pip
    fi

}

checkAndInstallAwscli()
{
    dpkg -s aws
    if [ $? -ne 0 ]; then
        pip3 install awscli --upgrade
    fi

}

checkAndInstallGcloud()
{
    dpkg -s gcloud
    if [ $? -ne 0 ]; then        

        # Create environment variable for correct distribution
        export CLOUD_SDK_REPO="cloud-sdk-cosmic"

        rm -rf /etc/apt/sources.list.d/google-cloud-sdk.list

        # Add the Cloud SDK distribution URI as a package source
        echo "deb http://packages.cloud.google.com/apt $CLOUD_SDK_REPO main" | tee -a /etc/apt/sources.list.d/google-cloud-sdk.list

        # Import the Google Cloud Platform public key
        curl https://packages.cloud.google.com/apt/doc/apt-key.gpg | apt-key add -

        # Update the package list and install the Cloud SDK
        yes | apt-get update
        yes | apt-get install google-cloud-sdk

        # on a development ubuntu laptop: 
        #apt-get install google-cloud-sdk

    fi

}

configAws()
{
    aws configure set aws_access_key_id YOUR_AWS_KEY_ID_HERE
    aws configure set aws_secret_access_key YOUR_AWS_SECRET_KEY_HERE
}

configgcloud()
{
    gcloud beta container clusters get-credentials dev-ww-cluster --region us-central1 --project next-gftn
}

checkAndInstallKops
checkAndInstallKubectl
checkAndInstallPip3
checkAndInstallAwscli
checkAndInstallGcloud
configgcloud
# configAws    