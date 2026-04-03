CREATE TABLE
	IF NOT EXISTS clients (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
		phone VARCHAR(20) NOT NULL UNIQUE,
		name VARCHAR(100) NOT NULL,
		bonus_balance BIGINT NOT NULL DEFAULT 0,
		is_active BOOLEAN NOT NULL DEFAULT TRUE,
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
		updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW ()
	);

COMMENT ON TABLE clients IS 'Platform clients who make purchases';

COMMENT ON COLUMN clients.bonus_balance IS 'Accumulated bonus points';

COMMENT ON COLUMN clients.is_active IS 'False if account is blocked';