#!/usr/bin/env sh
set -eu
: "${API_KEY:?Set API_KEY to the generated API key}"
curl -sS http://localhost:8080/v1/ping \
  -H "Authorization: Bearer ${API_KEY}"
printf '\n'
