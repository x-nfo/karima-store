-- +migrate Down
-- Rollback flash_sale_products schema fix

-- Drop indexes
DROP INDEX IF EXISTS idx_flash_sale_products_flash_sale_id;
DROP INDEX IF EXISTS idx_flash_sale_products_product_id;

-- Drop table
DROP TABLE IF EXISTS flash_sale_products CASCADE;
