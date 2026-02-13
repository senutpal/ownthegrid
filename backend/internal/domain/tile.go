package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrTileInvalid        = errors.New("tile ID is out of range")
	ErrTileAlreadyClaimed = errors.New("tile is already claimed")
)

type Tile struct {
	ID            int        `db:"id" json:"id"`
	X             int        `db:"x" json:"x"`
	Y             int        `db:"y" json:"y"`
	OwnerID       *uuid.UUID `db:"owner_id" json:"ownerId"`
	ClaimedAt     *time.Time `db:"claimed_at" json:"claimedAt"`
	OwnerUsername *string    `db:"owner_username" json:"ownerUsername,omitempty"`
	OwnerColor    *string    `db:"owner_color" json:"ownerColor,omitempty"`
}

type ClaimRequest struct {
	TileID int       `json:"tileId"`
	UserID uuid.UUID `json:"userId"`
}
