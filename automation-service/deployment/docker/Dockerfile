FROM golang:alpine as builder
ARG BUILD_VERSION
LABEL maintainer=noreply@your.domain

RUN echo $BUILD_VERSION

# Vars
ENV PARENT_REPO=github.com/GFTN/gftn-services
ENV SERVICE_NAME=automation-service

RUN apk add --update --no-cache git \
    && apk add --update --no-cache make\
    && apk add --no-cache openssh \
    && apk add --no-cache build-base

# Setup Runtime Dir Structure
RUN mkdir -p /go/src/$PARENT_REPO/$SERVICE_NAME \
        && mkdir -p /go/src/$PARENT_REPO/$SERVICE_NAME/vendor

# Copy code
COPY . /go/src/$PARENT_REPO/$SERVICE_NAME
COPY deployment/docker/vendor/gftn-services /go/src/$PARENT_REPO/

RUN cd /go/src/$PARENT_REPO/$SERVICE_NAME \
        && pwd \
        && make test \
        && make install \			
        && make build \	
        && mv $SERVICE_NAME /tmp/$SERVICE_NAME \
        && mv ./k8s /tmp/k8s \
        && mv ./msk /tmp/msk \
        && mv error-codes/ /tmp/error-codes \
        && rm -Rf /go/src/*

FROM alpine:latest
ARG BUILD_VERSION
RUN apk --no-cache add ca-certificates
COPY --from=builder /tmp/automation-service .
RUN pwd && chmod a+x automation-service

RUN apk update
RUN apk add --no-cache bash \
    && apk add --no-cache curl \
    && apk add --no-cache jq \
    && apk add --no-cache python3 \
    && apk add --no-cache py-pip \
    && apk add --no-cache sed \
    && set -ex && apk --no-cache add sudo

RUN curl -Lo kubectl https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/linux/amd64/kubectl \
    && chmod +x ./kubectl \
    && sudo mv ./kubectl /usr/local/bin/kubectl

RUN apk del sudo

RUN curl -L https://git.io/getLatestIstio | ISTIO_VERSION=1.3.1 sh -
ENV PATH=/istio-1.3.1/bin:$PATH
RUN rm -rf /istio-1.3.1/samples

# Setup automate deployment script and files
RUN mkdir -p /var/k8s
RUN mkdir -p /var/msk
RUN mkdir -p /var/files
RUN mkdir -p /var/logs
COPY --from=builder /tmp/k8s /var/k8s
COPY --from=builder /tmp/msk /var/msk
RUN touch /var/logs/logs.txt
RUN mkdir -p /opt/error-codes/
COPY --from=builder /tmp/error-codes /opt/error-codes/

# TLS & logs
RUN mkdir -p /etc/tls 
ENV ENV_BUILD_VERSION=$BUILD_VERSION

# Create a group and user
RUN addgroup -S wwgroup && adduser -S wwuser -G wwgroup
# Change owner for log file
RUN chown -R wwuser /var/k8s
RUN chown -R wwuser /var/msk
RUN chown -R wwuser /var/files
RUN chown wwuser /var/logs/logs.txt
# Tell docker that all future commands should run as the appuser wwuser
USER wwuser

RUN pip3 install awscli --upgrade --user
RUN chmod +x /home/wwuser/.local/bin/aws
ENV PATH=/home/wwuser/.local/bin:$PATH

EXPOSE 5566
CMD ["./automation-service"]
