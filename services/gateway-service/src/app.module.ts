import { MiddlewareConsumer, Module, NestModule, RequestMethod } from '@nestjs/common';
import { APP_GUARD } from '@nestjs/core';
import {
  ThrottlerGuard,
  ThrottlerModule,
} from '@nestjs/throttler';

import { AuthModule } from './auth/auth.module';
import { UserModule } from './user/user.module';
import { NotificationModule } from './notification/notification.module';

import { RequestIdMiddleware } from './common/middleware/request-id.middleware';
import { LoggerMiddleware } from './common/middleware/logger.middleware';
import { PostModule } from './post/post.module';

@Module({
  imports: [
    ThrottlerModule.forRoot([
      {
        ttl: 60000,
        limit: 100,
      },
    ]),
    AuthModule,
    UserModule,
    NotificationModule,
    PostModule,
  ],
  providers: [
    {
      provide: APP_GUARD,
      useClass: ThrottlerGuard,
    },
  ],
})
export class AppModule implements NestModule {
  configure(
    consumer: MiddlewareConsumer,
  ) {
    consumer
      .apply(
        RequestIdMiddleware,
        LoggerMiddleware,
      )
      .forRoutes({ path: '*path', method: RequestMethod.ALL });
  }
}