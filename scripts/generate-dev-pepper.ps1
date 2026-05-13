$ErrorActionPreference = "Stop"

$secretDir = Join-Path (Get-Location) ".secrets"
New-Item -ItemType Directory -Force -Path $secretDir | Out-Null

$bytes = New-Object byte[] 32
$rng = [System.Security.Cryptography.RandomNumberGenerator]::Create()
try {
    $rng.GetBytes($bytes)
}
finally {
    $rng.Dispose()
}

$pepperPath = Join-Path $secretDir "api_key_pepper"
[Convert]::ToBase64String($bytes) | Set-Content -Path $pepperPath -Encoding ascii -NoNewline
Write-Host "created .secrets/api_key_pepper"
