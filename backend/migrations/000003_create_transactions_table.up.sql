-- Create transactions table (shared by all types)
CREATE TABLE IF NOT EXISTS transactions (
    id VARCHAR(21) PRIMARY KEY,
    space_id UUID NOT NULL REFERENCES spaces (id) ON DELETE CASCADE,
    type VARCHAR(20) NOT NULL DEFAULT 'expense',
    title VARCHAR(100),
    date TIMESTAMPTZ NOT NULL DEFAULT NOW (),
    currency VARCHAR(3) NOT NULL DEFAULT 'TWD',
    total_amount DECIMAL(10, 2) NOT NULL DEFAULT 0,
    note TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW ()
);

CREATE INDEX IF NOT EXISTS idx_transactions_space_id ON transactions (space_id);
CREATE INDEX IF NOT EXISTS idx_transactions_date ON transactions (date);
CREATE INDEX IF NOT EXISTS idx_transactions_type ON transactions (type);

-- Create transaction_expenses table (1:1 with transactions where type=expense)
CREATE TABLE IF NOT EXISTS transaction_expenses (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    transaction_id VARCHAR(21) NOT NULL UNIQUE REFERENCES transactions (id) ON DELETE CASCADE,
    category VARCHAR(50),
    exchange_rate DECIMAL(12, 6) DEFAULT 1,
    billing_amount DECIMAL(10, 2) DEFAULT 0,
    handling_fee DECIMAL(10, 2) DEFAULT 0,
    payment_method VARCHAR(50)
);

-- Create transaction_expense_items table (1:N with transaction_expenses)
CREATE TABLE IF NOT EXISTS transaction_expense_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    expense_id UUID NOT NULL REFERENCES transaction_expenses (id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    unit_price DECIMAL(10, 2) NOT NULL DEFAULT 0,
    quantity DECIMAL(8, 2) NOT NULL DEFAULT 1,
    discount DECIMAL(10, 2) DEFAULT 0,
    amount DECIMAL(10, 2) NOT NULL DEFAULT 0
);

CREATE INDEX IF NOT EXISTS idx_transaction_expense_items_expense_id ON transaction_expense_items (expense_id);
