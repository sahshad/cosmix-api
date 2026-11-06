package main

import (
	"log"
	"net"
	"os"

	"communication-service/internal/app"
	"communication-service/internal/database"

	"cosmix/shared/core/eventbus"
	"cosmix/shared/core/logger"
	"cosmix/shared/core/rabbitmq"
	communicationpb "cosmix/shared/grpc/gen/go/communication"
	"cosmix/shared/grpc/interceptors"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func main() {
	godotenv.Load(".env.local")

	rabbitURL := os.Getenv("RABBITMQ_URL")
	grpcPort := os.Getenv("GRPC_PORT")

	if grpcPort == "" {
		grpcPort = "50056"
	}

	zapLogger, err := logger.New("communication-service")
	if err != nil {
		log.Fatalf("failed to initialize logger: %v", err)
	}
	defer zapLogger.Sync()

	eventbus.SetLogger(zapLogger)

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
			interceptors.LoggingInterceptor(zapLogger),
			interceptors.RecoveryInterceptor(zapLogger),
			interceptors.ErrorInterceptor,
		),
	)

	communicationpb.RegisterCommunicationServiceServer(
		grpcServer,
		container.CommunicationGrpcServer,
	)

	if err := app.RegisterSubscriptions(container); err != nil {
		log.Fatalf("failed to register subscriptions: %v", err)
	}

	log.Printf("Communication gRPC server running on :%s", grpcPort)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("grpc serve failed: %v", err)
	}
}
