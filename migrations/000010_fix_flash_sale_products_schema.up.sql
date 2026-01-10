-- +migrate Up
-- Fix flash_sale_products table schema to match GORM model

-- Drop existing table and recreate with correct schema
DROP TABLE IF EXISTS flash_sale_products CASCADE;

-- Recreate flash_sale_products table with correct schema
CREATE TABLE flash_sale_products (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    flash_sale_id BIGINT NOT NULL,
    product_id BIGINT NOT NULL,

    -- Flash sale specific settings for this product
    flash_sale_price DECIMAL(10, 2) NOT NULL,
    flash_sale_stock INTEGER NOT NULL,
    sold_count INTEGER NOT NULL DEFAULT 0,

    CONSTRAINT fk_flash_sale_products_flash_sale FOREIGN KEY (flash_sale_id) REFERENCES flash_sales(id) ON DELETE CASCADE,
    CONSTRAINT fk_flash_sale_products_product FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE,
    CONSTRAINT uq_flash_sale_products UNIQUE (flash_sale_id, product_id)
);

-- Create indexes for flash_sale_products
CREATE INDEX IF NOT EXISTS idx_flash_sale_products_flash_sale_id ON flash_sale_products(flash_sale_id);
CREATE INDEX IF NOT EXISTS idx_flash_sale_products_product_id ON flash_sale_products(product_id);
