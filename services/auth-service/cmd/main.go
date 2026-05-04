package main

import (
	"log"
	"os"
	"time"

	"auth-service/internal/app"
	"auth-service/internal/database"
	"auth-service/internal/messaging"
	"auth-service/internal/middlewares"
	"auth-service/internal/routes"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	godotenv.Load()

	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	port := os.Getenv("PORT")
	rabbitURL := os.Getenv("RABBITMQ_URL")

	db, err := database.ConnectDB()
	if err != nil {
		log.Fatalf("db connect failed: %v", err)
	}

	rabbit := messaging.Connect(rabbitURL)
	if rabbit == nil {
		log.Println("RabbitMQ unavailable, events will not be published")
	}

	container := app.NewContainer(db, rabbit)

	logger, _ := zap.NewProduction()
	router := gin.New()

	router.Use(middlewares.RequestIDMiddleware())
	router.Use(ginzap.GinzapWithConfig(logger, &ginzap.Config{
		TimeFormat: time.RFC3339,
		UTC:        true,
		Context: func(c *gin.Context) []zapcore.Field {
			if id, exists := c.Get("RequestID"); exists {
				return []zapcore.Field{zap.String("request_id", id.(string))}
			}
			return nil
		},
	}))
	router.Use(ginzap.RecoveryWithZap(logger, true))

	routes.RegisterRoutes(router, container)
	app.RegisterConsumers(container)

	log.Println("Auth service running on :" + port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
