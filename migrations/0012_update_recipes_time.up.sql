BEGIN;


-- 1. Change the column type to INTERVAL
ALTER TABLE recipes
	ALTER COLUMN preparation_time TYPE INTERVAL
		USING (preparation_time * INTERVAL '1 minute'),
	ALTER COLUMN cooking_time TYPE INTERVAL
		USING (cooking_time * INTERVAL '1 minute');

-- 2. Set a new interval-compatible default (e.g., '0 minutes')
ALTER TABLE recipes
	ALTER COLUMN preparation_time SET DEFAULT INTERVAL '10 minutes',
	ALTER COLUMN cooking_time SET DEFAULT INTERVAL '10 minutes';

COMMIT;
