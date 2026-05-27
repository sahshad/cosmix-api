package interceptors

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func RecoveryInterceptor(logger *zap.Logger) grpc.UnaryServerInterceptor {

	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (_ interface{}, err error) {

		defer func() {

			if r := recover(); r != nil {

				logger.Error(
					"panic recovered",
					zap.Any("panic", r),
				)

				err = status.Error(
					codes.Internal,
					"internal server error",
				)
			}
		}()

		return handler(ctx, req)
	}
}
