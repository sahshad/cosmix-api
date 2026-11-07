import {
  Controller,
  Get,
  Param,
  Patch,
  Query,
  UseGuards,
} from '@nestjs/common';
import {
  ApiBearerAuth,
  ApiOkResponse,
  ApiOperation,
  ApiParam,
  ApiTags,
  ApiUnauthorizedResponse,
} from '@nestjs/swagger';

import { NotificationGrpcService } from './notification-client.service';

import { AuthGuard } from '../common/guard/auth.guard';
import { CurrentUser } from '../common/decorators/current-user.decorator';

import { GetUserNotificationsDto } from './dto/get-user-notifications.dto';

@ApiTags('notifications')
@ApiBearerAuth('access-token')
@Controller('notifications')
export class NotificationController {
  constructor(private readonly notificationGrpc: NotificationGrpcService) {}

  @ApiOperation({ summary: 'List the notifications of the current user' })
  @ApiOkResponse({ description: 'Paginated notifications' })
  @ApiUnauthorizedResponse({ description: 'Missing or invalid access token' })
  @UseGuards(AuthGuard)
  @Get('me')
  async getMyNotifications(
    @CurrentUser()
    user: { userId: string },

    @Query()
    query: GetUserNotificationsDto,
  ) {
    return this.notificationGrpc.getUserNotifications(
      user.userId,
      query.page ?? 1,
      query.limit ?? 10,
    );
  }

  @ApiOperation({
    summary: 'Get the unread notification count of the current user',
  })
  @ApiOkResponse({ description: 'Unread notification count' })
  @ApiUnauthorizedResponse({ description: 'Missing or invalid access token' })
  @UseGuards(AuthGuard)
  @Get('unread-count')
  async getUnreadCount(
    @CurrentUser()
    user: {
      userId: string;
    },
  ) {
    return this.notificationGrpc.getUnreadCount(user.userId);
  }

  @ApiOperation({ summary: 'Mark a notification as read' })
  @ApiParam({
    name: 'id',
    description: 'Notification ID',
    example: '0f8fad5b-d9cb-469f-a165-70867728950e',
    format: 'uuid',
  })
  @ApiOkResponse({ description: 'Notification marked as read' })
  @ApiUnauthorizedResponse({ description: 'Missing or invalid access token' })
  @UseGuards(AuthGuard)
  @Patch(':id/read')
  async markAsRead(
    @CurrentUser()
    user: { userId: string },

    @Param('id') notificationId: string,
  ) {
    return this.notificationGrpc.markAsRead(notificationId, user.userId);
  }

  @ApiOperation({ summary: 'Mark all notifications as read' })
  @ApiOkResponse({ description: 'All notifications marked as read' })
  @ApiUnauthorizedResponse({ description: 'Missing or invalid access token' })
  @UseGuards(AuthGuard)
  @Patch('read-all')
  async markAllAsRead(
    @CurrentUser()
    user: {
      userId: string;
    },
  ) {
    return this.notificationGrpc.markAllAsRead(user.userId);
  }
}
