package main

import (
	"log"
	"os"
	"time"

	"user-service/internal/app"
	"user-service/internal/database"
	"cosmix/shared/core/rabbitmq"
	"user-service/internal/middlewares"
	"user-service/internal/routes"

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

	rabbitChannel := rabbitmq.Connect(rabbitURL)
	if rabbitChannel == nil {
		log.Println("RabbitMQ unavailable, running without consumer")
	}

	container := app.NewContainer(db, rabbitChannel)

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

	log.Println("User service running on :" + port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
