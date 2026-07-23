package interceptors

import (
	"context"
	"strings"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func LoggingInterceptor(logger *zap.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {

		start := time.Now()
		resp, err := handler(ctx, req)
		duration := time.Since(start)
		requestID, _ := ctx.Value(RequestIDKey).(string)

		service, method := splitFullMethod(info.FullMethod)

		fields := []zap.Field{
			zap.String("service", service),
			zap.String("method", method),
			zap.String("request_id", requestID),
			zap.Duration("duration", duration),
		}

		if err != nil {
			st, _ := status.FromError(err)

			fields = append(
				fields,
				zap.String("full_method", info.FullMethod),
				zap.String("error", st.Message()),
			)

			logger.Error("gRPC Request", fields...)
		} else {
			logger.Info("gRPC Request", fields...)
		}

		return resp, err
	}
}

// splitFullMethod turns "/pkg.Service/Method" into ("Service", "Method").
func splitFullMethod(full string) (service, method string) {
	trimmed := strings.TrimPrefix(full, "/")

	parts := strings.SplitN(trimmed, "/", 2)
	if len(parts) != 2 {
		return "", full
	}

	pkgService, methodName := parts[0], parts[1]

	if idx := strings.LastIndex(pkgService, "."); idx != -1 {
		return pkgService[idx+1:], methodName
	}

	return pkgService, methodName
}
