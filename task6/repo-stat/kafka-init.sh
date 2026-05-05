#!/bin/bash
echo "Initializing Kafka..."

echo "Creating Kafka topics..."
kafka-topics --bootstrap-server kafka:9092 --create --topic repo-requests --partitions 1 --replication-factor 1 || true
kafka-topics --bootstrap-server kafka:9092 --create --topic repo-responses --partitions 1 --replication-factor 1 || true
kafka-topics --bootstrap-server kafka:9092 --create --topic subscription-updates --partitions 1 --replication-factor 1 || true

echo "Resetting consumer groups to avoid offset issues..."
sleep 5

kafka-consumer-groups --bootstrap-server kafka:9092 --delete --group collector-task-group 2>/dev/null || true
kafka-consumer-groups --bootstrap-server kafka:9092 --delete --group processor-response-group 2>/dev/null || true
kafka-consumer-groups --bootstrap-server kafka:9092 --delete --group processor-subscription-group 2>/dev/null || true

echo "Kafka initialization completed."
exit 0
