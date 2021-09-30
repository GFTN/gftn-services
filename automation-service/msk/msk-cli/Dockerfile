FROM alpine:latest

RUN apk update

ENV KAFKA_VERSION 2.3.1
ENV SCALA_VERSION 2.12

LABEL name="kafka" version=${KAFKA_VERSION}

RUN apk add --no-cache openjdk8-jre bash docker coreutils su-exec
RUN apk add --no-cache -t .build-deps curl ca-certificates jq \
  && mkdir -p /opt \
  && mirror=$(curl --stderr /dev/null https://www.apache.org/dyn/closer.cgi\?as_json\=1 | jq -r '.preferred') \
  && curl -sSL "${mirror}kafka/${KAFKA_VERSION}/kafka_${SCALA_VERSION}-${KAFKA_VERSION}.tgz" \
  | tar -xzf - -C /opt \
  && mv /opt/kafka_${SCALA_VERSION}-${KAFKA_VERSION} /opt/kafka \
  && adduser -DH -s /sbin/nologin kafka \
  && chown -R kafka: /opt/kafka \
  && rm -rf /tmp/* \
  && apk del --purge .build-deps

ENV PATH /sbin:/opt/kafka/bin/:$PATH

WORKDIR /wkdir

COPY ./script /wkdir/

# Create a group and user
RUN addgroup -S wwgroup && adduser -S wwuser -G wwgroup
RUN chown -R wwuser /wkdir
# Tell docker that all future commands should run as the appuser wwuser
USER wwuser

CMD sh