CREATE TABLE IF NOT EXISTS inv_capital_oversea_futures_statements (
    year_month VARCHAR(7) PRIMARY KEY REFERENCES inv_settlements(year_month) ON DELETE CASCADE,
    twd_balance INTEGER NOT NULL DEFAULT 0,
    converted_net INTEGER NOT NULL DEFAULT 0,
    profit_loss INTEGER NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS inv_capital_oversea_futures_currencies (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    year_month VARCHAR(7) NOT NULL REFERENCES inv_capital_oversea_futures_statements(year_month) ON DELETE CASCADE,
    currency VARCHAR(10) NOT NULL DEFAULT 'USD',
    balance INTEGER NOT NULL DEFAULT 0,
    unrealized INTEGER NOT NULL DEFAULT 0,
    exchange_rate INTEGER NOT NULL DEFAULT 30
);
