ALTER TABLE files
ADD COLUMN deleted_at TIMESTAMPTZ;

ALTER TABLE folders
ADD COLUMN deleted_at TIMESTAMPTZ;

CREATE INDEX files_owner_deleted_at_idx
ON files (owner_id, deleted_at)
WHERE deleted_at IS NOT NULL;

CREATE INDEX folders_owner_deleted_at_idx
ON folders (owner_id, deleted_at)
WHERE deleted_at IS NOT NULL;

CREATE INDEX files_active_folder_idx
ON files (owner_id, folder_id)
WHERE deleted_at IS NULL;

CREATE INDEX folders_active_parent_idx
ON folders (owner_id, parent_id)
WHERE deleted_at IS NULL;
