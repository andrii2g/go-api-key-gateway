#!/usr/bin/env sh
set -eu
mkdir -p .secrets
openssl rand -base64 32 > .secrets/api_key_pepper
chmod 600 .secrets/api_key_pepper || true
printf 'created .secrets/api_key_pepper\n'
