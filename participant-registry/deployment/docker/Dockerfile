FROM golang:alpine as builder
ARG BUILD_VERSION

RUN echo $BUILD_VERSION

# Vars

ENV PARENT_REPO=github.com/GFTN
ENV PARENT_PACKAGE=$PARENT_REPO/gftn-services
ENV SERVICE_NAME=participant-registry

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

FROM alpine:latest	
ARG BUILD_VERSION
RUN apk --no-cache add ca-certificates	
COPY --from=builder /tmp/participant-registry .	
RUN pwd && chmod a+x participant-registry

ENV ENV_BUILD_VERSION=$BUILD_VERSION
# TLS & logs
RUN mkdir -p /etc/tls 

# Error codes and log file
RUN mkdir -p /opt/error-codes/
COPY error-codes/ /opt/error-codes/
RUN mkdir -p /opt/logs/
RUN touch /opt/logs/logs.txt

# Create a group and user
RUN addgroup -S wwgroup && adduser -S wwuser -G wwgroup
# Change owner for log file
RUN chown wwuser /opt/logs/logs.txt
# Tell docker that all future commands should run as the appuser wwuser
USER wwuser

EXPOSE 8080
CMD ["./participant-registry"]
