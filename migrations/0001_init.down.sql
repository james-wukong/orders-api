-- ============================================
-- DOWN MIGRATION: RESTAURANT & INVENTORY
-- ============================================
BEGIN;

-- 1. Drop Views first (they depend on tables)
DROP VIEW IF EXISTS order_summary CASCADE;
DROP VIEW IF EXISTS daily_inventory_usage CASCADE;
DROP VIEW IF EXISTS low_stock_items CASCADE;
DROP VIEW IF EXISTS inventory_value CASCADE;

-- 2. Drop Triggers
-- Note: Triggers are automatically dropped when the table is dropped,
-- but we drop them explicitly here for a "clean" script.
-- DROP TRIGGER IF EXISTS check_low_stock_trigger ON inventory_items CASCADE;
-- DROP TRIGGER IF EXISTS deduct_inventory_trigger ON orders CASCADE;
-- DROP TRIGGER IF EXISTS set_po_number ON purchase_orders CASCADE;
-- DROP TRIGGER IF EXISTS set_order_number ON orders CASCADE;
-- DROP TRIGGER IF EXISTS track_order_status_changes ON orders CASCADE;

-- 3. Drop Functions
DROP FUNCTION IF EXISTS check_low_stock() CASCADE;
DROP FUNCTION IF EXISTS deduct_inventory_on_order() CASCADE;
DROP FUNCTION IF EXISTS generate_po_number() CASCADE;
DROP FUNCTION IF EXISTS generate_order_number() CASCADE;
DROP FUNCTION IF EXISTS log_order_status_change() CASCADE;
DROP FUNCTION IF EXISTS update_updated_at_column() CASCADE;

-- 4. Drop Tables (Ordered by Dependency: Child Tables first)
-- Inventory Management
DROP TABLE IF EXISTS stock_alerts CASCADE;
DROP TABLE IF EXISTS inventory_adjustments CASCADE;
DROP TABLE IF EXISTS inventory_transactions CASCADE;
DROP TABLE IF EXISTS purchase_order_items CASCADE;
DROP TABLE IF EXISTS purchase_orders CASCADE;
DROP TABLE IF EXISTS recipe_ingredients CASCADE;
DROP TABLE IF EXISTS recipes CASCADE;
DROP TABLE IF EXISTS inventory_items CASCADE;
DROP TABLE IF EXISTS suppliers CASCADE;
DROP TABLE IF EXISTS ingredient_categories CASCADE;
DROP TABLE IF EXISTS units_of_measure CASCADE;

-- Orders & Menu
DROP TABLE IF EXISTS order_status_history CASCADE;
DROP TABLE IF EXISTS order_items CASCADE;
DROP TABLE IF EXISTS orders CASCADE;
DROP TABLE IF EXISTS cart_items CASCADE;
DROP TABLE IF EXISTS carts CASCADE;
DROP TABLE IF EXISTS menu_item_images CASCADE;
DROP TABLE IF EXISTS menu_items CASCADE;
DROP TABLE IF EXISTS categories CASCADE;

-- Core Identity & Restaurants
DROP TABLE IF EXISTS delivery_addresses CASCADE;
DROP TABLE IF EXISTS user_sessions CASCADE;
DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS restaurants CASCADE;

-- 5. Drop Sequences (The culprit for your previous error)
DROP SEQUENCE IF EXISTS po_number_seq CASCADE;
DROP SEQUENCE IF EXISTS order_number_seq CASCADE;

-- 6. Drop Custom Types (Enums)
DROP TYPE IF EXISTS alert_type_enum CASCADE;
DROP TYPE IF EXISTS po_status_enum CASCADE;
DROP TYPE IF EXISTS adjustment_reason_enum CASCADE;
DROP TYPE IF EXISTS transaction_type_enum CASCADE;
DROP TYPE IF EXISTS order_type_enum CASCADE;
DROP TYPE IF EXISTS payment_method_enum CASCADE;
DROP TYPE IF EXISTS payment_status_enum CASCADE;
DROP TYPE IF EXISTS order_status_enum CASCADE;
DROP TYPE IF EXISTS user_role_enum CASCADE;

COMMIT;

-- Success message
DO $$
BEGIN
    RAISE NOTICE '✅ Database schema reverted successfully!';
END $$;
