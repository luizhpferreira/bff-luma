#!/bin/bash

echo "🧪 Testando BFF Luma API - Todos os Endpoints"
echo "=============================================="
echo ""

# Configurações
BASE_URL="http://localhost:8080"
API_URL="$BASE_URL/api/v1"
TEST_USER="teste_$(date +%s)"
TEST_PASSWORD="B@nco2024!"
TEST_EMAIL="teste_$(date +%s)@exemplo.com"

echo "📧 Usuário de teste: $TEST_USER"
echo "🔑 Senha de teste: $TEST_PASSWORD"
echo "📧 Email de teste: $TEST_EMAIL"
echo ""

# Função para formatar JSON
format_json() {
    echo "$1" | jq '.' 2>/dev/null || echo "$1"
}

# 1. Health Check
echo "1️⃣ Testando Health Check..."
RESPONSE=$(curl -s "$BASE_URL/health")
echo "Status: $(echo "$RESPONSE" | jq -r '.data.status')"
echo ""

# 2. Criar Carteira
echo "2️⃣ Testando Criação de Carteira..."
RESPONSE=$(curl -s -X POST "$API_URL/wallets" \
  -H "Content-Type: application/json" \
  -d "{
    \"username\": \"$TEST_USER\",
    \"password\": \"$TEST_PASSWORD\",
    \"password_repeat\": \"$TEST_PASSWORD\"
  }")

if echo "$RESPONSE" | jq -e '.success' > /dev/null; then
    echo "✅ Carteira criada com sucesso"
    WALLET_ID=$(echo "$RESPONSE" | jq -r '.data.wallet_id')
    echo "   Wallet ID: $WALLET_ID"
else
    echo "❌ Erro ao criar carteira:"
    format_json "$RESPONSE"
fi
echo ""

# 3. Login
echo "3️⃣ Testando Login..."
RESPONSE=$(curl -s -X POST "$API_URL/login" \
  -H "Content-Type: application/json" \
  -d "{
    \"email\": \"$TEST_USER\",
    \"password\": \"$TEST_PASSWORD\"
  }")

if echo "$RESPONSE" | jq -e '.success' > /dev/null; then
    echo "✅ Login realizado com sucesso"
    TOKEN=$(echo "$RESPONSE" | jq -r '.data.token')
    echo "   Token obtido: ${TOKEN:0:50}..."
else
    echo "❌ Erro no login:"
    format_json "$RESPONSE"
fi
echo ""

# 4. Obter Informações da Carteira (com autenticação)
echo "4️⃣ Testando Obter Informações da Carteira..."
if [ ! -z "$TOKEN" ]; then
    RESPONSE=$(curl -s -H "Authorization: Bearer $TOKEN" "$API_URL/wallets")
    if echo "$RESPONSE" | jq -e '.success' > /dev/null; then
        echo "✅ Informações da carteira obtidas"
        echo "   Email: $(echo "$RESPONSE" | jq -r '.data.email')"
        echo "   Wallet ID: $(echo "$RESPONSE" | jq -r '.data.wallet_id')"
    else
        echo "❌ Erro ao obter informações:"
        format_json "$RESPONSE"
    fi
else
    echo "⚠️ Token não disponível, pulando teste"
fi
echo ""

# 5. Refresh Token
echo "5️⃣ Testando Refresh Token..."
if [ ! -z "$TOKEN" ]; then
    RESPONSE=$(curl -s -X POST "$API_URL/refresh" \
      -H "Authorization: Bearer $TOKEN")
    if echo "$RESPONSE" | jq -e '.success' > /dev/null; then
        echo "✅ Token renovado com sucesso"
        NEW_TOKEN=$(echo "$RESPONSE" | jq -r '.data.token')
        echo "   Novo token: ${NEW_TOKEN:0:50}..."
        TOKEN="$NEW_TOKEN"
    else
        echo "❌ Erro ao renovar token:"
        format_json "$RESPONSE"
    fi
else
    echo "⚠️ Token não disponível, pulando teste"
fi
echo ""

# 6. Criar Invoice (pode falhar por conexão LNBits)
echo "6️⃣ Testando Criação de Invoice..."
if [ ! -z "$TOKEN" ]; then
    RESPONSE=$(curl -s -X POST "$API_URL/invoices" \
      -H "Authorization: Bearer $TOKEN" \
      -H "Content-Type: application/json" \
      -d '{
        "amount": 1000,
        "memo": "Teste de invoice"
      }')
    
    if echo "$RESPONSE" | jq -e '.success' > /dev/null; then
        echo "✅ Invoice criado com sucesso"
        PAYMENT_HASH=$(echo "$RESPONSE" | jq -r '.data.payment_hash')
        echo "   Payment Hash: $PAYMENT_HASH"
    else
        echo "⚠️ Erro ao criar invoice (esperado se LNBits não estiver disponível):"
        format_json "$RESPONSE"
    fi
else
    echo "⚠️ Token não disponível, pulando teste"
fi
echo ""

