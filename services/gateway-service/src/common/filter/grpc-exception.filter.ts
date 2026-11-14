import {
  Catch,
  ExceptionFilter,
  ArgumentsHost,
  HttpStatus,
} from '@nestjs/common';
import { RpcException } from '@nestjs/microservices';
import type { Response } from 'express';
import * as grpc from '@grpc/grpc-js';

interface GrpcError {
  code?: number;
  details?: string;
  message?: string;
}

@Catch(RpcException)
export class GrpcExceptionFilter implements ExceptionFilter {
  catch(exception: RpcException, host: ArgumentsHost) {
    const error = exception.getError() as GrpcError;
    const ctx = host.switchToHttp();
    const response = ctx.getResponse<Response>();

    const grpcToHttp: Record<number, HttpStatus> = {
      [grpc.status.NOT_FOUND]: HttpStatus.NOT_FOUND,
      [grpc.status.ALREADY_EXISTS]: HttpStatus.CONFLICT,
      [grpc.status.INVALID_ARGUMENT]: HttpStatus.BAD_REQUEST,
      [grpc.status.UNAUTHENTICATED]: HttpStatus.UNAUTHORIZED,
      [grpc.status.PERMISSION_DENIED]: HttpStatus.FORBIDDEN,
      [grpc.status.RESOURCE_EXHAUSTED]: HttpStatus.TOO_MANY_REQUESTS,
      [grpc.status.UNAVAILABLE]: HttpStatus.SERVICE_UNAVAILABLE,
      [grpc.status.INTERNAL]: HttpStatus.INTERNAL_SERVER_ERROR,
    };

    const httpStatus =
      (error?.code !== undefined ? grpcToHttp[error.code] : undefined) ??
      HttpStatus.INTERNAL_SERVER_ERROR;
    const statusText = HttpStatus[httpStatus];
    const message = error?.details || error?.message || 'Internal server error';

    response.locals.errorMessage = message;

    response.status(httpStatus).json({
      statusCode: httpStatus,
      status: statusText,
      message,
    });
  }
}
