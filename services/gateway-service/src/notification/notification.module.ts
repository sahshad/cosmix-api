import { Module } from '@nestjs/common';
import { NotificationController } from './notification.controller';
import { NotificationGrpcService } from './notification-client.service';
import { NotificationService } from './notification.service';

@Module({
  controllers: [NotificationController],
  providers: [NotificationGrpcService, NotificationService]
})
export class NotificationModule {}
