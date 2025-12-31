-- Add pricing enhancements tables
-- Coupons table
CREATE TABLE IF NOT EXISTS coupons (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,

    -- Basic Information
    code VARCHAR(50) NOT NULL UNIQUE,
    name VARCHAR(200) NOT NULL,
    description TEXT,
    type VARCHAR(20) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'active',

    -- Discount
    discount_value DECIMAL(10, 2) NOT NULL,
    max_discount DECIMAL(10, 2),

    -- Usage Limits
    min_purchase_amount DECIMAL(10, 2) DEFAULT 0,
    max_usage_count INTEGER DEFAULT 0,
    usage_count INTEGER DEFAULT 0,
    max_usage_per_user INTEGER DEFAULT 0,

    -- Validity
    valid_from TIMESTAMP WITH TIME ZONE,
    valid_until TIMESTAMP WITH TIME ZONE,

    -- Customer Type Restrictions
    for_retail BOOLEAN DEFAULT TRUE,
    for_reseller BOOLEAN DEFAULT TRUE,

    -- Statistics
    total_discount_used DECIMAL(15, 2) DEFAULT 0,
    order_count INTEGER DEFAULT 0
);

-- Coupon usages table
CREATE TABLE IF NOT EXISTS coupon_usages (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),

    coupon_id INTEGER NOT NULL REFERENCES coupons(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL,
    order_id INTEGER NOT NULL,
    discount_amount DECIMAL(10, 2) NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_coupon_usages_coupon_id ON coupon_usages(coupon_id);
CREATE INDEX IF NOT EXISTS idx_coupon_usages_user_id ON coupon_usages(user_id);
CREATE INDEX IF NOT EXISTS idx_coupon_usages_order_id ON coupon_usages(order_id);

-- Taxes table
CREATE TABLE IF NOT EXISTS taxes (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,

    -- Basic Information
    name VARCHAR(200) NOT NULL,
    description TEXT,
    type VARCHAR(20) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'active',

    -- Tax Rate
    rate DECIMAL(5, 4) NOT NULL,

    -- Applicability
    is_default BOOLEAN DEFAULT FALSE,
    apply_to_shipping BOOLEAN DEFAULT FALSE,

    -- Validity
    valid_from TIMESTAMP WITH TIME ZONE,
    valid_until TIMESTAMP WITH TIME ZONE
);

-- Shipping zones table
CREATE TABLE IF NOT EXISTS shipping_zones (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,

    -- Basic Information
    name VARCHAR(200) NOT NULL,
    description TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'active',

    -- Free Shipping
    free_shipping_enabled BOOLEAN DEFAULT FALSE,
    free_shipping_threshold DECIMAL(10, 2) DEFAULT 0,

    -- Shipping Rates (base rates per provider)
    jne_base_rate DECIMAL(10, 2) DEFAULT 15000,
    tiki_base_rate DECIMAL(10, 2) DEFAULT 16000,
    pos_base_rate DECIMAL(10, 2) DEFAULT 14000,
    sicepat_base_rate DECIMAL(10, 2) DEFAULT 13000,

    -- Additional Costs
    handling_fee DECIMAL(10, 2) DEFAULT 0,
    minimum_cost DECIMAL(10, 2) DEFAULT 9000,

    -- Validity
    valid_from TIMESTAMP WITH TIME ZONE,
    valid_until TIMESTAMP WITH TIME ZONE
);

-- Create indexes for coupons
CREATE INDEX IF NOT EXISTS idx_coupons_code ON coupons(code);
CREATE INDEX IF NOT EXISTS idx_coupons_status ON coupons(status);
CREATE INDEX IF NOT EXISTS idx_coupons_valid_from ON coupons(valid_from);
CREATE INDEX IF NOT EXISTS idx_coupons_valid_until ON coupons(valid_until);
CREATE INDEX IF NOT EXISTS idx_coupons_deleted_at ON coupons(deleted_at);

-- Create indexes for taxes
CREATE INDEX IF NOT EXISTS idx_taxes_status ON taxes(status);
CREATE INDEX IF NOT EXISTS idx_taxes_is_default ON taxes(is_default);
CREATE INDEX IF NOT EXISTS idx_taxes_valid_from ON taxes(valid_from);
CREATE INDEX IF NOT EXISTS idx_taxes_valid_until ON taxes(valid_until);
CREATE INDEX IF NOT EXISTS idx_taxes_deleted_at ON taxes(deleted_at);

-- Create indexes for shipping_zones
CREATE INDEX IF NOT EXISTS idx_shipping_zones_status ON shipping_zones(status);
CREATE INDEX IF NOT EXISTS idx_shipping_zones_valid_from ON shipping_zones(valid_from);
CREATE INDEX IF NOT EXISTS idx_shipping_zones_valid_until ON shipping_zones(valid_until);
CREATE INDEX IF NOT EXISTS idx_shipping_zones_deleted_at ON shipping_zones(deleted_at);

-- Insert default tax (11% VAT for Indonesia)
INSERT INTO taxes (name, description, type, status, rate, is_default)
VALUES (
    'PPN (Pajak Pertambahan Nilai)',
    'Indonesian Value Added Tax (VAT)',
    'percentage',
    'active',
    0.11,
    TRUE
) ON CONFLICT DO NOTHING;

-- Insert default shipping zone (Indonesia)
INSERT INTO shipping_zones (name, description, status, free_shipping_enabled, free_shipping_threshold, handling_fee, minimum_cost)
VALUES (
    'Indonesia Nationwide',
    'Default shipping zone for all regions in Indonesia',
    'active',
    FALSE,
    0,
    0,
    9000
) ON CONFLICT DO NOTHING;
