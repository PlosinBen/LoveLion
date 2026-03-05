-- Create ledgers table (now representing a "Space")
CREATE TABLE IF NOT EXISTS ledgers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    type VARCHAR(20) NOT NULL DEFAULT 'personal', -- personal, trip, group, etc.
    base_currency VARCHAR(3) NOT NULL DEFAULT 'TWD',
    currencies JSONB DEFAULT '["TWD"]',
    members JSONB DEFAULT '[]', -- legacy field for quick member names, deprecated but kept for compatibility
    categories JSONB DEFAULT '[]',
    payment_methods JSONB DEFAULT '[]',
    start_date DATE,
    end_date DATE,
    cover_image TEXT,
    is_pinned BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW ()
);

-- Create ledger_members table
CREATE TABLE IF NOT EXISTS ledger_members (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    ledger_id UUID NOT NULL REFERENCES ledgers(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role VARCHAR(20) NOT NULL DEFAULT 'member', -- 'owner', 'member'
    alias VARCHAR(50), -- Added by owner to identify this member
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
    UNIQUE(ledger_id, user_id)
);

-- Create ledger_invites table
CREATE TABLE IF NOT EXISTS ledger_invites (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    ledger_id UUID NOT NULL REFERENCES ledgers(id) ON DELETE CASCADE,
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
CREATE INDEX IF NOT EXISTS idx_ledgers_user_id ON ledgers (user_id);
CREATE INDEX IF NOT EXISTS idx_ledgers_is_pinned ON ledgers (is_pinned) WHERE is_pinned = TRUE;
CREATE INDEX IF NOT EXISTS idx_ledger_members_ledger_id ON ledger_members (ledger_id);
CREATE INDEX IF NOT EXISTS idx_ledger_members_user_id ON ledger_members (user_id);
CREATE INDEX IF NOT EXISTS idx_ledger_invites_ledger_id ON ledger_invites (ledger_id);
CREATE INDEX IF NOT EXISTS idx_ledger_invites_token ON ledger_invites (token);