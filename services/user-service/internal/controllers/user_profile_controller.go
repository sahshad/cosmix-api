package controllers

import (

	"user-service/internal/dto"
	"user-service/internal/services"
	"cosmix/shared/core/httpx"

	"github.com/gin-gonic/gin"
	ampqp "github.com/rabbitmq/amqp091-go"
)

type UserProfileController struct {
	service  services.UserProfileServiceInterface
	rabbitCh *ampqp.Channel
}

func NewUserProfileController(service services.UserProfileServiceInterface, rabbitCh *ampqp.Channel) *UserProfileController {
	return &UserProfileController{service: service, rabbitCh: rabbitCh}
}

func (ctrl *UserProfileController) HealthCheck(c *gin.Context) (interface{}, error) {
	return map[string]string{"message": "user service is ok"}, nil
}

func (ctrl *UserProfileController) GetMe(c *gin.Context) (interface{}, error) {
	return ctrl.GetMyProfile(c)
}

func (ctrl *UserProfileController) GetMyProfile(c *gin.Context) (interface{}, error) {
	userID, err := httpx.ParseUserIDHeader(c)
	if err != nil {
		return nil, err
	}

	ctx := c.Request.Context()

	profile, err := ctrl.service.GetProfile(ctx, uint(userID))
	if err != nil {
		return nil, err
	}

	return profile, nil
}

func (ctrl *UserProfileController) GetProfileByID(c *gin.Context) (interface{}, error) {
	id, err := httpx.ParseParamID(c, "id")
	if err != nil {
		return nil, err
	}

	ctx := c.Request.Context()

	profile, err := ctrl.service.GetProfile(ctx, uint(id))
	if err != nil {
		return nil, err
	}

	return profile, nil
}

func (ctrl *UserProfileController) UpdateMe(c *gin.Context) (interface{}, error){
	return ctrl.UpdateMyProfile(c)
}

func (ctrl *UserProfileController) UpdateMyProfile(c *gin.Context) (interface{}, error){
	userID, err := httpx.ParseUserIDHeader(c)
	if err != nil {
		return nil, err
	}

	var input dto.UpdateProfileDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		return nil, err
	}

	ctx := c.Request.Context()

	profile, err := ctrl.service.UpdateProfile(ctx, uint(userID), input)
	if err != nil {
		return nil, err
	}

	return profile, nil
}

func (ctrl *UserProfileController) GetByUsername(c *gin.Context) (interface{}, error){
	username, err := httpx.ParseQueryString(c, "username")
	if err != nil {
		return nil, err
	}

	ctx := c.Request.Context()

	profile, err := ctrl.service.GetProfileByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	return profile, nil
}
