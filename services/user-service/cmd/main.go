package main

import (
	"log"
	"os"

	"user-service/internal/app"
	"user-service/internal/database"
	"user-service/internal/messaging"
	"user-service/internal/routes"

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
	rabbitChannel := messaging.Connect(rabbitURL)
	if rabbitChannel == nil {
		log.Println("RabbitMQ unavailable, running without consumer")
	}

	container := app.NewContainer(db, rabbitChannel)

	// HTTP router
	router := gin.Default()
	routes.RegisterRoutes(router, container)
	app.RegisterConsumers(container)

	log.Println("User service running on :" + port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
