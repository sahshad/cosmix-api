package services

import (
	"context"
	"errors"
	"strings"
	"time"

	"communication-service/internal/dto"
	"communication-service/internal/models"
	"communication-service/internal/repositories"

	apperrors "cosmix/shared/core/errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	maxGroupNameLength   = 100
	maxGroupParticipants = 100
)

type ConversationService struct {
	repo        *repositories.ConversationRepository
	userRepo    *repositories.CommunicationUserRepository
	receiptRepo *repositories.MessageReadReceiptRepository
}

func NewConversationService(
	repo *repositories.ConversationRepository,
	userRepo *repositories.CommunicationUserRepository,
	receiptRepo *repositories.MessageReadReceiptRepository,
) *ConversationService {
	return &ConversationService{
		repo:        repo,
		userRepo:    userRepo,
		receiptRepo: receiptRepo,
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

func (svc *ConversationService) CreateGroup(
	ctx context.Context,
	creatorID uuid.UUID,
	name string,
	participantIDs []uuid.UUID,
	avatarURL *string,
) (*dto.Conversation, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, apperrors.NewValidation("name", "name is required")
	}

	if len(name) > maxGroupNameLength {
		return nil, apperrors.NewValidation("name", "name must be at most 100 characters")
	}

	participantIDs = uniqueUserIDs(participantIDs, creatorID)
	if len(participantIDs) == 0 {
		return nil, apperrors.NewValidation("participantIds", "at least one participant is required")
	}

	if len(participantIDs)+1 > maxGroupParticipants {
		return nil, apperrors.NewValidation("participantIds", "a group can have at most 100 participants")
	}

	if _, err := svc.userRepo.FindByUserID(ctx, creatorID); err != nil {
		return nil, apperrors.NewNotFound("user")
	}

	users, err := svc.userRepo.FindByUserIDs(ctx, participantIDs)
	if err != nil {
		return nil, err
	}

	if len(users) != len(participantIDs) {
		return nil, apperrors.NewNotFound("user")
	}

	now := time.Now().UTC()
	conversation := &models.Conversation{
		Type:      models.ConversationTypeGroup,
		Name:      &name,
		AvatarURL: normalizedAvatarURL(avatarURL),
		CreatedBy: creatorID,
		CreatedAt: now,
	}

	participants := make([]models.ConversationParticipant, 0, len(participantIDs)+1)
	participants = append(participants, models.ConversationParticipant{
		UserID:   creatorID,
		Role:     models.ParticipantRoleOwner,
		JoinedAt: now,
	})

	for _, participantID := range participantIDs {
		participants = append(participants, models.ConversationParticipant{
			UserID:   participantID,
			Role:     models.ParticipantRoleMember,
			JoinedAt: now,
		})
	}

	if err := svc.repo.CreateWithParticipants(ctx, conversation, participants); err != nil {
		return nil, err
	}

	return svc.repo.GetForUser(ctx, conversation.ID, creatorID)
}

func (svc *ConversationService) UpdateGroup(
	ctx context.Context,
	conversationID uuid.UUID,
	requesterID uuid.UUID,
	name *string,
	avatarURL *string,
) (*dto.Conversation, error) {
	conversation, err := svc.findGroupConversation(ctx, conversationID)
	if err != nil {
		return nil, err
	}

	if err := svc.requireManager(ctx, conversation.ID, requesterID); err != nil {
		return nil, err
	}

	fields := make(map[string]any)

	if name != nil {
		trimmed := strings.TrimSpace(*name)
		if trimmed == "" {
			return nil, apperrors.NewValidation("name", "name is required")
		}

		if len(trimmed) > maxGroupNameLength {
			return nil, apperrors.NewValidation("name", "name must be at most 100 characters")
		}

		fields["name"] = trimmed
	}

	if avatarURL != nil {
		fields["avatar_url"] = normalizedAvatarURL(avatarURL)
	}

	if len(fields) > 0 {
		if err := svc.repo.UpdateFields(ctx, conversationID, fields); err != nil {
			return nil, err
		}
	}

	return svc.repo.GetForUser(ctx, conversationID, requesterID)
}

func (svc *ConversationService) AddParticipants(
	ctx context.Context,
	conversationID uuid.UUID,
	requesterID uuid.UUID,
	participantIDs []uuid.UUID,
) (*dto.Conversation, error) {
	conversation, err := svc.findGroupConversation(ctx, conversationID)
	if err != nil {
		return nil, err
	}

	if err := svc.requireManager(ctx, conversation.ID, requesterID); err != nil {
		return nil, err
	}

	participantIDs = uniqueUserIDs(participantIDs, requesterID)
	if len(participantIDs) == 0 {
		return nil, apperrors.NewValidation("participantIds", "at least one participant is required")
	}

	existingIDs, err := svc.repo.ParticipantIDs(ctx, conversationID)
	if err != nil {
		return nil, err
	}

	existing := make(map[uuid.UUID]struct{}, len(existingIDs))
	for _, existingID := range existingIDs {
		existing[existingID] = struct{}{}
	}

	newIDs := make([]uuid.UUID, 0, len(participantIDs))
	for _, participantID := range participantIDs {
		if _, ok := existing[participantID]; !ok {
			newIDs = append(newIDs, participantID)
		}
	}

	if len(newIDs) == 0 {
		return svc.repo.GetForUser(ctx, conversationID, requesterID)
	}

	if len(existingIDs)+len(newIDs) > maxGroupParticipants {
		return nil, apperrors.NewValidation("participantIds", "a group can have at most 100 participants")
	}

	users, err := svc.userRepo.FindByUserIDs(ctx, newIDs)
	if err != nil {
		return nil, err
	}

	if len(users) != len(newIDs) {
		return nil, apperrors.NewNotFound("user")
	}

	now := time.Now().UTC()
	participants := make([]models.ConversationParticipant, 0, len(newIDs))
	for _, participantID := range newIDs {
		participants = append(participants, models.ConversationParticipant{
			ConversationID: conversationID,
			UserID:         participantID,
			Role:           models.ParticipantRoleMember,
			JoinedAt:       now,
		})
	}

	if err := svc.repo.AddParticipants(ctx, participants); err != nil {
		return nil, err
	}

	return svc.repo.GetForUser(ctx, conversationID, requesterID)
}

func (svc *ConversationService) RemoveParticipant(
	ctx context.Context,
	conversationID uuid.UUID,
	requesterID uuid.UUID,
	participantID uuid.UUID,
) error {
	conversation, err := svc.findGroupConversation(ctx, conversationID)
	if err != nil {
		return err
	}

	requester, err := svc.repo.FindParticipant(ctx, conversation.ID, requesterID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperrors.NewNotFound("conversation")
		}

		return err
	}

	target, err := svc.repo.FindParticipant(ctx, conversation.ID, participantID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperrors.NewNotFound("participant")
		}

		return err
	}

	if requesterID != participantID {
		if requester.Role != models.ParticipantRoleOwner && requester.Role != models.ParticipantRoleAdmin {
			return apperrors.NewForbidden("you do not have permission to manage this group")
		}

		if target.Role == models.ParticipantRoleOwner {
			return apperrors.NewForbidden("the group owner cannot be removed")
		}
	}

	if err := svc.repo.RemoveParticipant(ctx, conversationID, participantID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperrors.NewNotFound("participant")
		}

		return err
	}

	return nil
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
	err := svc.receiptRepo.MarkConversationRead(ctx, conversationID, userID, time.Now().UTC())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperrors.NewNotFound("conversation")
		}

		return err
	}

	return nil
}

