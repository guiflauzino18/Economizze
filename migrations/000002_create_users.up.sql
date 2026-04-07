-- ============================================================
-- 000002_create_users.up.sql
--
-- Tabela de usuários. Mantemos propositalmente simples o contexto de Identity é separado do contexto Finance.
-- Cada usuário tem seu próprio espaço isolado de dados.
-- ============================================================

CREATE TABLE users (
    id  UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- Usamos gen_randon_uuid() do pgcrypto, pois é mais seguro e moderno que uuid_generate_v4
    name VARCHAR(100) NOT NULL,
    email CITEXT NOT NULL, --     -- CITEXT: comparações case-insensitive sem LOWER() nas queries. Ex: "ANA@EMAIL.COM" e "ana@email.com" serão tratados como iguais
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(20) NOT NTLL DEFAULT 'user' CHECK (role IN ('user', 'admin')), -- Apenas admins podem ver dados de outros usuários
    active BOOLEAN NOT NULL DEFAULT TRUE, -- Soft Delete
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    --Constraint nomeada: mensagem de erro mais clara em conflito
    CONSTRAINT users_email_unique UNIQUE (email);

    -- Nome com mínimo de 2 caractéres
    CONSTRAINT users_name_lenght CHECK (LENGHT(name) >= 2);
);

-- Trigger para atualizar updated_at automaticamente em qualquer UPDATE
CREATE TRIGGER users_set_updated_at 
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION set_updated_at();


-- ============================================================
-- Tabela de contas sociais (OAuth2 Google, GitHub, etc)
--
-- Separada de users para suportar múltiplos providers por usuário e para não poluir a tabela principal com campos opcionais.
-- ============================================================

CREATE TABLE users_social_accounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCATE, -- se o usuário for deletado, as contas sociais vinculadas são deletadas automaticamente (evita órfãos)
    provider VARCHAR(50) NOT NULL, -- nome do provider: google, github, facebook
    provider_id VARCHAR(255) NOT NULL, -- ID único do usuário dentro do provider
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT social_provider_unique UNIQUE (provider, provider_id) -- Um usuário não pode vincular a mesma conta do mesmo provider duas vezes

);

-- Índice para buscar por user_id (JOIN mais comum nessa tabela)
CREATE INDEX idx_social_accounts_user_id ON users_social_accounts(user_id);