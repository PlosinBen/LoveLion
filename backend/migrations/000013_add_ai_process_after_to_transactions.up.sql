ALTER TABLE transactions
    ADD COLUMN ai_process_after TIMESTAMPTZ;

CREATE INDEX IF NOT EXISTS idx_transactions_ai_pending_schedule
    ON transactions (ai_process_after, created_at)
    WHERE ai_status = 'pending';
