$ErrorActionPreference = "Stop"

$root = Resolve-Path "$PSScriptRoot\.."

protoc `
  --proto_path="$root\proto" `
  --proto_path="C:\Users\user\protoc\include" `
  --go_out="$root\gen\go" `
  --go_opt=paths=source_relative `
  --go-grpc_out="$root\gen\go" `
  --go-grpc_opt=paths=source_relative `
  auth/auth.proto `
  user/user.proto `
  notification/notification.proto `
  post/post.proto

Write-Host "Go protobuf files generated successfully"