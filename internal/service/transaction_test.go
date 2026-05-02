package service_test

import (
	"context"
	"testing"

	"github.com/GulzhanKarakul/payment-service/internal/domain"
	"github.com/GulzhanKarakul/payment-service/internal/service"
	"github.com/GulzhanKarakul/payment-service/internal/service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func newTransactionSvc(t *testing.T) (
	service.TransactionService,
	*mocks.MockTransactionRepository,
	*mocks.MockClientRepository,
	*mocks.MockBusinessRepository,
	*mocks.MockBonusSettingsRepository,
) {
	t.Helper()
	txRepo := mocks.NewMockTransactionRepository(t)
	cliRepo := mocks.NewMockClientRepository(t)
	bizRepo := mocks.NewMockBusinessRepository(t)
	bonusRepo := mocks.NewMockBonusSettingsRepository(t)
	svc := service.NewTransactionService(txRepo, cliRepo, bizRepo, bonusRepo, testLogger())
	return svc, txRepo, cliRepo, bizRepo, bonusRepo
}

const (
	testTxID = "tx-uuid"
	testClientID = "cli-uuid-0000"
	testBusinessID = "biz-uuid-0000"
	testAmount = int64(100000)
)

func TestTransactionService_Create(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name string
		description *string
		setup func(txRepo *mocks.MockTransactionRepository, clientRepo *mocks.MockClientRepository, bizRepo *mocks.MockBusinessRepository, bsRepo *mocks.MockBonusSettingsRepository)
		wantErr error
	} {
		{
			name: "success with bonus",
			setup: func(
				txRepo *mocks.MockTransactionRepository,
				clientRepo *mocks.MockClientRepository,
				bizRepo *mocks.MockBusinessRepository,
				bsRepo *mocks.MockBonusSettingsRepository,
			) {
				clientRepo.EXPECT().GetByID(mock.Anything, testClientID).
					Return(activeClient(), nil).Once()
				bizRepo.EXPECT().GetByID(mock.Anything, testBusinessID).
					Return(activeBusiness(), nil).Once()
				bsRepo.EXPECT().GetByBusinessID(mock.Anything, testBusinessID).
					Return(activeBonus(10.0), nil).Once()
				txRepo.EXPECT().CreateWithBonus(mock.Anything, testClientID, testBusinessID, testAmount, 10.0, mock.Anything).
					Return(testTransaction(), nil).Once()
			},
		},
		{
			name: "success no bonus settings - bonusPercent=0",
			setup: func(
				txRepo *mocks.MockTransactionRepository,
				clientRepo *mocks.MockClientRepository,
				bizRepo *mocks.MockBusinessRepository,
				bsRepo *mocks.MockBonusSettingsRepository,
			) {
				clientRepo.EXPECT().GetByID(mock.Anything, testClientID).
					Return(activeClient(), nil).Once()
				bizRepo.EXPECT().GetByID(mock.Anything, testBusinessID).
					Return(activeBusiness(), nil).Once()
				bsRepo.EXPECT().GetByBusinessID(mock.Anything, testBusinessID).
					Return(domain.BonusSettings{}, domain.ErrBonusSettingsNotFound).Once()
				txRepo.EXPECT().CreateWithBonus(mock.Anything, testClientID, testBusinessID, testAmount, 0.0, mock.Anything).
					Return(testTransaction(), nil).Once()
			},
		},
		{
			name: "success - bonus settings inactive - bonusPercent=0",
			setup: func(
				txRepo *mocks.MockTransactionRepository,
				clientRepo *mocks.MockClientRepository,
				bizRepo *mocks.MockBusinessRepository,
				bsRepo *mocks.MockBonusSettingsRepository,
			) {
				clientRepo.EXPECT().GetByID(mock.Anything, testClientID).
					Return(activeClient(), nil).Once()
				bizRepo.EXPECT().GetByID(mock.Anything, testBusinessID).
					Return(activeBusiness(), nil).Once()
				bsRepo.EXPECT().GetByBusinessID(mock.Anything, testBusinessID).
					Return(inactiveBonus(0.0), nil).Once()
				txRepo.EXPECT().CreateWithBonus(mock.Anything, testClientID, testBusinessID, testAmount, 0.0, mock.Anything).
					Return(testTransaction(), nil).Once()
			},
		},
		{
			name: "success - with description",
			description: strPtr("Оплата заказа"),
			setup: func(
				txRepo *mocks.MockTransactionRepository,
				clientRepo *mocks.MockClientRepository,
				bizRepo *mocks.MockBusinessRepository,
				bsRepo *mocks.MockBonusSettingsRepository,
			) {
				clientRepo.EXPECT().GetByID(mock.Anything, testClientID).
					Return(activeClient(), nil).Once()
				bizRepo.EXPECT().GetByID(mock.Anything, testBusinessID).
					Return(activeBusiness(), nil).Once()
				bsRepo.EXPECT().GetByBusinessID(mock.Anything, testBusinessID).
					Return(activeBonus(10.0), nil).Once()
				txRepo.EXPECT().CreateWithBonus(mock.Anything, testClientID, testBusinessID, testAmount, 10.0, mock.Anything).
					Return(testTransaction(), nil).Once()
			},
		},
		{
			name: "client not found",
			setup: func(
				txRepo *mocks.MockTransactionRepository,
				clientRepo *mocks.MockClientRepository,
				bizRepo *mocks.MockBusinessRepository,
				bsRepo *mocks.MockBonusSettingsRepository,
			) {
				clientRepo.EXPECT().GetByID(mock.Anything, testClientID).
					Return(domain.Client{}, domain.ErrClientNotFound).Once()
			},
			wantErr: domain.ErrClientNotFound,
		},
		{
			name: "client inactive",
			setup: func(
				txRepo *mocks.MockTransactionRepository,
				clientRepo *mocks.MockClientRepository,
				bizRepo *mocks.MockBusinessRepository,
				bsRepo *mocks.MockBonusSettingsRepository,
			) {
				clientRepo.EXPECT().GetByID(mock.Anything, testClientID).
					Return(inactiveClient(), nil).Once()
			},
			wantErr: domain.ErrClientIsNotActive,
		},
		{
			name: "business not found",
			setup: func(
				txRepo *mocks.MockTransactionRepository,
				clientRepo *mocks.MockClientRepository,
				bizRepo *mocks.MockBusinessRepository,
				bsRepo *mocks.MockBonusSettingsRepository,
			) {
				clientRepo.EXPECT().GetByID(mock.Anything, testClientID).
					Return(activeClient(), nil).Once()
				bizRepo.EXPECT().GetByID(mock.Anything, testBusinessID).
					Return(domain.Business{}, domain.ErrBusinessNotFound).Once()
			},
			wantErr: domain.ErrBusinessNotFound,
		},
		{
			name: "business inactive",
			setup: func(
				txRepo *mocks.MockTransactionRepository,
				clientRepo *mocks.MockClientRepository,
				bizRepo *mocks.MockBusinessRepository,
				bsRepo *mocks.MockBonusSettingsRepository,
			) {
				clientRepo.EXPECT().GetByID(mock.Anything, testClientID).
					Return(activeClient(), nil).Once()
				bizRepo.EXPECT().GetByID(mock.Anything, testBusinessID).
					Return(inactiveBusiness(), nil).Once()
			},
			wantErr: domain.ErrBusinessIsNotActive,
		},
		{
			name: "insufficent balance",
			setup: func(
				txRepo *mocks.MockTransactionRepository,
				clientRepo *mocks.MockClientRepository,
				bizRepo *mocks.MockBusinessRepository,
				bsRepo *mocks.MockBonusSettingsRepository,
			) {
				clientRepo.EXPECT().GetByID(mock.Anything, testClientID).
					Return(activeClient(), nil).Once()
				bizRepo.EXPECT().GetByID(mock.Anything, testBusinessID).
					Return(activeBusiness(), nil).Once()
				bsRepo.EXPECT().GetByBusinessID(mock.Anything, testBusinessID).
					Return(activeBonus(5.0), nil).Once()
				txRepo.EXPECT().CreateWithBonus(mock.Anything, testClientID, testBusinessID, testAmount, 5.0, mock.Anything).
					Return(domain.Transaction{}, domain.ErrInsufficientBonusBalance).Once()
			},
			wantErr: domain.ErrInsufficientBonusBalance,
		},
		{
			name: "bonus settings sys error - tx aborted",
			description: func() *string {description := "Оплата заказа"; return &description}(),
			setup: func(
				txRepo *mocks.MockTransactionRepository,
				clientRepo *mocks.MockClientRepository,
				bizRepo *mocks.MockBusinessRepository,
				bsRepo *mocks.MockBonusSettingsRepository,
			) {
				clientRepo.EXPECT().GetByID(mock.Anything, testClientID).
					Return(activeClient(), nil).Once()
				bizRepo.EXPECT().GetByID(mock.Anything, testBusinessID).
					Return(activeBusiness(), nil).Once()
				bsRepo.EXPECT().GetByBusinessID(mock.Anything, testBusinessID).
					Return(domain.BonusSettings{}, ErrConnectionRefused).Once()
			},
			wantErr: ErrConnectionRefused,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, txRepo, clientRepo, bizRepo, bsRepo := newTransactionSvc(t)
			tt.setup(txRepo, clientRepo, bizRepo, bsRepo)

			createdTx, err := svc.Create(ctx, testClientID, testBusinessID, testAmount, tt.description)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}

			require.NoError(t, err)
			assert.NotEmpty(t, createdTx.ID)
			assert.Equal(t, domain.StatusCompleted, createdTx.Status)
		})
	}
}

