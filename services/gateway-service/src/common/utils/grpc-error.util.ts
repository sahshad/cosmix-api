import {
  BadRequestException,
  ForbiddenException,
  InternalServerErrorException,
  NotFoundException,
  UnauthorizedException,
} from '@nestjs/common';

export function handleGrpcError(error: any): Error {
  switch (error.code) {
    case 3:
      return new BadRequestException(
        error.details || 'Invalid request',
      );

    case 5:
      return new NotFoundException(
        error.details || 'Resource not found',
      );

    case 7:
      return new ForbiddenException(
        error.details || 'Access denied',
      );

    case 16:
      return new UnauthorizedException(
        error.details || 'Unauthorized',
      );

    default:
      return new InternalServerErrorException(
        error.details || 'Internal server error',
      );
  }
}