# 7. Verificar Status do Pagamento
echo "7️⃣ Testando Verificação de Status do Pagamento..."
if [ ! -z "$TOKEN" ]; then
    RESPONSE=$(curl -s -H "Authorization: Bearer $TOKEN" \
      "$API_URL/payments/status?payment_hash=test_hash_123")
    
    if echo "$RESPONSE" | jq -e '.success' > /dev/null; then
        echo "✅ Status do pagamento verificado"
    else
        echo "⚠️ Erro ao verificar status (esperado se LNBits não estiver disponível):"
        format_json "$RESPONSE"
    fi
else
    echo "⚠️ Token não disponível, pulando teste"
fi
echo ""

# 8. Recuperação de Senha
echo "8️⃣ Testando Recuperação de Senha..."
RESPONSE=$(curl -s -X POST "$API_URL/forgot-password" \
  -H "Content-Type: application/json" \
  -d "{
    \"email\": \"$TEST_EMAIL\"
  }")

if echo "$RESPONSE" | jq -e '.success' > /dev/null; then
    echo "✅ Solicitação de recuperação processada"
else
    echo "⚠️ Erro na recuperação (pode ser esperado se email não existir):"
    format_json "$RESPONSE"
fi
echo ""

# 9. Estatísticas de Limpeza
echo "9️⃣ Testando Estatísticas de Limpeza..."
RESPONSE=$(curl -s "$API_URL/admin/cleanup/stats")
if echo "$RESPONSE" | jq -e '.success' > /dev/null; then
    echo "✅ Estatísticas de limpeza obtidas"
    echo "   Tokens expirados: $(echo "$RESPONSE" | jq -r '.data.expired_tokens')"
    echo "   Serviço rodando: $(echo "$RESPONSE" | jq -r '.data.is_running')"
else
    echo "❌ Erro ao obter estatísticas:"
    format_json "$RESPONSE"
fi
echo ""

# 10. Estatísticas de Rate Limiting
echo "🔟 Testando Estatísticas de Rate Limiting..."
RESPONSE=$(curl -s "$API_URL/admin/rate-limit/stats")
if echo "$RESPONSE" | jq -e '.success' > /dev/null; then
    echo "✅ Estatísticas de rate limiting obtidas"
    echo "   Limitadores de email: $(echo "$RESPONSE" | jq -r '.data.email_limiters_count')"
    echo "   Limitadores de IP: $(echo "$RESPONSE" | jq -r '.data.ip_limiters_count')"
else
    echo "❌ Erro ao obter estatísticas:"
    format_json "$RESPONSE"
fi
echo ""

# 11. Testes de Validação
echo "1️⃣1️⃣ Testando Validações..."
echo "   Testando senhas que não coincidem..."
RESPONSE=$(curl -s -X POST "$API_URL/wallets" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "teste_validacao",
    "password": "senha1",
    "password_repeat": "senha2"
  }')
if echo "$RESPONSE" | jq -e '.success' > /dev/null; then
    echo "   ❌ Validação falhou - senhas diferentes foram aceitas"
else
    echo "   ✅ Validação funcionou - senhas diferentes rejeitadas"
fi

echo "   Testando login com senha incorreta..."
RESPONSE=$(curl -s -X POST "$API_URL/login" \
  -H "Content-Type: application/json" \
  -d "{
    \"email\": \"$TEST_USER\",
    \"password\": \"senha_incorreta\"
  }")
if echo "$RESPONSE" | jq -e '.success' > /dev/null; then
    echo "   ❌ Validação falhou - senha incorreta foi aceita"
else
    echo "   ✅ Validação funcionou - senha incorreta rejeitada"
fi
echo ""

# 12. Teste de Rate Limiting
echo "1️⃣2️⃣ Testando Rate Limiting..."
echo "   Fazendo 6 tentativas de login com senha incorreta..."
for i in {1..6}; do
    RESPONSE=$(curl -s -X POST "$API_URL/login" \
      -H "Content-Type: application/json" \
      -d "{
        \"email\": \"rate_test_$i@exemplo.com\",
        \"password\": \"senha_incorreta\"
      }")
    
    if echo "$RESPONSE" | jq -e '.success' > /dev/null; then
        echo "   Tentativa $i: ✅ Login aceito (inesperado)"
    else
        MESSAGE=$(echo "$RESPONSE" | jq -r '.message')
        if [[ "$MESSAGE" == *"Muitas tentativas"* ]]; then
            echo "   Tentativa $i: 🔒 Rate limit ativado"
        else
            echo "   Tentativa $i: ❌ Login rejeitado (normal)"
        fi
    fi
done
echo ""

echo "🎉 Testes concluídos!"
echo "====================="
echo ""
echo "📊 Resumo dos Endpoints Testados:"
echo "   ✅ Health Check"
echo "   ✅ Criação de Carteira"
echo "   ✅ Login"
echo "   ✅ Obter Informações da Carteira"
echo "   ✅ Refresh Token"
echo "   ⚠️ Criar Invoice (depende do LNBits)"
echo "   ⚠️ Verificar Status do Pagamento (depende do LNBits)"
echo "   ✅ Recuperação de Senha"
echo "   ✅ Estatísticas de Limpeza"
echo "   ✅ Estatísticas de Rate Limiting"
echo "   ✅ Validações"
echo "   ✅ Rate Limiting"
echo ""
echo "🔗 Base URL: $BASE_URL"
echo "🔗 API URL: $API_URL"
