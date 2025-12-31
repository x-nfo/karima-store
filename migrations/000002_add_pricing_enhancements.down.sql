-- Drop pricing enhancements tables
DROP INDEX IF EXISTS idx_shipping_zones_deleted_at;
DROP INDEX IF EXISTS idx_shipping_zones_valid_until;
DROP INDEX IF EXISTS idx_shipping_zones_valid_from;
DROP INDEX IF EXISTS idx_shipping_zones_status;

DROP INDEX IF EXISTS idx_taxes_deleted_at;
DROP INDEX IF EXISTS idx_taxes_valid_until;
DROP INDEX IF EXISTS idx_taxes_valid_from;
DROP INDEX IF EXISTS idx_taxes_is_default;
DROP INDEX IF EXISTS idx_taxes_status;

DROP INDEX IF EXISTS idx_coupons_deleted_at;
DROP INDEX IF EXISTS idx_coupons_valid_until;
DROP INDEX IF EXISTS idx_coupons_valid_from;
DROP INDEX IF EXISTS idx_coupons_status;
DROP INDEX IF EXISTS idx_coupons_code;

DROP INDEX IF EXISTS idx_coupon_usages_order_id;
DROP INDEX IF EXISTS idx_coupon_usages_user_id;
DROP INDEX IF EXISTS idx_coupon_usages_coupon_id;

DROP TABLE IF EXISTS shipping_zones;
DROP TABLE IF EXISTS taxes;
DROP TABLE IF EXISTS coupon_usages;
DROP TABLE IF EXISTS coupons;
