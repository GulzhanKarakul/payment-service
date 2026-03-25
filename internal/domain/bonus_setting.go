package domain

import "time"

// BonusSettings represents bonus program configuration for a business
type BonusSettings struct {
	ID string
	BusinessID string
	BonusPercent float64
	IsActive bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
