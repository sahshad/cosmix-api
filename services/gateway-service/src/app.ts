import express from "express"
import helmet from "helmet"
import cors from "cors"
import morgan from "morgan"

import { apiRateLimiter, authRateLimiter, notificationRateLimiter, userRateLimiter, postRateLimiter } from "./middlewares/rateLimit.middleware"
import { errorHandler } from "./middlewares/error.middleware"

import { authRoutes } from "./routes/auth.routes"
import { healthRoutes } from "./routes/health.routes"
import { userRoutes } from "./routes/user.routes"
import { assignRequestID } from "./middlewares/requestId.middleware"
import { notificationRoutes } from "./routes/notification.routes"
import { postRoutes } from "./routes/post.routes"

export const app = express()

app.use(assignRequestID)
app.use(express.json({ limit: "10kb" }))
app.use(express.urlencoded({ extended: true, limit: "10kb" }))

morgan.token("id", (req: any) => req.headers["x-request-id"])
morgan.format("custom", ":id :method :url :status :response-time ms - :res[content-length]")
app.use(morgan("custom"))

app.use(helmet())
app.use(cors({
    origin: "http://localhost:3000",
    credentials: true,
    methods: ["GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"],
    allowedHeaders: ["Content-Type", "Authorization", "X-Requested-With"],
    exposedHeaders: ["Content-Length", "Content-Type", "Authorization", "X-Requested-With"],
}))

app.use("/api/health", apiRateLimiter, healthRoutes)
app.use("/api/auth", authRateLimiter, authRoutes)
app.use("/api/users", userRateLimiter, userRoutes)
app.use("/api/notifications", notificationRateLimiter, notificationRoutes)
app.use("/api/posts", postRateLimiter, postRoutes)

app.use(errorHandler)
