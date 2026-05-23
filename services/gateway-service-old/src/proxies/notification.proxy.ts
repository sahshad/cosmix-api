import { createProxyMiddleware } from "http-proxy-middleware"
import { services } from "../config/service"

export const notificationProxy = createProxyMiddleware({
    target: services.notification,
    changeOrigin: true,
    pathRewrite: {
        "^/notifications": ""
    },
    on: {
        proxyReq: (proxyReq, req: any) => {
            if (req.body && Object.keys(req.body).length > 0) {
                const bodyData = JSON.stringify(req.body)
                proxyReq.setHeader("Content-Type", "application/json")
                proxyReq.setHeader("Content-Length", Buffer.byteLength(bodyData))
                proxyReq.write(bodyData)
            }
        }
    }
})