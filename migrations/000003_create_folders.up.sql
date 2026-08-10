CREATE TABLE folders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    owner_id UUID NOT NULL
        REFERENCES users(id)
        ON DELETE CASCADE,

    parent_id UUID
        REFERENCES folders(id)
        ON DELETE CASCADE,

    name TEXT NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT folders_name_not_blank
        CHECK (length(trim(name)) > 0),

    CONSTRAINT folders_not_self_parent
        CHECK (parent_id IS NULL OR parent_id <> id)
);

CREATE INDEX folders_owner_id_idx
    ON folders (owner_id);

CREATE INDEX folders_parent_id_idx
    ON folders (parent_id);

CREATE INDEX folders_owner_parent_idx
    ON folders (owner_id, parent_id);

CREATE UNIQUE INDEX folders_unique_root_name_idx
    ON folders (owner_id, lower(name))
    WHERE parent_id IS NULL;

CREATE UNIQUE INDEX folders_unique_child_name_idx
    ON folders (owner_id, parent_id, lower(name))
    WHERE parent_id IS NOT NULL;

ALTER TABLE files
    ADD COLUMN folder_id UUID
        REFERENCES folders(id)
        ON DELETE CASCADE;

CREATE INDEX files_folder_id_idx
    ON files (folder_id);

CREATE INDEX files_owner_folder_idx
    ON files (owner_id, folder_id);
