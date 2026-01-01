ALTER TABLE users DROP COLUMN identity_id;
-- Note: We generally don't revert 'DROP NOT NULL' as it might break if nulls were introduced.
-- ALTER TABLE users ALTER COLUMN password SET NOT NULL; 
