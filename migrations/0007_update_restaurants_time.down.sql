BEGIN;

-- 1. Remove the interval default value
ALTER TABLE restaurants
ALTER COLUMN estimated_delivery_time DROP DEFAULT;

-- 2. Change the type back to INTEGER
-- We extract 'epoch' (total seconds) and divide by 60 to get minutes back
ALTER TABLE restaurants
ALTER COLUMN estimated_delivery_time TYPE INTEGER
USING (EXTRACT(EPOCH FROM estimated_delivery_time) / 60)::INTEGER;

-- 3. Set the original integer default (e.g., 0)
ALTER TABLE restaurants
ALTER COLUMN estimated_delivery_time SET DEFAULT 0;

COMMIT;
