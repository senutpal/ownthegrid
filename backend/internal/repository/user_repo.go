package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"ownthegrid/internal/domain"
)

type UserRepo struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	user := &domain.User{}
	query := `SELECT id, username, color, created_at, last_seen FROM users WHERE id = $1`
	err := r.db.QueryRowxContext(ctx, query, id).StructScan(user)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("GetByID: %w", err)
	}
	return user, nil
}

func (r *UserRepo) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	user := &domain.User{}
	query := `SELECT id, username, color, created_at, last_seen FROM users WHERE username = $1`
	err := r.db.QueryRowxContext(ctx, query, username).StructScan(user)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("GetByUsername: %w", err)
	}
	return user, nil
}

func (r *UserRepo) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	created := &domain.User{}
	query := `
        INSERT INTO users (username, color)
        VALUES ($1, $2)
        RETURNING id, username, color, created_at, last_seen
    `
	if err := r.db.QueryRowxContext(ctx, query, user.Username, user.Color).StructScan(created); err != nil {
		return nil, fmt.Errorf("Create: %w", err)
	}
	return created, nil
}

func (r *UserRepo) UpdateLastSeen(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE users SET last_seen = NOW() WHERE id = $1`
	if _, err := r.db.ExecContext(ctx, query, id); err != nil {
		return fmt.Errorf("UpdateLastSeen: %w", err)
	}
	return nil
}

func (r *UserRepo) GetByIDs(ctx context.Context, ids []uuid.UUID) ([]*domain.User, error) {
	if len(ids) == 0 {
		return []*domain.User{}, nil
	}
	query, args, err := sqlx.In(`SELECT id, username, color, created_at, last_seen FROM users WHERE id IN (?)`, ids)
	if err != nil {
		return nil, fmt.Errorf("GetByIDs: %w", err)
	}
	query = r.db.Rebind(query)
	users := []*domain.User{}
	if err := r.db.SelectContext(ctx, &users, query, args...); err != nil {
		return nil, fmt.Errorf("GetByIDs: %w", err)
	}
	return users, nil
}

type LeaderboardEntry struct {
	UserID    uuid.UUID `db:"id" json:"userId"`
	Username  string    `db:"username" json:"username"`
	Color     string    `db:"color" json:"color"`
	TileCount int       `db:"tile_count" json:"tileCount"`
	Rank      int       `json:"rank"`
}

func (r *UserRepo) GetLeaderboard(ctx context.Context, limit int) ([]LeaderboardEntry, error) {
	entries := []LeaderboardEntry{}
	query := `
        SELECT
            u.id, u.username, u.color,
            COUNT(t.id) AS tile_count
        FROM users u
        LEFT JOIN tiles t ON t.owner_id = u.id
        GROUP BY u.id
        ORDER BY tile_count DESC
        LIMIT $1
    `
	if err := r.db.SelectContext(ctx, &entries, query, limit); err != nil {
		return nil, fmt.Errorf("GetLeaderboard: %w", err)
	}
	for i := range entries {
		entries[i].Rank = i + 1
	}
	return entries, nil
}

func (r *UserRepo) CountUsers(ctx context.Context) (int, error) {
	var total int
	query := `SELECT COUNT(*) FROM users`
	if err := r.db.QueryRowxContext(ctx, query).Scan(&total); err != nil {
		return 0, fmt.Errorf("CountUsers: %w", err)
	}
	return total, nil
}
