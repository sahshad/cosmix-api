import { CanActivate, ExecutionContext, Injectable, UnauthorizedException } from '@nestjs/common';
import { Observable } from 'rxjs';
import * as jwt from 'jsonwebtoken';
import { AuthJwtPayload } from '../types/jwt-payload.type';

@Injectable()
export class AuthGuard implements CanActivate {
  canActivate(
    context: ExecutionContext,
  ): boolean | Promise<boolean> | Observable<boolean> {
    const request =
      context.switchToHttp().getRequest();

    const header =
      request.headers.authorization;

    if (!header) {
      throw new UnauthorizedException(
        'Authorization header missing',
      );
    }

    const token =
      header.split(' ')[1];

    if (!token) {
      throw new UnauthorizedException(
        'Invalid authorization header',
      );
    }

    try {
      const payload = jwt.verify(
        token,
        process.env.JWT_PUBLIC_KEY!,
      ) as AuthJwtPayload;

      request.user = payload;

      request.headers['x-user-id'] = payload.userId.toString()

      return true;
    } catch {
      throw new UnauthorizedException(
        'Invalid token',
      );
    }
  }
}
