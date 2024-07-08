package model

import (
	"time"
)

type Leave struct {
	ID        string
	UserID    string
	Type      string
	StartedAt time.Time
	EndedAt   time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}
