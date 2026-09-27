import { Inject, Injectable } from '@nestjs/common';

import {
  NOTIFICATION_PACKAGE_NAME,
  NOTIFICATION_SERVICE_NAME,
  NotificationServiceClient,
  UserNotificationsResponse,
  MarkAsReadResponse,
  MarkAllAsReadResponse,
  GetUnreadCountResponse,
} from '../generated/notification/notification';

import type { ClientGrpc } from '@nestjs/microservices';
import { firstValueFrom } from 'rxjs';

@Injectable()
export class NotificationGrpcService {
  private client!: NotificationServiceClient;

  constructor(
    @Inject(NOTIFICATION_PACKAGE_NAME)
    private readonly grpcClient: ClientGrpc,
  ) {}

  onModuleInit() {
    this.client = this.grpcClient.getService<NotificationServiceClient>(
      NOTIFICATION_SERVICE_NAME,
    );
  }

  getUserNotifications(
    userId: string,
    page = 1,
    limit = 10,
  ): Promise<UserNotificationsResponse> {
    return firstValueFrom(
      this.client.getUserNotifications({ userId, page, limit }),
    );
  }

  markAsRead(
    notificationId: string,
    userId: string,
  ): Promise<MarkAsReadResponse> {
    return firstValueFrom(this.client.markAsRead({ notificationId, userId }));
  }

  markAllAsRead(userId: string): Promise<MarkAllAsReadResponse> {
    return firstValueFrom(this.client.markAllAsRead({ userId }));
  }

  getUnreadCount(userId: string): Promise<GetUnreadCountResponse> {
    return firstValueFrom(this.client.getUnreadCount({ userId }));
  }
}
