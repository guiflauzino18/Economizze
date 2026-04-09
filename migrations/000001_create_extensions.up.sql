-- ============================================================
-- 000001_create_extensions.up.sql
--
-- Extensões do PostgreSQL que o sistema vai usar.
-- Devem ser criadas ANTES de qualquer tabela porque alguns tipos de coluna dependem delas (ex: UUID).
-- ============================================================

-- uuid-ossp: funções para gerar UUIDs (uuid_generate_v4())
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- pgcrypto: funções de criptografia e geração segura de bytes
-- gen_random_uuid() é mais moderno e seguro que uuid_generate_v4()
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- citext: tipo de texto case-insensitive (usado por exemplo no campo e-mail para user@x.com seja igual a user@X.com)
CREATE EXTENSION IF NOT EXISTS "citext";

-- pg_trgm: índices trigram para busca parcial eficiente
CREATE EXTENSION IF NOT EXISTS "pg_trgm";


-- ============================================================
-- Tipos ENUM
-- Centralizamos os valores válidos no banco como segunda linha
-- de defesa (a primeira é a validação no domínio Go).
-- ============================================================

-- Tipos de contas bancárias
CREATE TYPE account_type AS ENUM (
    'checking',     -- conta corrente
    'savings',      -- poupança
    'wallet',       -- carteira/dinheiro físico
    'credit_card',  -- cartão de crédito
    'investment'    -- investimentos
);

-- Tipos de transações financeiras
CREATE TYPE transaction_type AS ENUM (
    'income',   -- receita (entrada de dinheiro)
    'expense',  -- despesa (saída de dinheiro)
    'transfer'  -- transferência entre contas próprias
);

-- Frequências para transações recorrentes
CREATE TYPE recurrence_frequency AS ENUM (
    'daily',    -- diárias
    'wekkly',   -- semanal
    'monthly',  -- mensal
    'yearly'    -- anual
);


-- ============================================================
-- Função auxiliar para updated_at automático
--
-- Em vez de repetir a lógica em cada trigger, criamos uma função reutilizável. NEW é a linha que está sendo atualizada — modificamos o campo updated_at antes de salvar.
-- ============================================================
CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN 
    -- Now() retorna timestamp atual com timezone
    NEW.updated_at = Now();
    RETURN NEW;
END;

$$ LANGUAGE plpgsql