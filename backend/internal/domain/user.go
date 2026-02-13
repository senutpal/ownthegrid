package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID         uuid.UUID `db:"id" json:"id"`
	Username   string    `db:"username" json:"username"`
	Color      string    `db:"color" json:"color"`
	CreatedAt  time.Time `db:"created_at" json:"createdAt"`
	LastSeen   time.Time `db:"last_seen" json:"lastSeen"`
	ClaimCount int       `json:"claimCount,omitempty"`
	IsOnline   bool      `json:"isOnline,omitempty"`
	Token      string    `json:"token,omitempty"`
}

var ColorPalette = []string{
	"#FF6B6B", "#FF8E53", "#FFC107", "#CDDC39", "#66BB6A",
	"#26C6DA", "#42A5F5", "#7E57C2", "#EC407A", "#FF7043",
	"#29B6F6", "#26A69A", "#D4E157", "#FFCA28", "#78909C",
	"#EF5350", "#AB47BC", "#5C6BC0", "#00ACC1", "#43A047",
}
