package repositories

import (
	"context"

	"communication-service/internal/dto"
	"communication-service/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MessageRepository struct {
	*BaseRepository[models.Message]
}

func NewMessageRepository(
	db *gorm.DB,
) *MessageRepository {
	return &MessageRepository{
		NewBaseRepository[models.Message](db),
	}
}

func (repo *MessageRepository) FindByIDInConversation(ctx context.Context, messageID uuid.UUID, conversationID uuid.UUID) (*models.Message, error) {
	var message models.Message

	err := repo.db.
		WithContext(ctx).
		Where("id = ? AND conversation_id = ?", messageID, conversationID).
		First(&message).Error

	if err != nil {
		return nil, err
	}

	return &message, nil
}

func (repo *MessageRepository) CreateWithConversationUpdate(ctx context.Context, message *models.Message) error {
	return repo.Transaction(ctx, func(tx *gorm.DB) error {
		if err := tx.Create(message).Error; err != nil {
			return err
		}

		return tx.
			Model(&models.Conversation{}).
			Where("id = ?", message.ConversationID).
			Updates(map[string]any{
				"last_message_id": message.ID,
				"last_message_at": message.CreatedAt,
				"updated_at":      message.CreatedAt,
			}).Error
	})
}

func (repo *MessageRepository) ListByConversation(ctx context.Context, conversationID uuid.UUID, params dto.PaginationRequest) (*dto.MessageListResponse, error) {
	var totalCount int64

	err := repo.db.
		WithContext(ctx).
		Model(&models.Message{}).
		Where("conversation_id = ?", conversationID).
		Count(&totalCount).Error

	if err != nil {
		return nil, err
	}

	var messages []models.Message

	err = repo.db.
		WithContext(ctx).
		Preload("Sender").
		Where("conversation_id = ?", conversationID).
		Order("created_at DESC").
		Limit(int(params.Limit)).
		Offset(int((params.Page - 1) * params.Limit)).
		Find(&messages).Error

	if err != nil {
		return nil, err
	}

	messageList := make([]dto.Message, 0, len(messages))
	for i := range messages {
		messageList = append(messageList, *messageDTO(&messages[i]))
	}

	return &dto.MessageListResponse{
		Messages:   messageList,
		Pagination: paginationResponse(totalCount, params),
	}, nil
}

func (repo *MessageRepository) UnreadTotal(ctx context.Context, userID uuid.UUID) (uint, error) {
	var count int64

	err := repo.db.
		WithContext(ctx).
		Model(&models.Message{}).
		Joins("JOIN conversation_participants cp ON cp.conversation_id = messages.conversation_id").
		Joins("LEFT JOIN message_read_receipts mrr ON mrr.message_id = messages.id AND mrr.user_id = ?", userID).
		Where("cp.user_id = ? AND messages.sender_id <> ?", userID, userID).
		Where("messages.created_at >= cp.joined_at AND mrr.message_id IS NULL").
		Count(&count).Error

	return uint(count), err
}
