-- Inicialização do banco de dados BFF Luma
-- Este arquivo é executado automaticamente quando o container PostgreSQL é criado

-- Criação da tabela temp_passwords se não existir
CREATE TABLE IF NOT EXISTS temp_passwords (
    id SERIAL PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Criação do índice se não existir
CREATE INDEX IF NOT EXISTS idx_temp_passwords_email ON temp_passwords(email);

-- Comentário sobre a funcionalidade
COMMENT ON TABLE temp_passwords IS 'Tabela para armazenar senhas originais temporariamente durante o processo de confirmação de email';
