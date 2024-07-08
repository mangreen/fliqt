package model

import (
	"time"
)

type Performance struct {
	ID        string
	UserID    string
	Rating    string
	CreatedAt time.Time
	UpdatedAt time.Time
}
