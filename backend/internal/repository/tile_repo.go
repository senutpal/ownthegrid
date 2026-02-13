package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"ownthegrid/internal/domain"
)

type TileRepo struct {
	db *sqlx.DB
}

func NewTileRepo(db *sqlx.DB) *TileRepo {
	return &TileRepo{db: db}
}

func (r *TileRepo) GetAllTilesWithOwners(ctx context.Context) ([]*domain.Tile, error) {
	tiles := []*domain.Tile{}
	query := `
        SELECT
            t.id, t.x, t.y, t.owner_id, t.claimed_at,
            u.username AS owner_username,
            u.color    AS owner_color
        FROM tiles t
        LEFT JOIN users u ON u.id = t.owner_id
        ORDER BY t.id
    `
	if err := r.db.SelectContext(ctx, &tiles, query); err != nil {
		return nil, fmt.Errorf("GetAllTilesWithOwners: %w", err)
	}
	return tiles, nil
}

func (r *TileRepo) ClaimTile(ctx context.Context, tileID int, userID uuid.UUID) (*domain.Tile, error) {
	tile := &domain.Tile{}
	query := `
        UPDATE tiles
        SET owner_id = $1, claimed_at = NOW()
        WHERE id = $2
          AND owner_id IS NULL
        RETURNING
            id, x, y, owner_id, claimed_at,
            (SELECT username FROM users WHERE id = $1) AS owner_username,
            (SELECT color FROM users WHERE id = $1) AS owner_color
    `
	err := r.db.QueryRowxContext(ctx, query, userID, tileID).StructScan(tile)
	if err == sql.ErrNoRows {
		return nil, domain.ErrTileAlreadyClaimed
	}
	if err != nil {
		return nil, fmt.Errorf("ClaimTile: %w", err)
	}

	_, _ = r.db.ExecContext(ctx,
		`INSERT INTO tile_events (tile_id, user_id, event_type) VALUES ($1, $2, 'claim')`,
		tileID, userID,
	)

	return tile, nil
}

func (r *TileRepo) CountTiles(ctx context.Context) (int, int, error) {
	var total int
	var claimed int
	query := `SELECT COUNT(*) AS total, COUNT(owner_id) AS claimed FROM tiles`
	row := r.db.QueryRowxContext(ctx, query)
	if err := row.Scan(&total, &claimed); err != nil {
		return 0, 0, fmt.Errorf("CountTiles: %w", err)
	}
	return total, claimed, nil
}

func (r *TileRepo) LastActivity(ctx context.Context) (*time.Time, error) {
	var last sql.NullTime
	query := `SELECT MAX(created_at) FROM tile_events`
	if err := r.db.QueryRowxContext(ctx, query).Scan(&last); err != nil {
		return nil, fmt.Errorf("LastActivity: %w", err)
	}
	if !last.Valid {
		return nil, nil
	}
	return &last.Time, nil
}

func (r *TileRepo) SeedTiles(ctx context.Context, gridWidth, gridHeight int) error {
	query := `
		INSERT INTO tiles (id, x, y)
		SELECT
			y * $1 + x AS id,
			x,
			y
		FROM generate_series(0, $2 - 1) AS x,
		     generate_series(0, $3 - 1) AS y
		ON CONFLICT (id) DO NOTHING
	`
	_, err := r.db.ExecContext(ctx, query, gridWidth, gridWidth, gridHeight)
	if err != nil {
		return fmt.Errorf("SeedTiles: %w", err)
	}
	return nil
}

func (r *TileRepo) TileCount(ctx context.Context) (int, error) {
	var count int
	if err := r.db.QueryRowxContext(ctx, "SELECT COUNT(*) FROM tiles").Scan(&count); err != nil {
		return 0, fmt.Errorf("TileCount: %w", err)
	}
	return count, nil
}
