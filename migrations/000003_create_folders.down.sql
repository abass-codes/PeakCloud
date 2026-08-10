DROP INDEX IF EXISTS files_owner_folder_idx;
DROP INDEX IF EXISTS files_folder_id_idx;

ALTER TABLE files
    DROP COLUMN IF EXISTS folder_id;

DROP TABLE IF EXISTS folders;
