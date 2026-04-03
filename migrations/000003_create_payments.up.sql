CREATE TABLE
	IF NOT EXISTS payments (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
		client_id UUID NOT NULL REFERENCES clients (id),
		business_id UUID NOT NULL REFERENCES businesses (id),
		amount BIGINT NOT NULL,
		bonus_accrued BIGINT NOT NULL DEFAULT 0,
		currency VARCHAR(3) NOT NULL DEFAULT 'KZT',
		status VARCHAR(20) NOT NULL DEFAULT 'pending',
		description TEXT,
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
		updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW ()
	);

COMMENT ON TABLE payments IS 'financial transactions between clients and businesses';

COMMENT ON COLUMN payments.status IS 'pending|completed|failed|cancelled';

COMMENT ON COLUMN payments.currency IS 'ISO 4217 currency code (KZT, USD, RUB)';