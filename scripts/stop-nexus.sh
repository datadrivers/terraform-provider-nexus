#!/bin/bash
set -eo pipefail

echo "Stopping running Nexus containers..."

for container_id in $(docker ps -f label=terraform-provider-nexus=true \
        | tail -n +2 \
        | awk '{ print $1 }') ; do
    echo "Stopping ${container_id}"
    docker stop "${container_id}" > /dev/null
done

echo "Nexus containers stopped"
