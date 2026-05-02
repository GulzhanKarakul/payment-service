package repository_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/GulzhanKarakul/payment-service/internal/domain"
	"github.com/GulzhanKarakul/payment-service/internal/repository"
)

func TestBusinessRepo_Create(t *testing.T) {
	repo := repository.NewBusinessRepository(testDB)
	ctx := context.Background()

	tests := []struct {
		name         string
		businessName string
		ownerPhone   string
		wantErr      error
	}{
		{
			name:         "success",
			businessName: "Samsung",
			ownerPhone:   "+77771112233",
			wantErr:      nil,
		}, {
			name:         "same owner phone is allowed",
			businessName: "Samsung",
			ownerPhone:   "+77771112233",
			wantErr:      nil,
		},
	}

	for _, tt := range tests {
		cleanDB(t)
		t.Run(tt.name, func(t *testing.T) {
			business, err := repo.Create(ctx, tt.businessName, tt.ownerPhone)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.businessName, business.Name)
			assert.Equal(t, tt.ownerPhone, business.OwnerPhone)
			assert.True(t, business.IsActive)
			assert.NotZero(t, business.CreatedAt)
			assert.Zero(t, business.BonusBalance)
		})
	}
}

func TestBusinessRepo_GetByID(t *testing.T) {
	cleanDB(t)

	repo := repository.NewBusinessRepository(testDB)
	ctx := context.Background()

	created := createTestBusiness(t, "samsung", "+77771112233")

	t.Run("success", func(t *testing.T) {
		business, err := repo.GetByID(ctx, created.ID)
		require.NoError(t, err)

		assert.Equal(t, created.ID, business.ID)
		assert.Equal(t, created.OwnerPhone, business.OwnerPhone)
		assert.Equal(t, created.Name, business.Name)
		assert.Equal(t, created.BonusBalance, business.BonusBalance)
		assert.Equal(t, created.IsActive, business.IsActive)
	})

	t.Run("not found", func(t *testing.T) {
		_, err := repo.GetByID(ctx, "00000000-0000-0000-0000-000000000000")
		require.ErrorIs(t, err, domain.ErrBusinessNotFound)
	})
}

func TestBusinessRepo_UpdateBonusBalance(t *testing.T) {
	cleanDB(t)

	repo := repository.NewBusinessRepository(testDB)
	ctx := context.Background()

	t.Run("update balance", func(t *testing.T) {
		created := createTestBusiness(t, "Samsung", "+7777112233")
		
		updated, err := repo.UpdateBonusBalance(ctx, created.ID, 50000)
		require.NoError(t, err)

		assert.Equal(t, int64(50000), updated.BonusBalance)
	})

	t.Run("deduct balance", func(t *testing.T) {
		created := createTestBusiness(t, "Apple", "+7777112244")
		
		_, err := repo.UpdateBonusBalance(ctx, created.ID, 50000)
		require.NoError(t, err)

		updated, err := repo.UpdateBonusBalance(ctx, created.ID, -30000)
		require.NoError(t, err)

		assert.Equal(t, int64(20000), updated.BonusBalance)
	})

	t.Run("not found", func(t *testing.T) {
		_, err := repo.UpdateBonusBalance(ctx, "00000000-0000-0000-0000-000000000000", 1000)
		require.ErrorIs(t, err, domain.ErrBusinessNotFound)
	})
}
