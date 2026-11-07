package repositories

import (
	"communication-service/internal/dto"
	"communication-service/internal/models"

	"github.com/google/uuid"
)

func avatarURL(url *string) string {
	if url == nil {
		return ""
	}

	return *url
}

func userDTO(user *models.CommunicationUser) dto.User {
	return dto.User{
		ID:          user.UserID,
		Username:    user.Username,
		DisplayName: user.DisplayName,
		AvatarURL:   avatarURL(user.AvatarURL),
		IsVerified:  user.IsVerified,
	}
}

func participantDTO(participant *models.ConversationParticipant) dto.Participant {
	return dto.Participant{
		ID:          participant.UserID,
		Username:    participant.User.Username,
		DisplayName: participant.User.DisplayName,
		AvatarURL:   avatarURL(participant.User.AvatarURL),
		IsVerified:  participant.User.IsVerified,
		Role:        participant.Role,
		JoinedAt:    formatTime(participant.JoinedAt),
	}
}

func messageDTO(message *models.Message) *dto.Message {
	return &dto.Message{
		ID:             message.ID,
		ConversationID: message.ConversationID,
		SenderID:       message.SenderID,
		Content:        message.Content,
		CreatedAt:      formatTime(message.CreatedAt),
		UpdatedAt:      formatTimePtr(message.UpdatedAt),
		Sender:         userDTO(&message.Sender),
	}
}

func messageReadReceiptDTO(receipt *models.MessageReadReceipt) dto.MessageReadReceipt {
	return dto.MessageReadReceipt{
		User:   userDTO(&receipt.User),
		ReadAt: formatTime(receipt.ReadAt),
	}
}

func conversationDTO(conversation *models.Conversation, unreadCount uint, userID uuid.UUID) *dto.Conversation {
	participants := make([]dto.Participant, 0, len(conversation.Participants))
	currentUserRole := ""

	for i := range conversation.Participants {
		participants = append(participants, participantDTO(&conversation.Participants[i]))

		if conversation.Participants[i].UserID == userID {
			currentUserRole = conversation.Participants[i].Role
		}
	}

	name := ""
	if conversation.Name != nil {
		name = *conversation.Name
	}

	result := &dto.Conversation{
		ID:              conversation.ID,
		Type:            conversation.Type,
		Name:            name,
		AvatarURL:       avatarURL(conversation.AvatarURL),
		CreatedBy:       conversation.CreatedBy,
		CreatedAt:       formatTime(conversation.CreatedAt),
		UpdatedAt:       formatTimePtr(conversation.UpdatedAt),
		Participants:    participants,
		UnreadCount:     unreadCount,
		CurrentUserRole: currentUserRole,
	}

	if conversation.LastMessage != nil {
		result.LastMessage = messageDTO(conversation.LastMessage)
	}

	return result
}

func paginationResponse(totalCount int64, params dto.PaginationRequest) dto.PaginationResponse {
	totalPages := (uint(totalCount) + params.Limit - 1) / params.Limit

	return dto.PaginationResponse{
		TotalCount: uint(totalCount),
		Page:       params.Page,
		Limit:      params.Limit,
		TotalPages: totalPages,
	}
}
