CREATE TABLE
	IF NOT EXISTS businesses (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
		name VARCHAR(100) NOT NULL,
		owner_phone VARCHAR(20) NOT NULL UNIQUE,
		bonus_balance DECIMAL(10, 2) NOT NULL DEFAULT 0,
		is_active BOOLEAN NOT NULL DEFAULT TRUE,
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
		updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW ()
	);

COMMENT ON TABLE businesses IS 'Platform businesses who take purchases';

COMMENT ON COLUMN businesses.bonus_balance IS 'Balance to give to clients bonuses';

COMMENT ON COLUMN businesses.is_active IS 'False if account is blocked';