package dto

import (
	"time"

	"github.com/GulzhanKarakul/payment-service/internal/domain"
)


type TransactionResponse struct {
	ID           string                   `json:"id"`
	BusinessID   string                   `json:"business_id"`
	ClientID     string                   `json:"client_id"`
	Amount       int64                    `json:"amount"`
	BonusAccrued int64                    `json:"bonus_accrued"`
	Currency     string                   `json:"currency"`
	Status       domain.TransactionStatus `json:"status"`
	Description  *string                  `json:"description"`
	CreatedAt    time.Time                `json:"created_at"`
	UpdatedAt    time.Time                `json:"updated_at"`
	DeletedAt    *time.Time               `json:"deleted_at,omitempty"`
}

func ToTransactionResponse(t domain.Transaction) TransactionResponse {
	return TransactionResponse{
		ID:           t.ID,
		BusinessID:   t.BusinessID,
		ClientID:     t.ClientID,
		Amount:       t.Amount,
		BonusAccrued: t.BonusAccrued,
		Currency:     t.Currency,
		Status:       t.Status,
		Description:  t.Description,
		CreatedAt:    t.CreatedAt,
		UpdatedAt:    t.UpdatedAt,
		DeletedAt:    t.DeletedAt,
	}
}