package service

import (
	"context"

	"github.com/GulzhanKarakul/payment-service/internal/domain"
)

// ClientService defines business logic contract for clients
type ClientService interface {
	Create(ctx context.Context, name, phone string) (domain.Client, error)
	GetByID(ctx context.Context, id string) (domain.Client, error)
	GetByPhone(ctx context.Context, phone string) (domain.Client, error)
}

// BusinessService defines business logic contract for businesses
type BusinessService interface {
	Create(ctx context.Context, name, ownerPhone string) (domain.Business, error)
	GetByID(ctx context.Context, id string) (domain.Business, error)
}

// BonusSettingsService defines business logic contract for bonus settings
type BonusSettingsService interface {
	Upsert(ctx context.Context, businessId string, bonusPercent float64, isActive bool) (domain.BonusSettings, error)
	GetByBusinessID(ctx context.Context, businessId string) (domain.BonusSettings, error)
}

// TransactionService defines business logic contract for transactions
type TransactionService interface {
	Create(ctx context.Context, businessId, clientId string, amount int64, description *string) (domain.Transaction, error)
	GetByID(ctx context.Context, id string) (domain.Transaction, error)
	GetByClientID(ctx context.Context, clientId string, limit, offset int) ([]domain.Transaction, error)
	Cancel(ctx context.Context, id string) error
}
