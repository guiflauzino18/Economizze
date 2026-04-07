-- ============================================================
-- 000001_create_extensions.down.sql
--
-- Reverte EXATAMENTE o que o .up.sql criou, na ordem inversa.
-- Tipos ENUM precisam ser dropados antes das extensões.
-- ============================================================

DROP FUNCTION IF EXISTS set_updated_at();

-- Ordem importa: tipos derivados de outros devem vir primeiro
DROP TYPE IF EXISTS recurrence_frequency;
DROP TYPE IF EXISTS transaction_type;
DROP TYPE IF EXISTS account_type;

-- Não removemos as extensões no down porque:
--   1. Podem estar sendo usadas por outros schemas
--   2. Extensões são difíceis de reinstalar com dados existentes
--   3. O risco de remover é maior que o de manter