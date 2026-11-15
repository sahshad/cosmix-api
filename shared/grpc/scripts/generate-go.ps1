$ErrorActionPreference = "Stop"

$root = Resolve-Path "$PSScriptRoot\.."
$repoRoot = Resolve-Path "$root\..\.."

$targets = @(
  @{ Proto = "auth/auth.proto"; Service = "auth-service" },
  @{ Proto = "user/user.proto"; Service = "user-service" },
  @{ Proto = "notification/notification.proto"; Service = "notification-service" },
  @{ Proto = "post/post.proto"; Service = "post-service" },
  @{ Proto = "communication/communication.proto"; Service = "communication-service" }
)

foreach ($target in $targets) {
  $moduleRoot = "$repoRoot\services\$($target.Service)"
  New-Item -ItemType Directory -Path "$moduleRoot\internal\grpc\generated" -Force | Out-Null

  protoc `
    --proto_path="$root\proto" `
    --proto_path="C:\Users\user\protoc\include" `
    --go_out="$moduleRoot" `
    --go_opt=module="$($target.Service)" `
    --go-grpc_out="$moduleRoot" `
    --go-grpc_opt=module="$($target.Service)" `
    $target.Proto

  if ($LASTEXITCODE -ne 0) {
    throw "Failed to generate Go code for $($target.Proto)"
  }

  Write-Host "Generated $($target.Service)/internal/grpc/generated from $($target.Proto)"
}

Write-Host "Go protobuf files generated successfully"
