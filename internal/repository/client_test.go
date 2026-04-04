package repository_test

import (
	"context"
	"testing"

	"github.com/GulzhanKarakul/payment-service/internal/domain"
	"github.com/GulzhanKarakul/payment-service/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClientRepo_Create(t *testing.T) {
	cleanDB(t)

	repo := repository.NewClientRepository(testDB)
	ctx := context.Background()

	tests := []struct {
		name       string
		phone      string
		clientName string
		wantErr    bool
	}{
		{
			name:       "success",
			phone:      "+77771156580",
			clientName: "Gulzhan",
			wantErr:    false,
		},
		{
			name:       "empty phone",
			phone:      "",
			clientName: "Gulzhan",
			wantErr:    true, // NOT NULL constraint
		},
	}

	// act
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// arrange
			client, err := repo.Create(ctx, tt.phone, tt.clientName)

			if tt.wantErr {
				require.Error(t, err) // что делает метод Error
				return                // зачем ретерн если реквайер и так не даст след коду выполниться?
			}

			// assert
			require.NoError(t, err)
			assert.NotEmpty(t, client.ID)
			assert.Equal(t, tt.phone, client.Phone)
			assert.Zero(t, client.BonusBalance)
			assert.NotZero(t, client.CreatedAt)
		})
	}
}

// вопрос такой а че не сделать этот и первый тест вместе???? один же метод проверяем
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
		_, err := repo.GetByID(ctx, "non-existent-id")
		require.ErrorIs(t, err, domain.ErrClientNotFound)
	})
}

func TestClientRepo_GetByPhone(t *testing.T) {
	cleanDB(t)

	repo := repository.NewClientRepository(testDB)
	ctx := context.Background()

	created, err := repo.Create(ctx, "+7771156580", "Gulzhan")
	require.NoError(t, err)

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
		_, err := repo.GetByPhone(ctx, "+7777")
		require.Error(t, err)
		require.ErrorIs(t, err, domain.ErrClientNotFound)
	})
}

func TestClientRepo_UpdateBonusBalance(t *testing.T) {
	cleanDB(t)

	repo := repository.NewClientRepository(testDB)
	ctx := context.Background()

	created, err := repo.Create(ctx, "+77771156580", "Gulzhan")
	require.NoError(t, err)
	assert.Equal(t, int64(0), created.BonusBalance)

	t.Run("add bonus", func(t *testing.T) {
		updated, err := repo.UpdateBonusBalance(ctx, created.ID, 500)
		require.NoError(t, err)
		assert.Equal(t, int64(500), updated.BonusBalance)
	})

	t.Run("subtract bonus", func(t *testing.T) {
		client, err := repo.UpdateBonusBalance(ctx, created.ID, 10000)
		require.NoError(t, err)

		assert.Equal(t, int64(10000), client.BonusBalance)

		updated, err := repo.UpdateBonusBalance(ctx, created.ID, -3000)
		require.NoError(t, err)

		assert.Equal(t, int64(7000), updated.BonusBalance)
	})

	t.Run("not found", func(t *testing.T) {
		_, err := repo.UpdateBonusBalance(ctx, "not-existent-id", 4000)
		require.Error(t, err)
		require.ErrorIs(t, err, domain.ErrClientNotFound)
	})
}
