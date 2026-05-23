package interceptors

import (
	"context"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func LoggingInterceptor(
	logger *zap.Logger,
) grpc.UnaryServerInterceptor {

	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {

		start := time.Now()

		resp, err := handler(
			ctx,
			req,
		)

		duration := time.Since(start)

		requestID, _ :=
			ctx.Value(
				RequestIDKey,
			).(string)

		logger.Info(
			"grpc request",
			zap.String(
				"request_id",
				requestID,
			),
			zap.String(
				"method",
				info.FullMethod,
			),
			zap.Duration(
				"duration",
				duration,
			),
			zap.Error(err),
		)

		return resp, err
	}
}