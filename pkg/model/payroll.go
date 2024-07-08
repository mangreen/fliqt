package model

import (
	"time"
)

type Payroll struct {
	ID        string
	UserID    string
	Salary    int
	Insurance int
	Tax       int
	CreatedAt time.Time
	UpdatedAt time.Time
}
