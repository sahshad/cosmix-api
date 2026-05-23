import { Request, Response, NextFunction, response } from "express";
import { authClient } from "../grpc/auth.client";
import { LoginRequest, LoginResponse, RegisterRequest, RegisterResponse } from "../gen/auth/auth";
import { ServiceError } from "@grpc/grpc-js";

export function loginHandler(req: Request, res: Response) {

  const loginReq: LoginRequest = {
    email: req.body.email,
    password: req.body.password,
  }

  function handler(err: ServiceError | null, result: LoginResponse) {
    if (err)
      return err

    res.cookie("refresh_token", result.refreshToken,
      {
        httpOnly: true,
        sameSite: "strict",
      }
    );

    res.json(result);
  }

  authClient.login(loginReq, handler);
}

export function registerHandler(req: Request, res: Response) {

  const registerReq: RegisterRequest = {
    displayName: req.body.display_name,
    email: req.body.email,
    password: req.body.password
  }

  function handler(err: ServiceError | null, result: RegisterResponse) {
    if (err)
      return err;

    res.json(result)
  }

  authClient.register(registerReq, handler)
}