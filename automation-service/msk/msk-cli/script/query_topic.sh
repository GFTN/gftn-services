#!/usr/bin/env bash

BROKER=${BROKER:-kafka-1:9091,kafka-2:9092,kafka-3:9093}
PARTICIPANT_ID=${PARTICIPANT_ID:-participant-test-id}

G1_TOPICS="${PARTICIPANT_ID}_res ${PARTICIPANT_ID}_req"
for TOPIC in $G1_TOPICS
do
kafka-console-consumer.sh --bootstrap-server $BROKER --topic $TOPIC --from-beginning
done

G2_TOPICS="${PARTICIPANT_ID}_FEE ${PARTICIPANT_ID}_TRANSACTIONS ${PARTICIPANT_ID}_QUOTES ${PARTICIPANT_ID}_PAYMENT"
for TOPIC in $G2_TOPICS
do
kafka-console-consumer.sh --bootstrap-server $BROKER --topic $TOPIC --from-beginning
done
