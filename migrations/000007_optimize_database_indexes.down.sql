-- +migrate Down
-- Rollback Database Index Optimization Migration

-- Drop Covering Indexes
DROP INDEX CONCURRENTLY IF EXISTS idx_reviews_covering;
DROP INDEX CONCURRENTLY IF EXISTS idx_orders_covering;
DROP INDEX CONCURRENTLY IF EXISTS idx_products_covering;

-- Drop Unique Constraint Optimization Indexes
DROP INDEX CONCURRENTLY IF EXISTS uidx_users_email_active;
DROP INDEX CONCURRENTLY IF EXISTS uidx_products_sku_active;
DROP INDEX CONCURRENTLY IF EXISTS uidx_products_slug_active;

-- Drop Sorting and Pagination Indexes
DROP INDEX CONCURRENTLY IF EXISTS idx_orders_sort_created;
DROP INDEX CONCURRENTLY IF EXISTS idx_reviews_sort_created;
DROP INDEX CONCURRENTLY IF EXISTS idx_products_sort_rating;
DROP INDEX CONCURRENTLY IF EXISTS idx_products_sort_created;
DROP INDEX CONCURRENTLY IF EXISTS idx_products_sort_price;

-- Drop Performance Indexes for Join Operations
DROP INDEX CONCURRENTLY IF EXISTS idx_product_images_product_updated;
DROP INDEX CONCURRENTLY IF EXISTS idx_order_items_order_updated;
DROP INDEX CONCURRENTLY IF EXISTS idx_cart_items_cart_updated;

-- Drop Analytics and Reporting Indexes
DROP INDEX CONCURRENTLY IF EXISTS idx_users_stats;
DROP INDEX CONCURRENTLY IF EXISTS idx_orders_stats;
DROP INDEX CONCURRENTLY IF EXISTS idx_products_stats;

-- Drop Partial Indexes for Common Queries
DROP INDEX CONCURRENTLY IF EXISTS idx_orders_processing;
DROP INDEX CONCURRENTLY IF EXISTS idx_orders_pending;
DROP INDEX CONCURRENTLY IF EXISTS idx_products_low_stock;
DROP INDEX CONCURRENTLY IF EXISTS idx_products_available_stock;

-- Drop Full Text Search Indexes
DROP INDEX CONCURRENTLY IF EXISTS idx_products_description_search;
DROP INDEX CONCURRENTLY IF EXISTS idx_products_name_search;

-- Drop Review Images Table Composite Indexes
DROP INDEX CONCURRENTLY IF EXISTS idx_review_images_review;

-- Drop Product Images Table Composite Indexes
DROP INDEX CONCURRENTLY IF EXISTS idx_product_images_product_position;

-- Drop Product Variants Table Composite Indexes
DROP INDEX CONCURRENTLY IF EXISTS idx_product_variants_product_price;
DROP INDEX CONCURRENTLY IF EXISTS idx_product_variants_product_color;
DROP INDEX CONCURRENTLY IF EXISTS idx_product_variants_product_size;

-- Drop Flash Sale Products Table Composite Indexes
DROP INDEX CONCURRENTLY IF EXISTS idx_flash_sale_products_flash_product;

-- Drop Flash Sales Table Composite Indexes
DROP INDEX CONCURRENTLY IF EXISTS idx_flash_sales_active_time;
DROP INDEX CONCURRENTLY IF EXISTS idx_flash_sales_status_time;

-- Drop Wishlists Table Composite Indexes
DROP INDEX CONCURRENTLY IF EXISTS idx_wishlists_user_created;
DROP INDEX CONCURRENTLY IF EXISTS idx_wishlists_user_product;

-- Drop Reviews Table Composite Indexes
DROP INDEX CONCURRENTLY IF EXISTS idx_reviews_created_at;
DROP INDEX CONCURRENTLY IF EXISTS idx_reviews_product_approved;
DROP INDEX CONCURRENTLY IF EXISTS idx_reviews_product_rating;

-- Drop Order Items Table Composite Indexes
DROP INDEX CONCURRENTLY IF EXISTS idx_order_items_order_product;

-- Drop Cart Items Table Composite Indexes
DROP INDEX CONCURRENTLY IF EXISTS idx_cart_items_cart_variant;
DROP INDEX CONCURRENTLY IF EXISTS idx_cart_items_cart_product;

-- Drop Orders Table Composite Indexes
DROP INDEX CONCURRENTLY IF EXISTS idx_orders_user_created;
DROP INDEX CONCURRENTLY IF EXISTS idx_orders_created_at;
DROP INDEX CONCURRENTLY IF EXISTS idx_orders_payment_status;
DROP INDEX CONCURRENTLY IF EXISTS idx_orders_status_created;
DROP INDEX CONCURRENTLY IF EXISTS idx_orders_user_status;

-- Drop Users Table Composite Indexes
DROP INDEX CONCURRENTLY IF EXISTS idx_users_created_at;
DROP INDEX CONCURRENTLY IF EXISTS idx_users_role_active;
DROP INDEX CONCURRENTLY IF EXISTS idx_users_email_active;

-- Drop Products Table Composite Indexes
DROP INDEX CONCURRENTLY IF EXISTS idx_products_price_range;
DROP INDEX CONCURRENTLY IF EXISTS idx_products_rating;
DROP INDEX CONCURRENTLY IF EXISTS idx_products_created_at;
DROP INDEX CONCURRENTLY IF EXISTS idx_products_status_price;
DROP INDEX CONCURRENTLY IF EXISTS idx_products_category_price;
DROP INDEX CONCURRENTLY IF EXISTS idx_products_category_status;

-- Drop Time-based Query Indexes
DROP INDEX CONCURRENTLY IF EXISTS idx_orders_delivered;
DROP INDEX CONCURRENTLY IF EXISTS idx_orders_shipped;
DROP INDEX CONCURRENTLY IF EXISTS idx_orders_confirmed;
DROP INDEX CONCURRENTLY IF EXISTS idx_products_published;

-- Drop Status-based Query Indexes
DROP INDEX CONCURRENTLY IF EXISTS idx_reviews_by_status;
DROP INDEX CONCURRENTLY IF EXISTS idx_flash_sales_by_status;
DROP INDEX CONCURRENTLY IF EXISTS idx_products_by_status;
