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

type LikeService struct {
	repo         *repositories.LikeRepository
	postRepo     *repositories.PostRepository
	postUserRepo *repositories.PostUserRepository
	rabbitCh     *amqp.Channel
}

func NewLikeService(
	repo *repositories.LikeRepository,
	postRepo *repositories.PostRepository,
	postUserRepo *repositories.PostUserRepository,
	rabbitCh *amqp.Channel,
) *LikeService {
	return &LikeService{
		repo:         repo,
		postRepo:     postRepo,
		postUserRepo: postUserRepo,
		rabbitCh:     rabbitCh,
	}
}

func (svc *LikeService) LikePost(ctx context.Context, postID uuid.UUID, userID uuid.UUID) error {
	post, err := svc.postRepo.FindByID(ctx, postID)
	if err != nil {
		return errors.New("post not found")
	}

	exists, err := svc.repo.Exists(ctx, postID, userID)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("already liked")
	}

	like := &models.Like{
		PostID: postID,
		UserID: userID,
	}
	if err := svc.repo.Create(ctx, like); err != nil {
		return err
	}

	if err := svc.postRepo.IncrementLikesCount(ctx, postID, 1); err != nil {
		return err
	}

	if svc.rabbitCh != nil && post.UserID != userID {
		actor, _ := svc.postUserRepo.FindByUserID(ctx, userID)
		actorUsername := ""
		actorDisplayName := ""
		if actor != nil {
			actorUsername = actor.Username
			actorDisplayName = actor.DisplayName
		}

		event := postEvents.PostLiked{
			EventVersion:     postEvents.EventVersionOne,
			PostID:           postID,
			PostOwnerID:      post.UserID,
			ActorID:          userID,
			ActorUsername:    actorUsername,
			ActorDisplayName: actorDisplayName,
			PostTitle:        post.Content, // Use post content as title (first 100 chars)
			LikedAt:          time.Now(),
		}
		eventbus.Publish(ctx, svc.rabbitCh, "cosmix.events", "post.liked", event)
	}

	return nil
}

func (svc *LikeService) UnlikePost(ctx context.Context, postID uuid.UUID, userID uuid.UUID) error {
	rows, err := svc.repo.Delete(ctx, postID, userID)
	if err != nil {
		return err
	}
	if rows == 0 {
		return nil
	}

	return svc.postRepo.IncrementLikesCount(ctx, postID, -1)
}

func (svc *LikeService) GetLikesCount(ctx context.Context, postID uuid.UUID) (int64, error) {
	return svc.repo.CountByPostID(ctx, postID)
}
