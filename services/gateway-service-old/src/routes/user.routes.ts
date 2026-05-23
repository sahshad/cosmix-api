import { Router } from "express";
import { userProxy } from "../proxies/user.proxy";
import { authenticate } from "../middlewares/auth.middleware";

export const userRoutes = Router()
userRoutes.use("/", authenticate, userProxy)