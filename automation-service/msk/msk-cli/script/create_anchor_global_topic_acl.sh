#!/bin/bash

set -o errexit \
    -o verbose \
    -o nounset \
    -o xtrace

# create topic
PARTITIONS=${PARTITIONS:-3}
REPLICATION_FACTOR=${REPLICATION_FACTOR:-2}
ZOOKEEPER=${ZOOKEEPER:-zookeeper-1:32181}
PARTICIPANT_ID=${PARTICIPANT_ID:-participant-test-id}
DN=${DN:-*.worldwire.io}

# Create TOPICs
CONSUMER_GROUP1="G1"
G1_TOPICS="ww_res ww_req"

for TOPIC in $G1_TOPICS
do
kafka-topics.sh --create --zookeeper $ZOOKEEPER --replication-factor $REPLICATION_FACTOR  --partitions $PARTITIONS --topic $TOPIC
kafka-acls.sh --authorizer-properties zookeeper.connect=$ZOOKEEPER --add \
--allow-principal User:'*' --producer --topic $TOPIC
kafka-acls.sh --authorizer-properties zookeeper.connect=$ZOOKEEPER --add \
--allow-principal User:CN=$DN --consumer --topic $TOPIC --group $CONSUMER_GROUP1
done