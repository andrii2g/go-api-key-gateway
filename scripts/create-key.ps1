$ErrorActionPreference = "Stop"

$adminToken = $env:APIKEY_ADMIN_TOKEN
if ([string]::IsNullOrWhiteSpace($adminToken)) {
    $adminToken = "local-admin-token"
}

$body = @{
    app = "crm"
    env = "local"
    name = "CRM local integration"
    created_by = "admin@example.com"
    scopes = @("demo:read", "messages:write")
    expires_at = $null
} | ConvertTo-Json -Depth 8

$response = Invoke-RestMethod `
    -Method Post `
    -Uri "http://localhost:8080/admin/api-keys" `
    -Headers @{ "X-Admin-Token" = $adminToken } `
    -ContentType "application/json" `
    -Body $body

$response | ConvertTo-Json -Depth 20
