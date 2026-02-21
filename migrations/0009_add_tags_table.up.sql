BEGIN;

-- Tags
CREATE TABLE tags (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) UNIQUE NOT NULL,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Composite table menu_item_tags
CREATE TABLE menu_item_tags (
    menu_item_id UUID NOT NULL REFERENCES menu_items(id) ON DELETE CASCADE,
    tag_id UUID NOT NULL REFERENCES tags(id) ON DELETE CASCADE,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (menu_item_id, tag_id)
);

-- Indexes for better performance
CREATE INDEX idx_menu_item_tags_menu_item ON menu_item_tags(menu_item_id);
CREATE INDEX idx_menu_item_tags_tag ON menu_item_tags(tag_id);

-- remove tags field in menu_items
ALTER TABLE menu_items
	DROP COLUMN tags;

COMMIT;
