package db

import (
	"context"
	"log"

	"github.com/jmoiron/sqlx"
)

func Migrate(db *sqlx.DB) error {
	ctx := context.Background()

	migrations := []string{
		`CREATE EXTENSION IF NOT EXISTS "pgcrypto"`,
		`CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			username VARCHAR(32) NOT NULL UNIQUE,
			color CHAR(7) NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			last_seen TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS tiles (
			id INTEGER PRIMARY KEY,
			x INTEGER NOT NULL,
			y INTEGER NOT NULL,
			owner_id UUID REFERENCES users(id) ON DELETE SET NULL,
			claimed_at TIMESTAMPTZ,
			UNIQUE (x, y)
		)`,
		`CREATE INDEX IF NOT EXISTS idx_tiles_owner ON tiles(owner_id)`,
		`CREATE INDEX IF NOT EXISTS idx_tiles_claimed_at ON tiles(claimed_at)`,
		`CREATE TABLE IF NOT EXISTS tile_events (
			id BIGSERIAL PRIMARY KEY,
			tile_id INTEGER NOT NULL REFERENCES tiles(id),
			user_id UUID NOT NULL REFERENCES users(id),
			event_type VARCHAR(16) NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE INDEX IF NOT EXISTS idx_events_user ON tile_events(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_events_created ON tile_events(created_at DESC)`,
	}

	for _, m := range migrations {
		if _, err := db.ExecContext(ctx, m); err != nil {
			log.Printf("Migration warning: %v", err)
			continue
		}
	}

	log.Println("Database migration completed")
	return nil
}
