#!/usr/bin/env bash
set -eo pipefail

export $(cat .env | xargs)

until wget -t 1 http://127.0.0.1:${NEXUS_PORT} -O /dev/null -q; do
    >&2 echo "Waiting for Nexus..."
    sleep 5
done
echo "Nexus is started."
