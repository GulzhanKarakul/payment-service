package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/GulzhanKarakul/payment-service/internal/domain"
)

type clientService struct {
	repo ClientRepository
	log  *slog.Logger
}

func NewClientService(repo ClientRepository, log *slog.Logger) ClientService {
	return &clientService{repo: repo, log: log}
}

func (s *clientService) Create(ctx context.Context, phone, name string) (domain.Client, error) {
	_, err := s.repo.GetByPhone(ctx, phone)
	if err == nil {
		return domain.Client{}, domain.ErrClientAlreadyExist
	}
	if !errors.Is(err, domain.ErrClientNotFound) {
		return domain.Client{}, fmt.Errorf("clientService.Create: %w", err)
	}

	client, err := s.repo.Create(ctx, phone, name)
	if err != nil {
		return domain.Client{}, fmt.Errorf("clientService.Create: %w", err)
	}

	return client, nil
}

func (s *clientService) GetByID(ctx context.Context, id string) (domain.Client, error) {
	client, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return domain.Client{}, fmt.Errorf("clientService.GetByID: %w", err)
	}

	return client, nil
}

func (s *clientService) GetByPhone(ctx context.Context, phone string) (domain.Client, error) {
	client, err := s.repo.GetByPhone(ctx, phone)
	if err != nil {
		return domain.Client{}, fmt.Errorf("clientService.GetByPhone: %w", err)
	}

	return client, nil
}
