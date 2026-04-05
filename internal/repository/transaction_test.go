package repository_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/GulzhanKarakul/payment-service/internal/domain"
	"github.com/GulzhanKarakul/payment-service/internal/repository"
)

func TestTransactionRepo_CreateWithBonus_WithBonus(t *testing.T) {
	cleanDB(t)

	ctx := context.Background()
	repo := repository.NewTransactionRepository(testDB)
	clientRepo := repository.NewClientRepository(testDB)
	bizRepo := repository.NewBusinessRepository(testDB)

	client := createTestClient(t, "+77771156580", "Gulzhan")
	business := createTestBusinessWithBonusBalance(t, "samsung", "+77771112233", 100000)

	tx, err := repo.CreateWithBonus(ctx, client.ID, business.ID, 50000, 10.0, nil)
	require.NoError(t, err)

	assert.NotEmpty(t, tx.ID)
	assert.Equal(t, int64(5000), tx.BonusAccrued)
	assert.Equal(t, int64(50000), tx.Amount)
	assert.Equal(t, domain.StatusCompleted, tx.Status)

	updatedClient, err := clientRepo.GetByID(ctx, client.ID)
	require.NoError(t, err)
	assert.Equal(t, int64(5000), updatedClient.BonusBalance)

	updatedBiz, err := bizRepo.GetByID(ctx, business.ID)
	require.NoError(t, err)
	assert.Equal(t, int64(95000), updatedBiz.BonusBalance)

}

func TestTransactionRepo_CreateWithBonus_ZeroBonus(t *testing.T) {
	cleanDB(t)

	ctx := context.Background()
	repo := repository.NewTransactionRepository(testDB)
	clientRepo := repository.NewClientRepository(testDB)
	bizRepo := repository.NewBusinessRepository(testDB)

	client := createTestClient(t, "+77771156580", "Gulzhan")
	business := createTestBusinessWithBonusBalance(t, "samsung", "+77771112233", 100000)

	tx, err := repo.CreateWithBonus(ctx, client.ID, business.ID, 50000, 0.0, nil)
	require.NoError(t, err)

	assert.NotEmpty(t, tx.ID)
	assert.Equal(t, int64(0), tx.BonusAccrued)
	assert.Equal(t, int64(50000), tx.Amount)
	assert.Equal(t, domain.StatusCompleted, tx.Status)

	updatedClient, err := clientRepo.GetByID(ctx, client.ID)
	require.NoError(t, err)
	assert.Equal(t, int64(0), updatedClient.BonusBalance)

	updatedBiz, err := bizRepo.GetByID(ctx, business.ID)
	require.NoError(t, err)
	assert.Equal(t, int64(100000), updatedBiz.BonusBalance)
}

func TestTransactionRepo_CreateWithBonus_InsufficientBalance(t *testing.T) {
	cleanDB(t)

	ctx := context.Background()
	repo := repository.NewTransactionRepository(testDB)
	clientRepo := repository.NewClientRepository(testDB)

	client := createTestClient(t, "+77771112233", "Gulzhan")
	business := createTestBusiness(t, "Samsung", "+77771112234")

	_, err := repo.CreateWithBonus(ctx, client.ID, business.ID, 10000, 5.0, nil) 
	require.ErrorIs(t, err, domain.ErrInsufficientBonusBalance)

	updatedClient, err := clientRepo.GetByID(ctx, client.ID)
	require.NoError(t, err)
	assert.Equal(t, int64(0), updatedClient.BonusBalance)
}

func TestTransactionRepo_GetByID(t *testing.T) {
	cleanDB(t)

	ctx := context.Background()
	repo := repository.NewTransactionRepository(testDB)

	client := createTestClient(t, "+77771156580", "Gulzhan")
	business := createTestBusinessWithBonusBalance(t, "samsung", "+77771112233", 100000)

	t.Run("success", func(t *testing.T) {
		created := createTestTransaction(t, client.ID, business.ID, 50000, 0.0, nil)

		found, err := repo.GetByID(ctx, created.ID)
		require.NoError(t, err)

		assert.Equal(t, created.ID, found.ID)
		assert.Equal(t, created.ClientID, found.ClientID)
		assert.Equal(t, created.BusinessID, found.BusinessID)
		assert.Equal(t, created.BonusAccrued, found.BonusAccrued)
		assert.Equal(t, created.Amount, found.Amount)
		assert.Equal(t, created.CreatedAt, found.CreatedAt)
	})

	t.Run("not found", func(t *testing.T) {
		_, err := repo.GetByID(ctx, "00000000-0000-0000-0000-000000000000")
		require.ErrorIs(t, err, domain.ErrTransactionNotFound)
	})
}

func TestTransactionRepo_GetByClientID_Pagination(t *testing.T){
	cleanDB(t)

	ctx := context.Background()
	repo := repository.NewTransactionRepository(testDB)

	client := createTestClient(t, "+77771112233", "Gulzhan")
	business := createTestBusinessWithBonusBalance(t, "apple", "+77771112234", 1000000)

	for i := 0; i < 5; i++ {
		_ = createTestTransaction(t, client.ID, business.ID, 5000, 10.0, nil)
	}

	t.Run("first page", func(t *testing.T) {
		txList, err := repo.GetByClientID(ctx, client.ID, 2, 0)
		require.NoError(t, err)

		assert.Equal(t, client.ID, txList[0].ClientID)
		assert.Len(t, txList, 2)
		assert.True(t, txList[0].CreatedAt.After(txList[1].CreatedAt) || 
			txList[0].CreatedAt.Equal(txList[1].CreatedAt))
	})

	t.Run("second page", func(t *testing.T) {
		txList, err := repo.GetByClientID(ctx, client.ID, 2, 2)
		require.NoError(t, err)

		assert.Len(t, txList, 2)
	})

	t.Run("last page", func(t *testing.T) {
		txList, err := repo.GetByClientID(ctx, client.ID, 2, 4)
		require.NoError(t, err)

		assert.Len(t, txList, 1)
	})

	t.Run("beyond last page", func(t *testing.T) {
		txList, err := repo.GetByClientID(ctx, client.ID, 2, 10)
		require.NoError(t, err)

		assert.Len(t, txList, 0)
	})
}


func TestTransactionRepo_Cancel(t *testing.T) {
	cleanDB(t)

	ctx := context.Background()
	repo := repository.NewTransactionRepository(testDB)

	client := createTestClient(t, "+77771112233", "Gulzhan")
	business := createTestBusinessWithBonusBalance(t, "apple", "+77771112234", 10000)

	t.Run("success", func(t *testing.T) {
		tx := createTestTransaction(t, client.ID, business.ID, 5000, 10.0, nil)

		err := repo.Cancel(ctx, tx.ID)
		require.NoError(t, err)

		_, err = repo.GetByID(ctx, tx.ID)
		require.ErrorIs(t, err, domain.ErrTransactionNotFound)
	})

	t.Run("already cancelled", func(t *testing.T) {
		tx := createTestTransaction(t, client.ID, business.ID, 5000, 10.0, nil)

		err := repo.Cancel(ctx, tx.ID)
		require.NoError(t, err)

		err = repo.Cancel(ctx, tx.ID)
		require.ErrorIs(t, err, domain.ErrTransactionCancelled)
	})

	t.Run("not found", func(t *testing.T) {
		err := repo.Cancel(ctx, "00000000-0000-0000-0000-000000000000")
		require.ErrorIs(t, err, domain.ErrTransactionNotFound)
	})
}
