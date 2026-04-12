-- Add role column to users
ALTER TABLE users ADD COLUMN role VARCHAR(20) NOT NULL DEFAULT 'user';

-- Create announcements table
CREATE TABLE announcements (
    id VARCHAR(20) PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL DEFAULT '',
    status VARCHAR(20) NOT NULL DEFAULT 'draft',
    broadcast_start TIMESTAMPTZ,
    broadcast_end TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_announcements_status ON announcements (status);
CREATE INDEX idx_announcements_broadcast ON announcements (status, broadcast_start, broadcast_end);
