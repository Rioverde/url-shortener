-- Reverse of 0001_create_urls_table.up.sql.
-- Drop in reverse order of creation: index first, then table.
DROP INDEX IF EXISTS idx_urls_short_code;
DROP TABLE IF EXISTS urls;
