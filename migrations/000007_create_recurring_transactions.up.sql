-- 000007_create_recurring_transactions.up.sql
--
-- Transações recorrentes: salário mensal, aluguel, assinaturas.
-- O sistema usa esta tabela para gerar automaticamente
-- transações futuras com base na frequência configurada.
--
-- FLUXO:
--   1. Usuário cria uma RecurringTransaction
--   2. Worker job roda diariamente
--   3. Para cada RecurringTransaction ativa, verifica se é hora de gerar a próxima transação
--   4. Cria a Transaction e atualiza next_occurrence
-- ============================================================

CREATE TABLE recurring_transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    account_id UUID NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
    category_id UUID REFERENCES categories(id) ON DELETE SET NULL,

    description VARCHAR(255) NOT NULL,
    cents BIGINT NOT NULL CHECK (cents > 0),
    currency CHAR(3) NOT NULL DEFAULT 'BRL',
    type transaction_type NOT NULL,
    frequency recurrence_frequency NOT NULL,
    day_of_month SMALLINT CHECK (day_of_month BETWEEN 1 AND 28),
    starts_on DATE NOT NULL,
    ends_on DATE NOT NULL,
    next_occurrency DATE NOT NULL,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    -- Validações de consistência
    CONSTRAINT recurring_ends_after_starts CHECK (
        ends_on IS NULL OR ends_on > starts_on
    ),
    CONSTRAINT recurring_monthly_needs_day CHECK (
        frequency != 'monthly' OR day_of_month IS NOT NULL
    )
);

CREATE TRIGGER recurring_transactions_set_updated_at
    BEFORE UPDATE ON recurring_transactions
    FOR EACH ROW
    EXECUTE FUNCTION set_updated_at();