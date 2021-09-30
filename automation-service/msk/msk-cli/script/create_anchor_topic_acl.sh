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

CONSUMER_GROUP2="G2"
G2_TOPICS="${PARTICIPANT_ID}_TRANSACTIONS"

for TOPIC in $G2_TOPICS
do
kafka-topics.sh --create --zookeeper $ZOOKEEPER --replication-factor $REPLICATION_FACTOR  --partitions $PARTITIONS --topic $TOPIC
kafka-acls.sh --authorizer-properties zookeeper.connect=$ZOOKEEPER --add \
--allow-principal User:'*' --producer --topic $TOPIC
kafka-acls.sh --authorizer-properties zookeeper.connect=$ZOOKEEPER --add \
--allow-principal User:CN=$DN --consumer --topic $TOPIC --group $CONSUMER_GROUP2
done