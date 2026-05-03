package service_test

import (
	"errors"
	"log/slog"

	"github.com/GulzhanKarakul/payment-service/internal/domain"
)

var (
	ErrConnectionRefused = errors.New("sys error: connection refused")
	ErrDBNotAvailable = errors.New("db error")
)

func strPtr(s string) *string {
	return &s
}

func testLogger() *slog.Logger {
	return slog.New(slog.DiscardHandler)
}

func activeClient() domain.Client {
	return domain.Client{ID: "cli-uuid-0000", Name: "Gulzhan", Phone: "+77771156580", IsActive: true}
}

func inactiveClient() domain.Client {
	return domain.Client{ID: "cli-uuid-0000", Name: "Gulzhan", Phone: "+77771156580", IsActive: false}
}

func activeBusiness() domain.Business {
	return domain.Business{ID: "biz-uuid-0000", Name: "Apple", OwnerPhone: "+77771112233", IsActive: true}
}

func inactiveBusiness() domain.Business {
	return domain.Business{ID: "biz-uuid-0000", Name: "Apple", OwnerPhone: "+77771112233", IsActive: false}
}

func activeBonus(percent float64) domain.BonusSettings {
	return domain.BonusSettings{ID: "bonus-uuid-0000", BusinessID: "biz-uuid-0000", BonusPercent: percent, IsActive: true}
}

func inactiveBonus(percent float64) domain.BonusSettings {
	return domain.BonusSettings{ID: "bonus-uuid-0000", BusinessID: "biz-uuid-0000", BonusPercent: percent, IsActive: false}
}

func testTransaction() domain.Transaction {
	return domain.Transaction{ID: "tx-uuid", Status: domain.StatusCompleted}
}