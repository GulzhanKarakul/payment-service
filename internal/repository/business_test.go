package repository_test

import (
	"context"
	"testing"

	"github.com/GulzhanKarakul/payment-service/internal/domain"
	"github.com/GulzhanKarakul/payment-service/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBusinessRepo_Create(t *testing.T) {
	cleanDB(t)

	repo := repository.NewBusinessRepository(testDB)
	ctx := context.Background()

	tests := []struct {
		name         string
		businessName string
		ownerPhone   string
		wantErr      bool
	}{
		{
			name:         "success",
			businessName: "Samsung",
			ownerPhone:   "+77771112233",
			wantErr:      false,
		}, {
			name:         "empty phone",
			businessName: "Mi",
			ownerPhone:   "",
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			business, err := repo.Create(ctx, tt.businessName, tt.ownerPhone)

			if tt.wantErr {
				require.Error(t, err)
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

	created, err := repo.Create(ctx, "Samsung", "+77771112233")
	require.NoError(t, err)

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
		_, err := repo.GetByID(ctx, "no-existent-id")
		require.Error(t, err)
		require.ErrorIs(t, err, domain.ErrBusinessNotFound)
	})
}

func TestBusinessRepo_UpdateBonusBalance(t *testing.T) {
	cleanDB(t)

	repo := repository.NewBusinessRepository(testDB)
	ctx := context.Background()

	created, err := repo.Create(ctx, "samsung", "+77771112233")
	require.NoError(t, err)

	t.Run("success", func(t *testing.T) {
		updated, err := repo.UpdateBonusBalance(ctx, created.ID, 50000)
		require.NoError(t, err)

		assert.Equal(t, created.ID, updated.ID)
		assert.Equal(t, created.OwnerPhone, updated.OwnerPhone)
		assert.Equal(t, int64(50000), updated.BonusBalance)
	})

	t.Run("not found", func(t *testing.T) {
		_, err := repo.UpdateBonusBalance(ctx, "not-existent-id", 1000)
		require.Error(t, err)
		require.ErrorIs(t, err, domain.ErrBusinessNotFound)
	})
}
