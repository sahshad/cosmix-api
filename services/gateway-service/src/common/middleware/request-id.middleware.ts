import { Injectable, NestMiddleware } from '@nestjs/common';
import { Request, Response, NextFunction } from 'express';
import { v4 as uuidv4 } from 'uuid';

@Injectable()
export class RequestIdMiddleware implements NestMiddleware {
  use(req: Request, res: Response, next: NextFunction) {
    const header = req.headers['x-request-id'];
    const id = (Array.isArray(header) ? header[0] : header) || uuidv4();

    req.headers['x-request-id'] = id;
    res.setHeader('x-request-id', id);

    next();
  }
}
