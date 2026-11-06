package repositories

import (
	"communication-service/internal/dto"
	"communication-service/internal/models"
)

func userDTO(user *models.CommunicationUser) dto.User {
	avatarURL := ""
	if user.AvatarURL != nil {
		avatarURL = *user.AvatarURL
	}

	return dto.User{
		ID:          user.UserID,
		Username:    user.Username,
		DisplayName: user.DisplayName,
		AvatarURL:   avatarURL,
		IsVerified:  user.IsVerified,
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

func conversationDTO(conversation *models.Conversation, unreadCount uint) *dto.Conversation {
	participants := make([]dto.User, 0, len(conversation.Participants))
	for i := range conversation.Participants {
		participants = append(participants, userDTO(&conversation.Participants[i].User))
	}

	result := &dto.Conversation{
		ID:           conversation.ID,
		Type:         conversation.Type,
		CreatedAt:    formatTime(conversation.CreatedAt),
		UpdatedAt:    formatTimePtr(conversation.UpdatedAt),
		Participants: participants,
		UnreadCount:  unreadCount,
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
