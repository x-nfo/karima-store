ALTER TABLE users ADD COLUMN IF NOT EXISTS kratos_id VARCHAR(36);
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_kratos_id ON users(kratos_id);
