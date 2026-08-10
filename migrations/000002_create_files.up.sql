CREATE TABLE files (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    owner_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    object_key TEXT NOT NULL UNIQUE,
    original_name TEXT NOT NULL,
    content_type TEXT NOT NULL,
    size_bytes BIGINT NOT NULL,
    etag TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT files_original_name_not_blank
        CHECK (length(trim(original_name)) > 0),

    CONSTRAINT files_object_key_not_blank
        CHECK (length(trim(object_key)) > 0),

    CONSTRAINT files_content_type_not_blank
        CHECK (length(trim(content_type)) > 0),

    CONSTRAINT files_size_non_negative
        CHECK (size_bytes >= 0)
);

CREATE INDEX files_owner_id_idx
    ON files (owner_id);

CREATE INDEX files_owner_created_at_idx
    ON files (owner_id, created_at DESC);
