-- Create transaction_debts table (replaces transaction_splits)
CREATE TABLE IF NOT EXISTS transaction_debts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    transaction_id VARCHAR(21) NOT NULL REFERENCES transactions (id) ON DELETE CASCADE,
    payer_name VARCHAR(50) NOT NULL,
    payee_name VARCHAR(50) NOT NULL,
    amount DECIMAL(10, 2) NOT NULL DEFAULT 0,
    settled_amount DECIMAL(10, 2) NOT NULL DEFAULT 0,
    is_spot_paid BOOLEAN NOT NULL DEFAULT false
);

CREATE INDEX IF NOT EXISTS idx_transaction_debts_transaction_id ON transaction_debts (transaction_id);
