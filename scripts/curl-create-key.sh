#!/usr/bin/env sh
set -eu
curl -sS -X POST http://localhost:8080/admin/api-keys \
  -H "Content-Type: application/json" \
  -H "X-Admin-Token: ${APIKEY_ADMIN_TOKEN:-local-admin-token}" \
  -d '{
    "app": "crm",
    "env": "local",
    "name": "CRM local integration",
    "created_by": "admin@example.com",
    "scopes": ["demo:read", "messages:write"],
    "expires_at": null
  }'
printf '\n'
