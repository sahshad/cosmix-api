package main

import (
	"log"
	"net"
	"os"

	"notification-service/internal/app"
	"notification-service/internal/database"

	"cosmix/shared/core/rabbitmq"
	interceptors "cosmix/shared/grpc/interceptors"

	notificationpb "cosmix/shared/grpc/gen/go/notification"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	godotenv.Load()

	rabbitURL := os.Getenv("RABBITMQ_URL")
	grpcPort := os.Getenv("GRPC_PORT")

	if grpcPort == "" {
		grpcPort = "50053"
	}

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	db, err := database.ConnectDB()
	if err != nil {
		log.Fatalf("db connect failed: %v", err)
	}

	rabbit := rabbitmq.Connect(rabbitURL)
	if rabbit == nil {
		log.Println(
			"RabbitMQ unavailable, events will not be published",
		)
	}

	container := app.NewContainer(
		db,
		rabbit,
	)

	lis, err := net.Listen(
		"tcp",
		":"+grpcPort,
	)
	if err != nil {
		log.Fatalf(
			"grpc listen failed: %v",
			err,
		)
	}

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptors.RequestIDInterceptor,
			interceptors.LoggingInterceptor(logger),
			interceptors.RecoveryInterceptor(logger),
			interceptors.ErrorInterceptor,
		),
	)

	notificationpb.RegisterNotificationServiceServer(
		grpcServer,
		container.NotificationGrpcServer,
	)

	app.RegisterConsumers(container)

	log.Printf(
		"Auth gRPC server running on :%s",
		grpcPort,
	)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf(
			"grpc serve failed: %v",
			err,
		)
	}
}