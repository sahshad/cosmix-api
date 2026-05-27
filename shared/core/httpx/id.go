package httpx

import (
	"strconv"

	appErr "cosmix/shared/core/errors"

	"github.com/gin-gonic/gin"
)

func ParseUserIDHeader(c *gin.Context) (uint, error) {
	userIDStr := c.GetHeader("X-User-Id")
	if userIDStr == "" {
		return 0, appErr.NewUnauthorized(
			"user not authenticated",
		)
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		return 0, appErr.NewBadRequest(
			"X-User-Id",
			"invalid user id",
		)
	}

	return uint(userID), nil
}

func ParseParamID(c *gin.Context, param string) (uint, error) {
	idStr := c.Param(param)
	if idStr == "" {
		return 0, appErr.NewBadRequest(
			param,
			"id is required",
		)
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return 0, appErr.NewBadRequest(
			param,
			"invalid id",
		)
	}

	return uint(id), nil
}

func ParseParamIDWithDefault(c *gin.Context, param string, defaultValue uint) (uint, error) {
	idStr := c.Param(param)
	if idStr == "" {
		return defaultValue, nil
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return 0, appErr.NewBadRequest(
			param,
			"invalid id",
		)
	}

	return uint(id), nil
}

func ParseQueryString(c *gin.Context, param string) (string, error) {
	value := c.Query(param)
	if value == "" {
		return "", appErr.NewBadRequest(
			param,
			param+" is required",
		)
	}

	return value, nil
}

func ParseQueryIntWithDefault(c *gin.Context, param string, defaultValue int) (int64, error) {
	valueStr := c.Query(param)
	if valueStr == "" {
		return int64(defaultValue), nil
	}

	value, err := strconv.ParseInt(valueStr, 10, 32)
	if err != nil {
		return 0, appErr.NewBadRequest(
			param,
			"invalid id",
		)
	}

	return int64(value), nil
}
