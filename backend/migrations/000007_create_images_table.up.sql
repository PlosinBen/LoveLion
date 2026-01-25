CREATE TABLE IF NOT EXISTS images (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    entity_id VARCHAR(50) NOT NULL,
    entity_type VARCHAR(50) NOT NULL,
    file_path TEXT NOT NULL,
    sort_order INT DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Index for fast lookup by entity
CREATE INDEX IF NOT EXISTS idx_images_entity ON images(entity_id, entity_type);
