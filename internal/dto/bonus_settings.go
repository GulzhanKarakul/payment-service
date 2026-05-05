package dto

import (
	"time"

	"github.com/GulzhanKarakul/payment-service/internal/domain"
)


type BonusSettingsResponse struct {
	ID           string    `json:"id"`
	BusinessID   string    `json:"business_id"`
	BonusPercent float64   `json:"bonus_percent"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func ToBonusSettingsResponse(bs domain.BonusSettings) BonusSettingsResponse {
	return BonusSettingsResponse{
		ID:           bs.ID,
		BusinessID:   bs.BusinessID,
		BonusPercent: bs.BonusPercent,
		IsActive:     bs.IsActive,
		CreatedAt:    bs.CreatedAt,
		UpdatedAt:    bs.UpdatedAt,
	}
}