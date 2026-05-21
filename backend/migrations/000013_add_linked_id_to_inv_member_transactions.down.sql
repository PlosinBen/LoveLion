DROP INDEX IF EXISTS idx_inv_member_transactions_linked_id;
ALTER TABLE inv_member_transactions DROP COLUMN IF EXISTS linked_id;
ALTER TABLE inv_settlement_allocations DROP COLUMN IF EXISTS fee;
