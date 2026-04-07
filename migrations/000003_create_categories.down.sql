-- 000003_create_categories.down.sql
DROP TRIGGER IF EXISTS categories_set_updated_at ON categories;
DROP TABLE IF EXISTS categories;