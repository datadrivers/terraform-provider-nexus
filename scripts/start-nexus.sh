#!/bin/bash
set -eo pipefail

export NEXUS_VERSION=3.23.0
LOCAL_PORT=8082

echo "Starting Nexus v${NEXUS_VERSION} container..."
docker run \
    -d \
    --rm \
    --name nexus \
    -l terraform-provider-nexus=true \
    -p 127.0.0.1:${LOCAL_PORT}:8081 \
    -v "${PWD}/nexus.properties:/nexus-data/etc/nexus.properties" \
    "sonatype/nexus3:${NEXUS_VERSION}"

function wait_for_nexus {
    echo -n "Waiting for Nexus to be ready "
    i=1
    until wget -t 1 http://127.0.0.1:${LOCAL_PORT} -O /dev/null -q
    do
        sleep 1
        echo -n .
        if [[ $((i%5)) == 0 ]]; then echo -n ' '; fi
        (( i++ ))
    done
    echo ""
}

wait_for_nexus

echo "Getting admin password..."
NEXUS_ADMIN_PASSWORD=$(docker exec -ti nexus cat /nexus-data/admin.password)

echo "Setting admin password..."
curl -X PUT "http://127.0.0.1:${LOCAL_PORT}/service/rest/beta/security/users/admin/change-password" -H "accept: application/json" -H "Content-Type: text/plain" -d "admin123" -u "admin:${NEXUS_ADMIN_PASSWORD}"

echo "Nexus is started."
