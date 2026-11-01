package services

import (
	"context"
	"strings"

	"realtime-service/internal/realtime"

	postEvents "cosmix/shared/events/post"
	userEvents "cosmix/shared/events/user"

	"github.com/google/uuid"
)

const (
	maxContentLength = 20
	maxPostTitleLen  = 20
	maxPostLikeLen   = 40
)

type EventService struct {
	sender realtime.Sender
}

func NewEventService(sender realtime.Sender) *EventService {
	return &EventService{sender: sender}
}

func (svc *EventService) HandlePostLiked(ctx context.Context, event postEvents.PostLiked) error {
	if event.PostOwnerID == event.ActorID {
		return nil
	}

	actorName := nameOf(event.ActorDisplayName, event.ActorUsername)
	postTitle := truncate(event.PostTitle, maxPostLikeLen)

	svc.push(event.PostOwnerID, realtime.Notification{
		Type:       "post.like",
		Title:      actorName + " liked your post",
		Body:       "\"" + postTitle + "\"",
		EntityID:   event.PostID.String(),
		EntityType: "post",
		ActionURL:  "/posts/" + event.PostID.String(),
		Actor:      newActor(event.ActorID, event.ActorUsername, event.ActorDisplayName),
		CreatedAt:  event.LikedAt,
	})

	return nil
}

func (svc *EventService) HandlePostCommented(ctx context.Context, event postEvents.PostCommented) error {
	if event.PostOwnerID == event.ActorID {
		return nil
	}

	actorName := nameOf(event.ActorDisplayName, event.ActorUsername)
	commentContent := truncate(event.Content, maxContentLength)
	postTitle := truncate(event.PostTitle, maxPostTitleLen)

	svc.push(event.PostOwnerID, realtime.Notification{
		Type:       "post.comment",
		Title:      actorName + " commented on your post",
		Body:       "\"" + commentContent + "\" on \"" + postTitle + "\"",
		EntityID:   event.PostID.String(),
		EntityType: "post",
		ActionURL:  "/posts/" + event.PostID.String(),
		Actor:      newActor(event.ActorID, event.ActorUsername, event.ActorDisplayName),
		CreatedAt:  event.CreatedAt,
	})

	return nil
}

func (svc *EventService) HandleCommentReplied(ctx context.Context, event postEvents.CommentReplied) error {
	if event.ParentCommentOwnerID == event.ActorID {
		return nil
	}

	actorName := nameOf(event.ActorDisplayName, event.ActorUsername)
	replyContent := truncate(event.Content, maxContentLength)
	postTitle := truncate(event.PostTitle, maxPostTitleLen)

	svc.push(event.ParentCommentOwnerID, realtime.Notification{
		Type:       "comment.reply",
		Title:      actorName + " replied to your comment",
		Body:       "\"" + replyContent + "\" on \"" + postTitle + "\"",
		EntityID:   event.ReplyCommentID.String(),
		EntityType: "comment",
		ActionURL:  "/posts/" + event.PostID.String(),
		Actor:      newActor(event.ActorID, event.ActorUsername, event.ActorDisplayName),
		CreatedAt:  event.CreatedAt,
	})

	return nil
}

func (svc *EventService) HandleCommentLiked(ctx context.Context, event postEvents.CommentLiked) error {
	if event.CommentOwnerID == event.ActorID {
		return nil
	}

	actorName := nameOf(event.ActorDisplayName, event.ActorUsername)
	commentContent := truncate(event.CommentContent, maxContentLength)
	postTitle := truncate(event.PostTitle, maxPostTitleLen)

	svc.push(event.CommentOwnerID, realtime.Notification{
		Type:       "comment.like",
		Title:      actorName + " liked your comment",
		Body:       "\"" + commentContent + "\" on \"" + postTitle + "\"",
		EntityID:   event.CommentID.String(),
		EntityType: "comment",
		ActionURL:  "/posts/" + event.PostID.String(),
		Actor:      newActor(event.ActorID, event.ActorUsername, event.ActorDisplayName),
		CreatedAt:  event.LikedAt,
	})

	return nil
}

func (svc *EventService) HandleUserFollowed(ctx context.Context, event userEvents.UserFollowed) error {
	if event.FollowerID == event.FollowingID {
		return nil
	}

	actorName := nameOf(event.FollowerDisplayName, event.FollowerUsername)

	svc.push(event.FollowingID, realtime.Notification{
		Type:       "user.follow",
		Title:      actorName + " started following you",
		EntityID:   event.FollowerID.String(),
		EntityType: "user",
		ActionURL:  "/users/" + event.FollowerUsername,
		Actor:      newActor(event.FollowerID, event.FollowerUsername, event.FollowerDisplayName),
		CreatedAt:  event.CreatedAt,
	})

	return nil
}

func (svc *EventService) push(userID uuid.UUID, notification realtime.Notification) {
	svc.sender.SendToUser(userID, realtime.Message{
		Type: realtime.MessageTypeNotification,
		Data: notification,
	})
}

func newActor(id uuid.UUID, username, displayName string) *realtime.Actor {
	return &realtime.Actor{
		ID:          id.String(),
		Username:    username,
		DisplayName: displayName,
	}
}

func nameOf(displayName, username string) string {
	if displayName != "" {
		return displayName
	}

	if username != "" {
		return username
	}

	return "Someone"
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}

	return strings.TrimSpace(s[:maxLen]) + "..."
}
