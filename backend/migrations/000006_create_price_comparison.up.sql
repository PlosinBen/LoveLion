-- Create trip_comparison_stores table
CREATE TABLE IF NOT EXISTS trip_comparison_stores (
    id VARCHAR(21) PRIMARY KEY,
    trip_id VARCHAR(21) NOT NULL REFERENCES trips (id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    location TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW ()
);

-- Create trip_comparison_products table
CREATE TABLE IF NOT EXISTS trip_comparison_products (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    store_id VARCHAR(21) NOT NULL REFERENCES trip_comparison_stores (id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    currency VARCHAR(3) DEFAULT 'TWD',
    unit VARCHAR(20),
    note TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW ()
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_comparison_stores_trip_id ON trip_comparison_stores (trip_id);

CREATE INDEX IF NOT EXISTS idx_comparison_products_store_id ON trip_comparison_products (store_id);

CREATE INDEX IF NOT EXISTS idx_comparison_products_name ON trip_comparison_products (name);