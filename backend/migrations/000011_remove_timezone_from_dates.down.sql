ALTER TABLE transactions
    ALTER COLUMN date TYPE timestamp WITH TIME ZONE,
    ALTER COLUMN created_at TYPE timestamp WITH TIME ZONE,
    ALTER COLUMN updated_at TYPE timestamp WITH TIME ZONE;
