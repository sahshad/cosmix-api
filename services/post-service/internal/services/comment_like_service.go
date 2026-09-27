package services

import (
	"context"
	"errors"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"

	"post-service/internal/models"
	"post-service/internal/repositories"

	"cosmix/shared/core/eventbus"
	postEvents "cosmix/shared/events/post"

	"github.com/google/uuid"
)

type CommentLikeService struct {
	repo         *repositories.CommentLikeRepository
	commentRepo  *repositories.CommentRepository
	postUserRepo *repositories.PostUserRepository
	rabbitCh     *amqp.Channel
}

func NewCommentLikeService(
	repo *repositories.CommentLikeRepository,
	commentRepo *repositories.CommentRepository,
	postUserRepo *repositories.PostUserRepository,
	rabbitCh *amqp.Channel,
) *CommentLikeService {
	return &CommentLikeService{
		repo:         repo,
		commentRepo:  commentRepo,
		postUserRepo: postUserRepo,
		rabbitCh:     rabbitCh,
	}
}

func (svc *CommentLikeService) LikeComment(ctx context.Context, commentID uuid.UUID, userID uuid.UUID) error {
	comment, err := svc.commentRepo.FindByID(ctx, commentID)
	if err != nil {
		return errors.New("comment not found")
	}

	exists, err := svc.repo.Exists(ctx, commentID, userID)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("already liked")
	}

	like := &models.CommentLike{
		CommentID: commentID,
		UserID:    userID,
	}
	if err := svc.repo.Create(ctx, like); err != nil {
		return err
	}

	if err := svc.commentRepo.IncrementLikesCount(ctx, commentID, 1); err != nil {
		return err
	}

	if svc.rabbitCh != nil && comment.UserID != userID {
		actor, _ := svc.postUserRepo.FindByUserID(ctx, userID)
		actorUsername := ""
		actorDisplayName := ""
		if actor != nil {
			actorUsername = actor.Username
			actorDisplayName = actor.DisplayName
		}

		post, _ := svc.commentRepo.GetPostByID(ctx, comment.PostID)
		postTitle := ""
		if post != nil {
			postTitle = post.Content
		}

		event := postEvents.CommentLiked{
			EventVersion:       postEvents.EventVersionOne,
			CommentID:          commentID,
			CommentOwnerID:     comment.UserID,
			PostID:             comment.PostID,
			ActorID:            userID,
			ActorUsername:      actorUsername,
			ActorDisplayName:   actorDisplayName,
			CommentContent:     comment.Content,
			PostTitle:          postTitle,
			LikedAt:            time.Now(),
		}
		eventbus.Publish(ctx, svc.rabbitCh, "cosmix.events", "comment.liked", event)
	}

	return nil
}

func (svc *CommentLikeService) UnlikeComment(ctx context.Context, commentID uuid.UUID, userID uuid.UUID) error {
	rows, err := svc.repo.Delete(ctx, commentID, userID)
	if err != nil {
		return err
	}
	if rows == 0 {
		return nil
	}

	return svc.commentRepo.IncrementLikesCount(ctx, commentID, -1)
}
