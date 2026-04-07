-- 000006_create_budgets.down.sql
DROP TRIGGER IF EXISTS budgets_set_updated_at ON budgets;
DROP TABLE IF EXISTS budgets;