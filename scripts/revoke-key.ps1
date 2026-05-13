$ErrorActionPreference = "Stop"

if ([string]::IsNullOrWhiteSpace($env:API_KEY_ID)) {
    throw "Set API_KEY_ID to the numeric key id to revoke."
}

$adminToken = $env:APIKEY_ADMIN_TOKEN
if ([string]::IsNullOrWhiteSpace($adminToken)) {
    $adminToken = "local-admin-token"
}

$response = Invoke-RestMethod `
    -Method Post `
    -Uri "http://localhost:8080/admin/api-keys/$env:API_KEY_ID/revoke" `
    -Headers @{ "X-Admin-Token" = $adminToken }

$response | ConvertTo-Json -Depth 20
