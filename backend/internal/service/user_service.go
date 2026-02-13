package service

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"

	"ownthegrid/internal/domain"
	"ownthegrid/internal/repository"
)

var ErrUsernameTaken = errors.New("username already taken")

type UserService struct {
	repo      *repository.UserRepo
	redis     RedisStore
	jwtSecret string
	tokenTTL  time.Duration
}

func NewUserService(repo *repository.UserRepo, redis RedisStore, jwtSecret string, tokenTTL time.Duration) *UserService {
	return &UserService{
		repo:      repo,
		redis:     redis,
		jwtSecret: jwtSecret,
		tokenTTL:  tokenTTL,
	}
}

func (s *UserService) Register(ctx context.Context, username string) (*domain.User, error) {
	username = strings.TrimSpace(username)
	if username == "" {
		return nil, errors.New("username is required")
	}

	existing, err := s.repo.GetByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrUsernameTaken
	}

	color := pickColor()
	user := &domain.User{Username: username, Color: color}
	created, err := s.repo.Create(ctx, user)
	if err != nil {
		if isUniqueViolation(err) {
			return nil, ErrUsernameTaken
		}
		return nil, err
	}

	token, err := buildToken(created.ID, created.Username, s.jwtSecret, s.tokenTTL)
	if err != nil {
		return nil, fmt.Errorf("token: %w", err)
	}
	created.Token = token
	return created, nil
}

func (s *UserService) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *UserService) ValidateToken(token string) (*Claims, error) {
	return parseToken(token, s.jwtSecret)
}

func (s *UserService) UpdateLastSeen(ctx context.Context, id uuid.UUID) error {
	return s.repo.UpdateLastSeen(ctx, id)
}

func (s *UserService) GetLeaderboard(ctx context.Context, limit int) ([]repository.LeaderboardEntry, error) {
	return s.repo.GetLeaderboard(ctx, limit)
}

func (s *UserService) CountUsers(ctx context.Context) (int, error) {
	return s.repo.CountUsers(ctx)
}

func (s *UserService) OnlineCount(ctx context.Context) (int, error) {
	count, err := s.redis.SCard(ctx, "board:online")
	if err != nil {
		return 0, fmt.Errorf("online count: %w", err)
	}
	return int(count), nil
}

func (s *UserService) SetOnline(ctx context.Context, userID uuid.UUID) error {
	if err := s.redis.SAdd(ctx, "board:online", userID.String()); err != nil {
		return fmt.Errorf("set online: %w", err)
	}
	return nil
}

func (s *UserService) SetOffline(ctx context.Context, userID uuid.UUID) error {
	if err := s.redis.SRem(ctx, "board:online", userID.String()); err != nil {
		return fmt.Errorf("set offline: %w", err)
	}
	return nil
}

func (s *UserService) ListOnlineUsers(ctx context.Context) ([]*domain.User, error) {
	ids, err := s.redis.SMembers(ctx, "board:online")
	if err != nil {
		return nil, fmt.Errorf("online users: %w", err)
	}
	if len(ids) == 0 {
		return []*domain.User{}, nil
	}
	parsed := make([]uuid.UUID, 0, len(ids))
	for _, id := range ids {
		uid, err := uuid.Parse(id)
		if err != nil {
			continue
		}
		parsed = append(parsed, uid)
	}
	if len(parsed) == 0 {
		return []*domain.User{}, nil
	}
	return s.repo.GetByIDs(ctx, parsed)
}

func pickColor() string {
	if len(domain.ColorPalette) == 0 {
		return "#999999"
	}
	return domain.ColorPalette[rand.Intn(len(domain.ColorPalette))]
}

func isUniqueViolation(err error) bool {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		if pqErr.Code == "23505" {
			return true
		}
	}
	return false
}
