import { env } from "./env"

export const services = {
  auth: env.authServiceUrl,
  user: env.userServiceUrl,
  notification: env.notificationServiceUrl,
  post: env.postServiceUrl,
}
