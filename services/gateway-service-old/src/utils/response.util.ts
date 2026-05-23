import { Response } from "express"

export const sendResponse = (
  res: Response,
  statusCode: number,
  success: boolean,
  message: string,
  data: any = null
) => {
  return res.status(statusCode).json({
    success,
    message,
    ...(data && { data }),
  })
}

export const sendError = (res: Response, statusCode: number, message: string) => {
  return sendResponse(res, statusCode, false, message)
}

export const sendSuccess = (res: Response, data: any, message = "Success") => {
  return sendResponse(res, 200, true, message, data)
}
