package server

import (
	"context"

	"notification-service/internal/dto"
	"notification-service/internal/services"

	notificationpb "cosmix/shared/grpc/gen/go/notification"

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

func (s *NotificationServer) HealthCheck(
	ctx context.Context,
	req *notificationpb.HealthCheckRequest,
) (*notificationpb.HealthCheckResponse, error) {

	return &notificationpb.HealthCheckResponse{
		Status: "ok",
	}, nil
}

func (s *NotificationServer) GetUserNotifications(
	ctx context.Context,
	req *notificationpb.GetUserNotificationsRequest,
) (*notificationpb.UserNotificationsResponse, error) {

	result, err := s.notificationService.GetUserNotifications(
		ctx,
		uint(req.UserId),
		dto.PaginationRequest{
			Page:  uint(req.Page),
			Limit: uint(req.Limit),
		},
	)

	if err != nil {
		return nil, err
	}

	response :=
		&notificationpb.UserNotificationsResponse{
			Pagination:
				&notificationpb.Pagination{
					TotalCount:
						uint32(
							result.Pagination.TotalCount,
						),
					Page:
						uint32(
							result.Pagination.Page,
						),
					Limit:
						uint32(
							result.Pagination.Limit,
						),
					TotalPages:
						uint32(
							result.Pagination.TotalPages,
						),
				},
		}

	for _, item := range result.Notifications {

		response.Notifications =
			append(
				response.Notifications,
				mapNotification(item),
			)
	}

	return response, nil
}

func mapNotification(
	item dto.NotificationList,
) *notificationpb.Notification {

	var actorID uint64
	if item.ActorID != nil {
		actorID = uint64(*item.ActorID)
	}

	var entityID uint64
	if item.EntityID != nil {
		entityID = uint64(*item.EntityID)
	}

	var readAt *timestamppb.Timestamp
	if item.ReadAt != nil {
		readAt =
			timestamppb.New(
				*item.ReadAt,
			)
	}

	return &notificationpb.Notification{
		Id:               uint64(item.ID),
		UserId:           uint64(item.UserID),
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
		CreatedAt:        timestamppb.New(item.CreatedAt),
	}
}

// func stringValue(
// 	value *string,
// ) string {
// 	if value == nil {
// 		return ""
// 	}

// 	return *value
// }