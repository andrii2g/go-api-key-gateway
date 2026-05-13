# Quickstart

This quickstart is for Linux in WSL.

## Prerequisites

- Docker Desktop installed on Windows with WSL integration enabled
- WSL shell with `docker`, `curl`, and `openssl` available

## 1. Start the stack

From the repository root:

Make the helper scripts executable:

```sh
chmod +x ./scripts/generate-dev-pepper.sh ./scripts/curl-create-key.sh ./scripts/curl-protected-ping.sh ./scripts/curl-revoke-key.sh
```

Generate a local development pepper file at `.secrets/api_key_pepper`. This is the base64 secret the service uses as the server-side HMAC pepper:

```sh
./scripts/generate-dev-pepper.sh
```

Build and start MySQL, run migrations, and start the sample service:

```sh
docker compose up --build
```

Keep that shell running.

## 2. Create an API key

Open a second WSL shell in the same repository:

Call the admin endpoint that creates a new API key using the default local admin token:

```sh
./scripts/curl-create-key.sh
```

Copy the `api_key` value from the JSON response.

## 3. Call a protected endpoint

In the second shell:

Set the generated full API key in the environment for the next helper script:

```sh
export API_KEY='paste-the-api-key-here'
```

Call the protected `/v1/ping` endpoint with `Authorization: Bearer $API_KEY`:

```sh
./scripts/curl-protected-ping.sh
```

Expected response:

```json
{"message":"pong", ...}
```

## 4. Revoke the key

Use the `id` returned by the create response:

Set the numeric key id you want to revoke:

```sh
export API_KEY_ID='1'
```

Call the admin revoke endpoint for that key id:

```sh
./scripts/curl-revoke-key.sh
```

Expected response:

```json
{"revoked":true}
```

## 5. Stop the stack

In the first shell, press `Ctrl+C`, then run:

```sh
docker compose down -v
```

## Troubleshooting

- If `docker compose` fails in WSL, verify Docker Desktop WSL integration is enabled.
- If the service cannot start, check that `.secrets/api_key_pepper` exists.
- If `curl` scripts return unauthorized, confirm `APIKEY_ADMIN_TOKEN` is still the default `local-admin-token`.
