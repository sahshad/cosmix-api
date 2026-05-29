import { Module } from '@nestjs/common';
import { NotificationController } from './notification.controller';
import { NotificationGrpcService } from './notification-client.service';
import { NotificationService } from './notification.service';
import { ClientsModule, Transport } from '@nestjs/microservices';
import { NOTIFICATION_PACKAGE_NAME } from '../generated/notification/notification';
import { join } from 'path';

@Module({
  imports: [
    ClientsModule.register([
      {
        name: NOTIFICATION_PACKAGE_NAME,
        transport: Transport.GRPC,
        options: {
          url: process.env.NOTIFICATION_GRPC_ADDR ?? 'localhost:50054',
          package: 'notification',
          protoPath: join(
            __dirname,
            '../../../../shared/grpc/proto/notification/notification.proto',
          ),
        },
      },
    ]),
  ],
  controllers: [NotificationController],
  providers: [NotificationService, NotificationGrpcService],
  exports: [NotificationGrpcService],
})
export class NotificationModule { }
