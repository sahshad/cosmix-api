import { Inject, Injectable } from '@nestjs/common';

import {
  COMMUNICATION_PACKAGE_NAME,
  COMMUNICATION_SERVICE_NAME,
  CommunicationServiceClient,
  ConversationListResponse,
  ConversationResponse,
  GetUnreadCountResponse,
  MarkConversationReadResponse,
  MessageListResponse,
  MessageResponse,
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

  getUnreadCount(userId: string): Promise<GetUnreadCountResponse> {
    return firstValueFrom(this.client.getUnreadCount({ userId }));
  }
}
