import { CanActivate, ExecutionContext, Injectable } from '@nestjs/common';
import { Observable } from 'rxjs';
import * as jwt from 'jsonwebtoken';
import { AuthJwtPayload } from '../types/jwt-payload.type';

@Injectable()
export class OptionalAuthGuard implements CanActivate {
  canActivate(
    context: ExecutionContext,
  ): boolean | Promise<boolean> | Observable<boolean> {
    const request =
      context.switchToHttp().getRequest();

    const header =
      request.headers.authorization;

    if (!header) {
      return true;
    }

    const token =
      header.split(' ')[1];

    if (!token) {
      return true;
    }

    try {
      const payload = jwt.verify(
        token,
        process.env.JWT_PUBLIC_KEY!,
      ) as AuthJwtPayload;

      request.user = payload;

      request.headers['x-user-id'] = payload.userId
    } catch {
      // invalid/expired token on a public route: proceed anonymously
    }

    return true;
  }
}
