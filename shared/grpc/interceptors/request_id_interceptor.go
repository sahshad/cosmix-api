package interceptors

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const RequestIDKey = "x-request-id"

func RequestIDInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {

	requestID := uuid.NewString()

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if values := md.Get(RequestIDKey); len(values) > 0 {
			requestID = values[0]
		}
	}

	ctx = context.WithValue(
		ctx,
		RequestIDKey,
		requestID,
	)

	return handler(ctx, req)
}
