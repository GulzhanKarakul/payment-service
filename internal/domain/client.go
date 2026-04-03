package domain

import "time"

// Client represents application clients
type Client struct {
	ID           string
	Phone        string
	Name         string
	BonusBalance int64
	IsActive     bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
