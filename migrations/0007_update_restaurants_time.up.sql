BEGIN;

-- 1. Remove the old integer default value
ALTER TABLE restaurants
	ALTER COLUMN estimated_delivery_time DROP DEFAULT;

-- 2. Change the column type to INTERVAL
ALTER TABLE restaurants
	ALTER COLUMN estimated_delivery_time TYPE INTERVAL
		-- USING (estimated_delivery_time * INTERVAL '1 minute');
		USING (estimated_delivery_time || ' minutes')::INTERVAL;

-- 3. Set a new interval-compatible default (e.g., '0 minutes')
ALTER TABLE restaurants
	ALTER COLUMN estimated_delivery_time SET DEFAULT INTERVAL '10 minutes';

COMMIT;
