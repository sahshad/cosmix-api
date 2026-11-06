package grpc

import (
	"context"

	"communication-service/internal/dto"
	"communication-service/internal/services"

	apperrors "cosmix/shared/core/errors"
	communicationpb "cosmix/shared/grpc/gen/go/communication"

	"github.com/google/uuid"
)

type CommunicationServer struct {
	communicationpb.UnimplementedCommunicationServiceServer

	conversationService *services.ConversationService
	messageService      *services.MessageService
}

func NewCommunicationServer(
	conversationService *services.ConversationService,
	messageService *services.MessageService,
) *CommunicationServer {
	return &CommunicationServer{
		conversationService: conversationService,
		messageService:      messageService,
	}
}

func (srv *CommunicationServer) CreateDirectConversation(ctx context.Context, req *communicationpb.CreateDirectConversationRequest) (*communicationpb.ConversationResponse, error) {
	userID, err := parseUUID(req.UserId, "user_id")
	if err != nil {
		return nil, err
	}

	participantID, err := parseUUID(req.ParticipantId, "participant_id")
	if err != nil {
		return nil, err
	}

	conversation, err := srv.conversationService.CreateDirect(ctx, userID, participantID)
	if err != nil {
		return nil, err
	}

	return &communicationpb.ConversationResponse{
		Conversation: toPbConversation(conversation),
	}, nil
}

func (srv *CommunicationServer) GetConversations(ctx context.Context, req *communicationpb.GetConversationsRequest) (*communicationpb.ConversationListResponse, error) {
	userID, err := parseUUID(req.UserId, "user_id")
	if err != nil {
		return nil, err
	}

	params := normalizePagination(req.Page, req.Limit)

	result, err := srv.conversationService.List(ctx, userID, params)
	if err != nil {
		return nil, err
	}

	conversations := make([]*communicationpb.Conversation, 0, len(result.Conversations))
	for i := range result.Conversations {
		conversations = append(conversations, toPbConversation(&result.Conversations[i]))
	}

	return &communicationpb.ConversationListResponse{
		Conversations: conversations,
		Pagination:    toPbPagination(result.Pagination),
	}, nil
}

func (srv *CommunicationServer) GetConversation(ctx context.Context, req *communicationpb.GetConversationRequest) (*communicationpb.ConversationResponse, error) {
	conversationID, err := parseUUID(req.ConversationId, "conversation_id")
	if err != nil {
		return nil, err
	}

	userID, err := parseUUID(req.UserId, "user_id")
	if err != nil {
		return nil, err
	}

	conversation, err := srv.conversationService.Get(ctx, conversationID, userID)
	if err != nil {
		return nil, err
	}

	return &communicationpb.ConversationResponse{
		Conversation: toPbConversation(conversation),
	}, nil
}

func (srv *CommunicationServer) GetMessages(ctx context.Context, req *communicationpb.GetMessagesRequest) (*communicationpb.MessageListResponse, error) {
	conversationID, err := parseUUID(req.ConversationId, "conversation_id")
	if err != nil {
		return nil, err
	}

	userID, err := parseUUID(req.UserId, "user_id")
	if err != nil {
		return nil, err
	}

	params := normalizePagination(req.Page, req.Limit)

	result, err := srv.messageService.List(ctx, conversationID, userID, params)
	if err != nil {
		return nil, err
	}

	messages := make([]*communicationpb.Message, 0, len(result.Messages))
	for i := range result.Messages {
		messages = append(messages, toPbMessage(&result.Messages[i]))
	}

	return &communicationpb.MessageListResponse{
		Messages:   messages,
		Pagination: toPbPagination(result.Pagination),
	}, nil
}

func (srv *CommunicationServer) SendMessage(ctx context.Context, req *communicationpb.SendMessageRequest) (*communicationpb.MessageResponse, error) {
	conversationID, err := parseUUID(req.ConversationId, "conversation_id")
	if err != nil {
		return nil, err
	}

	senderID, err := parseUUID(req.SenderId, "sender_id")
	if err != nil {
		return nil, err
	}

	message, err := srv.messageService.Send(ctx, conversationID, senderID, req.Content)
	if err != nil {
		return nil, err
	}

	return &communicationpb.MessageResponse{
		Message: toPbMessage(message),
	}, nil
}

func (srv *CommunicationServer) MarkConversationRead(ctx context.Context, req *communicationpb.MarkConversationReadRequest) (*communicationpb.MarkConversationReadResponse, error) {
	conversationID, err := parseUUID(req.ConversationId, "conversation_id")
	if err != nil {
		return nil, err
	}

	userID, err := parseUUID(req.UserId, "user_id")
	if err != nil {
		return nil, err
	}

	if err := srv.conversationService.MarkRead(ctx, conversationID, userID); err != nil {
		return nil, err
	}

	return &communicationpb.MarkConversationReadResponse{
		Success: true,
	}, nil
}

func (srv *CommunicationServer) GetUnreadCount(ctx context.Context, req *communicationpb.GetUnreadCountRequest) (*communicationpb.GetUnreadCountResponse, error) {
	userID, err := parseUUID(req.UserId, "user_id")
	if err != nil {
		return nil, err
	}

	count, err := srv.messageService.UnreadCount(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &communicationpb.GetUnreadCountResponse{
		Count: uint32(count),
	}, nil
}

func parseUUID(value string, field string) (uuid.UUID, error) {
	id, err := uuid.Parse(value)
	if err != nil {
		return uuid.Nil, apperrors.NewBadRequest(field, "invalid "+field)
	}

	return id, nil
}

func normalizePagination(page uint32, limit uint32) dto.PaginationRequest {
	if page == 0 {
		page = 1
	}

	if limit == 0 {
		limit = 20
	}

	if limit > 100 {
		limit = 100
	}

	return dto.PaginationRequest{
		Page:  uint(page),
		Limit: uint(limit),
	}
}

func toPbUser(user *dto.User) *communicationpb.User {
	return &communicationpb.User{
		Id:          user.ID.String(),
		Username:    user.Username,
		DisplayName: user.DisplayName,
		AvatarUrl:   user.AvatarURL,
		IsVerified:  user.IsVerified,
	}
}

func toPbMessage(message *dto.Message) *communicationpb.Message {
	return &communicationpb.Message{
		Id:             message.ID.String(),
		ConversationId: message.ConversationID.String(),
		SenderId:       message.SenderID.String(),
		Content:        message.Content,
		CreatedAt:      message.CreatedAt,
		UpdatedAt:      message.UpdatedAt,
		Sender:         toPbUser(&message.Sender),
	}
}

func toPbConversation(conversation *dto.Conversation) *communicationpb.Conversation {
	participants := make([]*communicationpb.User, 0, len(conversation.Participants))
	for i := range conversation.Participants {
		participants = append(participants, toPbUser(&conversation.Participants[i]))
	}

	result := &communicationpb.Conversation{
		Id:           conversation.ID.String(),
		Type:         conversation.Type,
		CreatedAt:    conversation.CreatedAt,
		UpdatedAt:    conversation.UpdatedAt,
		Participants: participants,
		UnreadCount:  uint32(conversation.UnreadCount),
	}

	if conversation.LastMessage != nil {
		result.LastMessage = toPbMessage(conversation.LastMessage)
	}

	return result
}

func toPbPagination(pagination dto.PaginationResponse) *communicationpb.Pagination {
	return &communicationpb.Pagination{
		TotalCount: uint32(pagination.TotalCount),
		Page:       uint32(pagination.Page),
		Limit:      uint32(pagination.Limit),
		TotalPages: uint32(pagination.TotalPages),
	}
}
