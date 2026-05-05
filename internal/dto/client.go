package dto

import (
	"time"

	"github.com/GulzhanKarakul/payment-service/internal/domain"
)

type ClientResponse struct {
	ID           string    `json:"id"`
	Phone        string    `json:"phone"`
	Name         string    `json:"name"`
	BonusBalance int64     `json:"bonus_balance"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func ToClientResponse(c domain.Client) ClientResponse {
	return ClientResponse{
		ID:           c.ID,
		Phone:        c.Phone,
		Name:         c.Name,
		BonusBalance: c.BonusBalance,
		IsActive:     c.IsActive,
		CreatedAt:    c.CreatedAt,
		UpdatedAt:    c.UpdatedAt,
	}
}