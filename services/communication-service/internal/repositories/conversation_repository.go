package repositories

import (
	"context"

	"communication-service/internal/dto"
	"communication-service/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ConversationRepository struct {
	*BaseRepository[models.Conversation]
}

func NewConversationRepository(
	db *gorm.DB,
) *ConversationRepository {
	return &ConversationRepository{
		NewBaseRepository[models.Conversation](db),
	}
}

func (repo *ConversationRepository) FindDirectByKey(ctx context.Context, directKey string) (*models.Conversation, error) {
	var conversation models.Conversation

	err := repo.db.
		WithContext(ctx).
		Where("direct_key = ?", directKey).
		First(&conversation).Error

	if err != nil {
		return nil, err
	}

	return &conversation, nil
}

func (repo *ConversationRepository) CreateDirect(ctx context.Context, conversation *models.Conversation, participantIDs []uuid.UUID) error {
	participants := make([]models.ConversationParticipant, 0, len(participantIDs))
	for _, participantID := range participantIDs {
		participants = append(participants, models.ConversationParticipant{
			UserID:   participantID,
			Role:     models.ParticipantRoleMember,
			JoinedAt: conversation.CreatedAt,
		})
	}

	return repo.CreateWithParticipants(ctx, conversation, participants)
}

func (repo *ConversationRepository) CreateWithParticipants(ctx context.Context, conversation *models.Conversation, participants []models.ConversationParticipant) error {
	return repo.Transaction(ctx, func(tx *gorm.DB) error {
		if err := tx.Create(conversation).Error; err != nil {
			return err
		}

		if len(participants) == 0 {
			return nil
		}

		for i := range participants {
			participants[i].ConversationID = conversation.ID
			if participants[i].JoinedAt.IsZero() {
				participants[i].JoinedAt = conversation.CreatedAt
			}
		}

		return tx.Create(&participants).Error
	})
}

func (repo *ConversationRepository) UpdateFields(ctx context.Context, conversationID uuid.UUID, fields map[string]any) error {
	return repo.db.
		WithContext(ctx).
		Model(&models.Conversation{}).
		Where("id = ?", conversationID).
		Updates(fields).Error
}

func (repo *ConversationRepository) GetForUser(ctx context.Context, conversationID uuid.UUID, userID uuid.UUID) (*dto.Conversation, error) {
	var conversation models.Conversation

	err := repo.db.
		WithContext(ctx).
		Select("conversations.*").
		Joins("JOIN conversation_participants cp ON cp.conversation_id = conversations.id").
		Where("conversations.id = ? AND cp.user_id = ?", conversationID, userID).
		Preload("Participants.User").
		Preload("LastMessage.Sender").
		First(&conversation).Error

	if err != nil {
		return nil, err
	}

	unreadCount, err := repo.unreadCount(ctx, conversationID, userID)
	if err != nil {
		return nil, err
	}

	return conversationDTO(&conversation, unreadCount, userID), nil
}

func (repo *ConversationRepository) ListForUser(ctx context.Context, userID uuid.UUID, params dto.PaginationRequest) (*dto.ConversationListResponse, error) {
	var totalCount int64

	err := repo.db.
		WithContext(ctx).
		Model(&models.Conversation{}).
		Joins("JOIN conversation_participants cp ON cp.conversation_id = conversations.id").
		Where("cp.user_id = ?", userID).
		Count(&totalCount).Error

	if err != nil {
		return nil, err
	}

	var conversations []models.Conversation

	err = repo.db.
		WithContext(ctx).
		Select("conversations.*").
		Joins("JOIN conversation_participants cp ON cp.conversation_id = conversations.id").
		Where("cp.user_id = ?", userID).
		Preload("Participants.User").
		Preload("LastMessage.Sender").
		Order("conversations.last_message_at DESC NULLS LAST").
		Order("conversations.created_at DESC").
		Limit(int(params.Limit)).
		Offset(int((params.Page - 1) * params.Limit)).
		Find(&conversations).Error

	if err != nil {
		return nil, err
	}

	unreadCounts, err := repo.unreadCountsForUser(ctx, userID, conversations)
	if err != nil {
		return nil, err
	}

	conversationList := make([]dto.Conversation, 0, len(conversations))
	for i := range conversations {
		conversationList = append(conversationList, *conversationDTO(&conversations[i], unreadCounts[conversations[i].ID], userID))
	}

	return &dto.ConversationListResponse{
		Conversations: conversationList,
		Pagination:    paginationResponse(totalCount, params),
	}, nil
}

