ALTER TABLE payments
RENAME TO transactions;

COMMENT ON TABLE transactions IS 'financial transactions between clients and businesses';