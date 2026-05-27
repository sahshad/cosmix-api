$root = Split-Path -Parent (Split-Path -Parent $PSScriptRoot)

Set-Location "$root/services/gateway-service"

pnpm start:dev