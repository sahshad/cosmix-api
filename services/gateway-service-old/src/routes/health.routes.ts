import { Router } from "express"

export const healthRoutes = Router()

healthRoutes.get("/", (_, res) => {
  res.json({ status: "API Gateway OK" })
})
