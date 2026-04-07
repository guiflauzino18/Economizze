-- ============================================================
-- 000003_create_categories.up.sql
--
-- Categorias organizam as transações (ex: Alimentação, Salário).
-- Suportamos duas origens:
--   1. Categorias padrão do sistema (user_id IS NULL)
--   2. Categorias personalizadas por usuário (user_id NOT NULL)
--
-- Isso evita duplicar categorias comuns para cada usuário.
-- ============================================================

CREATE TABLE categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE -- NULL = categoria padrão do sistema (visível para todos) -- NOT NULL = categoria personalizada (visível só para o dono) / Se deletar User deleta categoria criada por ele
    name VARCHAR(100) NOT NULL,
    default_type transaction_type, -- Tipo de transação padrão para esta categoria (Ex.: Salário sempre Income, Aluguel sempre expense)
    active BOOLEAN NOT NULL DEFAULT TRUE, -- Soft Delete
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT categories_name_unique UNIQUE (user_id, name); -- Um usuário não pode ter duas categorias com mesmo nome. NULL != NULL no SQL, então categorias do sistema não conflitam com categorias de usuário de mesmo nome

)

-- Trigger para atualizar updated_at em todo update
CREATE TRIGGER categories_set_updated_at
    BEFORE UPDATE ON categories
    FOR EACH ROW
    EXECUTE FUNCTION set_updated_at();

