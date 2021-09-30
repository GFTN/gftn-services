FROM centos:7
LABEL maintainer=your.user@your.domain
ARG BUILD_VERSION

RUN echo $BUILD_VERSION

# Vars
ENV ENV_BUILD_VERSION=$BUILD_VERSION
ENV PARENT_REPO=github.com/GFTN
ENV PARENT_PACKAGE=$PARENT_REPO/gftn-services
ENV SERVICE_NAME=crypto-service
ENV SERVICE_PATH=/$SERVICE_NAME
ENV ENVIRONMENT=dev
ENV CERT_FILE=ED25519_client_$ENVIRONMENT

# update package
RUN yum -y update
RUN yum install -y wget unzip make gcc libtool-ltdl-devel

# Download cache cleanup
RUN yum clean all

# Setting up Golang
RUN wget https://dl.google.com/go/go1.13.1.linux-amd64.tar.gz
RUN tar -xzf go1.13.1.linux-amd64.tar.gz
RUN mv go /usr/local
RUN mkdir /go
ENV GOPATH=/go
ENV PATH=$PATH:/usr/local/go/bin:$GOPATH/bin


# Setup Runtime Dir Structure
# Setup directory for hsm & vault key files
RUN mkdir -p $SERVICE_PATH/hsm/driver
COPY deployment/docker/$ENVIRONMENT/$CERT_FILE.zip $SERVICE_PATH/hsm/driver
RUN unzip $SERVICE_PATH/hsm/driver/$CERT_FILE.zip -d $SERVICE_PATH/hsm/driver


# Setup necessary environment variables for HSM client
ENV CERT_PATH=$SERVICE_PATH/hsm/driver/$CERT_FILE
ENV PKCS11_LIB=$CERT_PATH/libs/64/libCryptoki2.so
ENV ChrystokiConfigurationPath=$CERT_PATH
RUN cd $CERT_PATH && source ./setenv

# Copy code
RUN mkdir -p /go/src/$PARENT_PACKAGE/$SERVICE_NAME \
    && mkdir -p /go/src/$PARENT_PACKAGE/$SERVICE_NAME/vendor
COPY deployment/docker/vendor/gftn-services /go/src/$PARENT_PACKAGE/
COPY . /go/src/$PARENT_PACKAGE/$SERVICE_NAME
RUN cd /go/src/$PARENT_PACKAGE/$SERVICE_NAME \
        && pwd \
        && make test \
        && make install \			
        && make build \
        && rm -Rf /go/src/*

# TLS & logs
RUN mkdir -p /etc/tls

# Error codes and log file	
RUN mkdir -p /opt/error-codes/	
COPY error-codes/ /opt/error-codes/	
RUN mkdir -p /opt/logs/	
RUN touch /opt/logs/logs.txt	

# Create a group and user
RUN groupadd -r wwgroup && useradd -r -g wwgroup wwuser
# Change owner for log file
RUN chown wwuser /opt/logs/logs.txt
# Tell docker that all future commands should run as the appuser wwuser
USER wwuser

EXPOSE 9999
CMD ["crypto-service"]
