package http

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"ownthegrid/internal/service"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Username string `json:"username"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid payload")
		return
	}

	user, err := h.userService.Register(r.Context(), payload.Username)
	if err != nil {
		if err == service.ErrUsernameTaken {
			respondJSON(w, http.StatusConflict, map[string]string{
				"error": "Username already taken",
				"code":  "USERNAME_TAKEN",
			})
			return
		}
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	if user.Token != "" {
		http.SetCookie(w, &http.Cookie{
			Name:     "otg_token",
			Value:    user.Token,
			Path:     "/",
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
			Secure:   r.TLS != nil,
		})
	}

	respondJSON(w, http.StatusCreated, user)
}

func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid user id")
		return
	}
	user, err := h.userService.GetByID(r.Context(), id)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to load user")
		return
	}
	if user == nil {
		respondError(w, http.StatusNotFound, "User not found")
		return
	}
	respondJSON(w, http.StatusOK, user)
}

func (h *UserHandler) GetLeaderboard(w http.ResponseWriter, r *http.Request) {
	leaderboard, err := h.userService.GetLeaderboard(r.Context(), 10)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to get leaderboard")
		return
	}
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"leaderboard": leaderboard,
	})
}

func (h *UserHandler) GetOnlineCount(w http.ResponseWriter, r *http.Request) {
	users, err := h.userService.ListOnlineUsers(r.Context())
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to get online users")
		return
	}
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"users":       users,
		"onlineCount": len(users),
	})
}
