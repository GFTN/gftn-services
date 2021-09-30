FROM golang:alpine as builder
ARG BUILD_VERSION

RUN echo $BUILD_VERSION

# Vars

ENV PARENT_REPO=github.com/GFTN
ENV PARENT_PACKAGE=$PARENT_REPO/gftn-services
ENV SERVICE_NAME=anchor-service

RUN apk update \
     && apk add --update --no-cache pkgconfig\
     && apk add --update --no-cache libxml2\
     && apk add --update --no-cache libxml2-dev\
     && apk add --update --no-cache 'librdkafka>=1.2.1-r0' --repository=http://dl-cdn.alpinelinux.org/alpine/edge/community \
     && apk add --update --no-cache 'librdkafka-dev>=1.2.1-r0' --repository=http://dl-cdn.alpinelinux.org/alpine/edge/community \
     && apk add --update --no-cache gcc\
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
COPY --from=builder /tmp/anchor-service .	
RUN pwd && chmod a+x anchor-service	

# Error codes and log file
RUN mkdir -p /opt/error-codes/
COPY error-codes/ /opt/error-codes/
RUN mkdir -p /opt/logs/
RUN touch /opt/logs/logs.txt
RUN mkdir -p /var/private/ssl/
RUN mkdir -p /opt/configuration/
COPY deployment/docker/vendor/gftn-services/utility/payment/xsd/ /opt/configuration/
ENV SERVICE_FILE=/opt/configuration/errors.json

RUN apk update \
     && apk add --update --no-cache pkgconfig\
     && apk add --update --no-cache libxml2\
     && apk add --update --no-cache libxml2-dev\
     && apk add --update --no-cache 'librdkafka>=1.2.1-r0' --repository=http://dl-cdn.alpinelinux.org/alpine/edge/community \
     && apk add --update --no-cache 'librdkafka-dev>=1.2.1-r0' --repository=http://dl-cdn.alpinelinux.org/alpine/edge/community

ENV ENV_BUILD_VERSION=$BUILD_VERSION
# TLS & logs
RUN mkdir -p /etc/tls 

# Create a group and user
RUN addgroup -S wwgroup && adduser -S wwuser -G wwgroup
# Change owner for log file
RUN chown wwuser /opt/logs/logs.txt
# Tell docker that all future commands should run as the appuser wwuser
USER wwuser

EXPOSE 8080
CMD ["./anchor-service"]
