import { Router } from "express"
import { loginHandler, registerHandler } from "../controllers/auth.controller"
// import { authProxy } from "../proxies/auth.proxy"

export const authRoutes = Router()
authRoutes.post("/login", loginHandler)
authRoutes.post("/register", registerHandler)
