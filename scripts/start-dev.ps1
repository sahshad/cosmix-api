$root = Split-Path -Parent $PSScriptRoot

Write-Host "Starting infrastructure..."

docker compose -f "$root/docker-compose.dev.yml" up -d

if ($LASTEXITCODE -ne 0) {
    Write-Host "Failed to start infrastructure."
    exit 1
}

Write-Host "Infrastructure started successfully."

wt `
new-tab --title auth-service powershell -NoExit -File "$root/services/auth-service/auth-service.ps1" `
`; `
new-tab --title user-service powershell -NoExit -File "$root/services/user-service/user-service.ps1" `
`; `
new-tab --title post-service powershell -NoExit -File "$root/services/post-service/post-service.ps1" `
`; `
new-tab --title notification-service powershell -NoExit -File "$root/services/notification-service/notification-service.ps1" `
`; `
new-tab --title gateway-service powershell -NoExit -File "$root/services/gateway-service/gateway-service.ps1"

Write-Host "Services started successfully."