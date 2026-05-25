package controllers

// import (
// 	"os"
// 	"time"

// 	"auth-service/internal/dto"
// 	publisher "auth-service/internal/messaging/publisher"
// 	"auth-service/internal/services"
// 	"auth-service/internal/utils"
// 	apperrors "cosmix/shared/core/errors"

// 	authEvents "cosmix/shared/events/auth"

// 	"github.com/gin-gonic/gin"
// 	amqp "github.com/rabbitmq/amqp091-go"
// )

// type AuthController struct {
// 	authService services.AuthServiceInterface
// 	rabbitCh    *amqp.Channel
// }

// func NewAuthController(authService services.AuthServiceInterface, rabbitCh *amqp.Channel) *AuthController {
// 	return &AuthController{
// 		authService: authService,
// 		rabbitCh:    rabbitCh,
// 	}
// }

// func (ctrl *AuthController) HealthCheck(c *gin.Context) (interface{}, error) {
// 	return map[string]interface{}{"message": "auth service is ok"}, nil
// }

// func (ctrl *AuthController) Register(c *gin.Context) (interface{}, error) {
// 	var registerDTO dto.RegisterDTO
// 	if err := c.ShouldBindJSON(&registerDTO); err != nil {
// 		return nil, err
// 	}

// 	ctx := c.Request.Context()

// 	user, err := ctrl.authService.Register(ctx, registerDTO)
// 	if err != nil {
// 		return nil, err
// 	}

// 	username := utils.GenerateUsername(registerDTO.DisplayName)

// 	// Publish user registered event
// 	publisher.PublishAuthUserRegistered(ctrl.rabbitCh, authEvents.AuthUserRegistered{
// 		EventVersion: authEvents.EventVersionOne,
// 		AuthUserID:   user.ID,
// 		Email:        user.Email,
// 		Username:     username,
// 		DisplayName:  registerDTO.DisplayName,
// 		CreatedAt:    time.Now().UTC(),
// 	})

// 	return map[string]interface{}{"id": user.ID}, nil
// }

// func (ctrl *AuthController) Login(c *gin.Context) (interface{}, error) {
// 	var input dto.LoginDTO
// 	if err := c.ShouldBindJSON(&input); err != nil {
// 		return nil, err
// 	}

// 	ctx := c.Request.Context()

// 	loginResponse, err := ctrl.authService.Login(ctx, input)
// 	if err != nil {
// 		return nil, err
// 	}

// 	secure := false
// 	domain := ""
// 	if os.Getenv("ENV") == "production" {
// 		secure = true
// 		domain = os.Getenv("COOKIE_DOMAIN")
// 	}

// 	c.SetCookie("refresh_token", loginResponse.RefreshToken, 60*60*24*30, "/", domain, secure, true)

// 	return map[string]interface{}{
// 		"access_token": loginResponse.AccessToken,
// 		"user":         loginResponse.AuthUser,
// 	}, nil
// }

// func (ctrl *AuthController) Refresh(c *gin.Context) (interface{}, error) {
// 	rt, err := c.Cookie("refresh_token")
// 	if err != nil || rt == "" {
// 		return nil, apperrors.NewUnauthorized("no refresh token")
// 	}

// 	ctx := c.Request.Context()

// 	refreshResponse, err := ctrl.authService.Refresh(ctx, rt)
// 	if err != nil {
// 		return nil, err
// 	}

// 	secure := false
// 	domain := ""
// 	if os.Getenv("ENV") == "production" {
// 		secure = true
// 		domain = os.Getenv("COOKIE_DOMAIN")
// 	}
// 	c.SetCookie("refresh_token", refreshResponse.RefreshToken, 60*60*24*30, "/", domain, secure, true)

// 	return map[string]interface{}{
// 		"access_token": refreshResponse.AccessToken,
// 	}, nil
// }

// func (ctrl *AuthController) Logout(c *gin.Context) (interface{}, error) {
// 	c.SetCookie("refresh_token", "", -1, "/", "", false, true)
// 	return map[string]interface{}{"message": "logged out"}, nil
// }

// func (ctrl *AuthController) UpdateUserPassword(c *gin.Context) (interface{}, error) {
// 	var input dto.UpdateUserPasswordDTO
// 	if err := c.ShouldBindJSON(&input); err != nil {
// 		return nil, err
// 	}

// 	ctx := c.Request.Context()

// 	err := ctrl.authService.UpdateUserPassword(ctx, input.UserID, input.NewPassword)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return map[string]interface{}{"message": "password updated successfully"}, nil
// }
