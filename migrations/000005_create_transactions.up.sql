-- ============================================================
-- 000005_create_transactions.up.sql
--
-- Transações financeiras — o registro de cada movimentação.
-- É a tabela com mais leituras e escritas do sistema.
--
-- CONVENÇÃO DE SINAL:
--   amount_cents > 0 = receita (dinheiro entrou)
--   amount_cents < 0 = despesa (dinheiro saiu)
--
-- Isso simplifica os cálculos de saldo:
--   SUM(amount_cents) = saldo do período
-- ============================================================

CREATE TABLE transtactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    account_id UUID NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
    category_id UUID REFERENCES categories(id) ON DELETE SET NULL, -- categoria é opcional / ON DELETE: ao deletar categoria seta o campo como null e mantém a transação
    transfer_peer_id UUID REFERENCES transactions(id) ON DELETE SET NULL, -- se for transferencia armaze a conta de destino
    cents BIGINT NOT NULL CHECK (cents != 0),
    currency CHAR(3) NOT NULL,
    type transaction_type NOT NULL, -- Tipo criado no migration 000001
    description VARCHAR(255) NOT NULL, -- Descrição livre
    notes TEXT, -- Notas adicionais opcionais. Texto mais longo,
    occurred_on DATE NOT NULL,
    recurring_id UUID REFERENCES recurring_transactions(id) ON DELETE SET NULL, -- Se for transação recorrente (automática) armazena a origem
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Garante consistência: Tipo e sinal do valor devem bater: income sempre positivo, expanse sempre negativo
    CONSTRAINT transactions_sign_matches_type CHECK(
        (type = 'income' AND cents > 0) OR
        (type = 'expense' AND cents < 0) OR
        (type = 'transfer') -- transferência pode ter qualquer sinal
    )

    CONSTRAINT transactions_currency_format CHECK (currency ~ '^[A-Z]{3}$')

);

-- Trigger para atualizar updated_at em cada update
CREATE TRIGGER transactions_set_updated_at 
    BEFORE UPDATE ON transactions
    FOR EACH ROW
    EXECUTE FUNCTION set_updated_at();