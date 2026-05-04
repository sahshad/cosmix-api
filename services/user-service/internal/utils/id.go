package utils

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ParseError struct {
	Field   string
	Message string
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}


// ParseUserIDHeader parses the "X-User-Id" header and returns a uint or ParseError
func ParseUserIDHeader(c *gin.Context) (uint, error) {
	userIDStr := c.GetHeader("X-User-Id")
	if userIDStr == "" {
		return 0, &ParseError{
			Field:   "X-User-Id",
			Message: "user not authenticated",
		}
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		return 0, &ParseError{
			Field:   "X-User-Id",
			Message: "invalid user id",
		}
	}

	return uint(userID), nil
}

// ParseParamID parses a URL parameter (like "/users/:id") and returns a uint or ParseError
func ParseParamID(c *gin.Context, param string) (uint, error) {
	idStr := c.Param(param)
	if idStr == "" {
		return 0, &ParseError{
			Field:   param,
			Message: "id is required",
		}
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return 0, &ParseError{
			Field:   param,
			Message: "invalid id",
		}
	}

	return uint(id), nil
}