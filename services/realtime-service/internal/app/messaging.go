package app

import (
	"cosmix/shared/core/eventbus"
	"cosmix/shared/core/rabbitmq"
)

func RegisterSubscriptions(container *Container) error {
	if err := rabbitmq.DeclareExchanges(container.Rabbit.Channel); err != nil {
		return err
	}

	if err := eventbus.Subscribe(
		container.Rabbit.Channel,
		rabbitmq.ExchangeEvents,
		rabbitmq.RealtimePostLikedQueue,
		rabbitmq.PostLiked,
		container.EventService.HandlePostLiked,
	); err != nil {
		return err
	}

	if err := eventbus.Subscribe(
		container.Rabbit.Channel,
		rabbitmq.ExchangeEvents,
		rabbitmq.RealtimePostCommentedQueue,
		rabbitmq.CommentCreated,
		container.EventService.HandlePostCommented,
	); err != nil {
		return err
	}

	if err := eventbus.Subscribe(
		container.Rabbit.Channel,
		rabbitmq.ExchangeEvents,
		rabbitmq.RealtimeCommentRepliedQueue,
		rabbitmq.CommentReplied,
		container.EventService.HandleCommentReplied,
	); err != nil {
		return err
	}

	if err := eventbus.Subscribe(
		container.Rabbit.Channel,
		rabbitmq.ExchangeEvents,
		rabbitmq.RealtimeCommentLikedQueue,
		rabbitmq.CommentLiked,
		container.EventService.HandleCommentLiked,
	); err != nil {
		return err
	}

	if err := eventbus.Subscribe(
		container.Rabbit.Channel,
		rabbitmq.ExchangeEvents,
		rabbitmq.RealtimeUserFollowedQueue,
		rabbitmq.UserFollowed,
		container.EventService.HandleUserFollowed,
	); err != nil {
		return err
	}

	return nil
}
