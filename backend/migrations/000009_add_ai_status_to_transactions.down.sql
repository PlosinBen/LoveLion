DROP INDEX IF EXISTS idx_transactions_ai_status;

ALTER TABLE transactions
    DROP COLUMN IF EXISTS ai_error,
    DROP COLUMN IF EXISTS ai_status;
