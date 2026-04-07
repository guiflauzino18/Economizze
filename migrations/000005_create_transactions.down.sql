-- 000005_create_transactions.down.sql
DROP TRIGGER IF EXISTS transactions_set_updated_at ON transactions;
DROP TABLE IF EXISTS transactions;