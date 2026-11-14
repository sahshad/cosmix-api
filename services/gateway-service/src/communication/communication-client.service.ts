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

import { RpcException, type ClientGrpc } from '@nestjs/microservices';
import { catchError, firstValueFrom, Observable } from 'rxjs';

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

  private call<T>(observable: Observable<T>): Promise<T> {
    return firstValueFrom(
      observable.pipe(
        catchError((err: unknown) => {
          throw new RpcException(err as string | object);
        }),
      ),
    );
  }

  createDirectConversation(
    userId: string,
    participantId: string,
  ): Promise<ConversationResponse> {
    return this.call(
      this.client.createDirectConversation({ userId, participantId }),
    );
  }

  createGroupConversation(
    userId: string,
    name: string,
    participantIds: string[],
    avatarUrl?: string,
  ): Promise<ConversationResponse> {
    return this.call(
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
    return this.call(
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
    return this.call(
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
    return this.call(
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
    return this.call(this.client.getConversations({ userId, page, limit }));
  }

  getConversation(
    conversationId: string,
    userId: string,
  ): Promise<ConversationResponse> {
    return this.call(this.client.getConversation({ conversationId, userId }));
  }

  getMessages(
    conversationId: string,
    userId: string,
    page: number,
    limit: number,
  ): Promise<MessageListResponse> {
    return this.call(
      this.client.getMessages({ conversationId, userId, page, limit }),
    );
  }

  sendMessage(
    conversationId: string,
    senderId: string,
    content: string,
  ): Promise<MessageResponse> {
    return this.call(
      this.client.sendMessage({ conversationId, senderId, content }),
    );
  }

  markConversationRead(
    conversationId: string,
    userId: string,
  ): Promise<MarkConversationReadResponse> {
    return this.call(
      this.client.markConversationRead({ conversationId, userId }),
    );
  }

  markMessageRead(
    conversationId: string,
    messageId: string,
    userId: string,
  ): Promise<MarkMessageReadResponse> {
    return this.call(
      this.client.markMessageRead({ conversationId, messageId, userId }),
    );
  }

  getMessageReadReceipts(
    conversationId: string,
    messageId: string,
    userId: string,
  ): Promise<MessageReadReceiptListResponse> {
    return this.call(
      this.client.getMessageReadReceipts({
        conversationId,
        messageId,
        userId,
      }),
    );
  }

  getUnreadCount(userId: string): Promise<GetUnreadCountResponse> {
    return this.call(this.client.getUnreadCount({ userId }));
  }
}
