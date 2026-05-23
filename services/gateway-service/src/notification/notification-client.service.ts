import { Injectable } from '@nestjs/common';
import * as grpc from '@grpc/grpc-js';

import {
  GetUserNotificationsRequest,
  NotificationServiceClient,
  UserNotificationsResponse,
} from '../generated/notification/notification';

import { grpcUnaryCall } from '../common/utils/grpc.util';

@Injectable()
export class NotificationGrpcService {
  private readonly client: NotificationServiceClient;

  constructor() {
    this.client =
      new NotificationServiceClient(
        process.env.NOTIFICATION_GRPC_ADDR ??
          'notification-service:50053',
        grpc.credentials.createInsecure(),
      );
  }

  getUserNotifications(
    userId: number,
    page = 1,
    limit = 10,
  ): Promise<UserNotificationsResponse> {
    const request: GetUserNotificationsRequest =
      {
        userId,
        page,
        limit,
      };

    return grpcUnaryCall(
      callback =>
        this.client.getUserNotifications(
          request,
          callback,
        ),
    );
  }
}