-- Script de migração para adicionar colunas de confirmação de email
-- Execute este script para atualizar a tabela wallets existente

-- Conecta ao banco de dados
\c bff_luma;

-- Adiciona as novas colunas para confirmação de email
ALTER TABLE wallets 
ADD COLUMN IF NOT EXISTS email_confirmed BOOLEAN DEFAULT FALSE,
ADD COLUMN IF NOT EXISTS email_confirmation_token TEXT,
ADD COLUMN IF NOT EXISTS email_confirmation_expires_at TIMESTAMP;

-- Cria índice para o token de confirmação
CREATE INDEX IF NOT EXISTS idx_wallets_email_confirmation_token ON wallets(email_confirmation_token);

-- Log de migração
DO $$
BEGIN
    RAISE NOTICE 'Migração de confirmação de email concluída com sucesso!';
    RAISE NOTICE 'Colunas adicionadas: email_confirmed, email_confirmation_token, email_confirmation_expires_at';
END $$;
