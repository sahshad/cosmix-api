import {
  Controller,
  Get,
  Query,
  UseGuards,
} from '@nestjs/common';

import { NotificationGrpcService } from './notification-client.service';

import { AuthGuard } from '../common/guard/auth.guard';
import { CurrentUser } from '../common/decorators/current-user.decorator';

import { GetUserNotificationsDto } from './dto/get-user-notifications.dto';

@Controller('notifications')
export class NotificationController {
  constructor(
    private readonly notificationGrpc:
      NotificationGrpcService,
  ) {}

  @UseGuards(AuthGuard)
  @Get('me')
  async getMyNotifications(
    @CurrentUser()
    user: { userId: number },

    @Query()
    query: GetUserNotificationsDto,
  ) {
    return this.notificationGrpc.getUserNotifications(
      user.userId,
      query.page ?? 1,
      query.limit ?? 10,
    );
  }
}