func TestTransactionService_GetByID(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name string
		setup func(repo *mocks.MockTransactionRepository)
		wantErr error
	}{
		{
			name: "success",
			setup: func(repo *mocks.MockTransactionRepository) {
				repo.EXPECT().GetByID(mock.Anything, testTxID).
				Return(testTransaction(), nil).Once()
			},
		},
		{
			name: "not found",
			setup: func(repo *mocks.MockTransactionRepository) {
				repo.EXPECT().GetByID(mock.Anything, testTxID).
				Return(domain.Transaction{}, domain.ErrTransactionNotFound).Once()
			},
			wantErr: domain.ErrTransactionNotFound,
		},
		{
			name: "system error",
			setup: func(repo *mocks.MockTransactionRepository) {
				repo.EXPECT().GetByID(mock.Anything, testTxID).
				Return(domain.Transaction{}, ErrDBNotAvailable).Once()
			},
			wantErr: ErrDBNotAvailable,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, txRepo, _, _, _ := newTransactionSvc(t)
			tt.setup(txRepo)

			tx, err := svc.GetByID(ctx, testTxID)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, testTxID, tx.ID)
			assert.Equal(t, domain.StatusCompleted, tx.Status)
		})
	}
}

func TestTransactionService_GetByClientID_Pagination(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name string
		limit int
		offset int
		setup func(repo *mocks.MockTransactionRepository)
		wantErr error
	}{
		{
			name: "valid limit and offset",
			limit: 5,
			offset: 10,
			setup: func(repo *mocks.MockTransactionRepository) {
				result := getResult(10)

				repo.EXPECT().GetByClientID(mock.Anything, testClientID, 5, 10).
				Return(result, nil)
			},
		},
		{
			name: "zero limit replaced with default 20",
			limit: 0,
			offset: 0,
			setup: func(repo *mocks.MockTransactionRepository) {
				result := getResult(10)

				repo.EXPECT().GetByClientID(mock.Anything, testClientID, 20, 0).
					Return(result, nil)
			},
		},
		{
			name: "negative limit replaced with default 20",
			limit: -5,
			offset: 0,
			setup: func(repo *mocks.MockTransactionRepository) {
				result := getResult(10)

				repo.EXPECT().GetByClientID(mock.Anything, testClientID, 20, 0).
					Return(result, nil)
			},
		},
		{
			name: "negative offset replaced with default 20",
			limit: -5,
			offset: -1,
			setup: func(repo *mocks.MockTransactionRepository) {
				result := getResult(10)

				repo.EXPECT().GetByClientID(mock.Anything, testClientID, 20, 0).
					Return(result, nil)
			},
		},
		{
			name: "empty result is valid",
			limit: -5,
			offset: -1,
			setup: func(repo *mocks.MockTransactionRepository) {
				repo.EXPECT().GetByClientID(mock.Anything, testClientID, 20, 0).
					Return(nil, nil)
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			svc, txRepo, _, _, _ := newTransactionSvc(t)
			tt.setup(txRepo)

			_, err := svc.GetByClientID(ctx, testClientID, tt.limit, tt.offset)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}

			require.NoError(t, err)
		})
	}
}

