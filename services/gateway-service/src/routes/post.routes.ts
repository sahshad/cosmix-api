import { Router } from "express";
import { postProxy } from "../proxies/post.proxy";
import { authenticate } from "../middlewares/auth.middleware";

export const postRoutes = Router()
postRoutes.use("/", authenticate, postProxy)
