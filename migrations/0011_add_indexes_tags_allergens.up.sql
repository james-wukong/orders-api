BEGIN;

CREATE INDEX idx_allergens_name ON allergens(name);
CREATE INDEX idx_tags_name ON tags(name);

COMMIT;
