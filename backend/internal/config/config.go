package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found, using environment variables")
	}
}

type Config struct {
	Port                string
	ClientOrigin        string
	DatabaseURL         string
	RedisURL            string
	JwtSecret           string
	GridWidth           int
	GridHeight          int
	TokenTTL            time.Duration
	LeaderboardInterval time.Duration
	LeaderboardLimit    int
}

func Load() Config {
	cfg := Config{
		Port:                getEnv("PORT", "8080"),
		ClientOrigin:        getEnv("CLIENT_ORIGIN", "http://localhost:5173"),
		DatabaseURL:         requireEnv("DATABASE_URL"),
		RedisURL:            requireEnv("REDIS_URL"),
		JwtSecret:           requireEnv("JWT_SECRET"),
		GridWidth:           getEnvInt("GRID_WIDTH", 50),
		GridHeight:          getEnvInt("GRID_HEIGHT", 40),
		TokenTTL:            time.Duration(getEnvInt("TOKEN_TTL_HOURS", 168)) * time.Hour,
		LeaderboardInterval: time.Duration(getEnvInt("LEADERBOARD_INTERVAL_SECONDS", 10)) * time.Second,
		LeaderboardLimit:    getEnvInt("LEADERBOARD_LIMIT", 10),
	}

	if len(cfg.JwtSecret) < 32 {
		log.Fatal("JWT_SECRET must be at least 32 characters")
	}

	if cfg.GridWidth <= 0 || cfg.GridHeight <= 0 {
		log.Fatal("GRID_WIDTH and GRID_HEIGHT must be positive")
	}

	return cfg
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func requireEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("%s is required", key)
	}
	return value
}

func getEnvInt(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		log.Fatalf("%s must be an integer", key)
	}
	return parsed
}
