ALTER TABLE transactions
    ALTER COLUMN date TYPE timestamp WITHOUT TIME ZONE,
    ALTER COLUMN created_at TYPE timestamp WITHOUT TIME ZONE,
    ALTER COLUMN updated_at TYPE timestamp WITHOUT TIME ZONE;
