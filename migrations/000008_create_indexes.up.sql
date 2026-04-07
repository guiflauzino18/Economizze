-- 000008_create_indexes.up.sql
--
-- Índices separados das tabelas por dois motivos:
--   1. Permite criar com CONCURRENTLY (não bloqueia writes)
--   2. Facilita rollback isolado sem afetar a estrutura
--
-- ESTRATÉGIA DE INDEXAÇÃO:
--   - Índices parciais (WHERE) são menores e mais rápidos
--   - Índices compostos seguem a ordem das queries mais comuns
--   - Analisamos as queries antes de criar índices
-- ============================================================


-- ── ACCOUNTS ─────────────────────────────────────────────────
-- Query mais comum: "me dê todas as contas ativas deste usuário"
-- ORDER BY é incluído no índice para evitar sort adicional
CREATE INDEX CONCURRENTLY idx_accounts_user_active
    ON accounts(user_id, created_at DESC)
    WHERE active = TRUE;
-- Índice parcial (WHERE active = TRUE) é menor e mais eficiente que um índice em toda a tabela



-- ── TRANSACTIONS ─────────────────────────────────────────────

-- Query mais comum: "me dê as transações desta conta no período X"
-- Composto: (account_id, occurred_on) porque filtramos por conta e ordenamos por data na maioria das queries

CREATE INDEX CONCURRENTLY idx_transactions_account_date
    ON transactions(account_id, occurred_on DESC)


-- Filtro por categoria: relatório "quanto gastei em Alimentação"
CREATE INDEX CONCURRENTLY idx_transactions_category
    ON transactions(category_id, occurred_on DESC)
    WHERE category_id IS NOT NULL;
--Índice parcial: ignora transações sem categoria (NULL)

-- Filtro por tipo: " me dê só as receitas do mês"
CREATE INDEX CONCURRENTLY idx_transactions_type_date
    ON transactions(account_id, type, occurred_on DESC)

-- Busca por recorrência: worker precisa encontrar transações geradas por uma recorrência específica
CREATE INDEX CONCURRENTLY idx_transactions_recurring
    ON transactions(recurring_id)
    WHERE recurring_id IS NOT NULL;

-- ── BUSCA FULL-TEXT ───────────────────────────────────────────

-- Índice GIN trigram para busca por descrição
-- Permite: WHERE description ILIKE '%almoço%' com boa performance
-- Sem este índice, ILIKE faz sequential scan (lento para tabelas grandes)

CREATE INDEX CONCURRENTLY idx_transactions_description_trgm
    ON transactions USING GIN (description gin_trgm_ops);




-- ── BUDGETS ───────────────────────────────────────────────────
-- -- Query do worker e do dashboard: "orçamentos ativos deste usuário neste período"
CREATE INDEX CONCURRENTLY idx_budgets_user_period
    ON budgets(user_id, period_start, period_end);




-- ── CATEGORIES ────────────────────────────────────────────────

-- Busca: "categorias disponíveis para este usuário"
-- inclui categorias do sistema (user_id IS NULL) e as do usuário
CREATE INDEX CONCURRENTLY idx_categories_user
    ON categories(user_id)
    WHERE active = TRUE;




-- ── RECURRING TRANSACTIONS ────────────────────────────────────

-- O worker roda diariamente buscando: "recorrências ativas com
-- next_occurrence <= hoje". Este índice é crítico para performance.
CREATE INDEX CONCURRENTLY idx_recurring_next_occurrence
    ON recurring_transactions(next_occurrence)
    WHERE active = TRUE;




-- ── USERS ─────────────────────────────────────────────────────

-- Login: busca por email. CITEXT já cria índice mas
-- índice parcial em active=TRUE é mais eficiente
CREATE INDEX CONCURRENTLY idx_users_email_active
    ON users(email)
    WHERE active = TRUE;