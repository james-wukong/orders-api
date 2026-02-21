BEGIN;

-- 1. Remove interval default
ALTER TABLE menu_items
ALTER COLUMN preparation_time DROP DEFAULT;

-- 2. Change the type back to INTEGER
-- We extract 'epoch' (total seconds) and divide by 60 to get minutes back
ALTER TABLE menu_items
ALTER COLUMN preparation_time TYPE INTEGER
USING (EXTRACT(EPOCH FROM preparation_time) / 60)::INTEGER;

COMMIT;
