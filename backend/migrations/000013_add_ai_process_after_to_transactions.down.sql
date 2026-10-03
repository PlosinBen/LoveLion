DROP INDEX IF EXISTS idx_transactions_ai_pending_schedule;

ALTER TABLE transactions
    DROP COLUMN IF EXISTS ai_process_after;
