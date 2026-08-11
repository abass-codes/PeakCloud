CREATE TYPE share_resource_type AS ENUM ('file', 'folder');
CREATE TYPE share_permission AS ENUM ('viewer', 'editor');

CREATE TABLE resource_shares (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    owner_id UUID NOT NULL
        REFERENCES users(id)
        ON DELETE CASCADE,

    recipient_id UUID NOT NULL
        REFERENCES users(id)
        ON DELETE CASCADE,

    resource_type share_resource_type NOT NULL,
    resource_id UUID NOT NULL,

    permission share_permission NOT NULL DEFAULT 'viewer',
    allow_download BOOLEAN NOT NULL DEFAULT TRUE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT resource_shares_not_self
        CHECK (owner_id <> recipient_id),

    CONSTRAINT resource_shares_unique_recipient
        UNIQUE (recipient_id, resource_type, resource_id)
);

CREATE INDEX resource_shares_owner_id_idx
    ON resource_shares (owner_id);

CREATE INDEX resource_shares_recipient_id_idx
    ON resource_shares (recipient_id);

CREATE INDEX resource_shares_resource_idx
    ON resource_shares (resource_type, resource_id);

CREATE TABLE public_share_links (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    owner_id UUID NOT NULL
        REFERENCES users(id)
        ON DELETE CASCADE,

    resource_type share_resource_type NOT NULL,
    resource_id UUID NOT NULL,

    token_hash TEXT NOT NULL UNIQUE,
    password_hash TEXT,

    permission share_permission NOT NULL DEFAULT 'viewer',
    allow_download BOOLEAN NOT NULL DEFAULT TRUE,

    expires_at TIMESTAMPTZ,
    revoked_at TIMESTAMPTZ,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT public_share_token_hash_not_blank
        CHECK (length(trim(token_hash)) > 0)
);

CREATE INDEX public_share_links_owner_id_idx
    ON public_share_links (owner_id);

CREATE INDEX public_share_links_resource_idx
    ON public_share_links (resource_type, resource_id);
