package http

import (
	"github.com/go-chi/chi/v5"

	"ownthegrid/internal/service"
)

func Mount(r chi.Router, tileService *service.TileService, userService *service.UserService) {
	tileHandler := NewTileHandler(tileService, userService)
	userHandler := NewUserHandler(userService)

	r.Route("/api", func(api chi.Router) {
		api.Route("/users", func(users chi.Router) {
			users.Post("/register", userHandler.Register)
			users.Get("/{id}", userHandler.GetByID)
			users.Get("/online", userHandler.GetOnlineCount)
			users.Get("/leaderboard", userHandler.GetLeaderboard)
		})

		api.Route("/board", func(board chi.Router) {
			board.Get("/", tileHandler.GetBoard)
			board.Get("/stats", tileHandler.GetStats)
		})
	})
}
