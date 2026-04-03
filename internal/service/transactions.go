package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/GulzhanKarakul/payment-service/internal/domain"
)

type transactionService struct {
	txRepo     TransactionRepository
	clientRepo ClientRepository
	bizRepo    BusinessRepository
	bonusRepo  BonusSettingsRepository
	log        *slog.Logger
}

func NewTransactionService(
	txRepo TransactionRepository,
	clientRepo ClientRepository,
	bizRepo BusinessRepository,
	bonusRepo BonusSettingsRepository,
	log *slog.Logger,
) TransactionService {
	return &transactionService{
		txRepo:     txRepo,
		clientRepo: clientRepo,
		bizRepo:    bizRepo,
		bonusRepo:  bonusRepo,
		log:        log,
	}
}

func (s *transactionService) Create(
	ctx context.Context,
	businessID,
	clientID string,
	amount int64,
	description *string,
) (domain.Transaction, error) {
	client, err := s.clientRepo.GetByID(ctx, clientID)
	if err != nil {
		return domain.Transaction{}, fmt.Errorf("transactionService.Create: %w", err)
	}

	if !client.IsActive {
		return domain.Transaction{}, domain.ErrClientIsNotActive
	}

	business, err := s.bizRepo.GetByID(ctx, businessID)
	if err != nil {
		return domain.Transaction{}, fmt.Errorf("transactionService.Create: %w", err)
	}

	if !business.IsActive {
		return domain.Transaction{}, domain.ErrBusinessIsNotActive
	}

	bonusSettings, err := s.bonusRepo.GetByBusinessID(ctx, businessID)
	if err != nil && !errors.Is(err, domain.ErrBonusSettingsNotFound) {
		return domain.Transaction{}, fmt.Errorf("transactionService.Create: %w", err)
	}

	var bonusPercent float64
	if err == nil && bonusSettings.IsActive {
		bonusPercent = bonusSettings.BonusPercent
	}

	tx, err := s.txRepo.CreateWithBonus(ctx, clientID, businessID, amount, bonusPercent, description)
	if err != nil {
		return domain.Transaction{}, fmt.Errorf("transactionService.Create: %w", err)
	}

	return tx, nil
}

func (s *transactionService) GetByID(ctx context.Context, id string) (domain.Transaction, error) {
	tx, err := s.txRepo.GetByID(ctx, id)
	if err != nil {
		return domain.Transaction{}, fmt.Errorf("transactionService.GetByID: %w", err)
	}

	return tx, nil
}

func (s *transactionService) GetByClientID(ctx context.Context, clientID string, limit, offset int) ([]domain.Transaction, error) {
	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	txSl, err := s.txRepo.GetByClientID(ctx, clientID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("transactionService.GetByClientID: %w", err)
	}

	return txSl, nil
}

func (s *transactionService) Cancel(ctx context.Context, id string) error {
	if err := s.txRepo.Cancel(ctx, id); err != nil {
		return fmt.Errorf("transactionService.Cancel: %w", err)
	}
	return nil
}
