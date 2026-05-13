# go-api-key-gateway

Copy-first Go API key component with a MySQL-backed store, sample `net/http` service, migrations, and local Docker Compose flow.

## API key format

Keys use:

```text
ak_{app}_{public_key}_{secret}
```

Example:

```text
ak_crm_7F3K9Q2M8N4P6R1T_VL8f2pQ9uYk6sN3xA0mB7cD4eF1gH5jKp4xR9zQ
```

## Security model

- Only an HMAC-SHA256 hash of the secret is stored.
- The pepper is loaded from base64 config or a mounted secret file.
- Full API keys are returned once on creation and should never be logged.
- Validation supports revocation, expiration, environment checks, and scopes.

## Local run

POSIX or WSL:

```sh
./scripts/generate-dev-pepper.sh
docker compose up --build
```

Windows PowerShell:

```powershell
powershell -ExecutionPolicy Bypass -File .\scripts\generate-dev-pepper.ps1
docker compose up --build
```

## Create a key

POSIX or WSL:

```sh
./scripts/curl-create-key.sh
```

Windows PowerShell:

```powershell
powershell -ExecutionPolicy Bypass -File .\scripts\create-key.ps1
```

## Call protected endpoints

POSIX or WSL:

```sh
API_KEY='ak_crm_...' ./scripts/curl-protected-ping.sh
```

Windows PowerShell:

```powershell
$env:API_KEY = 'ak_crm_...'
powershell -ExecutionPolicy Bypass -File .\scripts\protected-ping.ps1
```

## Revoke a key

POSIX or WSL:

```sh
API_KEY_ID='1' ./scripts/curl-revoke-key.sh
```

Windows PowerShell:

```powershell
$env:API_KEY_ID = '1'
powershell -ExecutionPolicy Bypass -File .\scripts\revoke-key.ps1
```

## Configuration

See [.env.example](/C:/github/a2g.name/go-api-key-gateway/.env.example) for the local environment variables. The main required values are:

- `APIKEY_ENV`
- `APIKEY_DB_DSN`
- one of `APIKEY_PEPPER_BASE64` or `APIKEY_PEPPER_FILE`
- `APIKEY_ADMIN_TOKEN`

## Copy-first integration

Copy the packages you need into your own Go service:

- `apikey/`
- `storage/sqlstore/`
- `config/`
- `httpapi/` if you want the sample HTTP layer

Then replace imports from `github.com/andrii2g/go-api-key-gateway/...` with your own module path.

## Production notes

- Terminate TLS before serving protected endpoints.
- Store the pepper in a secret manager or mounted secret file.
- Add rate limiting and stronger admin authentication before production use.
