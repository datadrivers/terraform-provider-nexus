#!/bin/bash
set -eo pipefail

export NEXUS_VERSION=3.22.1

echo "Starting Nexus v${NEXUS_VERSION} container..."
docker run \
    -d \
    --rm \
    --name nexus \
    -l terraform-provider-nexus=true \
    -p 127.0.0.1:8081:8081 \
    -v "${PWD}/nexus.properties:/nexus-data/etc/nexus.properties" \
    "sonatype/nexus3:${NEXUS_VERSION}"

function wait_for_nexus {
    echo -n "Waiting for Nexus to be ready "
    i=1
    until wget -t 1 http://127.0.0.1:8081 -O /dev/null -q
    do
        sleep 1
        echo -n .
        if [[ $((i%5)) == 0 ]]; then echo -n ' '; fi
        (( i++ ))
    done
    echo ""
}

wait_for_nexus

echo "Nexus is started."
