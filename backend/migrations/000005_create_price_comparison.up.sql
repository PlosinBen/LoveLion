-- Create comparison_stores table
CREATE TABLE IF NOT EXISTS comparison_stores (
    id VARCHAR(21) PRIMARY KEY,
    space_id UUID NOT NULL REFERENCES spaces (id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    google_map_url TEXT,
    location TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW ()
);

-- Create comparison_products table
CREATE TABLE IF NOT EXISTS comparison_products (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    store_id VARCHAR(21) NOT NULL REFERENCES comparison_stores (id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    currency VARCHAR(3) DEFAULT 'TWD',
    unit VARCHAR(20),
    note TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW ()
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_comparison_stores_space_id ON comparison_stores (space_id);
CREATE INDEX IF NOT EXISTS idx_comparison_products_store_id ON comparison_products (store_id);
CREATE INDEX IF NOT EXISTS idx_comparison_products_name ON comparison_products (name);
