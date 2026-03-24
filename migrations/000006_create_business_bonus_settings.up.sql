CREATE TABLE
	IF NOT EXISTS business_bonus_settings (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
		business_id UUID NOT NULL UNIQUE REFERENCES businesses (id),
		bonus_percent DECIMAL(5, 2) NOT NULL CHECK (bonus_percent > 0),
		is_active BOOLEAN NOT NULL DEFAULT FALSE,
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
		updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW ()
	);

COMMENT ON TABLE business_bonus_settings IS 'Bonus program settings configured per business';