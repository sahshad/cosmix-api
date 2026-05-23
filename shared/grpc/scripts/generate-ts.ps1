$ErrorActionPreference = "Stop"

$root = Resolve-Path "$PSScriptRoot\.."

protoc `
  --proto_path="$root\proto" `
  --proto_path="C:\Users\user\protoc\include" `
  --plugin=protoc-gen-ts_proto="$PSScriptRoot\ts_proto_plugin.cmd" `
  --ts_proto_out="$root\..\..\services\gateway-service\src\generated" `
  --ts_proto_opt=outputServices=grpc-js `
  auth/auth.proto `
  user/user.proto `
  notification/notification.proto `
  post/post.proto

Write-Host "TypeScript protobuf files generated successfully"