package ws

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"

	"ownthegrid/internal/domain"
	"ownthegrid/internal/pubsub"
	"ownthegrid/internal/service"
)

type Handler struct {
	hub       *Hub
	tileSvc   *service.TileService
	userSvc   *service.UserService
	publisher pubsub.Publisher
}

func NewHandler(hub *Hub, tileSvc *service.TileService, userSvc *service.UserService, publisher pubsub.Publisher) *Handler {
	return &Handler{
		hub:       hub,
		tileSvc:   tileSvc,
		userSvc:   userSvc,
		publisher: publisher,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	userIDParam := r.URL.Query().Get("userId")
	token := r.URL.Query().Get("token")
	if token == "" {
		if cookie, err := r.Cookie("otg_token"); err == nil {
			token = cookie.Value
		}
	}
	if token == "" {
		http.Error(w, "missing token", http.StatusUnauthorized)
		return
	}

	claims, err := h.userSvc.ValidateToken(token)
	if err != nil {
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}

	if userIDParam == "" {
		userIDParam = claims.UserID
	}
	if claims.UserID != userIDParam {
		http.Error(w, "token mismatch", http.StatusUnauthorized)
		return
	}

	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		http.Error(w, "invalid userId", http.StatusBadRequest)
		return
	}

	user, err := h.userSvc.GetByID(r.Context(), userID)
	if err != nil || user == nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WS upgrade error: %v", err)
		return
	}

	client := &Client{
		hub:      h.hub,
		conn:     conn,
		send:     make(chan []byte, 256),
		UserID:   user.ID.String(),
		Username: user.Username,
		onDisconnect: func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := h.userSvc.SetOffline(ctx, user.ID); err != nil {
				log.Printf("Set offline error: %v", err)
			}
			onlineCount, _ := h.userSvc.OnlineCount(ctx)
			payload := map[string]interface{}{
				"userId":      user.ID.String(),
				"username":    user.Username,
				"onlineCount": onlineCount,
			}
			if err := h.publisher.Publish(ctx, MsgTypeUserLeft, payload); err != nil {
				log.Printf("Publish user left failed: %v", err)
			}
		},
	}
	h.hub.register <- client

	if err := h.userSvc.SetOnline(r.Context(), user.ID); err != nil {
		log.Printf("Set online error: %v", err)
	}
	if err := h.userSvc.UpdateLastSeen(r.Context(), user.ID); err != nil {
		log.Printf("Update last seen error: %v", err)
	}

	onlineCount, _ := h.userSvc.OnlineCount(r.Context())
	tiles := h.getAllTiles(r.Context())
	gridWidth, gridHeight := h.tileSvc.GridSize()
	initPayload := map[string]interface{}{
		"tiles":       tiles,
		"user":        user,
		"onlineCount": onlineCount,
		"gridWidth":   gridWidth,
		"gridHeight":  gridHeight,
	}
	h.hub.SendToUser(user.ID.String(), MsgTypeInitBoard, initPayload)

	h.broadcastUserJoined(r, user, onlineCount)

	go client.writePump()
	go client.readPump(h.handleMessage)
}

func (h *Handler) handleMessage(c *Client, inbound InboundMessage) {
	switch inbound.Type {
	case "PING":
		h.hub.SendToUser(c.UserID, MsgTypePong, map[string]interface{}{})
	case "CLAIM_TILE":
		h.handleClaimTile(c, inbound.Payload)
	default:
		h.hub.SendToUser(c.UserID, MsgTypeError, map[string]interface{}{
			"code":    "UNKNOWN_MESSAGE",
			"message": "Unknown message type",
		})
	}
}

func (h *Handler) handleClaimTile(c *Client, payload json.RawMessage) {
	var request struct {
		TileID int `json:"tileId"`
	}
	if err := json.Unmarshal(payload, &request); err != nil {
		h.hub.SendToUser(c.UserID, MsgTypeError, map[string]interface{}{
			"code":    "BAD_PAYLOAD",
			"message": "Invalid payload",
		})
		return
	}

	userID, err := uuid.Parse(c.UserID)
	if err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	tile, err := h.tileSvc.ClaimTile(ctx, request.TileID, userID)
	if err != nil {
		h.handleClaimError(c, request.TileID, err)
		return
	}

	color := ""
	if tile.OwnerColor != nil {
		color = *tile.OwnerColor
	}
	payloadOut := map[string]interface{}{
		"tileId":        tile.ID,
		"x":             tile.X,
		"y":             tile.Y,
		"userId":        userID.String(),
		"username":      c.Username,
		"color":         color,
		"claimedAt":     tile.ClaimedAt,
		"previousOwner": nil,
	}
	if err := h.publisher.Publish(ctx, MsgTypeTileClaimed, payloadOut); err != nil {
		log.Printf("Publish claim failed: %v", err)
	}
}

func (h *Handler) handleClaimError(c *Client, tileID int, err error) {
	reason := "SERVER_ERROR"
	if errors.Is(err, domain.ErrTileInvalid) {
		reason = "INVALID_TILE"
	} else if errors.Is(err, domain.ErrTileAlreadyClaimed) {
		reason = "ALREADY_CLAIMED"
	}

	h.hub.SendToUser(c.UserID, MsgTypeClaimRejected, map[string]interface{}{
		"tileId": tileID,
		"reason": reason,
	})
}

func (h *Handler) broadcastUserJoined(r *http.Request, user *domain.User, onlineCount int) {
	payload := map[string]interface{}{
		"userId":      user.ID.String(),
		"username":    user.Username,
		"color":       user.Color,
		"onlineCount": onlineCount,
	}
	if err := h.publisher.Publish(r.Context(), MsgTypeUserJoined, payload); err != nil {
		log.Printf("Publish user joined failed: %v", err)
	}
}

func (h *Handler) getAllTiles(ctx context.Context) []domain.Tile {
	tiles, err := h.tileSvc.GetAllTiles(ctx)
	if err != nil {
		log.Printf("Failed to load tiles: %v", err)
		return []domain.Tile{}
	}
	out := make([]domain.Tile, 0, len(tiles))
	for _, tile := range tiles {
		if tile != nil {
			out = append(out, *tile)
		}
	}
	return out
}
