DROP INDEX IF EXISTS folders_active_parent_idx;
DROP INDEX IF EXISTS files_active_folder_idx;
DROP INDEX IF EXISTS folders_owner_deleted_at_idx;
DROP INDEX IF EXISTS files_owner_deleted_at_idx;

ALTER TABLE folders
DROP COLUMN IF EXISTS deleted_at;

ALTER TABLE files
DROP COLUMN IF EXISTS deleted_at;
