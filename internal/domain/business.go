package domain

import "time"

// Business represents application businesses
type Business struct {
	ID string
	Name string
	OwnerPhone string
	BonusBalance int64
	IsActive bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
