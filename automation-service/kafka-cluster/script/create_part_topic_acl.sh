#!/bin/bash

set -o errexit \
    -o verbose \
    -o nounset \
    -o xtrace

# create topic
PARTITIONS=${PARTITIONS:-3}
REPLICATION_FACTOR=${REPLICATION_FACTOR:-2}
ZOOKEEPER=${ZOOKEEPER:-zookeeper-1:32181}
CONSUMER_GROUP=${CONSUMER_GROUP:-1}
PARTICIPANT_ID=${PARTICIPANT_ID:-participant-test-id}
CN=${CN:-$PARTICIPANT_ID}

# Create TOPICs
TOPICS="${PARTICIPANT_ID}_res ${PARTICIPANT_ID}_req"
for TOPIC in $TOPICS
do
kafka-topics --create --zookeeper $ZOOKEEPER --replication-factor $REPLICATION_FACTOR  --partitions $PARTITIONS --topic $TOPIC
kafka-acls --authorizer-properties zookeeper.connect=$ZOOKEEPER --add \
--allow-principal User:'*' --producer --topic $TOPIC
kafka-acls --authorizer-properties zookeeper.connect=$ZOOKEEPER --add \
--allow-principal User:CN=$CN --consumer --topic $TOPIC --group $CONSUMER_GROUP
done
