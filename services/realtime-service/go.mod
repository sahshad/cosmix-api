module realtime-service

go 1.25.1

require (
	cosmix/shared/core v0.0.0
	cosmix/shared/events v0.0.0
	github.com/golang-jwt/jwt/v5 v5.3.0
	github.com/google/uuid v1.6.0
	github.com/gorilla/websocket v1.5.3
	github.com/joho/godotenv v1.5.1
	go.uber.org/zap v1.28.0
)

require (
	github.com/rabbitmq/amqp091-go v1.11.0 // indirect
	go.uber.org/multierr v1.10.0 // indirect
)

replace cosmix/shared/events => ../../shared/events

replace cosmix/shared/core => ../../shared/core
