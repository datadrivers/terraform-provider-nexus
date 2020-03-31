#!/bin/bash
set -eo pipefail

echo "Starting Nexus container..."
docker run -d --rm --name nexus -p 127.0.0.1:8081:8081 sonatype/nexus3:3.22.0

function wait_for_nexus {
    echo -n "Waiting for Nexus to be ready "
    i=1
    until wget -t 1 http://127.0.0.1:8081 -O /dev/null -q
    do
        sleep 1
        echo -n .
        if [[ $((i%3)) == 0 ]]; then echo -n ' '; fi
        (( i++ ))
    done
}

wait_for_nexus

echo "Getting admin password..."
NEXUS_ADMIN_PASSWORD=$(docker exec -ti nexus cat /nexus-data/admin.password)

echo "Setting admin password..."
curl -X PUT "http://127.0.0.1:8081/service/rest/beta/security/users/admin/change-password" -H "accept: application/json" -H "Content-Type: text/plain" -d "admin123" -u "admin:${NEXUS_ADMIN_PASSWORD}"

echo "Enabling scripting related features and restarting Nexus. https://issues.sonatype.org/browse/NEXUS-23205"
docker exec nexus /bin/bash -c "echo 'nexus.scripts.allowCreation=true' >> /nexus-data/etc/nexus.properties" && docker restart nexus

wait_for_nexus