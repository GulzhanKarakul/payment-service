package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/GulzhanKarakul/payment-service/internal/domain"
)

type businessService struct {
	repo BusinessRepository
	log *slog.Logger
}

func NewBusinessService(repo BusinessRepository, log *slog.Logger) BusinessService {
	return &businessService{repo: repo, log: log}
}

func (s *businessService) Create(ctx context.Context, name, ownerPhone string) (domain.Business, error) {
	business, err := s.repo.Create(ctx, name, ownerPhone)
	if err != nil {
		return domain.Business{}, fmt.Errorf("businessService.Create: %w", err)
	}

	return business, nil
}

func (s *businessService) GetByID(ctx context.Context, id string) (domain.Business, error) {
	business, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return domain.Business{}, fmt.Errorf("businessService.GetByID: %w", err)
	}

	return business, nil
}
