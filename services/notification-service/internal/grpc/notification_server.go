package grpc

import (
	"context"

	"notification-service/internal/dto"
	"notification-service/internal/services"

	notificationpb "cosmix/shared/grpc/gen/go/notification"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type NotificationServer struct {
	notificationpb.UnimplementedNotificationServiceServer
	notificationService *services.NotificationService
}

func NewNotificationServer(
	notificationService *services.NotificationService,
) *NotificationServer {
	return &NotificationServer{
		notificationService: notificationService,
	}
}

func (srv *NotificationServer) HealthCheck(ctx context.Context, req *notificationpb.HealthCheckRequest) (*notificationpb.HealthCheckResponse, error) {
	return &notificationpb.HealthCheckResponse{
		Status: "ok",
	}, nil
}

func (srv *NotificationServer) GetUserNotifications(ctx context.Context, req *notificationpb.GetUserNotificationsRequest) (*notificationpb.UserNotificationsResponse, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, err
	}

	result, err := srv.notificationService.GetUserNotifications(
		ctx,
		userID,
		dto.PaginationRequest{
			Page:  uint(req.Page),
			Limit: uint(req.Limit),
		},
	)
	if err != nil {
		return nil, err
	}

	response := &notificationpb.UserNotificationsResponse{
		Pagination: &notificationpb.Pagination{
			TotalCount: uint32(result.Pagination.TotalCount),
			Page:       uint32(result.Pagination.Page),
			Limit:      uint32(result.Pagination.Limit),
			TotalPages: uint32(result.Pagination.TotalPages),
		},
	}

	for _, item := range result.Notifications {
		response.Notifications = append(
			response.Notifications,
			mapNotification(item),
		)
	}

	return response, nil
}

func mapNotification(item dto.NotificationList) *notificationpb.Notification {
	var actorID string
	if item.ActorID != nil {
		actorID = item.ActorID.String()
	}

	var entityID string
	if item.EntityID != nil {
		entityID = item.EntityID.String()
	}

	var readAt *timestamppb.Timestamp
	if item.ReadAt != nil {
		readAt = timestamppb.New(*item.ReadAt)
	}

	return &notificationpb.Notification{
		Id:               item.ID.String(),
		UserId:           item.UserID.String(),
		ActorId:          &actorID,
		ActorUsername:    item.ActorUsername,
		ActorDisplayName: item.ActorDisplayName,
		ActorAvatarUrl:   item.ActorAvatarURL,
		Type:             item.Type,
		EntityId:         &entityID,
		EntityType:       item.EntityType,
		Title:            item.Title,
		Body:             item.Body,
		ImageUrl:         item.ImageURL,
		ActionUrl:        item.ActionURL,
		IsRead:           item.IsRead,
		ReadAt:           readAt,
		CreatedAt:        item.CreatedAt,
	}
}
