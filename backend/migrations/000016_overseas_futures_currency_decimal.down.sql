ALTER TABLE inv_capital_oversea_futures_currencies
    ALTER COLUMN balance TYPE INTEGER USING balance::INTEGER,
    ALTER COLUMN unrealized TYPE INTEGER USING unrealized::INTEGER,
    ALTER COLUMN exchange_rate TYPE INTEGER USING exchange_rate::INTEGER;
