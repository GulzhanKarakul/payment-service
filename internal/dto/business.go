package dto

import (
	"time"

	"github.com/GulzhanKarakul/payment-service/internal/domain"
)


type BusinessResponse struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	OwnerPhone   string    `json:"owner_phone"`
	BonusBalance int64     `json:"bonus_balance"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func ToBusinessResponse(b domain.Business) BusinessResponse {
	return BusinessResponse{
		ID:           b.ID,
		Name:         b.Name,
		OwnerPhone:   b.OwnerPhone,
		BonusBalance: b.BonusBalance,
		IsActive:     b.IsActive,
		CreatedAt:    b.CreatedAt,
		UpdatedAt:    b.UpdatedAt,
	}
}
