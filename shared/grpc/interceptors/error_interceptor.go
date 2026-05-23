package interceptors

import (
	"context"
	"net/http"

	apperrors "cosmix/shared/core/errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ErrorInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {

	resp, err := handler(ctx, req)

	if err == nil {
		return resp, nil
	}

	if appErr, ok := err.(*apperrors.AppError); ok {

		switch appErr.Status {

		case http.StatusBadRequest:
			return nil, status.Error(
				codes.InvalidArgument,
				appErr.Message,
			)

		case http.StatusUnauthorized:
			return nil, status.Error(
				codes.Unauthenticated,
				appErr.Message,
			)

		case http.StatusForbidden:
			return nil, status.Error(
				codes.PermissionDenied,
				appErr.Message,
			)

		case http.StatusNotFound:
			return nil, status.Error(
				codes.NotFound,
				appErr.Message,
			)
		}
	}

	return nil, status.Error(
		codes.Internal,
		"internal server error",
	)
}