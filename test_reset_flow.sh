#!/bin/bash

# Script para testar o fluxo completo de reset de senha
# Testa: solicitação -> email -> página web -> deep link -> app

set -e

echo "🧪 Testando fluxo completo de reset de senha..."
echo "================================================"

# Configurações
API_BASE_URL="https://luma.app.br"
TEST_EMAIL="test@example.com"

echo "📧 1. Testando solicitação de reset de senha..."
echo "Email: $TEST_EMAIL"

# Solicitar reset de senha
RESPONSE=$(curl -s -X POST "$API_BASE_URL/api/v1/forgot-password" \
  -H "Content-Type: application/json" \
  -d "{\"email\": \"$TEST_EMAIL\"}")

echo "Resposta: $RESPONSE"

# Verificar se a resposta foi bem-sucedida
if echo "$RESPONSE" | grep -q '"success":true'; then
    echo "✅ Solicitação de reset enviada com sucesso"
else
    echo "❌ Erro na solicitação de reset"
    echo "$RESPONSE"
    exit 1
fi

echo ""
echo "🔗 2. Testando página de reset de senha..."
echo "URL: $API_BASE_URL/reset-password?token=test-token"

# Testar página de reset
PAGE_RESPONSE=$(curl -s "$API_BASE_URL/reset-password?token=test-token")

if echo "$PAGE_RESPONSE" | grep -q "Redefinir Senha"; then
    echo "✅ Página de reset carregada corretamente"
else
    echo "❌ Erro ao carregar página de reset"
    exit 1
fi

echo ""
echo "🔑 3. Testando validação de token..."
echo "Token: test-token"

# Testar validação de token
VALIDATION_RESPONSE=$(curl -s -X POST "$API_BASE_URL/api/v1/validate-reset-token" \
  -H "Content-Type: application/json" \
  -d '{"token": "test-token"}')

echo "Resposta: $VALIDATION_RESPONSE"

# Verificar se a resposta indica token inválido (esperado)
if echo "$VALIDATION_RESPONSE" | grep -q '"success":false'; then
    echo "✅ Validação de token funcionando (token inválido rejeitado)"
else
    echo "❌ Erro na validação de token"
    echo "$VALIDATION_RESPONSE"
    exit 1
fi

echo ""
echo "📱 4. Verificando deep link na página..."
echo "Procurando por: bffluma://reset-password"

# Verificar se o deep link está na página
if echo "$PAGE_RESPONSE" | grep -q "bffluma://reset-password"; then
    echo "✅ Deep link encontrado na página"
else
    echo "❌ Deep link não encontrado na página"
    exit 1
fi

echo ""
echo "🌐 5. Verificando configurações do domínio..."

# Verificar se está usando HTTPS
if echo "$PAGE_RESPONSE" | grep -q "https://luma.app.br"; then
    echo "✅ Usando HTTPS e domínio correto"
else
    echo "❌ Não está usando HTTPS ou domínio incorreto"
    exit 1
fi

echo ""
echo "✅ Todos os testes passaram!"
echo "================================================"
echo ""
echo "📋 Resumo do fluxo:"
echo "• Backend configurado com domínio correto"
echo "• Página de reset funcionando"
echo "• Deep link configurado corretamente"
echo "• API de validação funcionando"
echo ""
echo "🚀 Próximos passos:"
echo "1. Solicite reset de senha no app"
echo "2. Verifique o email recebido"
echo "3. Clique no link do email"
echo "4. Verifique se a página web carrega"
echo "5. Clique em 'Abrir App' na página web"
echo "6. Verifique se o app abre na tela de reset"
echo ""
echo "🔗 Link de teste:"
echo "https://luma.app.br/reset-password?token=test-token"
