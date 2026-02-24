BEGIN;

-- 1. Remove the interval default value
ALTER TABLE recipes
	ALTER COLUMN preparation_time DROP DEFAULT,
	ALTER COLUMN cooking_time DROP DEFAULT;

-- 2. Change the type back to INTEGER
-- We extract 'epoch' (total seconds) and divide by 60 to get minutes back
ALTER TABLE recipes
	ALTER COLUMN preparation_time TYPE INTEGER
		USING (EXTRACT(EPOCH FROM preparation_time) / 60)::INTEGER,
	ALTER COLUMN cooking_time TYPE INTEGER
		USING (EXTRACT(EPOCH FROM cooking_time) / 60)::INTEGER;


COMMIT;
