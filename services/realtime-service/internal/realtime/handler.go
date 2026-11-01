package realtime

import (
	"net/http"
	"strings"

	"realtime-service/internal/auth"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type Handler struct {
	hub      *Hub
	verifier *auth.Verifier
	logger   *zap.Logger
	upgrader websocket.Upgrader
}

func NewHandler(
	hub *Hub,
	verifier *auth.Verifier,
	logger *zap.Logger,
	allowedOrigins []string,
) *Handler {
	handler := &Handler{
		hub:      hub,
		verifier: verifier,
		logger:   logger,
	}

	handler.upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return isOriginAllowed(allowedOrigins, r)
		},
	}

	return handler
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !websocket.IsWebSocketUpgrade(r) {
		http.Error(w, "websocket upgrade required", http.StatusUpgradeRequired)
		return
	}

	token := extractToken(r)
	if token == "" {
		http.Error(w, "missing access token", http.StatusUnauthorized)
		return
	}

	claims, err := h.verifier.Verify(token)
	if err != nil {
		h.logger.Debug("websocket authentication failed", zap.Error(err))
		http.Error(w, "invalid access token", http.StatusUnauthorized)
		return
	}

	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.logger.Debug("websocket upgrade failed", zap.Error(err))
		return
	}

	client := NewClient(h.hub, conn, claims.UserID, h.logger)

	h.hub.Register(client)
	client.SendMessage(Message{
		Type: MessageTypeConnected,
		Data: ConnectedData{UserID: claims.UserID.String()},
	})

	go client.WritePump()
	go client.ReadPump()

	h.logger.Info(
		"websocket connected",
		zap.String("user_id", claims.UserID.String()),
		zap.String("remote_addr", r.RemoteAddr),
	)
}

func extractToken(r *http.Request) string {
	if token := strings.TrimSpace(r.URL.Query().Get("token")); token != "" {
		return token
	}

	header := strings.TrimSpace(r.Header.Get("Authorization"))
	if len(header) > 7 && strings.EqualFold(header[:7], "bearer ") {
		return strings.TrimSpace(header[7:])
	}

	return ""
}

func ParseAllowedOrigins(value string) []string {
	if strings.TrimSpace(value) == "" {
		return nil
	}

	parts := strings.Split(value, ",")
	origins := make([]string, 0, len(parts))

	for _, part := range parts {
		if origin := strings.TrimSpace(part); origin != "" {
			origins = append(origins, origin)
		}
	}

	return origins
}

func isOriginAllowed(allowedOrigins []string, r *http.Request) bool {
	origin := r.Header.Get("Origin")
	if origin == "" || len(allowedOrigins) == 0 {
		return true
	}

	for _, allowed := range allowedOrigins {
		if allowed == "*" || strings.EqualFold(allowed, origin) {
			return true
		}
	}

	return false
}
