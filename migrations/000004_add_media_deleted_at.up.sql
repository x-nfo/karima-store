ALTER TABLE media ADD COLUMN deleted_at TIMESTAMP;
CREATE INDEX idx_media_deleted_at ON media(deleted_at);
