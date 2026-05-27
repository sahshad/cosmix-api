param(
    [Parameter(Mandatory = $true)]
    [string]$Service,

    [Parameter(Mandatory = $true)]
    [string]$Name
)

$root = Split-Path -Parent $PSScriptRoot

$servicePath = "$root/services/$Service"
$migrationsPath = "$servicePath/migrations"

if (!(Test-Path $servicePath)) {
    Write-Host "Service '$Service' does not exist."
    exit 1
}

if (!(Test-Path $migrationsPath)) {
    Write-Host "Migrations folder does not exist for '$Service'."
    exit 1
}

Push-Location $migrationsPath

goose create $Name sql

Pop-Location