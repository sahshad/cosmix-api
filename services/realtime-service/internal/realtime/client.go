package realtime

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 1024
	sendBufferSize = 32
)

type Client struct {
	hub    *Hub
	conn   *websocket.Conn
	logger *zap.Logger

	UserID uuid.UUID
	send   chan []byte

	closeOnce sync.Once
	closed    chan struct{}
}

func NewClient(hub *Hub, conn *websocket.Conn, userID uuid.UUID, logger *zap.Logger) *Client {
	return &Client{
		hub:    hub,
		conn:   conn,
		logger: logger,
		UserID: userID,
		send:   make(chan []byte, sendBufferSize),
		closed: make(chan struct{}),
	}
}

func (c *Client) Send(payload []byte) {
	select {
	case <-c.closed:
		return
	default:
	}

	select {
	case c.send <- payload:
	default:
		c.logger.Debug("dropping slow websocket client", zap.String("user_id", c.UserID.String()))
		c.Close()
	}
}

func (c *Client) SendMessage(message Message) {
	payload, err := json.Marshal(message)
	if err != nil {
		return
	}

	c.Send(payload)
}

func (c *Client) Close() {
	c.closeOnce.Do(func() {
		close(c.closed)
		_ = c.conn.Close()
	})
}

func (c *Client) ReadPump() {
	defer c.disconnect()

	c.conn.SetReadLimit(maxMessageSize)
	_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		return c.conn.SetReadDeadline(time.Now().Add(pongWait))
	})

	for {
		if _, _, err := c.conn.ReadMessage(); err != nil {
			if websocket.IsUnexpectedCloseError(
				err,
				websocket.CloseGoingAway,
				websocket.CloseAbnormalClosure,
				websocket.CloseNormalClosure,
			) {
				c.logger.Debug(
					"websocket read failed",
					zap.String("user_id", c.UserID.String()),
					zap.Error(err),
				)
			}

			return
		}
	}
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
		c.disconnect()
	}()

	for {
		select {
		case payload := <-c.send:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))

			if err := c.conn.WriteMessage(websocket.TextMessage, payload); err != nil {
				return
			}
		case <-ticker.C:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))

			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		case <-c.closed:
			return
		}
	}
}

func (c *Client) disconnect() {
	c.Close()
	c.hub.Unregister(c)
}
