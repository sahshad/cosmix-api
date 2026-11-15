package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"auth-service/internal/app"
	"auth-service/internal/database"
	"auth-service/internal/grpc"

	"cosmix/shared/core/eventbus"
	"cosmix/shared/core/logger"
	"cosmix/shared/core/rabbitmq"

	"github.com/joho/godotenv"
	grpcpkg "google.golang.org/grpc"
)

func main() {
	godotenv.Load(".env.local")

	rabbitURL := os.Getenv("RABBITMQ_URL")
	grpcPort := os.Getenv("GRPC_PORT")

	if grpcPort == "" {
		grpcPort = "50051"
	}

	zapLogger, err := logger.New("auth-service")
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

	grpcServer := grpc.NewServer(
		container.AuthGrpcServer,
		zapLogger,
	)

	if err := app.RegisterSubscriptions(container); err != nil {
		log.Fatalf("failed to register subscriptions: %v", err)
	}

	log.Printf("Auth gRPC server running on :%s", grpcPort)

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	errCh := make(chan error, 1)

	go func() {
		errCh <- grpcServer.Start(grpcPort)
	}()

	select {
	case err := <-errCh:
		if err != nil && !errors.Is(err, grpcpkg.ErrServerStopped) {
			log.Printf("gRPC server stopped: %v", err)
		}

	case <-ctx.Done():
		log.Println("Shutting down gRPC server...")

		done := make(chan struct{})

		go func() {
			grpcServer.GracefulStop()
			close(done)
		}()

		select {
		case <-done:
			log.Println("gRPC server stopped gracefully")

		case <-time.After(10 * time.Second):
			log.Println("Graceful shutdown timed out. Forcing shutdown...")
			grpcServer.Stop()
		}
	}
}
