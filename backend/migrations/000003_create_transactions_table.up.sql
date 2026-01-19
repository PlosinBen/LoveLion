-- Create transactions table
CREATE TABLE IF NOT EXISTS transactions (
    id VARCHAR(21) PRIMARY KEY,
    ledger_id UUID NOT NULL REFERENCES ledgers (id) ON DELETE CASCADE,
    payer VARCHAR(50),
    date TIMESTAMPTZ NOT NULL DEFAULT NOW (),
    currency VARCHAR(3) NOT NULL DEFAULT 'TWD',
    total_amount DECIMAL(10, 2) NOT NULL DEFAULT 0,
    exchange_rate DECIMAL(12, 6) DEFAULT 1,
    billing_amount DECIMAL(10, 2) DEFAULT 0,
    handling_fee DECIMAL(10, 2) DEFAULT 0,
    category VARCHAR(50),
    payment_method VARCHAR(50),
    note TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW ()
);

-- Create transaction_items table
CREATE TABLE IF NOT EXISTS transaction_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    transaction_id VARCHAR(21) NOT NULL REFERENCES transactions (id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    unit_price DECIMAL(10, 2) NOT NULL DEFAULT 0,
    quantity DECIMAL(8, 2) NOT NULL DEFAULT 1,
    discount DECIMAL(10, 2) DEFAULT 0,
    amount DECIMAL(10, 2) NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW ()
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_transactions_ledger_id ON transactions (ledger_id);

CREATE INDEX IF NOT EXISTS idx_transactions_date ON transactions (date);

CREATE INDEX IF NOT EXISTS idx_transaction_items_transaction_id ON transaction_items (transaction_id);