
-- ============================================================
-- 000004_create_accounts.up.sql
--
-- Contas financeiras do usuário (corrente, poupança, etc).
-- Esta tabela é o coração do sistema — tudo gira em torno do saldo das contas.
--
-- saldo armazenado em CENTAVOS (BIGINT)
-- Nunca usar DECIMAL/FLOAT para dinheiro:
--   - Float tem problemas de arredondamento (0.1 + 0.2 = 0.30000000000001)
--   - DECIMAL é preciso mas mais lento
--   - BIGINT em centavos é rápido, exato e simples
-- ============================================================

CREATE TABLE accounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE, -- Cada conta pertence exatamente a um usuário
    name VARCHAR(100) NOT NULL,
    account_type account_type NOT NULL, -- ENUM criado na migration 000001
    balance_cents BIGINT NOT NULL DEFAULT 0,
    currency CHAR(3) NOT NULL DEFAULT 'BRL', -- ISO 4217: BRL, USD, EUR, etc.
    is_default BOOLEAN NOT NULL DEFAULT FALSE, -- Flag para conta padrão
    active BOOLEAN NOT NULL DEFAULT TRUE -- Soft Delete
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT accounts_currency_format CHECK (currency ~ '^[A-Z]{3}$'), -- Garante formato da moeada com 3 letras maiúscula

);

-- Trigger para atualizar updated_at em cada update
CREATE TRIGGER accounts_set_updated_at
    BEFORE UPDATE ON accounts
    FOR EACH ROW
    EXECUTE FUNCTION set_updated_at();