func (svc *ConversationService) findGroupConversation(ctx context.Context, conversationID uuid.UUID) (*models.Conversation, error) {
	conversation, err := svc.repo.FindByID(ctx, conversationID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewNotFound("conversation")
		}

		return nil, err
	}

	if conversation.Type != models.ConversationTypeGroup {
		return nil, apperrors.NewBadRequest("conversation", "only group conversations support this action")
	}

	return conversation, nil
}

func (svc *ConversationService) requireManager(ctx context.Context, conversationID uuid.UUID, userID uuid.UUID) error {
	participant, err := svc.repo.FindParticipant(ctx, conversationID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperrors.NewNotFound("conversation")
		}

		return err
	}

	if participant.Role != models.ParticipantRoleOwner && participant.Role != models.ParticipantRoleAdmin {
		return apperrors.NewForbidden("you do not have permission to manage this group")
	}

	return nil
}

func normalizedAvatarURL(avatarURL *string) *string {
	if avatarURL == nil {
		return nil
	}

	trimmed := strings.TrimSpace(*avatarURL)
	if trimmed == "" {
		return nil
	}

	return &trimmed
}

func uniqueUserIDs(userIDs []uuid.UUID, exclude uuid.UUID) []uuid.UUID {
	seen := make(map[uuid.UUID]struct{}, len(userIDs))
	result := make([]uuid.UUID, 0, len(userIDs))

	for _, userID := range userIDs {
		if userID == uuid.Nil || userID == exclude {
			continue
		}

		if _, ok := seen[userID]; ok {
			continue
		}

		seen[userID] = struct{}{}
		result = append(result, userID)
	}

	return result
}

func directKey(firstUserID uuid.UUID, secondUserID uuid.UUID) string {
	first := firstUserID.String()
	second := secondUserID.String()

	if first > second {
		first, second = second, first
	}

	return first + ":" + second
}
