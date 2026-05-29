package main

import (
	"log"
	"net"
	"os"
	"strconv"

	"notification-service/internal/app"
	"notification-service/internal/database"
	"notification-service/internal/email"

	"cosmix/shared/core/eventbus"
	"cosmix/shared/core/rabbitmq"
	"cosmix/shared/grpc/interceptors"

	notificationpb "cosmix/shared/grpc/gen/go/notification"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	godotenv.Load(".env.local")

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

	eventbus.SetLogger(logger)

	db, err := database.ConnectDB()
	if err != nil {
		log.Fatalf("db connect failed: %v", err)
	}

	rabbit := rabbitmq.Connect(rabbitURL)
	if rabbit == nil {
		log.Println("RabbitMQ unavailable, events will not be published")
	}

	smtpPort, err := strconv.Atoi(
		os.Getenv("SMTP_PORT"),
	)
	if err != nil {
		log.Fatal("invalid SMTP_PORT")
	}

	emailService := email.NewMailDispatcher(
		os.Getenv("SMTP_HOST"),
		smtpPort,
		os.Getenv("SMTP_FROM"),
		"internal/email/templates",
	)

	container := app.NewContainer(db, rabbit, emailService)

	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatalf("grpc listen failed: %v", err)
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

	if err := app.RegisterSubscriptions(container); err != nil {
		log.Fatalf("failed to register subscriptions: %v", err)
	}

	log.Printf("Notification gRPC server running on :%s", grpcPort)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("grpc serve failed: %v", err)
	}
}
