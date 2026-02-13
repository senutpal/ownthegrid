package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"ownthegrid/internal/config"
	"ownthegrid/internal/db"
	httphandler "ownthegrid/internal/handler/http"
	"ownthegrid/internal/handler/ws"
	"ownthegrid/internal/pubsub"
	"ownthegrid/internal/repository"
	"ownthegrid/internal/service"
)

func main() {
	cfg := config.Load()

	pgDB := db.NewPostgres(cfg.DatabaseURL)
	defer pgDB.Close()

	redisClient := db.NewRedis(cfg.RedisURL)
	defer redisClient.Close()

	tileRepo := repository.NewTileRepo(pgDB)
	userRepo := repository.NewUserRepo(pgDB)

	redisStore := service.NewRedisStore(redisClient)
	tileService := service.NewTileService(tileRepo, redisStore, cfg.GridWidth, cfg.GridHeight)
	userService := service.NewUserService(userRepo, redisStore, cfg.JwtSecret, cfg.TokenTTL)

	if err := tileService.SeedIfNeeded(context.Background()); err != nil {
		log.Printf("Warning: Failed to seed tiles: %v", err)
	}

	hub := ws.NewHub()
	go hub.Run()

	publisher := pubsub.NewRedisPublisher(redisClient, pubsub.BoardEventsChannel)
	subscriber := pubsub.NewRedisSubscriber(redisClient, hub, pubsub.BoardEventsChannel)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go subscriber.Subscribe(ctx)

	go startLeaderboardTicker(ctx, userService, publisher, cfg.LeaderboardInterval, cfg.LeaderboardLimit)

	go startStaleConnectionCleanup(ctx, hub, userService, publisher)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{cfg.ClientOrigin},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	httphandler.Mount(r, tileService, userService)
	r.Get("/ws", ws.NewHandler(hub, tileService, userService, publisher).ServeHTTP)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("Server starting on port %s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	<-quit
	log.Println("Shutting down server...")
	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("Shutdown error: %v", err)
	}
}

func startLeaderboardTicker(
	ctx context.Context,
	userService *service.UserService,
	publisher pubsub.Publisher,
	interval time.Duration,
	limit int,
) {
	if interval <= 0 {
		return
	}
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			leaderboard, err := userService.GetLeaderboard(ctx, limit)
			if err != nil {
				log.Printf("Leaderboard update failed: %v", err)
				continue
			}
			payload := map[string]interface{}{
				"leaderboard": leaderboard,
			}
			if err := publisher.Publish(ctx, ws.MsgTypeLeaderboardUpdate, payload); err != nil {
				log.Printf("Leaderboard publish failed: %v", err)
			}
		}
	}
}

func startStaleConnectionCleanup(
	ctx context.Context,
	hub *ws.Hub,
	userService *service.UserService,
	publisher pubsub.Publisher,
) {
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			cleanupStaleRedisEntries(ctx, hub, userService, publisher)
		}
	}
}

func cleanupStaleRedisEntries(
	ctx context.Context,
	hub *ws.Hub,
	userService *service.UserService,
	publisher pubsub.Publisher,
) {
	redisUsers, err := userService.ListOnlineUsers(ctx)
	if err != nil {
		log.Printf("Stale cleanup: failed to list online users: %v", err)
		return
	}

	if len(redisUsers) == 0 {
		return
	}

	connectedIDs := make(map[string]bool)
	for _, id := range hub.GetConnectedUserIDs() {
		connectedIDs[id] = true
	}

	for _, user := range redisUsers {
		if !connectedIDs[user.ID.String()] {
			log.Printf("Stale cleanup: removing stale user %s (%s) from online set", user.ID, user.Username)
			if err := userService.SetOffline(ctx, user.ID); err != nil {
				log.Printf("Stale cleanup: failed to set offline: %v", err)
				continue
			}
			onlineCount, _ := userService.OnlineCount(ctx)
			payload := map[string]interface{}{
				"userId":      user.ID.String(),
				"username":    user.Username,
				"onlineCount": onlineCount,
			}
			if err := publisher.Publish(ctx, ws.MsgTypeUserLeft, payload); err != nil {
				log.Printf("Stale cleanup: failed to publish USER_LEFT: %v", err)
			}
		}
	}
}
