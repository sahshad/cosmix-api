package main
import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"realtime-service/internal/app"
	"realtime-service/internal/auth"
	"realtime-service/internal/realtime"

	"cosmix/shared/core/eventbus"
	"cosmix/shared/core/logger"
	"cosmix/shared/core/rabbitmq"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	godotenv.Load(".env.local")

	port := os.Getenv("PORT")
	if port == "" {
		port = "50055"
	}

	zapLogger, err := logger.New("realtime-service")
	if err != nil {
		log.Fatalf("failed to initialize logger: %v", err)
	}
	defer zapLogger.Sync()

	eventbus.SetLogger(zapLogger)

	rabbit := rabbitmq.Connect(os.Getenv("RABBITMQ_URL"))
	if rabbit == nil {
		log.Fatal("RabbitMQ unavailable, realtime events cannot be consumed")
	}

	verifier := auth.NewVerifier(auth.SecretFromEnv())
	allowedOrigins := realtime.ParseAllowedOrigins(os.Getenv("CORS_ALLOWED_ORIGINS"))

	container := app.NewContainer(rabbit, zapLogger, verifier, allowedOrigins)

	if err := app.RegisterSubscriptions(container); err != nil {
		log.Fatalf("failed to register subscriptions: %v", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/ws", container.WsHandler)
	mux.HandleFunc("/health", container.HealthHandler)

	server := &http.Server{
		Addr:              ":" + port,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		zapLogger.Info("Realtime server listening", zap.String("addr", server.Addr))

		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("http serve failed: %v", err)
		}
	}()

	shutdownCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	<-shutdownCtx.Done()

	zapLogger.Info("Shutting down realtime server")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		zapLogger.Error("http shutdown failed", zap.Error(err))
	}

	container.Hub.CloseAll()
	_ = rabbit.Connection.Close()
}

