BEGIN;

-- 1. Change the column type to INTERVAL
ALTER TABLE menu_items
	ALTER COLUMN preparation_time TYPE INTERVAL
	-- USING 	(preparation_time * INTERVAL '1 minute');
	USING (preparation_time || ' minutes')::INTERVAL;

-- 1. Set a new interval-compatible default (e.g., '5 minutes')
ALTER TABLE menu_items
ALTER	 COLUMN preparation_time SET DEFAULT INTERVAL '5 minutes';

COMMIT;
