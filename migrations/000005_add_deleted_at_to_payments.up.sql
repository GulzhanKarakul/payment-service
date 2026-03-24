ALTER TABLE payments
ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMPTZ;

COMMENT ON COLUMN payments.deleted_at IS 'timestamp when record was soft-deleted';