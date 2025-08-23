-- Script para forçar a criação das tabelas
-- Execute este script diretamente no SQLite

-- Criar tabela wallets
CREATE TABLE IF NOT EXISTS wallets (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    cpf TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    wallet_id TEXT NOT NULL UNIQUE,
    admin_key TEXT NOT NULL,
    invoice_key TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Criar tabela reset_tokens
CREATE TABLE IF NOT EXISTS reset_tokens (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    email TEXT NOT NULL,
    token TEXT NOT NULL UNIQUE,
    expires_at DATETIME NOT NULL,
    used BOOLEAN DEFAULT FALSE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Criar índices
CREATE INDEX IF NOT EXISTS idx_wallets_cpf ON wallets(cpf);
CREATE INDEX IF NOT EXISTS idx_wallets_email ON wallets(email);
CREATE INDEX IF NOT EXISTS idx_wallets_wallet_id ON wallets(wallet_id);
CREATE INDEX IF NOT EXISTS idx_reset_tokens_token ON reset_tokens(token);
CREATE INDEX IF NOT EXISTS idx_reset_tokens_email ON reset_tokens(email);

-- Verificar se as tabelas foram criadas
SELECT 'Tabelas criadas com sucesso!' as status;
SELECT name FROM sqlite_master WHERE type='table';
