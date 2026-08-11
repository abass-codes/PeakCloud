CREATE TABLE file_versions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    file_id UUID NOT NULL
        REFERENCES files(id)
        ON DELETE CASCADE,

    version_number INTEGER NOT NULL,

    object_key TEXT NOT NULL,

    size_bytes BIGINT NOT NULL,

    content_type TEXT NOT NULL,

    etag TEXT,

    created_by UUID NOT NULL
        REFERENCES users(id)
        ON DELETE RESTRICT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT file_versions_version_positive
        CHECK (version_number > 0),

    CONSTRAINT file_versions_size_non_negative
        CHECK (size_bytes >= 0),

    CONSTRAINT file_versions_object_key_not_blank
        CHECK (length(TRIM(object_key)) > 0),

    CONSTRAINT file_versions_content_type_not_blank
        CHECK (length(TRIM(content_type)) > 0),

    CONSTRAINT file_versions_file_version_unique
        UNIQUE (file_id, version_number),

    CONSTRAINT file_versions_object_key_unique
        UNIQUE (object_key)
);

CREATE INDEX file_versions_file_id_idx
    ON file_versions(file_id);

CREATE INDEX file_versions_file_created_at_idx
    ON file_versions(file_id, created_at DESC);

INSERT INTO file_versions (
    file_id,
    version_number,
    object_key,
    size_bytes,
    content_type,
    etag,
    created_by,
    created_at
)
SELECT
    id,
    1,
    object_key,
    size_bytes,
    content_type,
    etag,
    owner_id,
    created_at
FROM files;
