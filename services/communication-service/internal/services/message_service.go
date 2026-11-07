package services

import (
	"context"
	stderrors "errors"
	"strings"
	"time"

	"communication-service/internal/dto"
	"communication-service/internal/models"
	"communication-service/internal/repositories"

	"cosmix/shared/core/errors"
	"cosmix/shared/core/eventbus"
	"cosmix/shared/core/rabbitmq"
	chatEvents "cosmix/shared/events/chat"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
	"gorm.io/gorm"
)

const maxMessageLength = 5000

type MessageService struct {
	messageRepo      *repositories.MessageRepository
	conversationRepo *repositories.ConversationRepository
	userRepo         *repositories.CommunicationUserRepository
	receiptRepo      *repositories.MessageReadReceiptRepository
	rabbitCh         *amqp.Channel
}

func NewMessageService(
	messageRepo *repositories.MessageRepository,
	conversationRepo *repositories.ConversationRepository,
	userRepo *repositories.CommunicationUserRepository,
	receiptRepo *repositories.MessageReadReceiptRepository,
	rabbitCh *amqp.Channel,
) *MessageService {
	return &MessageService{
		messageRepo:      messageRepo,
		conversationRepo: conversationRepo,
		userRepo:         userRepo,
		receiptRepo:      receiptRepo,
		rabbitCh:         rabbitCh,
	}
}

func (svc *MessageService) Send(ctx context.Context, conversationID uuid.UUID, senderID uuid.UUID, content string) (*dto.Message, error) {
	content = strings.TrimSpace(content)
	if content == "" {
		return nil, errors.NewValidation("content", "content is required")
	}

	if len(content) > maxMessageLength {
		return nil, errors.NewValidation("content", "content must be at most 5000 characters")
	}

	isParticipant, err := svc.conversationRepo.IsParticipant(ctx, conversationID, senderID)
	if err != nil {
		return nil, err
	}

	if !isParticipant {
		return nil, errors.NewNotFound("conversation")
	}

	sender, err := svc.userRepo.FindByUserID(ctx, senderID)
	if err != nil {
		return nil, errors.NewNotFound("user")
	}

	message := &models.Message{
		ConversationID: conversationID,
		SenderID:       senderID,
		Content:        content,
		CreatedAt:      time.Now().UTC(),
	}

	if err := svc.messageRepo.CreateWithConversationUpdate(ctx, message); err != nil {
		return nil, err
	}

	svc.publishMessageSent(ctx, message, sender)

	return &dto.Message{
		ID:             message.ID,
		ConversationID: message.ConversationID,
		SenderID:       message.SenderID,
		Content:        message.Content,
		CreatedAt:      message.CreatedAt.Format(time.RFC3339),
		Sender: dto.User{
			ID:          sender.UserID,
			Username:    sender.Username,
			DisplayName: sender.DisplayName,
			AvatarURL:   avatarURL(sender.AvatarURL),
			IsVerified:  sender.IsVerified,
		},
	}, nil
}

func (svc *MessageService) List(ctx context.Context, conversationID uuid.UUID, userID uuid.UUID, params dto.PaginationRequest) (*dto.MessageListResponse, error) {
	isParticipant, err := svc.conversationRepo.IsParticipant(ctx, conversationID, userID)
	if err != nil {
		return nil, err
	}

	if !isParticipant {
		return nil, errors.NewNotFound("conversation")
	}

	result, err := svc.messageRepo.ListByConversation(ctx, conversationID, params)
	if err != nil {
		return nil, err
	}

	messageIDs := make([]uuid.UUID, 0, len(result.Messages))
	for i := range result.Messages {
		messageIDs = append(messageIDs, result.Messages[i].ID)
	}

	readCounts, err := svc.receiptRepo.CountsForMessages(ctx, messageIDs)
	if err != nil {
		return nil, err
	}

	for i := range result.Messages {
		result.Messages[i].ReadCount = readCounts[result.Messages[i].ID]
	}

	return result, nil
}

