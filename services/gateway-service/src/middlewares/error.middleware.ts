import { Request, Response, NextFunction } from "express"
import { sendError } from "../utils/response.util"

export interface ApiError extends Error {
  statusCode?: number
  details?: unknown
}

export function errorHandler(
  err: ApiError,
  _req: Request,
  res: Response,
  _next: NextFunction
) {
  const statusCode = err.statusCode ?? 500

  console.error("API Gateway Error:", {
    message: err.message,
    statusCode,
    details: err.details,
    stack: err.stack
  })

  return sendError(
    res,
    statusCode,
    statusCode >= 500 ? "Internal server error" : err.message
  )
}
