CREATE TABLE IF NOT EXISTS inv_members (
    id VARCHAR(21) PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    user_id UUID REFERENCES users (id) ON DELETE SET NULL,
    is_owner BOOLEAN NOT NULL DEFAULT FALSE,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    sort_order INTEGER NOT NULL DEFAULT 0,
    net_investment INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_inv_members_user_id ON inv_members (user_id);

CREATE TABLE IF NOT EXISTS inv_settlements (
    year_month VARCHAR(7) PRIMARY KEY,
    status VARCHAR(10) NOT NULL DEFAULT 'draft',
    total_profit_loss INTEGER NOT NULL DEFAULT 0,
    total_weight INTEGER NOT NULL DEFAULT 0,
    profit_loss_per_weight INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS inv_member_transactions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    member_id VARCHAR(21) NOT NULL REFERENCES inv_members (id) ON DELETE CASCADE,
    date DATE NOT NULL,
    type VARCHAR(20) NOT NULL,
    amount INTEGER NOT NULL DEFAULT 0,
    note TEXT
);

CREATE INDEX IF NOT EXISTS idx_inv_member_transactions_member_id ON inv_member_transactions (member_id);
CREATE INDEX IF NOT EXISTS idx_inv_member_transactions_date ON inv_member_transactions (date);

CREATE TABLE IF NOT EXISTS inv_settlement_allocations (
    year_month VARCHAR(7) NOT NULL REFERENCES inv_settlements (year_month) ON DELETE CASCADE,
    member_id VARCHAR(21) NOT NULL REFERENCES inv_members (id) ON DELETE CASCADE,
    weight INTEGER NOT NULL DEFAULT 0,
    amount INTEGER NOT NULL DEFAULT 0,
    deposit INTEGER NOT NULL DEFAULT 0,
    withdrawal INTEGER NOT NULL DEFAULT 0,
    balance INTEGER NOT NULL DEFAULT 0,
    PRIMARY KEY (year_month, member_id)
);

CREATE TABLE IF NOT EXISTS inv_futures_statements (
    year_month VARCHAR(7) PRIMARY KEY REFERENCES inv_settlements (year_month) ON DELETE CASCADE,
    ending_equity INTEGER NOT NULL DEFAULT 0,
    floating_profit_loss INTEGER NOT NULL DEFAULT 0,
    realized_profit_loss INTEGER NOT NULL DEFAULT 0,
    deposit INTEGER NOT NULL DEFAULT 0,
    withdrawal INTEGER NOT NULL DEFAULT 0,
    profit_loss INTEGER NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS inv_stock_statements (
    year_month VARCHAR(7) PRIMARY KEY REFERENCES inv_settlements (year_month) ON DELETE CASCADE,
    account_balance INTEGER NOT NULL DEFAULT 0,
    market_value INTEGER NOT NULL DEFAULT 0,
    deposit INTEGER NOT NULL DEFAULT 0,
    withdrawal INTEGER NOT NULL DEFAULT 0,
    profit_loss INTEGER NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS inv_stock_holdings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    year_month VARCHAR(7) NOT NULL REFERENCES inv_stock_statements (year_month) ON DELETE CASCADE,
    symbol VARCHAR(20) NOT NULL,
    shares INTEGER NOT NULL DEFAULT 0,
    closing_price DECIMAL(10, 2) NOT NULL DEFAULT 0,
    market_value INTEGER NOT NULL DEFAULT 0
);

CREATE INDEX IF NOT EXISTS idx_inv_stock_holdings_year_month ON inv_stock_holdings (year_month);

CREATE TABLE IF NOT EXISTS inv_stock_trades (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    trade_date DATE NOT NULL,
    symbol VARCHAR(20) NOT NULL,
    shares INTEGER NOT NULL DEFAULT 0,
    price DECIMAL(10, 2) NOT NULL DEFAULT 0,
    fee INTEGER NOT NULL DEFAULT 0,
    tax INTEGER NOT NULL DEFAULT 0
);

CREATE INDEX IF NOT EXISTS idx_inv_stock_trades_trade_date ON inv_stock_trades (trade_date);
