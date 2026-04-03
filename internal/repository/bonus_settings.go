package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/GulzhanKarakul/payment-service/internal/domain"
)

type postgresBonusSettingsRepo struct {
	db *sql.DB
}

func NewBonusSettingsRepository(db *sql.DB) *postgresBonusSettingsRepo {
	return &postgresBonusSettingsRepo{db: db}
}

func (r *postgresBonusSettingsRepo) Upsert(
	ctx context.Context,
	businessId string,
	bonusPercent float64,
	isActive bool,
) (domain.BonusSettings, error) {
	var bs domain.BonusSettings
	err := r.db.QueryRowContext(
		ctx,
		`INSERT INTO business_bonus_settings (business_id, bonus_percent, is_active)
		VALUES ($1, $2, $3)
		ON CONFLICT(business_id)
		DO UPDATE SET bonus_percent = EXCLUDED.bonus_percent,
			is_active = EXCLUDED.is_active,
			updated_at = NOW()
		RETURNING id, business_id, bonus_percent, is_active, created_at, updated_at`,
		businessId, bonusPercent, isActive,
	).Scan(
		&bs.ID, &bs.BusinessID, &bs.BonusPercent, &bs.IsActive, &bs.CreatedAt, &bs.UpdatedAt,
	)
	if err != nil {
		return domain.BonusSettings{}, fmt.Errorf("create and update bonus settings: %w", err)
	}

	return bs, nil
}

func (r *postgresBonusSettingsRepo) GetByBusinessID(
	ctx context.Context,
	businessId string,
) (domain.BonusSettings, error) {
	var bs domain.BonusSettings
	err := r.db.QueryRowContext(
		ctx,
		`SELECT id, business_id, bonus_percent, is_active, created_at, updated_at
		FROM business_bonus_settings
		WHERE business_id=$1`,
		businessId,
	).Scan(
		&bs.ID, &bs.BusinessID, &bs.BonusPercent, &bs.IsActive, &bs.CreatedAt, &bs.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.BonusSettings{}, domain.ErrBonusSettingsNotFound
		}
		return domain.BonusSettings{}, fmt.Errorf("get bonus settings by business id: %w", err)
	}

	return bs, nil
}
