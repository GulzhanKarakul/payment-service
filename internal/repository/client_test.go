package repository_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/GulzhanKarakul/payment-service/internal/domain"
	"github.com/GulzhanKarakul/payment-service/internal/repository"
)

func TestClientRepo_Create(t *testing.T) {
	repo := repository.NewClientRepository(testDB)
	ctx := context.Background()

	client, err := repo.Create(ctx, "+77771156580", "Gulzhan")
	require.NoError(t, err)

	// assert
	require.NoError(t, err)
	assert.NotEmpty(t, client.ID)
	assert.Equal(t, "+77771156580", client.Phone)
	assert.Equal(t, "Gulzhan", client.Name)
	assert.Equal(t, int64(0), client.BonusBalance)
	assert.True(t, client.IsActive)
	assert.NotZero(t, client.CreatedAt)
	assert.NotZero(t, client.UpdatedAt)
}

func TestClientRepo_Create_DuplicatePhone(t *testing.T) {
	cleanDB(t)

	repo := repository.NewClientRepository(testDB)
	ctx := context.Background()

	_, err := repo.Create(ctx, "+77771156580", "Gulzhan")
	require.NoError(t, err)

	_, err = repo.Create(ctx, "+77771156580", "another")
	require.ErrorIs(t, err, domain.ErrClientAlreadyExist)
}

func TestClientRepo_GetByID(t *testing.T) {
	cleanDB(t)

	repo := repository.NewClientRepository(testDB)
	ctx := context.Background()

	// arrange
	created, err := repo.Create(ctx, "+7771156580", "Gulzhan")
	require.NoError(t, err)

	t.Run("success", func(t *testing.T) {
		client, err := repo.GetByID(ctx, created.ID)
		require.NoError(t, err)

		assert.Equal(t, created.ID, client.ID)
		assert.Equal(t, created.Phone, client.Phone)
		assert.Equal(t, created.Name, client.Name)
		assert.Equal(t, created.BonusBalance, client.BonusBalance)
		assert.Equal(t, created.IsActive, client.IsActive)
		assert.Equal(t, created.CreatedAt, client.CreatedAt)
	})

	t.Run("not_found", func(t *testing.T) {
		_, err := repo.GetByID(ctx, "00000000-0000-0000-0000-000000000000")
		require.ErrorIs(t, err, domain.ErrClientNotFound)
	})
}

func TestClientRepo_GetByPhone(t *testing.T) {
	cleanDB(t)

	repo := repository.NewClientRepository(testDB)
	ctx := context.Background()

	created := createTestClient(t, "+77771156580", "Gulzhan")

	t.Run("success", func(t *testing.T) {
		client, err := repo.GetByPhone(ctx, created.Phone)
		require.NoError(t, err)

		assert.Equal(t, created.ID, client.ID)
		assert.Equal(t, created.Phone, client.Phone)
		assert.Equal(t, created.Name, client.Name)
		assert.Equal(t, created.BonusBalance, client.BonusBalance)
		assert.Equal(t, created.IsActive, client.IsActive)
		assert.Equal(t, created.CreatedAt, client.CreatedAt)
	})

	t.Run("not_found", func(t *testing.T) {
		_, err := repo.GetByPhone(ctx, "+777700000000")
		require.ErrorIs(t, err, domain.ErrClientNotFound)
	})
}

func TestClientRepo_UpdateBonusBalance(t *testing.T) {
	cleanDB(t)

	repo := repository.NewClientRepository(testDB)
	ctx := context.Background()

	t.Run("add bonus", func(t *testing.T) {
		created := createTestClient(t, "+77771156580", "Gulzhan")

		updated, err := repo.UpdateBonusBalance(ctx, created.ID, 500)
		require.NoError(t, err)

		assert.Equal(t, int64(500), updated.BonusBalance)
		assert.True(t, updated.UpdatedAt.After(created.UpdatedAt))
	})

	t.Run("subtract bonus", func(t *testing.T) {
		created := createTestClient(t, "+77771156581", "Gulzhan")

		_, err := repo.UpdateBonusBalance(ctx, created.ID, 10000)
		require.NoError(t, err)

		updated, err := repo.UpdateBonusBalance(ctx, created.ID, -3000)
		require.NoError(t, err)

		assert.Equal(t, int64(7000), updated.BonusBalance)
	})

	t.Run("not found", func(t *testing.T) {
		_, err := repo.UpdateBonusBalance(ctx, "00000000-0000-0000-0000-000000000000", 4000)
		require.ErrorIs(t, err, domain.ErrClientNotFound)
	})
}
