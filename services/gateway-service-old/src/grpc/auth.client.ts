import * as grpc from "@grpc/grpc-js";

import {
  AuthServiceClient
} from "../gen/auth/auth";

export const authClient =
  new AuthServiceClient(
    "auth-service:50051",
    grpc.credentials.createInsecure()
  );