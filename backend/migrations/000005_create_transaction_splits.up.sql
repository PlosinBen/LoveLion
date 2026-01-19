-- Create transaction_splits table for bill splitting
CREATE TABLE IF NOT EXISTS transaction_splits (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    transaction_id VARCHAR(21) NOT NULL REFERENCES transactions (id) ON DELETE CASCADE,
    member_id UUID NOT NULL REFERENCES trip_members (id) ON DELETE CASCADE,
    amount DECIMAL(10, 2) NOT NULL DEFAULT 0,
    is_payer BOOLEAN DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW ()
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_transaction_splits_transaction_id ON transaction_splits (transaction_id);

CREATE INDEX IF NOT EXISTS idx_transaction_splits_member_id ON transaction_splits (member_id);