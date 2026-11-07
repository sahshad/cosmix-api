$root = Split-Path -Parent $PSScriptRoot
$servicesRoot = "$root/services"

function Get-MigrationServices {
    Get-ChildItem -Path $servicesRoot -Directory | Where-Object {
        $files = @(Get-ChildItem -Path (Join-Path $_.FullName "migrations") -Filter "*.sql" -File -ErrorAction SilentlyContinue)
        $files.Count -gt 0
    } | Select-Object -ExpandProperty Name | Sort-Object
}

function Read-ServiceSelection {
    param([string[]]$Services)

    $index = 0

    while ($true) {
        Clear-Host
        Write-Host "Select a service to migrate UP:" -ForegroundColor Cyan
        Write-Host ""

        for ($i = 0; $i -lt $Services.Count; $i++) {
            if ($i -eq $index) {
                Write-Host "  > $($Services[$i])" -ForegroundColor Green
            } else {
                Write-Host "    $($Services[$i])"
            }
        }

        Write-Host ""
        Write-Host "Use Up/Down arrows and press Enter (Esc to cancel)." -ForegroundColor DarkGray

        $key = [Console]::ReadKey($true)

        switch ($key.Key) {
            "UpArrow" {
                if ($index -gt 0) { $index-- }
            }
            "DownArrow" {
                if ($index -lt $Services.Count - 1) { $index++ }
            }
            "Enter" {
                return $Services[$index]
            }
            "Escape" {
                return $null
            }
        }
    }
}

function Get-DatabaseUrl {
    param([string]$ServicePath)

    $envFile = Join-Path $ServicePath ".env.local"

    if (!(Test-Path $envFile)) {
        $envFile = Join-Path $ServicePath ".env.docker"
    }

    if (!(Test-Path $envFile)) {
        return $null
    }

    foreach ($line in Get-Content $envFile) {
        if ($line -match '^\s*DATABASE_URL\s*=\s*(.+?)\s*$') {
            $value = $Matches[1].Trim()

            if (($value.StartsWith('"') -and $value.EndsWith('"')) -or ($value.StartsWith("'") -and $value.EndsWith("'"))) {
                $value = $value.Substring(1, $value.Length - 2)
            }

            return $value
        }
    }

    return $null
}

if (!(Get-Command goose -ErrorAction SilentlyContinue)) {
    Write-Host "goose is not installed or not available on PATH." -ForegroundColor Red
    exit 1
}

$services = @(Get-MigrationServices)

if ($services.Count -eq 0) {
    Write-Host "No services with migrations found." -ForegroundColor Red
    exit 1
}

$service = Read-ServiceSelection -Services $services

if (!$service) {
    Write-Host "Cancelled."
    exit 0
}

$servicePath = Join-Path $servicesRoot $service
$migrationsPath = Join-Path $servicePath "migrations"

$databaseUrl = Get-DatabaseUrl -ServicePath $servicePath

if (!$databaseUrl) {
    Write-Host "DATABASE_URL not found for '$service' (.env.local or .env.docker)." -ForegroundColor Red
    exit 1
}

Clear-Host
Write-Host "Running migrations UP for '$service'..." -ForegroundColor Cyan

$env:GOOSE_DRIVER = "postgres"
$env:GOOSE_DBSTRING = $databaseUrl

Push-Location $migrationsPath

try {
    goose up

    if ($LASTEXITCODE -ne 0) {
        Write-Host "Migration UP failed for '$service'." -ForegroundColor Red
        exit 1
    }

    Write-Host "Migration UP completed for '$service'." -ForegroundColor Green
} finally {
    Pop-Location
}
