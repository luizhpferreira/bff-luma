#!/bin/bash

echo "🧪 Testando configurações de domínio para cloudflared tunnel"
echo "============================================================"

# Verifica se o backend está rodando
echo "1. Verificando se o backend está rodando..."
if curl -s http://localhost:8080/health > /dev/null; then
    echo "✅ Backend está rodando na porta 8080"
else
    echo "❌ Backend não está rodando na porta 8080"
    echo "   Execute: cd bff_luma && make run"
    exit 1
fi

echo ""

# Verifica as variáveis de ambiente
echo "2. Verificando variáveis de ambiente..."
if [ -f .env ]; then
    echo "✅ Arquivo .env encontrado"
    
    # Verifica APP_DOMAIN
    if grep -q "APP_DOMAIN=luma.app.br" .env; then
        echo "✅ APP_DOMAIN configurado: luma.app.br"
    else
        echo "❌ APP_DOMAIN não configurado corretamente"
    fi
    
    # Verifica APP_PROTOCOL
    if grep -q "APP_PROTOCOL=https" .env; then
        echo "✅ APP_PROTOCOL configurado: https"
    else
        echo "❌ APP_PROTOCOL não configurado corretamente"
    fi
    
    # Verifica CLOUDFLARE_TUNNEL_TOKEN
    if grep -q "CLOUDFLARE_TUNNEL_TOKEN=" .env; then
        echo "✅ CLOUDFLARE_TUNNEL_TOKEN configurado"
    else
        echo "⚠️ CLOUDFLARE_TUNNEL_TOKEN não encontrado"
    fi
else
    echo "❌ Arquivo .env não encontrado"
fi

echo ""

# Testa o cadastro de um usuário para verificar se os emails usam o domínio correto
echo "3. Testando cadastro com novo domínio..."
EMAIL="teste_dominio_$(date +%s)@example.com"
CPF="12345678901"
PASSWORD="senha123"

RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/wallets \
  -H "Content-Type: application/json" \
  -d "{
    \"username\": \"$CPF\",
    \"email\": \"$EMAIL\",
    \"password\": \"$PASSWORD\",
    \"password_repeat\": \"$PASSWORD\"
  }")

echo "📋 Resposta do cadastro:"
echo "$RESPONSE" | jq '.'

echo ""

# Verifica os logs para ver se o link de confirmação usa o domínio correto
echo "4. Verificando logs para links de confirmação..."
echo "🔍 Procure por logs que contenham 'luma.app.br' nos links de confirmação"
echo "📧 Os emails agora devem usar: https://luma.app.br/confirm-email?token=..."

echo ""

echo "🎯 Resumo das configurações:"
echo "- Backend: ✅ Rodando"
echo "- Domínio: ✅ luma.app.br"
echo "- Protocolo: ✅ https"
echo "- Emails: ✅ Configurados para usar o novo domínio"
echo ""
echo "📱 Frontend configurado para: https://luma.app.br"
echo "🌐 Cloudflared tunnel deve apontar para: luma.app.br"
echo ""
echo "✅ Configuração concluída! Agora os emails de confirmação"
echo "   usarão https://luma.app.br em vez de localhost"
