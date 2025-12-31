-- Migration: Add media table
-- Description: Create media table for storing product images and videos

CREATE TABLE IF NOT EXISTS media (
    id SERIAL PRIMARY KEY,
    product_id INTEGER NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    type VARCHAR(20) NOT NULL CHECK (type IN ('image', 'video')),
    file_name VARCHAR(255) NOT NULL,
    file_path VARCHAR(500) NOT NULL,
    file_size BIGINT NOT NULL,
    content_type VARCHAR(100),
    width INTEGER,
    height INTEGER,
    alt_text VARCHAR(255),
    position INTEGER DEFAULT 0,
    is_primary BOOLEAN DEFAULT FALSE,
    status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('active', 'inactive', 'processing', 'deleted')),
    storage_provider VARCHAR(50) DEFAULT 'local',
    storage_path VARCHAR(500),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for performance
CREATE INDEX idx_media_product_id ON media(product_id);
CREATE INDEX idx_media_type ON media(type);
CREATE INDEX idx_media_status ON media(status);
CREATE INDEX idx_media_is_primary ON media(is_primary);

-- Add comment to table
COMMENT ON TABLE media IS 'Stores product media files (images and videos) with metadata and storage information';

-- Add comments to columns
COMMENT ON COLUMN media.product_id IS 'Reference to the product this media belongs to';
COMMENT ON COLUMN media.type IS 'Media type: image or video';
COMMENT ON COLUMN media.file_name IS 'Original filename of the uploaded file';
COMMENT ON COLUMN media.file_path IS 'Path where the file is stored';
COMMENT ON COLUMN media.file_size IS 'File size in bytes';
COMMENT ON COLUMN media.content_type IS 'MIME type of the file (e.g., image/jpeg)';
COMMENT ON COLUMN media.width IS 'Image width in pixels (for images)';
COMMENT ON COLUMN media.height IS 'Image height in pixels (for images)';
COMMENT ON COLUMN media.alt_text IS 'Alternative text for accessibility';
COMMENT ON COLUMN media.position IS 'Display order position';
COMMENT ON COLUMN media.is_primary IS 'Whether this is the primary image for the product';
COMMENT ON COLUMN media.status IS 'Media status: active, inactive, processing, deleted';
COMMENT ON COLUMN media.storage_provider IS 'Storage provider: local, r2, s3';
COMMENT ON COLUMN media.storage_path IS 'Path in the storage provider';
COMMENT ON COLUMN media.created_at IS 'Timestamp when the media was created';
COMMENT ON COLUMN media.updated_at IS 'Timestamp when the media was last updated';
