# USAGE: 
#   full install = $ sudo bash install_dependencies.sh --gcloud --node

isTravis="false"

normal=$(tput sgr0)         # back to normal
bold=$(tput bold)           # bold mode
dim=$(tput dim)             # dim (half-bright) mode
underline_on=$(tput smul)   # underline mode
underline_off=$(tput rmul)  # underline mode off
standout_on=$(tput smso)    # standout (bold) mode on 
standout_off=$(tput rmso)   # standout (bold) mode off 
black=$(tput setaf 0)       # black     COLOR_BLACK     0,0,0
red=$(tput setaf 1)         # red       COLOR_RED       1,0,0
green=$(tput setaf 2)       # green     COLOR_GREEN     0,1,0
yellow=$(tput setaf 3)      # yellow    COLOR_YELLOW    1,1,0
blue=$(tput setaf 4)        # blue      COLOR_BLUE      0,0,1
magenta=$(tput setaf 5)     # magenta   COLOR_MAGENTA   1,0,1
cyan=$(tput setaf 6)        # cyan      COLOR_CYAN      0,1,1
white=$(tput setaf 7)       # white     COLOR_WHITE     1,1,1

echo ""
echo ""
echo "${green}${bold}Step 1${normal}: installing global dependencies for auth-service & portal"
echo ""
echo ""
echo "${yellow}You may have to enter your ${bold}password ${normal}${yellow}(if prompted)${normal}:"

installKubectl()
{
    dpkg -s kubectl
    if [ $? -ne 0 ]; then
        curl -Lo kubectl https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/linux/amd64/kubectl
        chmod +x ./kubectl
        mv ./kubectl /usr/local/bin/kubectl
    fi

}

# needed to install awscli
installPip3()
{
    dpkg -s pip3
    if [ $? -ne 0 ]; then
       apt-get -y install python3-pip
    fi

}

installAws()
{
    dpkg -s aws
    if [ $? -ne 0 ]; then
       pip3 install --upgrade --user awscli
    fi

}

installGcloud()
{
    # UNABLE TO GET TO WORK FROM SCRIPT ON TRAVIS 

    # # https://gist.github.com/mjackson/5887963e7d8b8fb0615416c510ae8857
    gcloud version || true
    
    if [ ! -d "$HOME/google-cloud-sdk/bin" ]; then rm -rf $HOME/google-cloud-sdk; export CLOUDSDK_CORE_DISABLE_PROMPTS=1; curl https://sdk.cloud.google.com | bash; fi
    
    # add gcloud to path 
    # NOTE: the following errors out on Docker image because
    # usually source or sudo is not available in running container
    if [ $isTravis = "true" ]
    then
        echo "========= set travis path =========="
        source /home/travis/google-cloud-sdk/path.bash.inc
    else
        echo "========= set local path =========="
        source ${HOME}/google-cloud-sdk/path.bash.inc
        # export PATH=$PATH:${HOME}/google-cloud-sdk/bin
    fi
    
    gcloud version

}

installNode()
{
    dpkg -s npm
    if [ $? -ne 0 ]; then
        curl -sL https://deb.nodesource.com/setup_10.x | sudo -E bash -
        sudo apt-get install nodejs
        node -v
        npm -v 
    fi

}

installDockerCompose()
{
    sudo curl -L "https://github.com/docker/compose/releases/download/1.24.1/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
    sudo chmod +x /usr/local/bin/docker-compose
    sudo ln -s /usr/local/bin/docker-compose /usr/bin/docker-compose
    docker-compose --version
}

show_help()
{ 
cat << EOF

Deploy auth-service in nodejs:

--travis OPTIONAL
   Denote if this is install is for travis

--docker-compose OPTIONAL
   Install docker compose

--gcloud OPTIONAL
   Install google cloud sdk

--node OPTIONAL
   Install nodejs

--aws OPTIONAL
   Install aws cli

--kubectl OPTIONAL
   Install kubernetes cli

--pip3 OPTIONAL
   Install python 3

EOF

}

# flags usage - https://archive.is/5jGpl#selection-709.0-715.1 or https://archive.is/TRzn4
while :; do
    case $1 in
          -h|-\?|--help)   # Call a "show_help" function to display a synopsis, then exit.
              show_help
              exit
              ;;
          --travis) 
             isTravis="true"
              ;;
          --docker-compose) 
             installDockerCompose
              ;;
          --gcloud) 
             installGcloud
              ;;
          --node) 
             installNode
              ;;
          --aws) 
             installAws
              ;;
          --kubectl) 
             installKubectl
              ;;
          --pip3) 
             installPip3
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


echo ""
echo ""
npm i -g tsoa@2.4.0 ts-node@8.2.0 @angular/cli@7.3.8 firebase-tools@7.3.0 firebase-bolt@0.8.4 typescript@3.4.5 mocha@6.1.4 chai@4.2.0 ts-node@8.4.1
echo ""
echo ""
echo "${green}${bold}Step 2${normal}: installing world wire auth-service project dependencies";
echo ""
echo ""
echo "${bold}${red}IMPORTANT: ${normal}Attention all IBM World Wire ${green}Superheros! ${normal}"
echo "${bold}${red}PREREQUISITE: ${normal}Manual installation required first before installing ${bold}${blue}auth-service ${normal}dependencies:"
echo "${bold} * vscode ${normal}      see https://code.visualstudio.com/Download"
echo "${bold} * nodejs ${normal}      see https://nodejs.org/en/download/"
echo "${bold} * gcloud sdk ${normal}  see https://cloud.google.com/sdk/install"
echo "${yellow} You can try to install the above packages on linux or mac by running '\$ ${magenta}bash install_dependencies.sh --gcloud --node ${yellow}', if this does not work click the links above one by one and install manually.${normal}"
echo ""
echo ""
npm i ; # node-types
cd authentication ; npm i ; cd .. 
cd project-provisioning ; npm i ; cd .. 
cd deployment ; npm i ;  cd ..
cd secret-mgmt ; npm i ;  cd ..
bash extract-debug-creds.sh
# cd test-example ; npm i ;  cd ..
echo ""
echo ""
echo "${green}Deps installed. Now you're cooking!"${normal};
echo ""
echo ""