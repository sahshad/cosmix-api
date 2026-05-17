package main

import (
	"cosmix/shared/core/middleware"
	"cosmix/shared/core/rabbitmq"
	"log"
	"notification-service/internal/app"
	"notification-service/internal/database"
	"notification-service/internal/routes"
	"os"
	"time"

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

	rabbit := rabbitmq.Connect(rabbitURL)
	if rabbit == nil {
		log.Fatal("rabbitmq connection failed")
	}

	container := app.NewContainer(db, rabbit)

	logger, _ := zap.NewProduction()
	router := gin.New()

	router.Use(middleware.RequestIDMiddleware())
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

	app.RegisterConsumers(container)
	routes.RegisterRoutes(router, container)

	log.Println("Notification Service running on :", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
