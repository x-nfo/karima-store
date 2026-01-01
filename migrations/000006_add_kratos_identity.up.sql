ALTER TABLE users ADD COLUMN identity_id VARCHAR(40);
CREATE INDEX idx_users_identity_id ON users(identity_id);
ALTER TABLE users ALTER COLUMN password DROP NOT NULL;
