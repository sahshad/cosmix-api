import { Injectable, NestMiddleware } from '@nestjs/common';

import { Request, Response, NextFunction } from 'express';
import { logger } from '../logger/app-logger';

const SERVICE_BY_PREFIX: [string, string][] = [
  ['/api/auth', 'auth-service'],
  ['/api/users', 'user-service'],
  ['/api/posts', 'post-service'],
  ['/api/conversations', 'communication-service'],
  ['/api/notifications', 'notification-service'],
];

function resolveTargetService(path: string): string {
  const match = SERVICE_BY_PREFIX.find(
    ([prefix]) => path === prefix || path.startsWith(`${prefix}/`),
  );

  return match?.[1] ?? 'gateway-service';
}

@Injectable()
export class LoggerMiddleware implements NestMiddleware {
  use(req: Request, res: Response, next: NextFunction) {
    const start = Date.now();

    res.on('finish', () => {
      const requestId = req.headers['x-request-id'];
      const user = (req as Request & { user?: { userId?: string } }).user;
      const locals = res.locals as {
        errorMessage?: string;
        errorStack?: string;
      };

      const fields = {
        method: req.method,
        path: req.originalUrl,
        service: resolveTargetService(req.path),
        status: res.statusCode,
        request_id: Array.isArray(requestId) ? requestId[0] : requestId,
        user_id: user?.userId,
        duration_ms: Date.now() - start,
        error: locals.errorMessage,
        stack: res.statusCode >= 500 ? locals.errorStack : undefined,
      };

      if (res.statusCode >= 500) {
        logger.error('HTTP Request', fields);
      } else if (res.statusCode >= 400) {
        logger.warn('HTTP Request', fields);
      } else {
        logger.info('HTTP Request', fields);
      }
    });

    next();
  }
}
