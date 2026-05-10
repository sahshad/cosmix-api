package controllers

import (
	"cosmix/shared/core/httpx"
	"notification-service/internal/dto"
	"notification-service/internal/services"

	"github.com/gin-gonic/gin"
)

type NotificationController struct {
	NotificationSvc services.NotificationServiceInterface
}

func NewNotificationController(notificationSvc services.NotificationServiceInterface) *NotificationController {
	return &NotificationController{
		NotificationSvc: notificationSvc,
	}
}

func (ctrl *NotificationController) HealthCheck(c *gin.Context) (interface{}, error) {
	return gin.H{
		"status": "ok",
	}, nil
}

func (ctrl *NotificationController) GetUserNotification(c *gin.Context) (interface{}, error) {
	page, err := httpx.ParseParamIDWithDefault(c, "page", 1)
	if err != nil {
		return nil, err
	}

	limit, err := httpx.ParseParamIDWithDefault(c, "limit", 10)
	if err != nil {
		return nil, err
	}

	userID, err := httpx.ParseUserIDHeader(c)
	if err != nil {
		return nil, err
	}

	ctx := c.Request.Context()
	reqParam := dto.PaginationRequest{
		Page:  uint(page),
		Limit: uint(limit),
	}

	userNotifications, err := ctrl.NotificationSvc.GetUserNotifications(ctx, uint(userID), reqParam)
	if err != nil {
		return nil, err
	}

	return userNotifications, nil
}
