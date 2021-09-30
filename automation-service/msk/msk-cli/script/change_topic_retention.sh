#!/bin/bash

set -o errexit \
    -o verbose \
    -o nounset \
    -o xtrace

# create topic
ZOOKEEPER=${ZOOKEEPER:-zookeeper-1:32181}
TOPIC=${TOPIC:-test}
RETENTIONTIME=${RETENTIONTIME:-1000}

kafka-topics.sh --zookeeper $ZOOKEEPER --alter --topic $TOPIC --config retention.ms=$RETENTIONTIME

