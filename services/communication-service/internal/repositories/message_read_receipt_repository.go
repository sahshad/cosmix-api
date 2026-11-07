package repositories

import (
	"context"
	"time"

	"communication-service/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type MessageReadReceiptRepository struct {
	*BaseRepository[models.MessageReadReceipt]
}

func NewMessageReadReceiptRepository(
	db *gorm.DB,
) *MessageReadReceiptRepository {
	return &MessageReadReceiptRepository{
		NewBaseRepository[models.MessageReadReceipt](db),
	}
}

func (repo *MessageReadReceiptRepository) Create(ctx context.Context, receipt *models.MessageReadReceipt) error {
	return repo.db.
		WithContext(ctx).
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(receipt).Error
}

func (repo *MessageReadReceiptRepository) MarkConversationRead(ctx context.Context, conversationID uuid.UUID, userID uuid.UUID, readAt time.Time) error {
	return repo.Transaction(ctx, func(tx *gorm.DB) error {
		var participant models.ConversationParticipant

		if err := tx.
			WithContext(ctx).
			Where("conversation_id = ? AND user_id = ?", conversationID, userID).
			First(&participant).Error; err != nil {
			return err
		}

		if err := tx.
			WithContext(ctx).
			Exec(`
				INSERT INTO message_read_receipts (message_id, conversation_id, user_id, read_at)
				SELECT messages.id, messages.conversation_id, ?, ?
				FROM messages
				WHERE messages.conversation_id = ?
				  AND messages.sender_id <> ?
				  AND messages.created_at > COALESCE(?, '-infinity'::timestamptz)
				ON CONFLICT (message_id, user_id) DO NOTHING
			`, userID, readAt, conversationID, userID, participant.LastReadAt).Error; err != nil {
			return err
		}

		return tx.
			WithContext(ctx).
			Model(&models.ConversationParticipant{}).
			Where("conversation_id = ? AND user_id = ?", conversationID, userID).
			Update("last_read_at", readAt).Error
	})
}

type messageReadCountRow struct {
	MessageID uuid.UUID
	ReadCount int64
}

func (repo *MessageReadReceiptRepository) CountsForMessages(ctx context.Context, messageIDs []uuid.UUID) (map[uuid.UUID]uint, error) {
	counts := make(map[uuid.UUID]uint, len(messageIDs))
	if len(messageIDs) == 0 {
		return counts, nil
	}

	var rows []messageReadCountRow

	err := repo.db.
		WithContext(ctx).
		Model(&models.MessageReadReceipt{}).
		Select("message_id, COUNT(*) AS read_count").
		Where("message_id IN ?", messageIDs).
		Group("message_id").
		Scan(&rows).Error

	if err != nil {
		return nil, err
	}

	for _, row := range rows {
		counts[row.MessageID] = uint(row.ReadCount)
	}

	return counts, nil
}

func (repo *MessageReadReceiptRepository) ListForMessage(ctx context.Context, messageID uuid.UUID) ([]models.MessageReadReceipt, error) {
	var receipts []models.MessageReadReceipt

	err := repo.db.
		WithContext(ctx).
		Preload("User").
		Where("message_id = ?", messageID).
		Order("read_at ASC").
		Find(&receipts).Error

	if err != nil {
		return nil, err
	}

	return receipts, nil
}
