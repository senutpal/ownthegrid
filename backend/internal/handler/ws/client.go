package ws

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 1024
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Client struct {
	hub          *Hub
	conn         *websocket.Conn
	send         chan []byte
	UserID       string
	Username     string
	onDisconnect func()
	lastPong     time.Time
	lastPongMu   sync.Mutex
}

func (c *Client) LastPong() time.Time {
	c.lastPongMu.Lock()
	defer c.lastPongMu.Unlock()
	if c.lastPong.IsZero() {
		return time.Now()
	}
	return c.lastPong
}

func (c *Client) updateLastPong() {
	c.lastPongMu.Lock()
	defer c.lastPongMu.Unlock()
	c.lastPong = time.Now()
}

type InboundMessage struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type MessageHandler func(c *Client, inbound InboundMessage)

func (c *Client) readPump(handler MessageHandler) {
	defer func() {
		c.hub.unregister <- c
		_ = c.conn.Close()
		if c.onDisconnect != nil {
			c.onDisconnect()
		}
	}()

	c.conn.SetReadLimit(maxMessageSize)
	_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.updateLastPong()
	c.conn.SetPongHandler(func(string) error {
		c.updateLastPong()
		_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WS read error for %s: %v", c.UserID, err)
			}
			break
		}

		var inbound InboundMessage
		if err := json.Unmarshal(message, &inbound); err != nil {
			log.Printf("WS parse error: %v", err)
			continue
		}

		handler(c, inbound)
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		_ = c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}

		case <-ticker.C:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
