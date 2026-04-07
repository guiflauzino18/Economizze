-- ============================================================
-- 000006_create_budgets.up.sql
--
-- Orçamentos mensais por categoria.
-- Permite definir quanto o usuário PLANEJA gastar em cada categoria em um período. O sistema alerta quando excede.
--
-- Ex: Orçamento de R$ 500,00 para "Alimentação" em Janeiro/2026
-- ============================================================

CREATE TABLE budgets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    category_id UUID NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
    period_start DATE NOT NULL,
    period_end DATE NOT NULL,
    limit_cents BIGINT NOT NULL CHECK (limit_cents > 0), -- Limite planejado em centavos
    spent_cents BIGINT not null default 0,
    currency CHAR(3) NOT NULL DEFAULT 'BRL',
    notify_when_exceeded BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Um usuário não pode ter dois orçamentos para a mesma categoria no mesmo período (evita confusão)
    CONSTRAINT budgets_unique_period UNIQUE (user_id, category_id, period_start, period_end),

    -- Valida período: fim maior que inicio
    CONSTRAINT budgets_periodo_valid CHECK (period_end > period_start),

    -- Gasto não pode ser negativo
    CONSTRAINT budgets_spent_non_negative CHECK (spent_cents >= 0),

);

CREATE TRIGGER budgets_set_update_at
    BEFORE UPDATE ON budgets
    FOR EACH ROW
    EXECUTE FUNCTION set_updated_at();