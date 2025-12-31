-- Migration: Drop media table
-- Description: Remove media table and its indexes

DROP INDEX IF EXISTS idx_media_product_id;
DROP INDEX IF EXISTS idx_media_type;
DROP INDEX IF EXISTS idx_media_status;
DROP INDEX IF EXISTS idx_media_is_primary;

DROP TABLE IF EXISTS media;
