import { Inject, Injectable } from '@nestjs/common';

import {
  NOTIFICATION_PACKAGE_NAME,
  NOTIFICATION_SERVICE_NAME,
  NotificationServiceClient,
  UserNotificationsResponse,
} from '../generated/notification/notification';

import type { ClientGrpc } from '@nestjs/microservices';
import { firstValueFrom } from 'rxjs';

@Injectable()
export class NotificationGrpcService {
  private client!: NotificationServiceClient;

  constructor(
    @Inject(NOTIFICATION_PACKAGE_NAME)
    private readonly grpcClient: ClientGrpc,
  ) { }

  onModuleInit() {
    this.client =
      this.grpcClient.getService<NotificationServiceClient>(
        NOTIFICATION_SERVICE_NAME,
      );
  }

  getUserNotifications(userId: number, page = 1, limit = 10): Promise<UserNotificationsResponse> {
    return firstValueFrom(this.client.getUserNotifications({ userId, page, limit }))
  }
}