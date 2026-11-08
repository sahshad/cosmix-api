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

function Get-DevServicePorts {
    return @(
        [pscustomobject]@{ Service = 'gateway-service'; Port = 5001 }
        [pscustomobject]@{ Service = 'auth-service'; Port = 50051 }
        [pscustomobject]@{ Service = 'user-service'; Port = 50052 }
        [pscustomobject]@{ Service = 'notification-service'; Port = 50053 }
        [pscustomobject]@{ Service = 'post-service'; Port = 50054 }
        [pscustomobject]@{ Service = 'realtime-service'; Port = 50055 }
        [pscustomobject]@{ Service = 'communication-service'; Port = 50056 }
    )
}

function Get-DevPortListeners {
    $listeners = @()
    $seenProcesses = @{}

    foreach ($servicePort in Get-DevServicePorts) {
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

    return @($listeners)
}

function Test-DevToolProcess {
    param([string]$Name)

    return ($Name -match '^(air|node|pnpm|go)(\.exe)?$') -or ($Name -match '-service(\.exe)?$')
}

function Test-DevHostProcess {
    param([string]$Name)

    return $Name -match '^(windowsterminal|windowsterminalpreview|wt|openconsole|conhost|explorer|code|cursor|devenv|opencode|idea64|webstorm64|goland64|rider64|pycharm64)(\.exe)?$'
}

function Get-DevProcessAncestors {
    param([int]$ProcessId)

    $ancestors = @()
    $currentId = $ProcessId

    for ($i = 0; $i -lt 12; $i++) {
        $info = Get-CimInstance Win32_Process -Filter "ProcessId = $currentId" -ErrorAction SilentlyContinue
        if (-not $info -or -not $info.ParentProcessId) {
            break
        }

        $parentId = [int]$info.ParentProcessId
        $parent = Get-Process -Id $parentId -ErrorAction SilentlyContinue
        if (-not $parent -or (Test-DevHostProcess -Name $parent.Name)) {
            break
        }

        $ancestors += $parent
        $currentId = $parentId
    }

    return @($ancestors)
}

function Stop-DevProcessTree {
    param([int]$ProcessId)

    if (-not (Get-Process -Id $ProcessId -ErrorAction SilentlyContinue)) {
        return $true
    }

    & taskkill.exe /PID $ProcessId /T /F 2>&1 | Out-Null
    Start-Sleep -Milliseconds 200

    return -not (Get-Process -Id $ProcessId -ErrorAction SilentlyContinue)
}

function Resolve-DevProcessRoot {
    param([int]$ProcessId)

    $process = Get-Process -Id $ProcessId -ErrorAction SilentlyContinue
    if (-not $process) {
        return $ProcessId
    }

    $chain = @($process) + @(Get-DevProcessAncestors -ProcessId $ProcessId)
    $devProcesses = @($chain | Where-Object { Test-DevToolProcess -Name $_.Name })

    if ($devProcesses.Count -gt 0) {
        return $devProcesses[-1].Id
    }

    return $ProcessId
}

function Get-DevTabProcessId {
    param([int]$ProcessId)

    $tab = @(Get-DevProcessAncestors -ProcessId $ProcessId |
        Where-Object { $_.Name -match '^(powershell|pwsh)(\.exe)?$' } |
        Select-Object -Last 1)
    if ($tab.Count -eq 0) {
        return $null
    }

    $tabProcessId = $tab[0].Id
    $tabInfo = Get-CimInstance Win32_Process -Filter "ProcessId = $tabProcessId" -ErrorAction SilentlyContinue
    $commandLine = [string]$tabInfo.CommandLine

    if ($commandLine -match '(?i)(-NoExit.*(-File|-EncodedCommand)|services[\\/].*\.ps1)') {
        return $tabProcessId
    }

    return $null
}

Write-Host "Stopping dev services..."

$stopped = @()
$pendingStops = @()

$sharedState = Invoke-DevProcessState {
    param($sharedState)
    $sharedState.stopRequested = $true
    $sharedState.processes = @($sharedState.processes | Where-Object { Test-DevProcessStateRecord -Record $_ })
}

foreach ($state in @($sharedState.processes)) {
    $pendingStops += [pscustomobject]@{
        Service   = [string]$state.service
        Pid       = [int]$state.pid
        Name      = [string]$state.name
        StartTime = [string]$state.startTime
        Storage   = 'shared'
    }
}

if (Test-Path -LiteralPath $legacyStateDir) {
    foreach ($file in Get-ChildItem -LiteralPath $legacyStateDir -Filter '*.json' -File -ErrorAction SilentlyContinue) {
        try {
            $state = Get-Content -LiteralPath $file.FullName -Raw | ConvertFrom-Json
        } catch {
            Remove-Item -LiteralPath $file.FullName -Force -ErrorAction SilentlyContinue
            continue
        }

        if (-not (Test-DevProcessStateRecord -Record $state)) {
            Remove-Item -LiteralPath $file.FullName -Force -ErrorAction SilentlyContinue
            continue
        }

        $serviceName = if ($state.service) { [string]$state.service } else { $file.BaseName }
        $stopFile = Join-Path $legacyStateDir "$serviceName.stop"

        try {
            Set-Content -LiteralPath $stopFile -Value '' -Encoding ASCII -Force -ErrorAction Stop
        } catch {
            Write-Warning "Failed to signal $serviceName to stop."
        }

        $pendingStops += [pscustomobject]@{
            Service   = $serviceName
            Pid       = [int]$state.pid
            Name      = [string]$state.name
            StartTime = [string]$state.startTime
            Storage   = 'legacy'
            StateFile = $file.FullName
            StopFile  = $stopFile
        }
    }
}

for ($i = 0; $i -lt 40 -and $pendingStops.Count -gt 0; $i++) {
    $stillRunning = @()

    foreach ($pending in $pendingStops) {
        if (Test-DevProcessStateRecord -Record $pending) {
            $stillRunning += $pending
        } else {
            $stopped += $pending.Service
            if ($pending.Storage -eq 'shared') {
                Remove-DevProcessStateRecord -ProcessId $pending.Pid -StartTime $pending.StartTime
            } else {
                Remove-Item -LiteralPath $pending.StateFile -Force -ErrorAction SilentlyContinue
                Remove-Item -LiteralPath $pending.StopFile -Force -ErrorAction SilentlyContinue
            }
        }
    }

    $pendingStops = @($stillRunning)

    if ($pendingStops.Count -gt 0) {
        Start-Sleep -Milliseconds 250
    }
}

foreach ($pending in $pendingStops) {
    if (Stop-DevProcessTree -ProcessId $pending.Pid) {
        $stopped += $pending.Service
        if ($pending.Storage -eq 'shared') {
            Remove-DevProcessStateRecord -ProcessId $pending.Pid -StartTime $pending.StartTime
        } else {
            Remove-Item -LiteralPath $pending.StateFile -Force -ErrorAction SilentlyContinue
        }
    } else {
        Write-Warning "Failed to stop $($pending.Service) (PID $($pending.Pid))."
    }

    if ($pending.Storage -eq 'legacy') {
        Remove-Item -LiteralPath $pending.StopFile -Force -ErrorAction SilentlyContinue
    }
}

foreach ($listener in @(Get-DevPortListeners)) {
    $tabProcessId = Get-DevTabProcessId -ProcessId $listener.Pid
    $rootProcessId = if ($tabProcessId) { $tabProcessId } else { Resolve-DevProcessRoot -ProcessId $listener.Pid }

    if (Stop-DevProcessTree -ProcessId $rootProcessId) {
        $stopped += "$($listener.Service) (port $($listener.Port))"
    } else {
        Write-Warning "Failed to stop process on port $($listener.Port) (PID $($listener.Pid))."
    }
}

for ($i = 0; $i -lt 10; $i++) {
    $remainingPorts = @(Get-DevPortListeners)
    if ($remainingPorts.Count -eq 0) {
        break
    }

    foreach ($listener in $remainingPorts) {
        $tabProcessId = Get-DevTabProcessId -ProcessId $listener.Pid
        $rootProcessId = if ($tabProcessId) { $tabProcessId } else { Resolve-DevProcessRoot -ProcessId $listener.Pid }
        $null = Stop-DevProcessTree -ProcessId $rootProcessId
    }

    Start-Sleep -Milliseconds 500
}

if (Test-Path -LiteralPath $legacyStateDir) {
    foreach ($file in Get-ChildItem -LiteralPath $legacyStateDir -Filter '*.json' -File -ErrorAction SilentlyContinue) {
        try {
            $state = Get-Content -LiteralPath $file.FullName -Raw | ConvertFrom-Json
            if (-not (Test-DevProcessStateRecord -Record $state)) {
                Remove-Item -LiteralPath $file.FullName -Force -ErrorAction SilentlyContinue
            }
        } catch {
            Remove-Item -LiteralPath $file.FullName -Force -ErrorAction SilentlyContinue
        }
    }

    Get-ChildItem -LiteralPath $legacyStateDir -Filter '*.tmp' -File -ErrorAction SilentlyContinue | Remove-Item -Force -ErrorAction SilentlyContinue
    Get-ChildItem -LiteralPath $legacyStateDir -Filter '*.stop' -File -ErrorAction SilentlyContinue | Remove-Item -Force -ErrorAction SilentlyContinue

    if (-not (Get-ChildItem -LiteralPath $legacyStateDir -ErrorAction SilentlyContinue)) {
        Remove-Item -LiteralPath $legacyStateDir -Force -ErrorAction SilentlyContinue
    }
}

$null = Invoke-DevProcessState {
    param($sharedState)
    $sharedState.processes = @($sharedState.processes | Where-Object { Test-DevProcessStateRecord -Record $_ })
}

$remaining = @(Get-DevPortListeners)

if ($remaining.Count -gt 0) {
    Write-Warning "Some dev service ports are still in use:"
    $remaining | ForEach-Object { Write-Warning "  - $($_.Service) on port $($_.Port) (PID $($_.Pid), $($_.Name))" }
} elseif ($stopped.Count -gt 0) {
    Write-Host "Stopped: $(($stopped | Select-Object -Unique) -join ', ')"
} else {
    Write-Host "No running dev services were found."
}

Write-Host "Stopping infrastructure..."
docker compose -f "$root/docker-compose.dev.yml" down
Write-Host "Infrastructure stopped."
