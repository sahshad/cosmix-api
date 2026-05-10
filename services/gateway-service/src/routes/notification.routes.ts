import { Router } from "express"
import { notificationProxy } from "../proxies/notification.proxy"
import { authenticate } from "../middlewares/auth.middleware"

export const notificationRoutes = Router()
notificationRoutes.use("/", authenticate, notificationProxy)
