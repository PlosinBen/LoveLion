-- Add name column
ALTER TABLE transaction_splits
ADD COLUMN name VARCHAR(50);

-- Migrate existing data
UPDATE transaction_splits
SET
    name = (
        SELECT
            name
        FROM
            trip_members
        WHERE
            id = transaction_splits.member_id
    );

-- Fallback for any orphans (shouldn't exist due to previous FK, but safe)
UPDATE transaction_splits
SET
    name = 'Unknown'
WHERE
    name IS NULL;

-- Make name required
ALTER TABLE transaction_splits
ALTER COLUMN name
SET
    NOT NULL;

-- Make member_id optional
ALTER TABLE transaction_splits
ALTER COLUMN member_id
DROP NOT NULL;

-- Drop Foreign Key
ALTER TABLE transaction_splits
DROP CONSTRAINT IF EXISTS transaction_splits_member_id_fkey;