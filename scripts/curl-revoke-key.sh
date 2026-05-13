#!/usr/bin/env sh
set -eu
: "${API_KEY_ID:?Set API_KEY_ID to revoke}"
curl -sS -X POST "http://localhost:8080/admin/api-keys/${API_KEY_ID}/revoke" \
  -H "X-Admin-Token: ${APIKEY_ADMIN_TOKEN:-local-admin-token}"
printf '\n'
