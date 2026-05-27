$root = Split-Path -Parent $PSScriptRoot

Write-Host "Stopping infrastructure..."

docker compose -f "$root/docker-compose.dev.yml" down

Write-Host "Infrastructure stopped."
Write-Host "Close the Cosmix terminal window manually if needed."