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

func newBonusSettingsService(t *testing.T) (service.BonusSettingsService, *mocks.MockBonusSettingsRepository, *mocks.MockBusinessRepository) {
	repo := mocks.NewMockBonusSettingsRepository(t)
	bizRepo := mocks.NewMockBusinessRepository(t)
	svc := service.NewBonusSettingsService(repo, bizRepo, testLogger())
	return svc, repo, bizRepo
}

// Upsert - 
func TestBonusSettingsService_Upsert_CreateThenUpdate(t *testing.T) {
	svc, bsRepo, bizRepo := newBonusSettingsService(t)
	bizRepo.EXPECT().GetByID(mock.Anything, activeBusiness().ID).
		Return(activeBusiness(), nil).Twice()
	bsRepo.EXPECT().Upsert(mock.Anything, activeBusiness().ID, float64(12), true).
		Return(activeBonus(12.0), nil).Once()
	bsRepo.EXPECT().Upsert(mock.Anything, activeBusiness().ID, float64(7.5), false).
		Return(inactiveBonus(7.5), nil).Once()

	created, err := svc.Upsert(context.Background(), activeBusiness().ID, 12.0, true)

	require.NoError(t, err)
	assert.Equal(t, float64(12), created.BonusPercent)
	assert.True(t, created.IsActive)

	updated, err := svc.Upsert(context.Background(), activeBusiness().ID, 7.5, false)

	require.NoError(t, err)
	assert.Equal(t, float64(7.5), updated.BonusPercent)
	assert.False(t, updated.IsActive)
}

func TestBonusSettingsService_Upsert_BusinessNotFound(t *testing.T) {
	svc, bsRepo, bizRepo := newBonusSettingsService(t)
	bizRepo.EXPECT().GetByID(mock.Anything, mock.Anything).
		Return(domain.Business{}, domain.ErrBusinessNotFound).Once()
	
	_, err := svc.Upsert(context.Background(), "biz-uuid-0000", float64(6.5), true)

	require.ErrorIs(t, err, domain.ErrBusinessNotFound)
	bsRepo.AssertNotCalled(t, "Upsert", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

// не знаю есть ли смысл проверять все равно что активен бизнес что неактивен все равно настройки отобразятся
func TestBonusSettingsService_Upsert_BusinessInactive(t *testing.T) {
	svc, bsRepo, bizRepo := newBonusSettingsService(t)
	bizRepo.EXPECT().GetByID(mock.Anything, inactiveBusiness().ID).
		Return(inactiveBusiness(), nil).Once()
	bsRepo.EXPECT().Upsert(mock.Anything, activeBusiness().ID, float64(6.5), true).
		Return(activeBonus(6.5), nil).Once()

	bs, err := svc.Upsert(context.Background(), inactiveBusiness().ID, 6.5, true)

	require.NoError(t, err)
	assert.Equal(t, inactiveBusiness().ID, bs.BusinessID)
	assert.Equal(t, float64(6.5), bs.BonusPercent)
}

func TestBonusSettingsService_Upsert_InvalidZeroPercent(t *testing.T) {
	svc, bsRepo, bizRepo := newBonusSettingsService(t)

	_, err := svc.Upsert(context.Background(), activeBusiness().ID, 0.0, true)

	require.ErrorIs(t, err, domain.ErrInvalidInput)
	bsRepo.AssertNotCalled(t, "Upsert", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	bizRepo.AssertNotCalled(t, "GetByID", mock.Anything, mock.Anything)
}

func TestBonusSettingsService_Upsert_InvalidNegativePercent(t *testing.T) {
	svc, bsRepo, bizRepo := newBonusSettingsService(t)

	_, err := svc.Upsert(context.Background(), activeBusiness().ID, -3.0, true)

	require.ErrorIs(t, err, domain.ErrInvalidInput)
	bsRepo.AssertNotCalled(t, "Upsert", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	bizRepo.AssertNotCalled(t, "GetByID", mock.Anything, mock.Anything)
}

func TestBonusSettingsService_Upsert_PropagateSystemError(t *testing.T) {
	svc, bsRepo, bizRepo := newBonusSettingsService(t)
	bizRepo.EXPECT().GetByID(mock.Anything, mock.Anything).Return(domain.Business{}, ErrConnectionRefused).Once()

	_, err := svc.Upsert(context.Background(), "biz-uuid-0000", float64(6.0), true)

	require.ErrorIs(t, err, ErrConnectionRefused)
	bsRepo.AssertNotCalled(t, "Upsert", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

// GetByBusinessID
func TestBonusSettingsService_GetByBusinessID(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name string
		businessID string
		setup func(repo *mocks.MockBonusSettingsRepository)
		wantActive *bool
		wantErr error
	} {
		{
			name: "success",
			businessID: "biz-uuid-0000",
			setup: func(repo *mocks.MockBonusSettingsRepository) {
				repo.EXPECT().GetByBusinessID(mock.Anything, "biz-uuid-0000").
				Return(activeBonus(6.0), nil).
				Once()
			},
		},
		{
			name: "not found",
			businessID: "biz-uuid-0000",
			setup: func(repo *mocks.MockBonusSettingsRepository) {
				repo.EXPECT().GetByBusinessID(mock.Anything, "biz-uuid-0000").
				Return(domain.BonusSettings{}, domain.ErrBonusSettingsNotFound).
				Once()
			},
			wantErr: domain.ErrBonusSettingsNotFound,
		},
		{
			name: "inactive settings returns as is",
			businessID: "biz-uuid-0000",
			setup: func(repo *mocks.MockBonusSettingsRepository) {
				repo.EXPECT().GetByBusinessID(mock.Anything, "biz-uuid-0000").
				Return(inactiveBonus(6.0), nil).
				Once()
			},
			wantActive: func() *bool {b := false; return &b}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo, _ := newBonusSettingsService(t)
			tt.setup(repo)

			bs, err := svc.GetByBusinessID(ctx, tt.businessID)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}

			require.NoError(t, err)
			assert.NotEmpty(t, bs.ID)
			assert.Equal(t, tt.businessID, bs.BusinessID)
			if tt.wantActive != nil {
				assert.Equal(t, *tt.wantActive, bs.IsActive)
			}
		})
	}
}