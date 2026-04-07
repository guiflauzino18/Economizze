-- 000007_create_recurring_transactions.down.sql
DROP TRIGGER IF EXISTS recurring_transactions_set_updated_at ON recurring_transactions;
DROP TABLE IF EXISTS recurring_transactions;