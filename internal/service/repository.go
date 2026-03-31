package service

import (
	"context"

	"github.com/GulzhanKarakul/payment-service/internal/domain"
)

// ClientRepository defines data access contract for clients.
type ClientRepository interface {
	Create(ctx context.Context, name, phone string,)(domain.Client, error) 
	GetByID(ctx context.Context, id string) (domain.Client, error)
	GetByPhone(ctx context.Context, phone string) (domain.Client, error)
	UpdateBonusBalance(ctx context.Context, id string, bonus int64) (domain.Client, error)
}

// BusinessRepository defines data access contract for businesses
type BusinessRepository interface {
	Create(ctx context.Context, name, ownerPhone string) (domain.Business, error)
	GetByID(ctx context.Context, id string) (domain.Business, error)
	UpdateBonusBalance(ctx context.Context, id string, delta int64) (domain.Business, error)
}

// BonusSettingsrepository defines data access contract for bonus settings
type BonusSettingsRepository interface {
	Upsert(ctx context.Context, businessID string, bonusPercent float64, isActive bool) (domain.BonusSettings, error)
	GetByBusinessID(ctx context.Context, businessID string) (domain.BonusSettings, error)
}

// TransactionRepository defines data access contract for transactions
type TransactionRepository interface {
	CreateWithBonus(ctx context.Context, clientId, businessId string, amount int64, bonusPercent float64, description *string) (domain.Transaction, error)
	GetByID(ctx context.Context, id string) (domain.Transaction, error)
	GetByClientID(ctx context.Context, clientId string, limit, offset int) ([]domain.Transaction, error)
	Cancel(ctx context.Context, id string) error
}