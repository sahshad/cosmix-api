import rateLimit from "express-rate-limit"

const commonOptions = {
  standardHeaders: true,
  legacyHeaders: false,
  message: {
    success: false,
    message: "Too many requests, please try again later.",
  },
}

export const apiRateLimiter = rateLimit({
  ...commonOptions,
  windowMs: 1 * 60 * 1000,
  max: 100,
})

export const authRateLimiter = rateLimit({
  ...commonOptions,
  windowMs: 15 * 60 * 1000,
  max: 20,
})

export const userRateLimiter = rateLimit({
  ...commonOptions,
  windowMs: 1 * 60 * 1000,
  max: 100,
})