DROP INDEX IF EXISTS idx_media_deleted_at;
ALTER TABLE media DROP COLUMN IF EXISTS deleted_at;
