package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/GulzhanKarakul/payment-service/internal/domain"
)

type postgresClientRepo struct {
	db *sql.DB
}

func NewClientRepository(db *sql.DB) *postgresClientRepo {
	return &postgresClientRepo{ db: db }
}

func (r *postgresClientRepo) Create(
	ctx context.Context, 
	phone, name string,
	)(domain.Client, error) {
	var c domain.Client
	err := r.db.QueryRowContext(
		ctx,
		`INSERT INTO clients (phone, name)
		VALUES ($1, $2)
		RETURNING id, phone, name, bonus_balance, is_active, created_at, updated_at`,
		phone, name,
	).Scan(
		&c.ID, &c.Phone, &c.Name, &c.BonusBalance, &c.IsActive, &c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		return domain.Client{}, fmt.Errorf("create client: %w", err)
	}

	return c, nil
}

func (r *postgresClientRepo) GetByID(
	ctx context.Context, 
	clientId string,
	)(domain.Client, error) {
	var c domain.Client
	err := r.db.QueryRowContext(
		ctx,
		`SELECT id, phone, name, bonus_balance, is_active, created_at, updated_at
		FROM clients
		WHERE id=$1`,
		clientId,
	).Scan(
		&c.ID, &c.Phone, &c.Name, &c.BonusBalance, &c.IsActive, &c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Client{}, domain.ErrClientNotFound
		}
		return domain.Client{}, fmt.Errorf("get client by id: %w", err)
	}

	return c, nil
}

func (r *postgresClientRepo) GetByPhone(
	ctx context.Context,
	phone string,
)(domain.Client, error) {
	var c domain.Client
	err := r.db.QueryRowContext(
		ctx,
		`SELECT id, phone, name, bonus_balance, is_active, created_at, updated_at
		FROM clients
		WHERE phone=$1`,
		phone,
	).Scan(
		&c.ID, &c.Phone, &c.Name, &c.BonusBalance, &c.IsActive, &c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Client{}, domain.ErrClientNotFound 
		}
		return domain.Client{}, fmt.Errorf("get client by phone: %w", err)
	}

	return c, nil
}

func (r *postgresClientRepo) UpdateBonusBalance(
	ctx context.Context,
	clientId string,
	bonus int64,
)(domain.Client, error) {
	var c domain.Client
	err := r.db.QueryRowContext(
		ctx,
		`UPDATE clients
		SET bonus_balance = bonus_balance + $1,
			updated_at = NOW()
		WHERE id=$2
		RETURNING id, phone, name, bonus_balance, is_active, created_at, updated_at`,
		bonus, clientId,
	).Scan(
		&c.ID, &c.Phone, &c.Name, &c.BonusBalance, &c.IsActive, &c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Client{}, domain.ErrClientNotFound 
		}
		return domain.Client{}, fmt.Errorf("update client bonus: %w", err)
	}

	return c, nil
}
