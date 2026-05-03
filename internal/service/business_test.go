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

func newBusinessService(t *testing.T) (service.BusinessService, *mocks.MockBusinessRepository) {
	repo := mocks.NewMockBusinessRepository(t)
	svc := service.NewBusinessService(repo, testLogger())
	return svc, repo
}

// Create - BusinessService
func TestBusinessService_Create_Success(t *testing.T) {
	svc, repo := newBusinessService(t)
	repo.EXPECT().Create(mock.Anything, activeBusiness().Name, activeBusiness().OwnerPhone).
		Return(activeBusiness(), nil).Once()

	business, err := svc.Create(context.Background(), "Apple", "+77771112233")

	require.NoError(t, err)
	assert.Equal(t, activeBusiness().ID, business.ID)
	assert.Equal(t, "+77771112233", business.OwnerPhone)
	assert.True(t, business.IsActive)
}

func TestBusinessService_Create_MultipleBusinesses(t *testing.T) {
	svc, repo := newBusinessService(t)
	repo.EXPECT().Create(mock.Anything, activeBusiness().Name, activeBusiness().OwnerPhone).
		Return(activeBusiness(), nil).Once()
	repo.EXPECT().Create(mock.Anything, "Samsung", "+77772223344").
		Return(domain.Business{ID: "biz-uuid-1111", Name: "Samsung", OwnerPhone: "+77772223344", IsActive: true}, nil).Once()

	firstBiz, err := svc.Create(context.Background(), "Apple", "+77771112233")

	require.NoError(t, err)
	assert.Equal(t, activeBusiness().ID, firstBiz.ID)
	assert.Equal(t, "+77771112233", firstBiz.OwnerPhone)
	assert.True(t, firstBiz.IsActive)

	secondBiz, err := svc.Create(context.Background(), "Samsung", "+77772223344")

	require.NoError(t, err)
	assert.Equal(t, "biz-uuid-1111", secondBiz.ID)
	assert.Equal(t, "+77772223344", secondBiz.OwnerPhone)
	assert.True(t, secondBiz.IsActive)
}

func TestBusinessService_Create_RepoSystemError(t *testing.T) {
	svc, repo := newBusinessService(t)

	repo.EXPECT().Create(mock.Anything, activeBusiness().Name, activeBusiness().OwnerPhone).
		Return(domain.Business{}, ErrConnectionRefused).Once()

	_, err := svc.Create(context.Background(), activeBusiness().Name, activeBusiness().OwnerPhone)

	require.ErrorIs(t, err, ErrConnectionRefused)
}

// GetByID
func TestBusinessService_GetByID(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name string
		id string
		setup func(repo *mocks.MockBusinessRepository)
		wantErr error
		wantActive *bool
		} {
		{
			name: "success",
			id: activeBusiness().ID,
			setup: func(repo *mocks.MockBusinessRepository) {
				repo.EXPECT().GetByID(mock.Anything, activeBusiness().ID).
				Return(activeBusiness(), nil).
				Once()
			},
		},
		{
			name: "not found",
			id: activeBusiness().ID,
			setup: func(repo *mocks.MockBusinessRepository) {
				repo.EXPECT().GetByID(mock.Anything, activeBusiness().ID).
				Return(domain.Business{}, domain.ErrBusinessNotFound).
				Once()
			},
			wantErr: domain.ErrBusinessNotFound,
		},
		{
			name: "inactive business — service returns it anyway",
			id: activeBusiness().ID,
			setup: func(repo *mocks.MockBusinessRepository) {
				repo.EXPECT().GetByID(mock.Anything, activeBusiness().ID).
				Return(inactiveBusiness(), nil).
				Once()
			},
			wantActive: func() *bool {b := false; return &b}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo := newBusinessService(t)
			tt.setup(repo)

			business, err := svc.GetByID(ctx, tt.id)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.id, business.ID)
			assert.Equal(t, activeBusiness().OwnerPhone, business.OwnerPhone)
			if tt.wantActive != nil {
				assert.Equal(t, *tt.wantActive, business.IsActive)
			}
		})
	}
}

func TestBusinessService_UpdateBonusBalance_TopUp(t *testing.T) {
	svc, repo := newBusinessService(t)
	repo.EXPECT().UpdateBonusBalance(mock.Anything, "biz-uuid-0000", int64(10000)).
		Return(domain.Business{ID: "biz-uuid-0000", BonusBalance: 10000}, nil).Once()

	business, err := svc.UpdateBonusBalance(context.Background(), "biz-uuid-0000", int64(10000))

	require.NoError(t, err)
	assert.Equal(t, "biz-uuid-0000", business.ID)
	assert.Equal(t, int64(10000), business.BonusBalance)
}

func TestBusinessService_UpdateBonusBalance_Deduct(t *testing.T) {
	svc, repo := newBusinessService(t)
	repo.EXPECT().UpdateBonusBalance(mock.Anything, "biz-uuid-0000", int64(-3000)).
		Return(domain.Business{ID: "biz-uuid-0000", BonusBalance: 7000}, nil).Once()
	
	business, err := svc.UpdateBonusBalance(context.Background(), "biz-uuid-0000", int64(-3000))

	require.NoError(t, err)
	assert.Equal(t, "biz-uuid-0000", business.ID)
	assert.Equal(t, int64(7000), business.BonusBalance)
}

func TestBusinessService_UpdateBonusBalance_BusinessNotFound(t *testing.T) {
	svc, repo := newBusinessService(t)
	repo.EXPECT().UpdateBonusBalance(mock.Anything, mock.Anything, mock.Anything).
		Return(domain.Business{}, domain.ErrBusinessNotFound).Once()

	_, err := svc.UpdateBonusBalance(context.Background(), "biz-uuid-0000", 10000)

	require.ErrorIs(t, err, domain.ErrBusinessNotFound)
}

func TestBusinessService_UpdateBonusBalance_ZeroDeltaAllowed(t *testing.T) {
	svc, repo := newBusinessService(t)
	repo.EXPECT().UpdateBonusBalance(mock.Anything, "biz-uuid-0000", int64(0)).
		Return(domain.Business{ID:"biz-uuid-0000", BonusBalance: 0}, nil).Once()

	business, err := svc.UpdateBonusBalance(context.Background(), "biz-uuid-0000", int64(0))

	require.NoError(t, err)
	assert.Equal(t, int64(0), business.BonusBalance)
}
