package eventbus

import "context"

type contextKey string

const CorrelationIDKey contextKey = "correlation_id"
const RequestIDKey = "x-request-id"

func WithCorrelationID(
	ctx context.Context,
	id string,
) context.Context {
	return context.WithValue(
		ctx,
		CorrelationIDKey,
		id,
	)
}

func CorrelationID(
	ctx context.Context,
) string {

	id, ok := ctx.Value(
		CorrelationIDKey,
	).(string)

	if !ok {
		return ""
	}

	return id
}

func RequestID(
	ctx context.Context,
) string {

	id, ok := ctx.Value(
		RequestIDKey,
	).(string)

	if !ok {
		return ""
	}

	return id
}
