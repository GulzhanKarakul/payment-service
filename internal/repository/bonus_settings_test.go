package repository_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/GulzhanKarakul/payment-service/internal/domain"
	"github.com/GulzhanKarakul/payment-service/internal/repository"
)

func TestBonusSettingsRepo_Upsert_Create(t *testing.T) {
	cleanDB(t)

	ctx := context.Background()
	repo := repository.NewBonusSettingsRepository(testDB)
	business := createTestBusiness(t, "Samsung", "+77771156580")

	t.Run("create new settings", func(t *testing.T) {
		created, err := repo.Upsert(ctx, business.ID, 5.0, true)
		require.NoError(t, err)

		assert.NotEmpty(t, created.ID)
		assert.Equal(t, 5.0, created.BonusPercent)
		assert.True(t, created.IsActive)
		assert.NotZero(t, created.CreatedAt)
	})

	t.Run("update existing", func(t *testing.T) {
		updated, err := repo.Upsert(ctx, business.ID, 7.5, false)
		require.NoError(t, err)

		assert.Equal(t, 7.5, updated.BonusPercent)
		assert.False(t, updated.IsActive)
		assert.NotZero(t, updated.UpdatedAt)

		var count int
		err = testDB.QueryRow(`
		SELECT COUNT(*) FROM business_bonus_settings WHERE business_id = $1
		`, business.ID).Scan(&count)
		require.NoError(t, err)
		assert.Equal(t, 1, count)
	})

	t.Run("business not found", func(t *testing.T) {
		_, err := repo.Upsert(ctx, "00000000-0000-0000-0000-000000000000", 5.0, false)
		require.Error(t, err)
	})
}

func TestBonusSettingsRepo_GetByBusinessID(t *testing.T) {
	cleanDB(t)

	ctx := context.Background()
	bonusRepo := repository.NewBonusSettingsRepository(testDB)
	business := createTestBusiness(t, "samsung", "+77771112233")

	created, err := bonusRepo.Upsert(ctx, business.ID, 6.0, true)
	require.NoError(t, err)

	t.Run("success", func(t *testing.T) {
		found, err := bonusRepo.GetByBusinessID(ctx, business.ID)
		require.NoError(t, err)

		assert.Equal(t, created.ID, found.ID)
		assert.Equal(t, 6.0, found.BonusPercent)
		assert.True(t, found.IsActive)
	})

	t.Run("not found", func(t *testing.T) {
		_, err := bonusRepo.GetByBusinessID(ctx, "00000000-0000-0000-0000-000000000000")
		require.ErrorIs(t, err, domain.ErrBonusSettingsNotFound)
	})
}
