#!/usr/bin/env bash
set -eo pipefail

source ./script.source

echo "Stopping the services..."
docker-compose down
