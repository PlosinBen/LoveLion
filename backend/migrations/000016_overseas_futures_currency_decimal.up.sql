ALTER TABLE inv_capital_oversea_futures_currencies
    ALTER COLUMN balance TYPE DECIMAL(12,2) USING balance::DECIMAL(12,2),
    ALTER COLUMN unrealized TYPE DECIMAL(12,2) USING unrealized::DECIMAL(12,2),
    ALTER COLUMN exchange_rate TYPE DECIMAL(12,2) USING exchange_rate::DECIMAL(12,2);
