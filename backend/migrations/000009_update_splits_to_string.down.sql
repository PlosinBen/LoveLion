-- Restore FK (Warning: This might fail if we have splits with NULL member_id)
-- We attempt to clean up or this down migration assumes no data loss.
-- Re-add FK
ALTER TABLE transaction_splits ADD CONSTRAINT transaction_splits_member_id_fkey FOREIGN KEY (member_id) REFERENCES trip_members (id) ON DELETE CASCADE;

-- Make member_id required again
DELETE FROM transaction_splits
WHERE
    member_id IS NULL;

-- Destructive but necessary for NOT NULL
ALTER TABLE transaction_splits
ALTER COLUMN member_id
SET
    NOT NULL;

-- Drop name column
ALTER TABLE transaction_splits
DROP COLUMN name;