package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"ownthegrid/internal/domain"
	"ownthegrid/internal/repository"
)

type TileService struct {
	repo       *repository.TileRepo
	redis      RedisStore
	gridWidth  int
	gridHeight int
}

func NewTileService(
	repo *repository.TileRepo,
	redis RedisStore,
	gridWidth int,
	gridHeight int,
) *TileService {
	return &TileService{
		repo:       repo,
		redis:      redis,
		gridWidth:  gridWidth,
		gridHeight: gridHeight,
	}
}

func (s *TileService) ClaimTile(ctx context.Context, tileID int, userID uuid.UUID) (*domain.Tile, error) {
	if tileID < 0 || tileID >= s.gridWidth*s.gridHeight {
		return nil, domain.ErrTileInvalid
	}

	tile, err := s.repo.ClaimTile(ctx, tileID, userID)
	if err != nil {
		return nil, err
	}

	if err := s.redis.ZIncrBy(ctx, "board:leaderboard", 1, userID.String()); err != nil {
		return nil, fmt.Errorf("leaderboard update: %w", err)
	}

	return tile, nil
}

func (s *TileService) GetAllTiles(ctx context.Context) ([]*domain.Tile, error) {
	return s.repo.GetAllTilesWithOwners(ctx)
}

func (s *TileService) GridSize() (int, int) {
	return s.gridWidth, s.gridHeight
}

func (s *TileService) SeedIfNeeded(ctx context.Context) error {
	count, err := s.repo.TileCount(ctx)
	if err != nil {
		return fmt.Errorf("SeedIfNeeded: %w", err)
	}
	if count == 0 {
		if err := s.repo.SeedTiles(ctx, s.gridWidth, s.gridHeight); err != nil {
			return fmt.Errorf("SeedIfNeeded: %w", err)
		}
	}
	return nil
}

func (s *TileService) GetBoardStats(ctx context.Context, onlineCount int, totalUsers int) (map[string]interface{}, error) {
	_, claimed, err := s.repo.CountTiles(ctx)
	if err != nil {
		return nil, err
	}
	lastActivity, err := s.repo.LastActivity(ctx)
	if err != nil {
		return nil, err
	}
	// Use grid dimensions from config, not DB count (migration may have different size)
	total := s.gridWidth * s.gridHeight
	unclaimed := total - claimed
	// Ensure unclaimed doesn't go negative if DB has more claimed tiles than expected
	if unclaimed < 0 {
		unclaimed = 0
	}
	payload := map[string]interface{}{
		"totalTiles":     total,
		"claimedTiles":   claimed,
		"unclaimedTiles": unclaimed,
		"onlineUsers":    onlineCount,
		"totalUsers":     totalUsers,
		"lastActivity":   lastActivity,
	}
	return payload, nil
}
