param([string]$Tab)

$root = Split-Path -Parent $PSScriptRoot
$legacyStateDir = Join-Path $root 'tmp/dev-pids'
$script:DevProcessStateRoot = [System.IO.Path]::GetFullPath($root).TrimEnd('\')
$stateHashProvider = [System.Security.Cryptography.SHA256]::Create()
try {
    $stateHashBytes = $stateHashProvider.ComputeHash([System.Text.Encoding]::UTF8.GetBytes($script:DevProcessStateRoot.ToLowerInvariant()))
    $script:DevProcessStateKey = ([System.BitConverter]::ToString($stateHashBytes) -replace '-', '').Substring(0, 16)
} finally {
    $stateHashProvider.Dispose()
}

$script:DevProcessStateFile = Join-Path ([System.IO.Path]::GetTempPath()) "cosmix-api-dev-$($script:DevProcessStateKey).json"
$script:DevProcessStateMutexName = "Local\CosmixApiDev_$($script:DevProcessStateKey)"

function Invoke-DevProcessStateLock {
    param([scriptblock]$Action)

    $mutex = New-Object -TypeName System.Threading.Mutex -ArgumentList @($false, $script:DevProcessStateMutexName)
    $hasLock = $false

    try {
        try {
            $hasLock = $mutex.WaitOne([TimeSpan]::FromSeconds(30))
        } catch [System.Threading.AbandonedMutexException] {
            $hasLock = $true
        }

        if (-not $hasLock) {
            throw "Timed out waiting for the dev process state lock."
        }

        & $Action
    } finally {
        if ($hasLock) {
            $mutex.ReleaseMutex()
        }
        $mutex.Dispose()
    }
}

function Invoke-DevProcessState {
    param([scriptblock]$Update)

    Invoke-DevProcessStateLock {
        $document = if (Test-Path -LiteralPath $script:DevProcessStateFile -PathType Leaf) {
            Get-Content -LiteralPath $script:DevProcessStateFile -Raw -ErrorAction Stop | ConvertFrom-Json -ErrorAction Stop
        } else {
            [pscustomobject]@{ stopRequested = $false; processes = @() }
        }

        if ($null -eq $document.stopRequested -or $null -eq $document.processes) {
            throw "Invalid dev process state file: $script:DevProcessStateFile"
        }

        if ($Update) {
            $null = & $Update $document
            $document | ConvertTo-Json -Depth 4 | Set-Content -LiteralPath $script:DevProcessStateFile -Encoding UTF8 -ErrorAction Stop
        }

        return $document
    }
}

function Test-DevProcessStateRecord {
    param([object]$Record)

    try {
        $process = Get-Process -Id ([int]$Record.pid) -ErrorAction SilentlyContinue
        if (-not $process -or $process.Name -ne [string]$Record.name) {
            return $false
        }

        $recordedStart = [datetime]$Record.startTime
        return $process.StartTime.ToUniversalTime().ToString('o') -eq $recordedStart.ToUniversalTime().ToString('o')
    } catch {
        return $false
    }
}

function Remove-DevProcessStateRecord {
    param(
        [int]$ProcessId,
        [string]$StartTime
    )

    $null = Invoke-DevProcessState {
        param($state)
        $state.processes = @($state.processes | Where-Object {
            [int]$_.pid -ne $ProcessId -or [string]$_.startTime -ne $StartTime
        })
    }
}

function Start-DevTab {
    param([string]$Service)

    $serviceScript = Join-Path $root "services/$Service/$Service.ps1"
    $process = Get-Process -Id $PID
    $serviceProcess = $null

    $state = [ordered]@{
        service   = $Service
        pid       = $PID
        name      = $process.Name
        startTime = $process.StartTime.ToString('o')
        startedAt = (Get-Date).ToString('o')
    }

    try {
        $null = Invoke-DevProcessState {
            param($sharedState)
            $sharedState.processes = @($sharedState.processes | Where-Object {
                [int]$_.pid -ne $PID -or [string]$_.startTime -ne [string]$state.startTime
            }) + @($state)
        }

        $argumentList = "-NoProfile -ExecutionPolicy Bypass -File `"$serviceScript`""
        $serviceProcess = Start-Process -FilePath 'powershell.exe' -ArgumentList $argumentList -NoNewWindow -PassThru

        while (-not $serviceProcess.HasExited) {
            if ((Invoke-DevProcessState).stopRequested) {
                & taskkill.exe /PID $serviceProcess.Id /T /F 2>&1 | Out-Null
                return 'stop'
            }

            Start-Sleep -Milliseconds 250
        }
    }
    finally {
        Remove-DevProcessStateRecord -ProcessId $PID -StartTime ([string]$state.startTime)
    }
}

if ($Tab) {
    if ($Tab -notmatch '^[a-z][a-z0-9-]*$') {
        Write-Host "Invalid service name: $Tab"
        return
    }

    $result = Start-DevTab -Service $Tab
    return $result
}

$records = @()

$sharedState = Invoke-DevProcessState {
    param($sharedState)
    $sharedState.processes = @($sharedState.processes | Where-Object { Test-DevProcessStateRecord -Record $_ })
}

foreach ($state in @($sharedState.processes)) {
    $records += [pscustomobject]@{
        Service = [string]$state.service
        Pid     = [int]$state.pid
    }
}

if (Test-Path -LiteralPath $legacyStateDir) {
    foreach ($file in Get-ChildItem -LiteralPath $legacyStateDir -Filter '*.json' -File -ErrorAction SilentlyContinue) {
        try {
            $state = Get-Content -LiteralPath $file.FullName -Raw | ConvertFrom-Json
            if (Test-DevProcessStateRecord -Record $state) {
                $records += [pscustomobject]@{
                    Service = if ($state.service) { [string]$state.service } else { $file.BaseName }
                    Pid     = [int]$state.pid
                }
            } else {
                Remove-Item -LiteralPath $file.FullName -Force -ErrorAction SilentlyContinue
            }
        } catch {
            Remove-Item -LiteralPath $file.FullName -Force -ErrorAction SilentlyContinue
        }
    }
}

$servicePorts = @(
    [pscustomobject]@{ Service = 'gateway-service'; Port = 5001 }
    [pscustomobject]@{ Service = 'auth-service'; Port = 50051 }
    [pscustomobject]@{ Service = 'user-service'; Port = 50052 }
    [pscustomobject]@{ Service = 'notification-service'; Port = 50053 }
    [pscustomobject]@{ Service = 'post-service'; Port = 50054 }
    [pscustomobject]@{ Service = 'realtime-service'; Port = 50055 }
    [pscustomobject]@{ Service = 'communication-service'; Port = 50056 }
)

$listeners = @()
$seenProcesses = @{}

foreach ($servicePort in $servicePorts) {
    foreach ($connection in @(Get-NetTCPConnection -LocalPort $servicePort.Port -State Listen -ErrorAction SilentlyContinue)) {
        $processId = [int]$connection.OwningProcess
        if ($processId -le 0 -or $seenProcesses.ContainsKey($processId)) {
            continue
        }

        $process = Get-Process -Id $processId -ErrorAction SilentlyContinue
        if (-not $process) {
            continue
        }

        $seenProcesses[$processId] = $true
        $listeners += [pscustomobject]@{
            Service = $servicePort.Service
            Port    = $servicePort.Port
            Pid     = $processId
            Name    = $process.Name
        }
    }
}

if ($records.Count -gt 0 -or $listeners.Count -gt 0) {
    Write-Host "Dev services are already running."

    if ($records.Count -gt 0) {
        $records | ForEach-Object { Write-Host "  - $($_.Service) (PID $($_.Pid))" }
    }

    if ($listeners.Count -gt 0) {
        $listeners | ForEach-Object { Write-Host "  - $($_.Service) listening on port $($_.Port) (PID $($_.Pid), $($_.Name))" }
    }

    Write-Host "Run scripts/stop-dev.ps1 first, or close their terminal tabs."
    exit 1
}

$null = Invoke-DevProcessState {
    param($sharedState)
    $sharedState.stopRequested = $false
}

Write-Host "Starting infrastructure..."

docker compose -f "$root/docker-compose.dev.yml" up -d

if ($LASTEXITCODE -ne 0) {
    Write-Host "Failed to start infrastructure."
    exit 1
}

Write-Host "Infrastructure started successfully."

$services = @(
    'auth-service'
    'user-service'
    'post-service'
    'notification-service'
    'communication-service'
    'realtime-service'
    'gateway-service'
)

$wtArgs = @('-w', '-1')

for ($i = 0; $i -lt $services.Count; $i++) {
    $service = $services[$i]
    $tabCommand = "`$result = & '$($PSCommandPath.Replace("'", "''"))' -Tab $service; if (`$result -eq 'stop') { exit 0 }"
    $encodedTabCommand = [Convert]::ToBase64String([Text.Encoding]::Unicode.GetBytes($tabCommand))
    $wtArgs += @('new-tab', '--title', "cosmix-dev: $service", '--suppressApplicationTitle', 'powershell', '-NoExit', '-EncodedCommand', $encodedTabCommand)

    if ($i -lt $services.Count - 1) {
        $wtArgs += ';'
    }
}

try {
    & wt @wtArgs
} catch {
    Write-Host "Failed to open the Windows Terminal window."
    exit 1
}

Write-Host "Services started successfully."
