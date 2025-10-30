-- Drop indexes first
DROP INDEX IF EXISTS idx_files_mime_type;
DROP INDEX IF EXISTS idx_files_created_at;
DROP INDEX IF EXISTS idx_files_is_public;
DROP INDEX IF EXISTS idx_files_category;
DROP INDEX IF EXISTS idx_files_user_id;
DROP INDEX IF EXISTS idx_files_created_by;
DROP INDEX IF EXISTS idx_files_updated_by;
DROP INDEX IF EXISTS idx_files_deleted_by;

-- Drop the files table
DROP TABLE IF EXISTS files;