package main

import (
	"log"
	"notification-service/internal/app"
	"notification-service/internal/database"
	"notification-service/internal/messaging"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	rabbitURL := os.Getenv("RABBITMQ_URL")

	db, err := database.ConnectDB()
	if err != nil {
		log.Fatalf("db connect failed: %v", err)
	}

	rabbit := messaging.Connect(rabbitURL)
	if rabbit == nil {
		log.Fatal("rabbitmq connection failed")
	}

	container := app.NewContainer(db, rabbit)

	app.RegisterConsumers(container)

	log.Println("notification-service started")

	select {}
}
