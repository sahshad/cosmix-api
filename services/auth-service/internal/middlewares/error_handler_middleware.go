package middlewares

import (
	"net/http"
	"auth-service/internal/apperrors"

	"github.com/gin-gonic/gin"
)

func ErrorHandler(handler func(*gin.Context) (interface{}, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		result, err := handler(c)
		if err != nil {
			switch e := err.(type) {
			case *apperrors.BadRequestError:
				c.JSON(http.StatusBadRequest, gin.H{"error": e.Error()})
			case *apperrors.UnauthorizedError:
				c.JSON(http.StatusUnauthorized, gin.H{"error": e.Error()})
			case *apperrors.NotFoundError:
				c.JSON(http.StatusNotFound, gin.H{"error": e.Error()})
			case *apperrors.InternalServerError:
				c.JSON(http.StatusInternalServerError, gin.H{"error": e.Error()})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"error": e.Error()})
			}
			return
		}
		c.JSON(http.StatusOK, result)
	}
}