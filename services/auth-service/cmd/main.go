package main

import (
	"log"
	"os"

	"auth-service/internal/app"
	"auth-service/internal/database"
	"auth-service/internal/messaging"
	"auth-service/internal/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	godotenv.Load()

	// Gin mode
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Configuration
	port := os.Getenv("PORT")
	rabbitURL := os.Getenv("RABBITMQ_URL")

	// Database
	db, err := database.ConnectDB()
	if err != nil {
		log.Fatalf("db connect failed: %v", err)
	}

	// RabbitMQ
	rabbit := messaging.Connect(rabbitURL)
	if rabbit == nil {
		log.Println("RabbitMQ unavailable, events will not be published")
	}

	container := app.NewContainer(db, rabbit)

	// HTTP router
	router := gin.Default()

	routes.RegisterRoutes(router, container)
	app.RegisterConsumers(container)

	log.Println("Auth service running on :" + port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
