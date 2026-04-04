package repository_test

import (
	"context"
	"testing"

	"github.com/GulzhanKarakul/payment-service/internal/domain"
	"github.com/GulzhanKarakul/payment-service/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBonusSettingsRepo_Upsert_Create(t *testing.T) {
	cleanDB(t)

	bonusRepo := repository.NewBonusSettingsRepository(testDB)
	bizRepo := repository.NewBusinessRepository(testDB)
	ctx := context.Background()

	business, err := bizRepo.Create(ctx, "Samsung", "+77771156580")
	require.NoError(t, err)
	t.Run("success", func(t *testing.T) {
		created, err := bonusRepo.Upsert(ctx, business.ID, float64(5), true)
		require.NoError(t, err)

		assert.NotEmpty(t, created.ID)
		assert.Equal(t, business.ID, created.BusinessID)
		assert.Equal(t, float64(5), created.BonusPercent)
		assert.Equal(t, true, created.IsActive)
		assert.NotZero(t, created.CreatedAt)
	})

	t.Run("business not found", func(t *testing.T) {
		_, err := bonusRepo.Upsert(ctx, "not-existent-id", float64(5), false)
		require.Error(t, err)
		require.ErrorIs(t, err, domain.ErrBusinessNotFound)
	})
}

func TestBonusSettingsRepo_Upsert_Update(t *testing.T) {
	cleanDB(t)

	bonusRepo := repository.NewBonusSettingsRepository(testDB)
	bizRepo := repository.NewBusinessRepository(testDB)
	ctx := context.Background()

	business, err := bizRepo.Create(ctx, "Samsung", "+77771112233")
	require.NoError(t, err)

	t.Run("update bonus", func(t *testing.T) {
		created, err := bonusRepo.Upsert(ctx, business.ID, float64(5), true)
		require.NoError(t, err)

		assert.Equal(t, float64(5), created.BonusPercent)
		assert.Equal(t, true, created.IsActive)

		updated, err := bonusRepo.Upsert(ctx, business.ID, float64(7.5), false)
		require.NoError(t, err)

		assert.Equal(t, float64(7.5), updated.BonusPercent)
		assert.Equal(t, false, updated.IsActive)
		assert.NotEqual(t, updated.CreatedAt, updated.UpdatedAt)
	})
}

func TestBonusSettingsRepo_GetByBusinessID(t *testing.T) {
	cleanDB(t)

	bonusRepo := repository.NewBonusSettingsRepository(testDB)
	bizRepo := repository.NewBusinessRepository(testDB)
	ctx := context.Background()

	business, err := bizRepo.Create(ctx, "samsung", "+77771112233")
	require.NoError(t, err)
	created, err := bonusRepo.Upsert(ctx, business.ID, float64(6), true)
	require.NoError(t, err)

	t.Run("success", func(t *testing.T) {
		found, err := bonusRepo.GetByBusinessID(ctx, business.ID)
		require.NoError(t, err)

		assert.Equal(t, created.ID, found.ID)
		assert.Equal(t, float64(6), found.BonusPercent)
		assert.Equal(t, true, found.IsActive)
		assert.NotZero(t, found.CreatedAt)
		assert.NotZero(t, found.UpdatedAt)
	})

	t.Run("not found", func(t *testing.T) {
		_, err := bonusRepo.GetByBusinessID(ctx, "not-existent-id")
		require.Error(t, err)
		require.ErrorIs(t, err, domain.ErrBonusSettingsNotFound)
	})
}
