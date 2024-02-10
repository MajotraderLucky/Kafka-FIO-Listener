#!/bin/bash

# Setting environment variables
export KAFKA_HOME=/opt/kafka
export PATH=$PATH:$KAFKA_HOME/bin

# Create the topics
kafka-topics.sh --create \
  --zookeeper zookeeper:2181 \
  --replication-factor 1 \
  --partitions 1 \
  --topic my-topic