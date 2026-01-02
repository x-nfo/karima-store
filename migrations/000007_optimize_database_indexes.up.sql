-- +migrate Up
-- Database Index Optimization Migration
-- This migration adds composite indexes and optimizes existing indexes for better query performance

-- Products Table - Composite Indexes
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_products_category_status ON products(category, status) WHERE deleted_at IS NULL;
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_products_category_price ON products(category, price) WHERE deleted_at IS NULL;
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_products_status_price ON products(status, price) WHERE deleted_at IS NULL;
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_products_created_at ON products(created_at DESC) WHERE deleted_at IS NULL;
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_products_rating ON products(rating DESC, review_count DESC) WHERE deleted_at IS NULL;
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_products_price_range ON products(price) WHERE deleted_at IS NULL AND status = 'available';

-- Users Table - Composite Indexes
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_users_email_active ON users(email) WHERE is_active = TRUE AND deleted_at IS NULL;
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_users_role_active ON users(role, is_active) WHERE deleted_at IS NULL;
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_users_created_at ON users(created_at DESC) WHERE deleted_at IS NULL;

-- Orders Table - Composite Indexes
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_orders_user_status ON orders(user_id, status) WHERE deleted_at IS NULL;
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_orders_status_created ON orders(status, created_at DESC) WHERE deleted_at IS NULL;
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_orders_payment_status ON orders(payment_status, created_at DESC) WHERE deleted_at IS NULL;
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_orders_created_at ON orders(created_at DESC) WHERE deleted_at IS NULL;
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_orders_user_created ON orders(user_id, created_at DESC) WHERE deleted_at IS NULL;

-- Cart Items Table - Composite Indexes
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_cart_items_cart_product ON cart_items(cart_id, product_id);
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_cart_items_cart_variant ON cart_items(cart_id, product_variant_id);

-- Order Items Table - Composite Indexes
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_order_items_order_product ON order_items(order_id, product_id);

-- Reviews Table - Composite Indexes
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_reviews_product_rating ON reviews(product_id, rating DESC) WHERE deleted_at IS NULL AND is_approved = TRUE;
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_reviews_product_approved ON reviews(product_id, is_approved, created_at DESC) WHERE deleted_at IS NULL;
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_reviews_created_at ON reviews(created_at DESC) WHERE deleted_at IS NULL AND is_approved = TRUE;

-- Wishlists Table - Composite Indexes
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_wishlists_user_product ON wishlists(user_id, product_id) WHERE deleted_at IS NULL;
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_wishlists_user_created ON wishlists(user_id, created_at DESC) WHERE deleted_at IS NULL;

-- Flash Sales Table - Composite Indexes
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_flash_sales_status_time ON flash_sales(status, start_time, end_time) WHERE deleted_at IS NULL;
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_flash_sales_active_time ON flash_sales(start_time, end_time) WHERE status = 'active' AND deleted_at IS NULL;

-- Flash Sale Products Table - Composite Indexes
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_flash_sale_products_flash_product ON flash_sale_products(flash_sale_id, product_id);

-- Product Variants Table - Composite Indexes
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_product_variants_product_size ON product_variants(product_id, size);
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_product_variants_product_color ON product_variants(product_id, color);
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_product_variants_product_price ON product_variants(product_id, price);

-- Product Images Table - Composite Indexes
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_product_images_product_position ON product_images(product_id, position);

-- Review Images Table - Composite Indexes
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_review_images_review ON review_images(review_id);

-- Full Text Search Indexes (PostgreSQL specific)
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_products_name_search ON products USING gin(to_tsvector('english', name)) WHERE deleted_at IS NULL;
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_products_description_search ON products USING gin(to_tsvector('english', description)) WHERE deleted_at IS NULL;

-- Partial Indexes for Common Queries
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_products_available_stock ON products(id, stock) WHERE status = 'available' AND stock > 0 AND deleted_at IS NULL;
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_products_low_stock ON products(id, stock) WHERE stock < 10 AND deleted_at IS NULL;
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_orders_pending ON orders(id, user_id) WHERE status = 'pending' AND deleted_at IS NULL;
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_orders_processing ON orders(id, user_id) WHERE status = 'processing' AND deleted_at IS NULL;

-- Indexes for Analytics and Reporting
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_products_stats ON products(view_count, sold_count, rating) WHERE deleted_at IS NULL;
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_orders_stats ON orders(created_at, total_amount, status) WHERE deleted_at IS NULL;
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_users_stats ON users(created_at, is_active) WHERE deleted_at IS NULL;

-- Performance Indexes for Join Operations
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_cart_items_cart_updated ON cart_items(cart_id, updated_at DESC);
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_order_items_order_updated ON order_items(order_id, updated_at DESC);
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_product_images_product_updated ON product_images(product_id, updated_at DESC);

-- Indexes for Sorting and Pagination
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_products_sort_price ON products(price DESC) WHERE deleted_at IS NULL;
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_products_sort_created ON products(created_at DESC) WHERE deleted_at IS NULL;
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_products_sort_rating ON products(rating DESC, review_count DESC) WHERE deleted_at IS NULL;
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_reviews_sort_created ON reviews(created_at DESC) WHERE deleted_at IS NULL AND is_approved = TRUE;
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_orders_sort_created ON orders(created_at DESC) WHERE deleted_at IS NULL;

-- Indexes for Unique Constraints Optimization
CREATE UNIQUE INDEX CONCURRENTLY IF NOT EXISTS uidx_products_slug_active ON products(slug) WHERE deleted_at IS NULL;
CREATE UNIQUE INDEX CONCURRENTLY IF NOT EXISTS uidx_products_sku_active ON products(sku) WHERE deleted_at IS NULL;
CREATE UNIQUE INDEX CONCURRENTLY IF NOT EXISTS uidx_users_email_active ON users(email) WHERE deleted_at IS NULL;

-- Covering Indexes for Frequently Queried Columns
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_products_covering ON products(category, status, price, name, slug, thumbnail) WHERE deleted_at IS NULL;
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_orders_covering ON orders(user_id, status, created_at, total_amount, order_number) WHERE deleted_at IS NULL;
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_reviews_covering ON reviews(product_id, rating, is_approved, created_at, user_id) WHERE deleted_at IS NULL;

-- Indexes for Time-based Queries
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_products_published ON products(published_at DESC) WHERE published_at IS NOT NULL AND deleted_at IS NULL;
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_orders_confirmed ON orders(confirmed_at DESC) WHERE confirmed_at IS NOT NULL AND deleted_at IS NULL;
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_orders_shipped ON orders(shipped_at DESC) WHERE shipped_at IS NOT NULL AND deleted_at IS NULL;
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_orders_delivered ON orders(delivered_at DESC) WHERE delivered_at IS NOT NULL AND deleted_at IS NULL;

-- Indexes for Status-based Queries
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_products_by_status ON products(status, created_at DESC) WHERE deleted_at IS NULL;
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_flash_sales_by_status ON flash_sales(status, start_time) WHERE deleted_at IS NULL;
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_reviews_by_status ON reviews(is_approved, created_at DESC) WHERE deleted_at IS NULL;
