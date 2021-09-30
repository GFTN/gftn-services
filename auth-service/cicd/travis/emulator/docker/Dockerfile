# USAGE: 
#   This docker container can be used to simulate an fresh ci/cd pipleline
#   Build image: $ sudo docker build -f cicd/travis/emulator/docker/Dockerfile -t cicd-service-travis-emulator --build-arg CACHEBUST=$(date +%s) --build-arg githubUsername=YOUR_GITHUB_USERNAME --build-arg githubPersonalAccessToken=YOUR_PERSONAL_ACCESS_TOKEN .
#   Run image:   $ sudo docker run --env branch=development -it cicd-service-travis-emulator:latest

FROM ubuntu:xenial

WORKDIR /emulator

# pre-installed travis dependencies
RUN apt-get update
# curl
RUN apt-get -y install curl
# nodejs
RUN curl -sL https://deb.nodesource.com/setup_10.x | bash -
RUN apt-get -y install nodejs
# git
RUN apt-get -y install git-core
# docker
# COPY cicd/travis/emulator/deps .
# RUN sh get-docker.sh
# RUN apt-get -y install apt-transport-https ca-certificates gnupg-agent software-properties-common
# RUN curl -fsSL https://download.docker.com/linux/ubuntu/gpg | apt-key add -
# RUN add-apt-repository "deb [arch=amd64] https://download.dockerlinux/$(lsb_release -cs) stable"
# RUN apt-get install docker-ce docker-ce-cli containerd.io

# check dep versions
RUN echo "node version" ; node -v
RUN echo "npm version" ; npm -v
RUN echo git --version
# RUN docker -v

# ARG githubUsername=empty
# ARG githubPersonalAccessToken=empty

# # clone repo in image upon build so run is faster
# RUN git clone https://${githubUsername}:${githubPersonalAccessToken}@github.com/GFTN/gftn-services.git

# below command will run without cache using CACHEBUST
ARG CACHEBUST=1

COPY google-cloud-sdk ../google-cloud-sdk
# RUN export CLOUDSDK_CORE_DISABLE_PROMPTS=1; ../google-cloud-sdk

# clone the gftn-services repo
# RUN git clone https://PERSONAL_ACCESS_TOKEN@github.com/GFTN/gftn-services.git

RUN pwd ; ls


COPY install_dependencies.sh .


# run operations in order cicd pipeline will run
# ENTRYPOINT [ "sh", "entry.sh" ]

# run image with access to cmd line:
ENTRYPOINT [ "sh" ]