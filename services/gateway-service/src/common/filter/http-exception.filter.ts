import {
  ArgumentsHost,
  Catch,
  ExceptionFilter,
  HttpException,
  HttpStatus,
} from '@nestjs/common';

@Catch()
export class HttpExceptionFilter
  implements ExceptionFilter {
  catch(
    exception: unknown,
    host: ArgumentsHost,
  ) {
    const ctx =
      host.switchToHttp();

    const response =
      ctx.getResponse();

    const statusCode =
      exception instanceof HttpException
        ? exception.getStatus()
        : 500;

    const message =
      exception instanceof HttpException
        ? exception.message
        : 'Internal server error';

    console.error(
      'API Gateway Error:',
      exception,
    );

    response.status(statusCode).json({
      statusCode,
      status: HttpStatus[statusCode],
      message,
    });
  }
}