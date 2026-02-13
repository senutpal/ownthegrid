package http

import (
	"net/http"

	"ownthegrid/internal/service"
)

type TileHandler struct {
	tileService *service.TileService
	userService *service.UserService
}

func NewTileHandler(tileService *service.TileService, userService *service.UserService) *TileHandler {
	return &TileHandler{tileService: tileService, userService: userService}
}

func (h *TileHandler) GetBoard(w http.ResponseWriter, r *http.Request) {
	tiles, err := h.tileService.GetAllTiles(r.Context())
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to load board")
		return
	}
	gridWidth, gridHeight := h.tileService.GridSize()
	total := len(tiles)
	claimed := 0
	for _, tile := range tiles {
		if tile != nil && tile.OwnerID != nil {
			claimed++
		}
	}
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"tiles":        tiles,
		"gridWidth":    gridWidth,
		"gridHeight":   gridHeight,
		"totalTiles":   total,
		"claimedTiles": claimed,
	})
}

func (h *TileHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	onlineCount, _ := h.userService.OnlineCount(r.Context())
	totalUsers, _ := h.userService.CountUsers(r.Context())
	stats, err := h.tileService.GetBoardStats(r.Context(), onlineCount, totalUsers)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to get stats")
		return
	}
	respondJSON(w, http.StatusOK, stats)
}
