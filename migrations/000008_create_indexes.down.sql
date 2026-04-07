-- ============================================================
-- 000008_create_indexes.down.sql
--
-- CONCURRENTLY também funciona no DROP — não bloqueia queries
-- ============================================================

DROP INDEX CONCURRENTLY IF EXISTS idx_users_email_active;
DROP INDEX CONCURRENTLY IF EXISTS idx_recurring_next_occurrence;
DROP INDEX CONCURRENTLY IF EXISTS idx_categories_user;
DROP INDEX CONCURRENTLY IF EXISTS idx_budgets_user_period;
DROP INDEX CONCURRENTLY IF EXISTS idx_transactions_description_trgm;
DROP INDEX CONCURRENTLY IF EXISTS idx_transactions_recurring;
DROP INDEX CONCURRENTLY IF EXISTS idx_transactions_type_date;
DROP INDEX CONCURRENTLY IF EXISTS idx_transactions_category;
DROP INDEX CONCURRENTLY IF EXISTS idx_transactions_account_date;
DROP INDEX CONCURRENTLY IF EXISTS idx_accounts_user_active;