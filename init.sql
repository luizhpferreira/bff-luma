-- Script de inicialização do banco de dados PostgreSQL para BFF Luma
-- Este script é executado automaticamente quando o container PostgreSQL é criado pela primeira vez

-- Cria o banco de dados se não existir (já é criado via variável de ambiente)
-- CREATE DATABASE IF NOT EXISTS bff_luma;

-- Conecta ao banco de dados
\c bff_luma;

-- Cria as tabelas (as tabelas serão criadas automaticamente pelo aplicativo Go)
-- Este script serve apenas como backup e documentação

-- Tabela de carteiras
CREATE TABLE IF NOT EXISTS wallets (
    id SERIAL PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    wallet_id TEXT NOT NULL UNIQUE,
    admin_key TEXT NOT NULL,
    invoice_key TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Tabela de tokens de reset
CREATE TABLE IF NOT EXISTS reset_tokens (
    id SERIAL PRIMARY KEY,
    email TEXT NOT NULL,
    token TEXT NOT NULL UNIQUE,
    expires_at TIMESTAMP NOT NULL,
    used BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Índices para melhor performance
CREATE INDEX IF NOT EXISTS idx_wallets_email ON wallets(email);
CREATE INDEX IF NOT EXISTS idx_wallets_wallet_id ON wallets(wallet_id);
CREATE INDEX IF NOT EXISTS idx_reset_tokens_token ON reset_tokens(token);
CREATE INDEX IF NOT EXISTS idx_reset_tokens_email ON reset_tokens(email);

-- Log de inicialização
DO $$
BEGIN
    RAISE NOTICE 'Banco de dados BFF Luma inicializado com sucesso!';
END $$;