func (svc *MessageService) MarkMessageRead(ctx context.Context, conversationID uuid.UUID, messageID uuid.UUID, userID uuid.UUID) error {
	isParticipant, err := svc.conversationRepo.IsParticipant(ctx, conversationID, userID)
	if err != nil {
		return err
	}

	if !isParticipant {
		return errors.NewNotFound("conversation")
	}

	message, err := svc.messageRepo.FindByIDInConversation(ctx, messageID, conversationID)
	if err != nil {
		if stderrors.Is(err, gorm.ErrRecordNotFound) {
			return errors.NewNotFound("message")
		}

		return err
	}

	if message.SenderID == userID {
		return nil
	}

	return svc.receiptRepo.Create(ctx, &models.MessageReadReceipt{
		MessageID:      message.ID,
		ConversationID: conversationID,
		UserID:         userID,
		ReadAt:         time.Now().UTC(),
	})
}

func (svc *MessageService) GetReadReceipts(ctx context.Context, conversationID uuid.UUID, messageID uuid.UUID, userID uuid.UUID) (*dto.MessageReadReceiptListResponse, error) {
	isParticipant, err := svc.conversationRepo.IsParticipant(ctx, conversationID, userID)
	if err != nil {
		return nil, err
	}

	if !isParticipant {
		return nil, errors.NewNotFound("conversation")
	}

	if _, err := svc.messageRepo.FindByIDInConversation(ctx, messageID, conversationID); err != nil {
		if stderrors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NewNotFound("message")
		}

		return nil, err
	}

	receipts, err := svc.receiptRepo.ListForMessage(ctx, messageID)
	if err != nil {
		return nil, err
	}

	result := make([]dto.MessageReadReceipt, 0, len(receipts))
	for i := range receipts {
		result = append(result, dto.MessageReadReceipt{
			User: dto.User{
				ID:          receipts[i].User.UserID,
				Username:    receipts[i].User.Username,
				DisplayName: receipts[i].User.DisplayName,
				AvatarURL:   avatarURL(receipts[i].User.AvatarURL),
				IsVerified:  receipts[i].User.IsVerified,
			},
			ReadAt: receipts[i].ReadAt.Format(time.RFC3339),
		})
	}

	return &dto.MessageReadReceiptListResponse{Receipts: result}, nil
}

func (svc *MessageService) UnreadCount(ctx context.Context, userID uuid.UUID) (uint, error) {
	return svc.messageRepo.UnreadTotal(ctx, userID)
}

func (svc *MessageService) publishMessageSent(ctx context.Context, message *models.Message, sender *models.CommunicationUser) {
	if svc.rabbitCh == nil {
		return
	}

	participantIDs, err := svc.conversationRepo.ParticipantIDs(ctx, message.ConversationID)
	if err != nil {
		return
	}

	recipientIDs := make([]uuid.UUID, 0, len(participantIDs))
	for _, participantID := range participantIDs {
		if participantID != message.SenderID {
			recipientIDs = append(recipientIDs, participantID)
		}
	}

	event := chatEvents.ChatMessageSent{
		EventVersion:      chatEvents.EventVersionOne,
		MessageID:         message.ID,
		ConversationID:    message.ConversationID,
		SenderID:          message.SenderID,
		SenderUsername:    sender.Username,
		SenderDisplayName: sender.DisplayName,
		SenderAvatarURL:   sender.AvatarURL,
		RecipientIDs:      recipientIDs,
		Content:           message.Content,
		CreatedAt:         message.CreatedAt,
	}

	_ = eventbus.Publish(ctx, svc.rabbitCh, rabbitmq.ExchangeEvents, rabbitmq.ChatMessageSent, event)
}

func avatarURL(url *string) string {
	if url == nil {
		return ""
	}

	return *url
}
