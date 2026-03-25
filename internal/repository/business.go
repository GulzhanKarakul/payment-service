package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/GulzhanKarakul/payment-service/internal/domain"
)

type postgresBusinessRepo struct {
	db *sql.DB
}

func NewBusinessRepository(db *sql.DB) *postgresBusinessRepo {
	return &postgresBusinessRepo{ db: db }
}

func (r *postgresBusinessRepo) Create(
	ctx context.Context,
	name, ownerPhone string,
)( domain.Business, error ) {
	var b domain.Business
	err := r.db.QueryRowContext(
		ctx,
		`INSERT INTO businesses (name, owner_phone)
		VALUES ($1, $2)
		RETURNING id, name, owner_phone, bonus_balance, is_active, created_at, updated_at`,
		name, ownerPhone,
	).Scan(
		&b.ID, &b.Name, &b.OwnerPhone, &b.BonusBalance, &b.IsActive, &b.CreatedAt, &b.UpdatedAt,
	)
	if err != nil {
		return domain.Business{}, fmt.Errorf("create business: %w", err)
	}

	return b, nil
}

func (r *postgresBusinessRepo) GetByID(
	ctx context.Context,
	businessId string,
)(domain.Business, error) {
	var b domain.Business
	err := r.db.QueryRowContext(
		ctx,
		`SELECT id, name, owner_phone, bonus_balance, is_active, created_at, updated_at
		FROM businesses
		WHERE id=$1`,
		businessId,
	).Scan(
		&b.ID, &b.Name, &b.OwnerPhone, &b.BonusBalance, &b.IsActive, &b.CreatedAt, &b.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Business{}, domain.ErrBusinessNotFound
		}
		return domain.Business{}, fmt.Errorf("get business by id: %w", err)
	}

	return b, nil
}

func (r *postgresBusinessRepo) UpdateBonusBalance(
	ctx context.Context,
	businessId string,
	bonus int64,
) (domain.Business, error) {
	var b domain.Business
	err := r.db.QueryRowContext(
		ctx,
		`UPDATE businesses
		SET bonus_balance = bonus_balance + $1,
			updated_at = NOW()
		WHERE id=$2
		RETURNING id, name, owner_phone, bonus_balance, is_active, created_at, updated_at`,
		bonus, businessId,
	).Scan(
		&b.ID, &b.Name, &b.OwnerPhone, &b.BonusBalance, &b.IsActive, &b.CreatedAt, &b.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Business{}, domain.ErrBusinessNotFound
		}
		return domain.Business{}, fmt.Errorf("update business bonus balance: %w", err)
	}

	return b, nil
}
