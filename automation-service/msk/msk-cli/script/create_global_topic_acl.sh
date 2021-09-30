#!/bin/bash

set -o errexit \
    -o verbose \
    -o nounset \
    -o xtrace

ZOOKEEPER=${ZOOKEEPER:-zookeeper-1:32181}
DN=${DN:-*.worldwire.io}

kafka-acls.sh --authorizer-properties zookeeper.connect=$ZOOKEEPER --add \
--allow-principal User:CN=$DN --producer --topic '*'