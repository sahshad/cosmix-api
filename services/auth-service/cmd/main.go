package main

import (
	"log"
	"net"
	"os"

	"auth-service/internal/app"
	"auth-service/internal/database"

	"cosmix/shared/core/eventbus"
	"cosmix/shared/core/rabbitmq"
	"cosmix/shared/grpc/interceptors"

	authpb "cosmix/shared/grpc/gen/go/auth"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	godotenv.Load()

	rabbitURL := os.Getenv("RABBITMQ_URL")
	grpcPort := os.Getenv("GRPC_PORT")

	if grpcPort == "" {
		grpcPort = "50051"
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

	container := app.NewContainer(db, rabbit)

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

	authpb.RegisterAuthServiceServer(
		grpcServer,
		container.AuthGrpcServer,
	)

	if err := app.RegisterSubscriptions(container); err != nil {
		log.Fatalf("failed to register subscriptions: %v", err)
	}

	log.Printf("Auth gRPC server running on :%s", grpcPort)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("grpc serve failed: %v", err)
	}
}
