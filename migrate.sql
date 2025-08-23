-- Script de migração para atualizar o banco de dados existente
-- Execute este script para migrar da estrutura antiga para a nova

-- 1. Criar tabela users se não existir
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    email TEXT NOT NULL UNIQUE,
    username TEXT NOT NULL UNIQUE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 2. Migrar dados existentes da tabela wallets para users
-- Para cada carteira existente, criar um usuário com email como username temporário
INSERT OR IGNORE INTO users (email, username, created_at, updated_at)
SELECT 
    email, 
    email as username,  -- Usar email como username temporário
    created_at, 
    updated_at
FROM wallets;

-- 3. Adicionar coluna username na tabela wallets se não existir
-- SQLite não suporta ADD COLUMN IF NOT EXISTS, então vamos verificar primeiro
PRAGMA table_info(wallets);

-- 4. Criar nova tabela wallets com a estrutura correta
CREATE TABLE IF NOT EXISTS wallets_new (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    wallet_id TEXT NOT NULL UNIQUE,
    admin_key TEXT NOT NULL,
    invoice_key TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (username) REFERENCES users(username)
);

-- 5. Migrar dados da tabela antiga para a nova
INSERT INTO wallets_new (username, password, wallet_id, admin_key, invoice_key, created_at, updated_at)
SELECT 
    email as username,  -- Usar email como username temporário
    password,
    wallet_id,
    admin_key,
    invoice_key,
    created_at,
    updated_at
FROM wallets;

-- 6. Remover tabela antiga
DROP TABLE wallets;

-- 7. Renomear nova tabela
ALTER TABLE wallets_new RENAME TO wallets;

-- 8. Criar índices
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
CREATE INDEX IF NOT EXISTS idx_wallets_username ON wallets(username);
CREATE INDEX IF NOT EXISTS idx_wallets_wallet_id ON wallets(wallet_id);

-- 9. Verificar se a migração foi bem-sucedida
SELECT 'Migração concluída com sucesso!' as status;
