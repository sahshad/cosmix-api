package middleware

import (
	apperrors "cosmix/shared/core/errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type HandlerFunc func(*gin.Context) (interface{}, error)

func ErrorHandler(handler HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {

		result, err := handler(c)
		if err != nil {
			if appErr, ok := err.(*apperrors.AppError); ok {
				c.JSON(appErr.Status, gin.H{
					"code":    appErr.Code,
					"message": appErr.Message,
					"field":   appErr.Field,
				})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    apperrors.CodeInternalServerError,
				"message": "internal server error",
			})
			return
		}
		c.JSON(http.StatusOK, result)
	}
}
