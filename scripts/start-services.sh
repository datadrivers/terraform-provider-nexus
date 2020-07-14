#!/usr/bin/env bash
set -eo pipefail

echo "Starting Nexus and other required services..."
docker-compose up -d

./wait-for-nexus.sh
