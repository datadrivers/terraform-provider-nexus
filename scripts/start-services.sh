#!/usr/bin/env bash
set -eo pipefail

source ./script.source

echo "Starting Nexus ${NEXUS_TYPE} and other required services..."
docker-compose up -d

./wait-for-nexus.sh
