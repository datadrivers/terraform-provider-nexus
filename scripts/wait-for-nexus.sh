#!/usr/bin/env bash
set -eo pipefail

source .env

IP=$(./detect-docker-env-ip.sh)

if command -v wget; then
  checkcmd="wget -t 1 http://${IP}:${NEXUS_PORT} -O /dev/null -q"
else
  checkcmd="curl -I -s --connect-timeout 1 http://${IP}:${NEXUS_PORT} -o /dev/null"
fi

until ${checkcmd}; do
    >&2 echo "Waiting for Nexus..."
    sleep 5
done
echo "Nexus is started."
echo "http://${IP}:${NEXUS_PORT}"
