-- ============================================================
-- 000002_create_users.down.sql
-- ============================================================

-- Triggers são dropados junto com a tabela, mas é boa prática
-- dropar explicitamente para deixar o down completo e claro
DROP TRIGGER IF EXISTS users_set_updated_at ON users;
DROP TABLE IF EXISTS user_social_accounts;
DROP TABLE IF EXISTS users;