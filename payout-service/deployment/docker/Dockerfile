FROM golang:alpine

ARG BUILD_VERSION

RUN echo $BUILD_VERSION

# Vars
ENV ENV_BUILD_VERSION=$BUILD_VERSION
ENV PARENT_REPO=github.com/GFTN
ENV PARENT_PACKAGE=$PARENT_REPO/gftn-services
ENV SERVICE_NAME=payout-service
# ENV SERVICE_ERROR_CODES_FILE=/opt/$SERVICE_NAME/configuration/error-codes/payoutErrorCodes.toml

RUN apk add --update --no-cache git \
    && apk add --update --no-cache make\
    && apk add --no-cache openssh

# Setup Runtime Dir Structure
RUN mkdir -p /packages \
    && mkdir -p /go/src/$PARENT_PACKAGE/$SERVICE_NAME \
    && mkdir -p /go/src/$PARENT_PACKAGE/$SERVICE_NAME/vendor \
    && mkdir -p /opt/$SERVICE_NAME/configuration/error-codes/

# Error codes and log file
RUN mkdir -p /opt/error-codes/
COPY error-codes/ /opt/error-codes/
RUN mkdir -p /opt/logs/
RUN touch /opt/logs/logs.txt

# install packages for neo4j
RUN apk add --update bash && rm -rf /var/cache/apk/*
RUN apk add --update cmake build-base gcc openssl-dev libressl-dev
ENV OPENSSL_ROOT_DIR=/usr/lib
RUN apk add --update pkgconfig
RUN git clone -b 1.7 https://github.com/neo4j-drivers/seabolt.git /seabolt
RUN /bin/bash /seabolt/make_debug.sh

# export env variables for neo4j
ENV PKG_CONFIG_PATH=/seabolt/build/dist/share/pkgconfig
ENV DYLD_LIBRARY_PATH=/seabolt/build/dist/lib64
ENV LD_LIBRARY_PATH=/seabolt/build/dist/lib64

# Copy code

COPY deployment/docker/vendor/gftn-services /go/src/$PARENT_PACKAGE/
COPY . /go/src/$PARENT_PACKAGE/$SERVICE_NAME
COPY error-codes/ /opt/$SERVICE_NAME/configuration/error-codes/
RUN cd /go/src/$PARENT_PACKAGE/$SERVICE_NAME \
        && pwd \
        && make test \
        && make build \
        && make install \
        && rm -Rf /go/src/*

# TLS & logs
RUN mkdir -p /etc/tls 

# Volumes
#VOLUME /opt/$SERVICE_NAME/configuration
#VOLUME /opt/$SERVICE_NAME/logs

# Create a group and user
RUN addgroup -S wwgroup && adduser -S wwuser -G wwgroup
# Change owner for log file
RUN chown wwuser /opt/logs/logs.txt
# Tell docker that all future commands should run as the appuser wwuser
USER wwuser

EXPOSE 8080
CMD ["payout-service"]