func TestTransactionService_Cancel(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name string
		txID string
		setup func(repo *mocks.MockTransactionRepository)
		wantErr error
	}{
		{
			name: "success",
			setup: func(repo *mocks.MockTransactionRepository) {
				repo.EXPECT().Cancel(mock.Anything, testTxID).
				Return(nil).Once()
			},
		},
		{
			name: "not found",
			setup: func(repo *mocks.MockTransactionRepository) {
				repo.EXPECT().Cancel(mock.Anything, testTxID).
					Return(domain.ErrTransactionNotFound).Once()
			},
			wantErr: domain.ErrTransactionNotFound,
		},
		{
			name: "already cancelled",
			setup: func(repo *mocks.MockTransactionRepository) {
				repo.EXPECT().Cancel(mock.Anything, testTxID).
					Return(domain.ErrTransactionCancelled).Once()
			},
			wantErr: domain.ErrTransactionCancelled,
		},
		{
			name: "system error",
			setup: func(repo *mocks.MockTransactionRepository) {
				repo.EXPECT().Cancel(mock.Anything, testTxID).
					Return(ErrConnectionRefused).Once()
			},
			wantErr: ErrConnectionRefused,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, txRepo, _, _, _ := newTransactionSvc(t)
			tt.setup(txRepo)

			err := svc.Cancel(ctx, testTxID)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}

			require.NoError(t, err)
		})
	}
}

func getResult(n int) []domain.Transaction {
	result := make([]domain.Transaction, 0, n)
		for i:=0; i<n; i++ {
			result = append(result, testTransaction())
		}
	return result
}