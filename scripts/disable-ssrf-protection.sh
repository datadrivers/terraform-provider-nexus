#!/usr/bin/env bash
set -eo pipefail

source .env

IP=$(./detect-docker-env-ip.sh)

# Since 3.93.0 Nexus ships SSRF protection that resolves and validates every
# outbound URL on save. This breaks the S3 blobstore tests (MinIO resolves to
# a private IP) and every proxy repository test using a placeholder remote
# URL. Disable it for the test instance. Best effort: older Nexus versions do
# not have this endpoint.
if curl -sf -u "admin:admin123" \
  -X PUT "http://${IP}:${NEXUS_PORT}/service/rest/v1/security/ssrf-protection" \
  -H "Content-Type: application/json" \
  -d '{"enabled": false, "allowedIPs": [], "allowedDomains": []}' \
  -o /dev/null; then
  echo "Disabled Nexus SSRF protection."
else
  >&2 echo "Could not update Nexus SSRF protection settings (endpoint may not exist on this version); skipping."
fi
