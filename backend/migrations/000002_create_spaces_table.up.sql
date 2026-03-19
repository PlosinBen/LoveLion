-- Create spaces table
CREATE TABLE IF NOT EXISTS spaces (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    type VARCHAR(20) NOT NULL DEFAULT 'personal', -- personal, trip, group, etc.
    base_currency VARCHAR(3) NOT NULL DEFAULT 'TWD',
    currencies JSONB DEFAULT '["TWD"]',
    split_members JSONB DEFAULT '[]',
    categories JSONB DEFAULT '[]',
    payment_methods JSONB DEFAULT '[]',
    start_date DATE,
    end_date DATE,
    is_pinned BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW ()
);

-- Create space_members table
CREATE TABLE IF NOT EXISTS space_members (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    space_id UUID NOT NULL REFERENCES spaces(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role VARCHAR(20) NOT NULL DEFAULT 'member', -- 'owner', 'member'
    alias VARCHAR(50),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
    UNIQUE(space_id, user_id)
);

-- Create space_invites table
CREATE TABLE IF NOT EXISTS space_invites (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    space_id UUID NOT NULL REFERENCES spaces(id) ON DELETE CASCADE,
    token VARCHAR(50) NOT NULL UNIQUE,
    is_one_time BOOLEAN NOT NULL DEFAULT TRUE,
    max_uses INTEGER NOT NULL DEFAULT 1, -- 0 for unlimited
    use_count INTEGER NOT NULL DEFAULT 0,
    expires_at TIMESTAMPTZ,
    created_by UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW ()
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_spaces_user_id ON spaces (user_id);
CREATE INDEX IF NOT EXISTS idx_spaces_is_pinned ON spaces (is_pinned) WHERE is_pinned = TRUE;
CREATE INDEX IF NOT EXISTS idx_space_members_space_id ON space_members (space_id);
CREATE INDEX IF NOT EXISTS idx_space_members_user_id ON space_members (user_id);
CREATE INDEX IF NOT EXISTS idx_space_invites_space_id ON space_invites (space_id);
CREATE INDEX IF NOT EXISTS idx_space_invites_token ON space_invites (token);
