package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/GulzhanKarakul/payment-service/internal/domain"
)

type bonusSettingsService struct {
	repo    BonusSettingsRepository
	bizRepo BusinessRepository
	log     *slog.Logger
}

func NewBonusSettingsService(repo BonusSettingsRepository, bizRepo BusinessRepository, log *slog.Logger) BonusSettingsService {
	return &bonusSettingsService{repo: repo, bizRepo: bizRepo, log: log}
}

func (s *bonusSettingsService) Upsert(
	ctx context.Context,
	businessID string,
	bonusPercent float64,
	isActive bool,
) (domain.BonusSettings, error) {
	_, err := s.bizRepo.GetByID(ctx, businessID)
	if err != nil {
		return domain.BonusSettings{}, fmt.Errorf("bonusSettingsService.Upsert: %w", err)
	}

	if bonusPercent <= 0 {
		return domain.BonusSettings{}, fmt.Errorf("bonusSettingsService.Upsert: %w", domain.ErrInvalidInput)
	}

	bs, err := s.repo.Upsert(ctx, businessID, bonusPercent, isActive)
	if err != nil {
		return domain.BonusSettings{}, fmt.Errorf("bonusSettingsService.Upsert: %w", err)
	}

	return bs, nil
}

func (s *bonusSettingsService) GetByBusinessID(ctx context.Context, businessID string) (domain.BonusSettings, error) {
	bs, err := s.repo.GetByBusinessID(ctx, businessID)
	if err != nil {
		return domain.BonusSettings{}, fmt.Errorf("bonusSettingsService.GetByID: %w", err)
	}

	return bs, nil
}
