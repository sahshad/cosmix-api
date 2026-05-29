$ErrorActionPreference = "Stop"

$root = Resolve-Path "$PSScriptRoot\.."

$tsProtoPlugin = Resolve-Path `
  "$root\..\..\services\gateway-service\node_modules\.bin\protoc-gen-ts_proto.cmd"

protoc `
  --proto_path="$root\proto" `
  --proto_path="C:\Users\user\protoc\include" `
  --plugin=protoc-gen-ts_proto="$tsProtoPlugin" `
  --ts_proto_out="$root\..\..\services\gateway-service\src\generated" `
  --ts_proto_opt=outputServices=grpc-js,nestJs=true `
  auth/auth.proto `
  user/user.proto `
  notification/notification.proto `
  post/post.proto

Write-Host "TypeScript protobuf files generated successfully"