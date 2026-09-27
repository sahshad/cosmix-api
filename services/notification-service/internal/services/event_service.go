package services

import (
	"context"
	"strings"

	"notification-service/internal/email"
	"notification-service/internal/models"

	authEvents "cosmix/shared/events/auth"
	postEvents "cosmix/shared/events/post"
	userEvents "cosmix/shared/events/user"
)

const (
	maxContentLength = 20
	maxPostTitleLen  = 20
	maxPostLikeLen   = 40
)

type EventService struct {
	EmailLogSvc               *EmailLogService
	NotificationSvc           *NotificationService
	NotificationPreferenceSvc *NotificationPreferenceService
	NotificationUserSvc       *NotificationUserService
	MailDispatcher            *email.MailDispatcher
}

func NewEventService(
	emailLogSvc *EmailLogService,
	notificationSvc *NotificationService,
	notificationPreferenceSvc *NotificationPreferenceService,
	notificationUserSvc *NotificationUserService,
	mailDispatcher *email.MailDispatcher,
) *EventService {
	return &EventService{
		EmailLogSvc:               emailLogSvc,
		NotificationSvc:           notificationSvc,
		NotificationPreferenceSvc: notificationPreferenceSvc,
		NotificationUserSvc:       notificationUserSvc,
		MailDispatcher:            mailDispatcher,
	}
}

func (svc *EventService) HandleUserEmailVerificationCompleted(ctx context.Context, event authEvents.AuthUserEmailVerificationCompleted) error {
	notificationUser := &models.NotificationUser{
		UserID:      event.AuthUserID,
		Username:    event.Username,
		DisplayName: event.DisplayName,
		CreatedAt:   event.CreatedAt,
	}
	if err := svc.NotificationUserSvc.Create(ctx, notificationUser); err != nil {
		return err
	}

	if err := svc.NotificationPreferenceSvc.CreateDefault(ctx, event.AuthUserID); err != nil {
		return err
	}

	notification := &models.Notification{
		UserID: event.AuthUserID,
		Type:   "system.welcome",
		Title:  "Welcome to Cosmix",
		Body:   "Your account has been created successfully.",
	}

	if err := svc.NotificationSvc.Create(ctx, notification); err != nil {
		return err
	}

	appLink := "http://localhost:3000"

	templateData := map[string]any{
		"Title":       "Welcome to Cosmix",
		"DisplayName": event.DisplayName,
		"AppLink":     appLink,
	}

	mailDto := email.SendEmailDTO{
		To:           event.Email,
		Subject:      "Welcome to Cosmix",
		TemplateName: email.TemplateAuthWelcome,
		TemplateData: templateData,
	}

	err := svc.MailDispatcher.Send(mailDto)
	if err != nil {
		return err
	}

	return nil
}

func (svc *EventService) HandlePostLiked(ctx context.Context, event postEvents.PostLiked) error {
	if event.PostOwnerID == event.ActorID {
		return nil
	}

	actorName := event.ActorDisplayName
	if actorName == "" {
		actorName = event.ActorUsername
	}
	if actorName == "" {
		actorName = "Someone"
	}

	postTitle := truncate(event.PostTitle, maxPostLikeLen)

	notification := &models.Notification{
		UserID:           event.PostOwnerID,
		ActorID:          &event.ActorID,
		Type:             "post.like",
		Title:            actorName + " liked your post",
		Body:             "\"" + postTitle + "\"",
		EntityID:         &event.PostID,
		EntityType:       strPtr("post"),
		ActionURL:        strPtr("/posts/" + event.PostID.String()),
	}

	return svc.NotificationSvc.Create(ctx, notification)
}

func (svc *EventService) HandlePostCommented(ctx context.Context, event postEvents.PostCommented) error {
	if event.PostOwnerID == event.ActorID {
		return nil
	}

	actorName := event.ActorDisplayName
	if actorName == "" {
		actorName = event.ActorUsername
	}
	if actorName == "" {
		actorName = "Someone"
	}

	commentContent := truncate(event.Content, maxContentLength)
	postTitle := truncate(event.PostTitle, maxPostTitleLen)

	notification := &models.Notification{
		UserID:           event.PostOwnerID,
		ActorID:          &event.ActorID,
		Type:             "post.comment",
		Title:            actorName + " commented on your post",
		Body:             "\"" + commentContent + "\" on \"" + postTitle + "\"",
		EntityID:         &event.PostID,
		EntityType:       strPtr("post"),
		ActionURL:        strPtr("/posts/" + event.PostID.String()),
	}

	return svc.NotificationSvc.Create(ctx, notification)
}

func (svc *EventService) HandleCommentReplied(ctx context.Context, event postEvents.CommentReplied) error {
	if event.ParentCommentOwnerID == event.ActorID {
		return nil
	}

	actorName := event.ActorDisplayName
	if actorName == "" {
		actorName = event.ActorUsername
	}
	if actorName == "" {
		actorName = "Someone"
	}

	replyContent := truncate(event.Content, maxContentLength)
	postTitle := truncate(event.PostTitle, maxPostTitleLen)

	notification := &models.Notification{
		UserID:           event.ParentCommentOwnerID,
		ActorID:          &event.ActorID,
		Type:             "comment.reply",
		Title:            actorName + " replied to your comment",
		Body:             "\"" + replyContent + "\" on \"" + postTitle + "\"",
		EntityID:         &event.ReplyCommentID,
		EntityType:       strPtr("comment"),
		ActionURL:        strPtr("/posts/" + event.PostID.String()),
	}

	return svc.NotificationSvc.Create(ctx, notification)
}

func (svc *EventService) HandleCommentLiked(ctx context.Context, event postEvents.CommentLiked) error {
	if event.CommentOwnerID == event.ActorID {
		return nil
	}

	actorName := event.ActorDisplayName
	if actorName == "" {
		actorName = event.ActorUsername
	}
	if actorName == "" {
		actorName = "Someone"
	}

	commentContent := truncate(event.CommentContent, maxContentLength)
	postTitle := truncate(event.PostTitle, maxPostTitleLen)

	notification := &models.Notification{
		UserID:           event.CommentOwnerID,
		ActorID:          &event.ActorID,
		Type:             "comment.like",
		Title:            actorName + " liked your comment",
		Body:             "\"" + commentContent + "\" on \"" + postTitle + "\"",
		EntityID:         &event.CommentID,
		EntityType:       strPtr("comment"),
		ActionURL:        strPtr("/posts/" + event.PostID.String()),
	}

	return svc.NotificationSvc.Create(ctx, notification)
}

func (svc *EventService) HandleUserFollowed(ctx context.Context, event userEvents.UserFollowed) error {
	actorName := event.FollowerDisplayName
	if actorName == "" {
		actorName = event.FollowerUsername
	}
	if actorName == "" {
		actorName = "Someone"
	}

	notification := &models.Notification{
		UserID:           event.FollowingID,
		ActorID:          &event.FollowerID,
		Type:             "user.follow",
		Title:            actorName + " started following you",
		Body:             "",
		EntityID:         &event.FollowerID,
		EntityType:       strPtr("user"),
		ActionURL:        strPtr("/users/" + event.FollowerUsername),
	}

	return svc.NotificationSvc.Create(ctx, notification)
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return strings.TrimSpace(s[:maxLen]) + "..."
}

func strPtr(s string) *string {
	return &s
}
