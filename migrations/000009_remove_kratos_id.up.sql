-- Drop kratos_id column and related index from users table
ALTER TABLE users DROP COLUMN IF EXISTS kratos_id;
