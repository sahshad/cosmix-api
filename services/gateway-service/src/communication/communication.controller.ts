import {
  Body,
  Controller,
  Delete,
  Get,
  Param,
  Patch,
  Post,
  Query,
  UseGuards,
} from '@nestjs/common';
import {
  ApiBadRequestResponse,
  ApiBearerAuth,
  ApiCreatedResponse,
  ApiOkResponse,
  ApiOperation,
  ApiParam,
  ApiQuery,
  ApiTags,
  ApiUnauthorizedResponse,
} from '@nestjs/swagger';

import { CommunicationGrpcService } from './communication-client.service';

import { AuthGuard } from '../common/guard/auth.guard';
import { CurrentUser } from '../common/decorators/current-user.decorator';

import { AddGroupParticipantsDto } from './dto/add-group-participants.dto';
import { CreateDirectConversationDto } from './dto/create-direct-conversation.dto';
import { CreateGroupConversationDto } from './dto/create-group-conversation.dto';
import { SendMessageDto } from './dto/send-message.dto';
import { UpdateGroupConversationDto } from './dto/update-group-conversation.dto';

@ApiTags('conversations')
@ApiBearerAuth('access-token')
@Controller('conversations')
export class CommunicationController {
  constructor(private readonly communicationGrpc: CommunicationGrpcService) {}

  @ApiOperation({ summary: 'Create or get a direct conversation' })
  @ApiCreatedResponse({ description: 'Direct conversation' })
  @ApiBadRequestResponse({ description: 'Invalid participant' })
  @ApiUnauthorizedResponse({ description: 'Missing or invalid access token' })
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

  @ApiOperation({ summary: 'Create a group conversation' })
  @ApiCreatedResponse({ description: 'Group conversation' })
  @ApiBadRequestResponse({ description: 'Invalid group data' })
  @ApiUnauthorizedResponse({ description: 'Missing or invalid access token' })
  @UseGuards(AuthGuard)
  @Post('group')
  createGroupConversation(
    @CurrentUser() user: { userId: string },
    @Body() body: CreateGroupConversationDto,
  ) {
    return this.communicationGrpc.createGroupConversation(
      user.userId,
      body.name,
      body.participantIds,
      body.avatarUrl,
    );
  }

  @ApiOperation({ summary: 'List the conversations of the current user' })
  @ApiQuery({
    name: 'page',
    required: false,
    description: 'Page number (1-based)',
    example: 1,
    type: Number,
  })
  @ApiQuery({
    name: 'limit',
    required: false,
    description: 'Number of conversations per page',
    example: 20,
    type: Number,
  })
  @ApiOkResponse({ description: 'Paginated conversations' })
  @ApiUnauthorizedResponse({ description: 'Missing or invalid access token' })
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

  @ApiOperation({ summary: 'Get the unread message count of the current user' })
  @ApiOkResponse({ description: 'Unread message count' })
  @ApiUnauthorizedResponse({ description: 'Missing or invalid access token' })
  @UseGuards(AuthGuard)
  @Get('unread-count')
  getUnreadCount(@CurrentUser() user: { userId: string }) {
    return this.communicationGrpc.getUnreadCount(user.userId);
  }

  @ApiOperation({ summary: 'Update a group conversation' })
  @ApiParam({
    name: 'id',
    description: 'Conversation ID',
    example: '0f8fad5b-d9cb-469f-a165-70867728950e',
    format: 'uuid',
  })
  @ApiOkResponse({ description: 'Updated group conversation' })
  @ApiBadRequestResponse({ description: 'Invalid group data' })
  @ApiUnauthorizedResponse({ description: 'Missing or invalid access token' })
  @UseGuards(AuthGuard)
  @Patch(':id')
  updateGroupConversation(
    @Param('id') id: string,
    @CurrentUser() user: { userId: string },
    @Body() body: UpdateGroupConversationDto,
  ) {
    return this.communicationGrpc.updateGroupConversation(
      id,
      user.userId,
      body.name,
      body.avatarUrl,
    );
  }

  @ApiOperation({ summary: 'Add participants to a group conversation' })
  @ApiParam({
    name: 'id',
    description: 'Conversation ID',
    example: '0f8fad5b-d9cb-469f-a165-70867728950e',
    format: 'uuid',
  })
  @ApiCreatedResponse({ description: 'Participants added' })
  @ApiBadRequestResponse({ description: 'Invalid participant data' })
  @ApiUnauthorizedResponse({ description: 'Missing or invalid access token' })
  @UseGuards(AuthGuard)
  @Post(':id/participants')
  addGroupParticipants(
    @Param('id') id: string,
    @CurrentUser() user: { userId: string },
    @Body() body: AddGroupParticipantsDto,
  ) {
    return this.communicationGrpc.addGroupParticipants(
      id,
      user.userId,
      body.participantIds,
    );
  }

  @ApiOperation({ summary: 'Remove a participant from a group conversation' })
  @ApiParam({
    name: 'id',
    description: 'Conversation ID',
    example: '0f8fad5b-d9cb-469f-a165-70867728950e',
    format: 'uuid',
  })
  @ApiParam({
    name: 'participantId',
    description: 'User ID of the participant to remove',
    example: '7c9e6679-7425-40de-944b-e07fc1f90ae7',
    format: 'uuid',
  })
  @ApiOkResponse({ description: 'Participant removed' })
  @ApiUnauthorizedResponse({ description: 'Missing or invalid access token' })
  @UseGuards(AuthGuard)
  @Delete(':id/participants/:participantId')
  removeGroupParticipant(
    @Param('id') id: string,
    @Param('participantId') participantId: string,
    @CurrentUser() user: { userId: string },
  ) {
    return this.communicationGrpc.removeGroupParticipant(
      id,
      user.userId,
      participantId,
    );
  }

