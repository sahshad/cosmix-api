import { Inject, Injectable } from '@nestjs/common';

import {
  COMMUNICATION_PACKAGE_NAME,
  COMMUNICATION_SERVICE_NAME,
  CommunicationServiceClient,
  ConversationListResponse,
  ConversationResponse,
  GetUnreadCountResponse,
  MarkConversationReadResponse,
  MarkMessageReadResponse,
  MessageListResponse,
  MessageReadReceiptListResponse,
  MessageResponse,
  RemoveGroupParticipantResponse,
} from '../generated/communication/communication';

import type { ClientGrpc } from '@nestjs/microservices';
import { firstValueFrom } from 'rxjs';

@Injectable()
export class CommunicationGrpcService {
  private client!: CommunicationServiceClient;

  constructor(
    @Inject(COMMUNICATION_PACKAGE_NAME)
    private readonly grpcClient: ClientGrpc,
  ) {}

  onModuleInit() {
    this.client = this.grpcClient.getService<CommunicationServiceClient>(
      COMMUNICATION_SERVICE_NAME,
    );
  }

  createDirectConversation(
    userId: string,
    participantId: string,
  ): Promise<ConversationResponse> {
    return firstValueFrom(
      this.client.createDirectConversation({ userId, participantId }),
    );
  }

  createGroupConversation(
    userId: string,
    name: string,
    participantIds: string[],
    avatarUrl?: string,
  ): Promise<ConversationResponse> {
    return firstValueFrom(
      this.client.createGroupConversation({
        userId,
        name,
        participantIds,
        avatarUrl,
      }),
    );
  }

  updateGroupConversation(
    conversationId: string,
    userId: string,
    name?: string,
    avatarUrl?: string,
  ): Promise<ConversationResponse> {
    return firstValueFrom(
      this.client.updateGroupConversation({
        conversationId,
        userId,
        name,
        avatarUrl,
      }),
    );
  }

  addGroupParticipants(
    conversationId: string,
    userId: string,
    participantIds: string[],
  ): Promise<ConversationResponse> {
    return firstValueFrom(
      this.client.addGroupParticipants({
        conversationId,
        userId,
        participantIds,
      }),
    );
  }

  removeGroupParticipant(
    conversationId: string,
    userId: string,
    participantId: string,
  ): Promise<RemoveGroupParticipantResponse> {
    return firstValueFrom(
      this.client.removeGroupParticipant({
        conversationId,
        userId,
        participantId,
      }),
    );
  }

  getConversations(
    userId: string,
    page: number,
    limit: number,
  ): Promise<ConversationListResponse> {
    return firstValueFrom(
      this.client.getConversations({ userId, page, limit }),
    );
  }

  getConversation(
    conversationId: string,
    userId: string,
  ): Promise<ConversationResponse> {
    return firstValueFrom(
      this.client.getConversation({ conversationId, userId }),
    );
  }

  getMessages(
    conversationId: string,
    userId: string,
    page: number,
    limit: number,
  ): Promise<MessageListResponse> {
    return firstValueFrom(
      this.client.getMessages({ conversationId, userId, page, limit }),
    );
  }

  sendMessage(
    conversationId: string,
    senderId: string,
    content: string,
  ): Promise<MessageResponse> {
    return firstValueFrom(
      this.client.sendMessage({ conversationId, senderId, content }),
    );
  }

  markConversationRead(
    conversationId: string,
    userId: string,
  ): Promise<MarkConversationReadResponse> {
    return firstValueFrom(
      this.client.markConversationRead({ conversationId, userId }),
    );
  }

  markMessageRead(
    conversationId: string,
    messageId: string,
    userId: string,
  ): Promise<MarkMessageReadResponse> {
    return firstValueFrom(
      this.client.markMessageRead({ conversationId, messageId, userId }),
    );
  }

  getMessageReadReceipts(
    conversationId: string,
    messageId: string,
    userId: string,
  ): Promise<MessageReadReceiptListResponse> {
    return firstValueFrom(
      this.client.getMessageReadReceipts({
        conversationId,
        messageId,
        userId,
      }),
    );
  }

  getUnreadCount(userId: string): Promise<GetUnreadCountResponse> {
    return firstValueFrom(this.client.getUnreadCount({ userId }));
  }
}
