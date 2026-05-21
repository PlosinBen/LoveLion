ALTER TABLE inv_member_transactions ADD COLUMN linked_id UUID;
CREATE INDEX idx_inv_member_transactions_linked_id ON inv_member_transactions (linked_id) WHERE linked_id IS NOT NULL;
ALTER TABLE inv_settlement_allocations ADD COLUMN fee INTEGER NOT NULL DEFAULT 0;
