package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/GulzhanKarakul/payment-service/internal/domain"
)

type postgresTransactionRepo struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *postgresTransactionRepo {
	return &postgresTransactionRepo{ db: db }
}

func (r *postgresTransactionRepo) CreateWithBonus(
	ctx context.Context,
	clientId, businessId string,
	amount int64,
	bonusPercent float64,
	description *string,
) (domain.Transaction, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return domain.Transaction{}, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	var t domain.Transaction

	err = tx.QueryRowContext(
		ctx,
		`INSERT INTO transactions (client_id, business_id, amount, description)
		VALUES ($1, $2, $3, $4)
		RETURNING 
			id, client_id, business_id, amount, bonus_accrued,
			currency, status, description, created_at, updated_at`,
			clientId, businessId, amount, description,
	).Scan(
		&t.ID, &t.ClientID, &t.BusinessID, &t.Amount, &t.BonusAccrued,
		&t.Currency, &t.Status, &t.Description, &t.CreatedAt, &t.UpdatedAt,
	)
	if err != nil {
		return domain.Transaction{}, fmt.Errorf("create transaction: %w", err)
	}

	bonus := int64(float64(t.Amount) * bonusPercent / 100)
	result, err := tx.ExecContext(
		ctx,
		`UPDATE clients
		SET bonus_balance = bonus_balance + $1,
			updated_at = NOW()
		WHERE id = $2`,
		bonus, clientId,
	)
	if err != nil {
		return domain.Transaction{}, fmt.Errorf("credit client bonus balance: %w", err)
	}
	if n, _ := result.RowsAffected(); n == 0 {
		return domain.Transaction{}, domain.ErrClientNotFound
	}

	result, err = tx.ExecContext(
		ctx,
		`UPDATE businesses
		SET bonus_balance = bonus_balance - $1,
			updated_at = NOW()
		WHERE id = $2
			AND bonus_balance >= $1`,
		bonus, businessId,
	)
	if err != nil {
		return domain.Transaction{}, fmt.Errorf("debit business bonus balance: %w", err)
	}
	if n, _ := result.RowsAffected(); n == 0 {
		return domain.Transaction{}, domain.ErrInsufficientBonusBalance
	}

	_, err = tx.ExecContext(ctx,
    `UPDATE transactions 
		SET status = $1, bonus_accrued = $2, updated_at = NOW() 
		WHERE id = $3`,
    domain.StatusCompleted, bonus, t.ID,
	)
	if err != nil {
		return domain.Transaction{}, fmt.Errorf("complete transaction: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return domain.Transaction{}, fmt.Errorf("commit tx: %w", err)
	}
	t.Status = domain.StatusCompleted
	t.BonusAccrued = bonus

	return t, nil
}


func (r *postgresTransactionRepo) GetByID(
	ctx context.Context,
	txId string,
) (domain.Transaction, error) {
	var t domain.Transaction
	err := r.db.QueryRowContext(
		ctx,
		`SELECT 
			id, client_id, business_id, amount, bonus_accrued,
			currency, status, description, created_at, updated_at
		FROM transactions
		WHERE id = $1
			AND deleted_at IS NULL`,
			txId,
	).Scan(
		&t.ID, &t.ClientID, &t.BusinessID, &t.Amount, &t.BonusAccrued,
		&t.Currency, &t.Status, &t.Description, &t.CreatedAt, &t.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Transaction{}, domain.ErrTransactionNotFound
		}
		return domain.Transaction{}, fmt.Errorf("get transaction by id: %w", err)
	}

	return t, nil
}

func (r *postgresTransactionRepo) GetByClientId(
	ctx context.Context,
	clientId string,
	limit, offset int,
) (
	[]domain.Transaction,
	error,
) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT 
			id, client_id, business_id, amount, bonus_accrued,
			currency, status, description, created_at, updated_at
		FROM transactions
		WHERE client_id = $1 
			AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`,
		clientId, limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("query transactions: %w", err)
	}
	defer rows.Close()

	var tl []domain.Transaction
	for rows.Next() {
		var t domain.Transaction
		if err = rows.Scan(
			&t.ID,
			&t.ClientID,
			&t.BusinessID,
			&t.Amount,
			&t.BonusAccrued,
			&t.Currency,
			&t.Status,
			&t.Description,
			&t.CreatedAt,
			&t.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("get transaction list by user id: %w", err)
		}
		tl = append(tl, t)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("scan transactions: %w", err)
	}

	return tl, nil
}

func (r *postgresTransactionRepo) Cancel(
	ctx context.Context,
	txId string,
) error {
	var status domain.TransactionStatus
	err := r.db.QueryRowContext(ctx,
			`SELECT status FROM transactions 
			 WHERE id = $1 AND deleted_at IS NULL`,
			txId,
	).Scan(&status)
	if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
					return domain.ErrTransactionNotFound
			}
			return fmt.Errorf("get transaction status: %w", err)
	}

	if status == domain.StatusCancelled {
			return domain.ErrTransactionCancelled
	}

	_, err = r.db.ExecContext(ctx,
			`UPDATE transactions
			 SET status = $1, deleted_at = NOW(), updated_at = NOW()
			 WHERE id = $2`,
			domain.StatusCancelled, txId,
	)
	if err != nil {
			return fmt.Errorf("cancel transaction: %w", err)
	}

	return nil
}