  @ApiOperation({
    summary: 'Leave a group conversation',
    description: 'Removes the current user from the group conversation.',
  })
  @ApiParam({
    name: 'id',
    description: 'Conversation ID',
    example: '0f8fad5b-d9cb-469f-a165-70867728950e',
    format: 'uuid',
  })
  @ApiCreatedResponse({ description: 'Left the group conversation' })
  @ApiUnauthorizedResponse({ description: 'Missing or invalid access token' })
  @UseGuards(AuthGuard)
  @Post(':id/leave')
  leaveGroupConversation(
    @Param('id') id: string,
    @CurrentUser() user: { userId: string },
  ) {
    return this.communicationGrpc.removeGroupParticipant(
      id,
      user.userId,
      user.userId,
    );
  }

  @ApiOperation({ summary: 'Get a conversation by ID' })
  @ApiParam({
    name: 'id',
    description: 'Conversation ID',
    example: '0f8fad5b-d9cb-469f-a165-70867728950e',
    format: 'uuid',
  })
  @ApiOkResponse({ description: 'Conversation details' })
  @ApiUnauthorizedResponse({ description: 'Missing or invalid access token' })
  @UseGuards(AuthGuard)
  @Get(':id')
  getConversation(
    @Param('id') id: string,
    @CurrentUser() user: { userId: string },
  ) {
    return this.communicationGrpc.getConversation(id, user.userId);
  }

  @ApiOperation({ summary: 'Get the paginated messages of a conversation' })
  @ApiParam({
    name: 'id',
    description: 'Conversation ID',
    example: '0f8fad5b-d9cb-469f-a165-70867728950e',
    format: 'uuid',
  })
  @ApiQuery({
    name: 'page',
    required: false,
    description: 'Page number (1-based)',
    example: 1,
    type: Number,
  })
  @ApiQuery({
    name: 'limit',
    required: false,
    description: 'Number of messages per page',
    example: 20,
    type: Number,
  })
  @ApiOkResponse({ description: 'Paginated messages' })
  @ApiUnauthorizedResponse({ description: 'Missing or invalid access token' })
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

  @ApiOperation({ summary: 'Send a message in a conversation' })
  @ApiParam({
    name: 'id',
    description: 'Conversation ID',
    example: '0f8fad5b-d9cb-469f-a165-70867728950e',
    format: 'uuid',
  })
  @ApiCreatedResponse({ description: 'Created message' })
  @ApiBadRequestResponse({ description: 'Invalid message data' })
  @ApiUnauthorizedResponse({ description: 'Missing or invalid access token' })
  @UseGuards(AuthGuard)
  @Post(':id/messages')
  sendMessage(
    @Param('id') id: string,
    @CurrentUser() user: { userId: string },
    @Body() body: SendMessageDto,
  ) {
    return this.communicationGrpc.sendMessage(id, user.userId, body.content);
  }

  @ApiOperation({ summary: 'Mark all messages in a conversation as read' })
  @ApiParam({
    name: 'id',
    description: 'Conversation ID',
    example: '0f8fad5b-d9cb-469f-a165-70867728950e',
    format: 'uuid',
  })
  @ApiOkResponse({ description: 'Conversation marked as read' })
  @ApiUnauthorizedResponse({ description: 'Missing or invalid access token' })
  @UseGuards(AuthGuard)
  @Patch(':id/read')
  markConversationRead(
    @Param('id') id: string,
    @CurrentUser() user: { userId: string },
  ) {
    return this.communicationGrpc.markConversationRead(id, user.userId);
  }

  @ApiOperation({ summary: 'Mark a single message as read' })
  @ApiParam({
    name: 'id',
    description: 'Conversation ID',
    example: '0f8fad5b-d9cb-469f-a165-70867728950e',
    format: 'uuid',
  })
  @ApiParam({
    name: 'messageId',
    description: 'Message ID',
    example: '7c9e6679-7425-40de-944b-e07fc1f90ae7',
    format: 'uuid',
  })
  @ApiCreatedResponse({ description: 'Message marked as read' })
  @ApiUnauthorizedResponse({ description: 'Missing or invalid access token' })
  @UseGuards(AuthGuard)
  @Post(':id/messages/:messageId/read')
  markMessageRead(
    @Param('id') id: string,
    @Param('messageId') messageId: string,
    @CurrentUser() user: { userId: string },
  ) {
    return this.communicationGrpc.markMessageRead(id, messageId, user.userId);
  }

  @ApiOperation({ summary: 'Get the read receipts of a message' })
  @ApiParam({
    name: 'id',
    description: 'Conversation ID',
    example: '0f8fad5b-d9cb-469f-a165-70867728950e',
    format: 'uuid',
  })
  @ApiParam({
    name: 'messageId',
    description: 'Message ID',
    example: '7c9e6679-7425-40de-944b-e07fc1f90ae7',
    format: 'uuid',
  })
  @ApiOkResponse({ description: 'Message read receipts' })
  @ApiUnauthorizedResponse({ description: 'Missing or invalid access token' })
  @UseGuards(AuthGuard)
  @Get(':id/messages/:messageId/receipts')
  getMessageReadReceipts(
    @Param('id') id: string,
    @Param('messageId') messageId: string,
    @CurrentUser() user: { userId: string },
  ) {
    return this.communicationGrpc.getMessageReadReceipts(
      id,
      messageId,
      user.userId,
    );
  }
}
