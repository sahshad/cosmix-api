package middlewares

import (
	"net/http"
	"user-service/internal/errors"

	"github.com/gin-gonic/gin"
)

func ErrorHandler(handler func(*gin.Context) (interface{}, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		result, err := handler(c)
		if err != nil {
			switch e := err.(type) {
			case *errors.BadRequestError:
				c.JSON(http.StatusBadRequest, gin.H{"error": e.Error()})
			case *errors.UnauthorizedError:
				c.JSON(http.StatusUnauthorized, gin.H{"error": e.Error()})
			case *errors.NotFoundError:
				c.JSON(http.StatusNotFound, gin.H{"error": e.Error()})
			case *errors.InternalServerError:
				c.JSON(http.StatusInternalServerError, gin.H{"error": e.Error()})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"error": e.Error()})
			}
			return
		}
		c.JSON(http.StatusOK, result)
	}
}