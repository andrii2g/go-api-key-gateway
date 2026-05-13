$ErrorActionPreference = "Stop"

if ([string]::IsNullOrWhiteSpace($env:API_KEY)) {
    throw "Set API_KEY to the generated API key."
}

$response = Invoke-RestMethod `
    -Method Get `
    -Uri "http://localhost:8080/v1/ping" `
    -Headers @{ "Authorization" = "Bearer $env:API_KEY" }

$response | ConvertTo-Json -Depth 20
