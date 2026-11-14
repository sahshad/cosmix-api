import {
  ArgumentsHost,
  Catch,
  ExceptionFilter,
  HttpException,
  HttpStatus,
} from '@nestjs/common';
import type { Response } from 'express';

@Catch()
export class HttpExceptionFilter implements ExceptionFilter {
  catch(exception: unknown, host: ArgumentsHost) {
    const ctx = host.switchToHttp();
    const response = ctx.getResponse<Response>();

    const statusCode =
      exception instanceof HttpException ? exception.getStatus() : 500;

    const message =
      exception instanceof HttpException
        ? exception.message
        : 'Internal server error';

    response.locals.errorMessage =
      exception instanceof Error ? exception.message : message;

    if (!(exception instanceof HttpException) && exception instanceof Error) {
      response.locals.errorStack = exception.stack;
    }

    response.status(statusCode).json({
      statusCode,
      status: HttpStatus[statusCode],
      message,
    });
  }
}
