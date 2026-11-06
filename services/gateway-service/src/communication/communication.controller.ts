import {
  Body,
  Controller,
  Get,
  Param,
  Patch,
  Post,
  Query,
  UseGuards,
} from '@nestjs/common';

import { CommunicationGrpcService } from './communication-client.service';

import { AuthGuard } from '../common/guard/auth.guard';
import { CurrentUser } from '../common/decorators/current-user.decorator';

import { CreateDirectConversationDto } from './dto/create-direct-conversation.dto';
import { SendMessageDto } from './dto/send-message.dto';

@Controller('conversations')
export class CommunicationController {
  constructor(private readonly communicationGrpc: CommunicationGrpcService) {}

  @UseGuards(AuthGuard)
  @Post('direct')
  createDirectConversation(
    @CurrentUser() user: { userId: string },
    @Body() body: CreateDirectConversationDto,
  ) {
    return this.communicationGrpc.createDirectConversation(
      user.userId,
      body.participantId,
    );
  }

  @UseGuards(AuthGuard)
  @Get()
  getConversations(
    @CurrentUser() user: { userId: string },
    @Query('page') page = '1',
    @Query('limit') limit = '20',
  ) {
    return this.communicationGrpc.getConversations(
      user.userId,
      Number(page),
      Number(limit),
    );
  }

  @UseGuards(AuthGuard)
  @Get('unread-count')
  getUnreadCount(@CurrentUser() user: { userId: string }) {
    return this.communicationGrpc.getUnreadCount(user.userId);
  }

  @UseGuards(AuthGuard)
  @Get(':id')
  getConversation(
    @Param('id') id: string,
    @CurrentUser() user: { userId: string },
  ) {
    return this.communicationGrpc.getConversation(id, user.userId);
  }

  @UseGuards(AuthGuard)
  @Get(':id/messages')
  getMessages(
    @Param('id') id: string,
    @CurrentUser() user: { userId: string },
    @Query('page') page = '1',
    @Query('limit') limit = '20',
  ) {
    return this.communicationGrpc.getMessages(
      id,
      user.userId,
      Number(page),
      Number(limit),
    );
  }

  @UseGuards(AuthGuard)
  @Post(':id/messages')
  sendMessage(
    @Param('id') id: string,
    @CurrentUser() user: { userId: string },
    @Body() body: SendMessageDto,
  ) {
    return this.communicationGrpc.sendMessage(id, user.userId, body.content);
  }

  @UseGuards(AuthGuard)
  @Patch(':id/read')
  markConversationRead(
    @Param('id') id: string,
    @CurrentUser() user: { userId: string },
  ) {
    return this.communicationGrpc.markConversationRead(id, user.userId);
  }
}
