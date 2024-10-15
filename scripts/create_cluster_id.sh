#!/bin/bash
set -x  # Enable debugging for troubleshooting

# Define the absolute path
file_path="/tmp/clusterID/clusterID"

# Ensure the /tmp/clusterID directory exists
mkdir -p /tmp/clusterID
chmod 777 /tmp/clusterID

# Create cluster ID if it doesn't exist
if [ ! -f "$file_path" ]; then
    /bin/kafka-storage random-uuid > "$file_path"
    echo "Cluster ID has been created at $file_path..."
else
    echo "Cluster ID already exists at $file_path."
fi

chmod 644 "$file_path"