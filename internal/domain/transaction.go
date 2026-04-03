package domain

import "time"

type TransactionStatus string

const (
	StatusPending   TransactionStatus = "pending"
	StatusCompleted TransactionStatus = "completed"
	StatusFailed    TransactionStatus = "failed"
	StatusCancelled TransactionStatus = "cancelled"
)

// Transaction represents financial transaction
type Transaction struct {
	ID           string
	ClientID     string
	BusinessID   string
	Amount       int64
	BonusAccrued int64
	Currency     string
	Status       TransactionStatus
	Description  *string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time
}
