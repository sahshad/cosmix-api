module notification-service

go 1.25.7

require (
	github.com/joho/godotenv v1.5.1
	github.com/pressly/goose/v3 v3.27.1
	github.com/rabbitmq/amqp091-go v1.11.0 // indirect
	gorm.io/driver/postgres v1.6.0
	gorm.io/gorm v1.31.1
)

require go.uber.org/zap v1.28.0

require cosmix/shared/events v0.0.0

require cosmix/shared/core v0.0.0

require (
	cosmix/shared/grpc v0.0.0
	google.golang.org/grpc v1.81.1
	google.golang.org/protobuf v1.36.11
	gopkg.in/gomail.v2 v2.0.0-20160411212932-81ebce5c23df
)

replace cosmix/shared/events => ../../shared/events

replace cosmix/shared/core => ../../shared/core

replace cosmix/shared/grpc => ../../shared/grpc

require (
	github.com/google/uuid v1.6.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/pgx/v5 v5.9.2 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/mfridman/interpolate v0.0.2 // indirect
	github.com/sethvargo/go-retry v0.3.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/net v0.53.0 // indirect
	golang.org/x/sync v0.20.0 // indirect
	golang.org/x/sys v0.44.0 // indirect
	golang.org/x/text v0.37.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260420184626-e10c466a9529 // indirect
	gopkg.in/alexcesaro/quotedprintable.v3 v3.0.0-20150716171945-2caba252f4dc // indirect
)