func (repo *ConversationRepository) IsParticipant(ctx context.Context, conversationID uuid.UUID, userID uuid.UUID) (bool, error) {
	var count int64

	err := repo.db.
		WithContext(ctx).
		Model(&models.ConversationParticipant{}).
		Where("conversation_id = ? AND user_id = ?", conversationID, userID).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (repo *ConversationRepository) FindParticipant(ctx context.Context, conversationID uuid.UUID, userID uuid.UUID) (*models.ConversationParticipant, error) {
	var participant models.ConversationParticipant

	err := repo.db.
		WithContext(ctx).
		Where("conversation_id = ? AND user_id = ?", conversationID, userID).
		First(&participant).Error

	if err != nil {
		return nil, err
	}

	return &participant, nil
}

func (repo *ConversationRepository) ParticipantIDs(ctx context.Context, conversationID uuid.UUID) ([]uuid.UUID, error) {
	var participantIDs []uuid.UUID

	err := repo.db.
		WithContext(ctx).
		Model(&models.ConversationParticipant{}).
		Where("conversation_id = ?", conversationID).
		Pluck("user_id", &participantIDs).Error

	return participantIDs, err
}

func (repo *ConversationRepository) AddParticipants(ctx context.Context, participants []models.ConversationParticipant) error {
	if len(participants) == 0 {
		return nil
	}

	return repo.db.
		WithContext(ctx).
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(&participants).Error
}

func (repo *ConversationRepository) RemoveParticipant(ctx context.Context, conversationID uuid.UUID, participantID uuid.UUID) error {
	return repo.Transaction(ctx, func(tx *gorm.DB) error {
		result := tx.
			WithContext(ctx).
			Where("conversation_id = ? AND user_id = ?", conversationID, participantID).
			Delete(&models.ConversationParticipant{})

		if result.Error != nil {
			return result.Error
		}

		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}

		var remaining int64
		if err := tx.
			WithContext(ctx).
			Model(&models.ConversationParticipant{}).
			Where("conversation_id = ?", conversationID).
			Count(&remaining).Error; err != nil {
			return err
		}

		if remaining == 0 {
			return tx.
				WithContext(ctx).
				Where("id = ?", conversationID).
				Delete(&models.Conversation{}).Error
		}

		var owners int64
		if err := tx.
			WithContext(ctx).
			Model(&models.ConversationParticipant{}).
			Where("conversation_id = ? AND role = ?", conversationID, models.ParticipantRoleOwner).
			Count(&owners).Error; err != nil {
			return err
		}

		if owners > 0 {
			return nil
		}

		var nextOwner models.ConversationParticipant
		if err := tx.
			WithContext(ctx).
			Where("conversation_id = ?", conversationID).
			Order("joined_at ASC").
			First(&nextOwner).Error; err != nil {
			return err
		}

		return tx.
			WithContext(ctx).
			Model(&models.ConversationParticipant{}).
			Where("conversation_id = ? AND user_id = ?", conversationID, nextOwner.UserID).
			Update("role", models.ParticipantRoleOwner).Error
	})
}

func (repo *ConversationRepository) unreadCount(ctx context.Context, conversationID uuid.UUID, userID uuid.UUID) (uint, error) {
	var count int64

	err := repo.db.
		WithContext(ctx).
		Model(&models.Message{}).
		Joins("JOIN conversation_participants cp ON cp.conversation_id = messages.conversation_id AND cp.user_id = ?", userID).
		Joins("LEFT JOIN message_read_receipts mrr ON mrr.message_id = messages.id AND mrr.user_id = ?", userID).
		Where("messages.conversation_id = ? AND messages.sender_id <> ?", conversationID, userID).
		Where("messages.created_at >= cp.joined_at AND mrr.message_id IS NULL").
		Count(&count).Error

	return uint(count), err
}

type unreadCountRow struct {
	ConversationID uuid.UUID
	UnreadCount    int64
}

func (repo *ConversationRepository) unreadCountsForUser(ctx context.Context, userID uuid.UUID, conversations []models.Conversation) (map[uuid.UUID]uint, error) {
	counts := make(map[uuid.UUID]uint, len(conversations))
	if len(conversations) == 0 {
		return counts, nil
	}

	conversationIDs := make([]uuid.UUID, 0, len(conversations))
	for i := range conversations {
		conversationIDs = append(conversationIDs, conversations[i].ID)
	}

	var rows []unreadCountRow

	err := repo.db.
		WithContext(ctx).
		Model(&models.Message{}).
		Select("messages.conversation_id, COUNT(*) AS unread_count").
		Joins("JOIN conversation_participants cp ON cp.conversation_id = messages.conversation_id AND cp.user_id = ?", userID).
		Joins("LEFT JOIN message_read_receipts mrr ON mrr.message_id = messages.id AND mrr.user_id = ?", userID).
		Where("messages.conversation_id IN ? AND messages.sender_id <> ?", conversationIDs, userID).
		Where("messages.created_at >= cp.joined_at AND mrr.message_id IS NULL").
		Group("messages.conversation_id").
		Scan(&rows).Error

	if err != nil {
		return nil, err
	}

	for _, row := range rows {
		counts[row.ConversationID] = uint(row.UnreadCount)
	}

	return counts, nil
}
