#!/usr/bin/env bash
set -eo pipefail

echo "Stopping the services..."
docker-compose down
