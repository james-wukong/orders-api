BEGIN;

-- Add tags field to menu_items
ALTER TABLE menu_items
	ADD COLUMN tags TEXT[];

-- Remove indexes from menu_item_tags table
DROP INDEX IF EXISTS idx_menu_item_tags_menu_item;
DROP INDEX IF EXISTS idx_menu_item_tags_tag;

-- Drop tables
DROP TABLE IF EXISTS menu_item_tags CASCADE;
DROP TABLE IF EXISTS tags CASCADE;

COMMIT;
