#!/bin/bash

set -o errexit \
    -o verbose \
    -o nounset \
    -o xtrace

# create topic
ZOOKEEPER=${ZOOKEEPER:-zookeeper-1:32181}
PARTICIPANT_ID=${PARTICIPANT_ID:-participant-test-id}
DN=${DN:-*.worldwire.io}

# Create TOPICs
CONSUMER_GROUP1="G1"
G1_TOPICS="${PARTICIPANT_ID}_res ${PARTICIPANT_ID}_req"

for TOPIC in $G1_TOPICS
do
kafka-topics.sh --delete --zookeeper $ZOOKEEPER --topic $TOPIC
kafka-acls.sh --authorizer-properties zookeeper.connect=$ZOOKEEPER --remove \
--allow-principal User:'*' --producer --topic $TOPIC
kafka-acls.sh --authorizer-properties zookeeper.connect=$ZOOKEEPER --remove \
--allow-principal User:CN=$DN --consumer --topic $TOPIC --group $CONSUMER_GROUP1
done

CONSUMER_GROUP2="G2"
G2_TOPICS="${PARTICIPANT_ID}_FEE ${PARTICIPANT_ID}_TRANSACTIONS ${PARTICIPANT_ID}_QUOTES ${PARTICIPANT_ID}_PAYMENT"

for TOPIC in $G2_TOPICS
do
kafka-topics.sh --delete --zookeeper $ZOOKEEPER --topic $TOPIC
kafka-acls.sh --authorizer-properties zookeeper.connect=$ZOOKEEPER --remove \
--allow-principal User:'*' --producer --topic $TOPIC
kafka-acls.sh --authorizer-properties zookeeper.connect=$ZOOKEEPER --remove \
--allow-principal User:CN=$DN --consumer --topic $TOPIC --group $CONSUMER_GROUP2
done