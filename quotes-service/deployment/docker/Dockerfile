FROM golang:alpine as builder
LABEL maintainer=your.user@your.domain
ARG BUILD_VERSION

RUN echo $BUILD_VERSION

# Vars
ENV ENV_BUILD_VERSION=$BUILD_VERSION
ENV PARENT_REPO=github.com/GFTN
ENV PARENT_PACKAGE=$PARENT_REPO/gftn-services
ENV SERVICE_NAME=quotes-service

# install necessary packages for kafka
RUN apk update \
    && apk add --update --no-cache pkgconfig \
    && apk add --update --no-cache libxml2\
    && apk add --update --no-cache libxml2-dev\
    && apk add --update --no-cache 'librdkafka>=1.2.1-r0' --repository=http://dl-cdn.alpinelinux.org/alpine/edge/community \
    && apk add --update --no-cache 'librdkafka-dev>=1.2.1-r0' --repository=http://dl-cdn.alpinelinux.org/alpine/edge/community \
    && apk add --update --no-cache gcc \
    && apk add --update --no-cache libc-dev

RUN apk add --update --no-cache git \
    && apk add --update --no-cache make\
    && apk add --no-cache openssh \
    && apk add --no-cache build-base

# Setup Runtime Dir Structure
RUN mkdir -p /go/src/$PARENT_PACKAGE/$SERVICE_NAME \
    && mkdir -p /go/src/$PARENT_PACKAGE/$SERVICE_NAME/vendor

# Copy code

COPY deployment/docker/vendor/gftn-services /go/src/$PARENT_PACKAGE/
COPY . /go/src/$PARENT_PACKAGE/$SERVICE_NAME
RUN cd /go/src/$PARENT_PACKAGE/$SERVICE_NAME \
        && pwd \
        && make test \
        && make install \			
        && make build \	
        && mv $SERVICE_NAME /tmp/$SERVICE_NAME \
        && rm -Rf /go/src/*

# Wrapper image	
FROM alpine:latest	
ARG BUILD_VERSION
RUN apk --no-cache add ca-certificates	
COPY --from=builder /tmp/quotes-service .	
RUN pwd && chmod a+x quotes-service	

RUN apk update \
     && apk add --update --no-cache pkgconfig\
     && apk add --update --no-cache libxml2\
     && apk add --update --no-cache libxml2-dev\
     && apk add --update --no-cache 'librdkafka>=1.2.1-r0' --repository=http://dl-cdn.alpinelinux.org/alpine/edge/community \
     && apk add --update --no-cache 'librdkafka-dev>=1.2.1-r0' --repository=http://dl-cdn.alpinelinux.org/alpine/edge/community

RUN echo $BUILD_VERSION	
ENV ENV_BUILD_VERSION=$BUILD_VERSION

RUN mkdir -p /opt/error-codes/
COPY error-codes/ /opt/error-codes/
RUN mkdir -p /opt/logs/
RUN touch /opt/logs/logs.txt

# TLS & logs
RUN mkdir -p /etc/tls 

# Create a group and user
RUN addgroup -S wwgroup && adduser -S wwuser -G wwgroup
# Change owner for log file
RUN chown wwuser /opt/logs/logs.txt
# Tell docker that all future commands should run as the appuser wwuser
USER wwuser

EXPOSE 11234
CMD ["./quotes-service"]
