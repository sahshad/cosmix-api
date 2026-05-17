import dotenv from "dotenv"

dotenv.config()

export const env = {
  port: Number(process.env.PORT) || 3000,
  authServiceUrl: process.env.AUTH_SERVICE_URL!,
  userServiceUrl: process.env.USER_SERVICE_URL!,
  notificationServiceUrl: process.env.NOTIFICATION_SERVICE_URL!,
  postServiceUrl: process.env.POST_SERVICE_URL!,
  jwtPublicKey: process.env.JWT_PUBLIC_KEY!
}
