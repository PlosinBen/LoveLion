ALTER TABLE transactions
    ADD COLUMN ai_status VARCHAR(20),
    ADD COLUMN ai_error  TEXT;

CREATE INDEX IF NOT EXISTS idx_transactions_ai_status
    ON transactions (ai_status)
    WHERE ai_status IS NOT NULL;
