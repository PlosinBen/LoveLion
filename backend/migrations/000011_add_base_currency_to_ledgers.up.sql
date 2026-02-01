ALTER TABLE ledgers ADD COLUMN base_currency VARCHAR(3) NOT NULL DEFAULT 'TWD';

-- Data Migration: Update ledgers.base_currency from linked trips.base_currency
UPDATE ledgers
SET base_currency = trips.base_currency
FROM trips
WHERE trips.ledger_id = ledgers.id;
