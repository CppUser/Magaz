#!/bin/sh

## Docker workaround: Remove check for KAFKA_ZOOKEEPER_CONNECT parameter
#sed -i '/KAFKA_ZOOKEEPER_CONNECT/d' /etc/confluent/docker/configure
#
## Docker workaround: Remove check for KAFKA_ADVERTISED_LISTENERS parameter
#sed -i '/dub ensure KAFKA_ADVERTISED_LISTENERS/d' /etc/confluent/docker/configure
#
## Docker workaround: Ignore cub zk-ready
#sed -i 's/cub zk-ready/echo ignore zk-ready/' /etc/confluent/docker/ensure

## Add debug output to check for syntax errors
#echo "Debugging /etc/confluent/docker/configure" >&2
#cat /etc/confluent/docker/configure >&2

file_path="/tmp/clusterID/clusterID"
interval=5  # wait interval in seconds

while [ ! -e "$file_path" ] || [ ! -s "$file_path" ]; do
  echo "Waiting for $file_path to be created..."
  sleep $interval
done

# Output the file content and use it for KRaft formatting
cat "$file_path"
echo "kafka-storage format --ignore-formatted -t $(cat "$file_path") -c /etc/kafka/kafka.properties" >> /etc/confluent/docker/ensure
