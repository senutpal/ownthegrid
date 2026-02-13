package ws

import (
	"encoding/json"
	"log"
	"sync"
	"time"
)

const (
	MsgTypeInitBoard         = "INIT_BOARD"
	MsgTypeTileClaimed       = "TILE_CLAIMED"
	MsgTypeClaimRejected     = "CLAIM_REJECTED"
	MsgTypeUserJoined        = "USER_JOINED"
	MsgTypeUserLeft          = "USER_LEFT"
	MsgTypeLeaderboardUpdate = "LEADERBOARD_UPDATE"
	MsgTypeError             = "ERROR"
	MsgTypePong              = "PONG"

	pingTimeout     = 90 * time.Second
	cleanupInterval = 30 * time.Second
)

type Message struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type Hub struct {
	clients    map[*Client]bool
	userMap    map[string]*Client
	register   chan *Client
	unregister chan *Client
	broadcast  chan []byte
	cleanup    chan *Client
	mu         sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		userMap:    make(map[string]*Client),
		register:   make(chan *Client, 256),
		unregister: make(chan *Client, 256),
		broadcast:  make(chan []byte, 1024),
		cleanup:    make(chan *Client, 256),
	}
}

func (h *Hub) Run() {
	cleanupTicker := time.NewTicker(cleanupInterval)
	defer cleanupTicker.Stop()

	for {
		select {
		case <-cleanupTicker.C:
			h.cleanupStaleClients()

		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.userMap[client.UserID] = client
			h.mu.Unlock()
			log.Printf("Client registered: %s (%s)", client.UserID, client.Username)

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				delete(h.userMap, client.UserID)
				close(client.send)
			}
			h.mu.Unlock()
			log.Printf("Client unregistered: %s", client.UserID)

		case client := <-h.cleanup:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				delete(h.userMap, client.UserID)
				close(client.send)
				log.Printf("Cleaned up stale client: %s", client.UserID)
			}
			h.mu.Unlock()

		case message := <-h.broadcast:
			h.mu.RLock()
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					go func(c *Client) { h.unregister <- c }(client)
				}
			}
			h.mu.RUnlock()
		}
	}
}

func (h *Hub) cleanupStaleClients() {
	h.mu.RLock()
	var staleClients []*Client
	for client := range h.clients {
		if time.Since(client.LastPong()) > pingTimeout {
			staleClients = append(staleClients, client)
		}
	}
	h.mu.RUnlock()

	if len(staleClients) > 0 {
		log.Printf("Cleanup: Found %d stale clients", len(staleClients))
	}

	for _, client := range staleClients {
		log.Printf("Detected stale client: %s (last pong: %v)", client.UserID, client.LastPong())
		if client.onDisconnect != nil {
			client.onDisconnect()
		}
		h.cleanup <- client
	}
}

func (h *Hub) GetConnectedUserIDs() []string {
	h.mu.RLock()
	defer h.mu.RUnlock()
	ids := make([]string, 0, len(h.userMap))
	for userID := range h.userMap {
		ids = append(ids, userID)
	}
	return ids
}

func (h *Hub) Broadcast(msgType string, payload interface{}) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Broadcast marshal error: %v", err)
		return
	}
	msg := Message{Type: msgType, Payload: payloadBytes}
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Broadcast marshal message error: %v", err)
		return
	}
	h.broadcast <- msgBytes
}

func (h *Hub) BroadcastRaw(message []byte) {
	h.broadcast <- message
}

func (h *Hub) SendToUser(userID string, msgType string, payload interface{}) {
	h.mu.RLock()
	client, ok := h.userMap[userID]
	h.mu.RUnlock()
	if !ok {
		return
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return
	}
	msg := Message{Type: msgType, Payload: payloadBytes}
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return
	}
	select {
	case client.send <- msgBytes:
	default:
	}
}
