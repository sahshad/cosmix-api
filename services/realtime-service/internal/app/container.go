package app

import (
	"encoding/json"
	"net/http"

	"realtime-service/internal/auth"
	"realtime-service/internal/realtime"
	"realtime-service/internal/services"

	"cosmix/shared/core/rabbitmq"

	"go.uber.org/zap"
)

type Container struct {
	Rabbit *rabbitmq.Rabbit
	Hub    *realtime.Hub

	WsHandler    *realtime.Handler
	EventService *services.EventService
}

func NewContainer(
	rabbit *rabbitmq.Rabbit,
	logger *zap.Logger,
	verifier *auth.Verifier,
	allowedOrigins []string,
) *Container {
	hub := realtime.NewHub()

	return &Container{
		Rabbit:       rabbit,
		Hub:          hub,
		WsHandler:    realtime.NewHandler(hub, verifier, logger, allowedOrigins),
		EventService: services.NewEventService(hub),
	}
}

func (container *Container) HealthHandler(w http.ResponseWriter, r *http.Request) {
	users, connections := container.Hub.Stats()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(map[string]any{
		"status":      "ok",
		"onlineUsers": users,
		"connections": connections,
	})
}
