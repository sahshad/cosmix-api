package services

import (
	"context"
	"errors"
	"time"

	"communication-service/internal/dto"
	"communication-service/internal/models"
	"communication-service/internal/repositories"

	apperrors "cosmix/shared/core/errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ConversationService struct {
	repo     *repositories.ConversationRepository
	userRepo *repositories.CommunicationUserRepository
}

func NewConversationService(
	repo *repositories.ConversationRepository,
	userRepo *repositories.CommunicationUserRepository,
) *ConversationService {
	return &ConversationService{
		repo:     repo,
		userRepo: userRepo,
	}
}

func (svc *ConversationService) CreateDirect(ctx context.Context, requesterID uuid.UUID, participantID uuid.UUID) (*dto.Conversation, error) {
	if requesterID == participantID {
		return nil, apperrors.NewBadRequest("participantId", "you cannot create a conversation with yourself")
	}

	if _, err := svc.userRepo.FindByUserID(ctx, requesterID); err != nil {
		return nil, apperrors.NewNotFound("user")
	}

	if _, err := svc.userRepo.FindByUserID(ctx, participantID); err != nil {
		return nil, apperrors.NewNotFound("user")
	}

	key := directKey(requesterID, participantID)

	existing, err := svc.repo.FindDirectByKey(ctx, key)
	if err == nil {
		return svc.repo.GetForUser(ctx, existing.ID, requesterID)
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	conversation := &models.Conversation{
		Type:      models.ConversationTypeDirect,
		DirectKey: &key,
		CreatedBy: requesterID,
		CreatedAt: time.Now().UTC(),
	}

	if err := svc.repo.CreateDirect(ctx, conversation, []uuid.UUID{requesterID, participantID}); err != nil {
		// Another request may have created the same direct conversation.
		existing, findErr := svc.repo.FindDirectByKey(ctx, key)
		if findErr != nil {
			return nil, err
		}

		return svc.repo.GetForUser(ctx, existing.ID, requesterID)
	}

	return svc.repo.GetForUser(ctx, conversation.ID, requesterID)
}

func (svc *ConversationService) List(ctx context.Context, userID uuid.UUID, params dto.PaginationRequest) (*dto.ConversationListResponse, error) {
	return svc.repo.ListForUser(ctx, userID, params)
}

func (svc *ConversationService) Get(ctx context.Context, conversationID uuid.UUID, userID uuid.UUID) (*dto.Conversation, error) {
	conversation, err := svc.repo.GetForUser(ctx, conversationID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewNotFound("conversation")
		}

		return nil, err
	}

	return conversation, nil
}

func (svc *ConversationService) MarkRead(ctx context.Context, conversationID uuid.UUID, userID uuid.UUID) error {
	err := svc.repo.MarkRead(ctx, conversationID, userID, time.Now().UTC())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperrors.NewNotFound("conversation")
		}

		return err
	}

	return nil
}

func directKey(firstUserID uuid.UUID, secondUserID uuid.UUID) string {
	first := firstUserID.String()
	second := secondUserID.String()

	if first > second {
		first, second = second, first
	}

	return first + ":" + second
}
