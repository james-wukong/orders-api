BEGIN;

-- Add tags field to menu_items
ALTER TABLE menu_items
	ADD COLUMN allergens TEXT[];

-- Remove indexes from menu_item_tags table
DROP INDEX IF EXISTS idx_menu_item_allergens_menu_item;
DROP INDEX IF EXISTS idx_menu_item_allergens_allergen;

-- Drop tables
DROP TABLE IF EXISTS menu_item_allergens CASCADE;
DROP TABLE IF EXISTS allergens CASCADE;

COMMIT;
