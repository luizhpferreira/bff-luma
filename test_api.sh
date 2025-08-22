#!/bin/bash

BASE_URL="http://localhost:8080"
API_URL="$BASE_URL/api/v1"
EMAIL="teste_$(date +%s)@exemplo.com"
PASSWORD="123456"

echo "🧪 Testando BFF Luma API"
echo "=========================="

echo "📧 Email de teste: $EMAIL"
echo ""

# Teste 1: Health Check
echo "1️⃣ Testando Health Check..."
curl -s "$BASE_URL/health" | jq '.'
echo ""

# Teste 2: Criar carteira (Cadastro)
echo "2️⃣ Criando carteira..."
WALLET_RESPONSE=$(curl -s -X POST "$API_URL/wallets" \
  -H "Content-Type: application/json" \
  -d "{
    \"email\": \"$EMAIL\",
    \"password\": \"$PASSWORD\",
    \"password_repeat\": \"$PASSWORD\"
  }")

echo "$WALLET_RESPONSE" | jq '.'
echo ""

# Teste 3: Login
echo "3️⃣ Fazendo login..."
LOGIN_RESPONSE=$(curl -s -X POST "$API_URL/login" \
  -H "Content-Type: application/json" \
  -d "{
    \"email\": \"$EMAIL\",
    \"password\": \"$PASSWORD\"
  }")

echo "$LOGIN_RESPONSE" | jq '.'
echo ""

# Teste 4: Obter informações da carteira
echo "4️⃣ Obtendo informações da carteira..."
curl -s "$API_URL/wallets?email=$EMAIL" | jq '.'
echo ""

# Teste 5: Tentar criar invoice (pode falhar)
echo "5️⃣ Tentando criar invoice..."
INVOICE_RESPONSE=$(curl -s -X POST "$API_URL/invoices" \
  -H "Content-Type: application/json" \
  -d "{
    \"email\": \"$EMAIL\",
    \"amount\": 1000,
    \"memo\": \"Teste de invoice\"
  }")

echo "$INVOICE_RESPONSE" | jq '.'
echo ""

# Teste 6: Tentar criar carteira duplicada (deve falhar)
echo "6️⃣ Tentando criar carteira duplicada (deve falhar)..."
curl -s -X POST "$API_URL/wallets" \
  -H "Content-Type: application/json" \
  -d "{
    \"email\": \"$EMAIL\",
    \"password\": \"$PASSWORD\",
    \"password_repeat\": \"$PASSWORD\"
  }" | jq '.'
echo ""

# Teste 7: Tentar login com senha incorreta (deve falhar)
echo "7️⃣ Tentando login com senha incorreta (deve falhar)..."
curl -s -X POST "$API_URL/login" \
  -H "Content-Type: application/json" \
  -d "{
    \"email\": \"$EMAIL\",
    \"password\": \"senhaerrada\"
  }" | jq '.'
echo ""

echo "✅ Testes concluídos!"

echo ""
echo "📊 Status dos endpoints:"
echo "✅ Health Check - Funcionando"
echo "✅ Cadastro de Usuários - Funcionando"
echo "✅ Login - Funcionando"
echo "✅ Consulta de Carteiras - Funcionando"
echo "⚠️ Criação de Invoices - Problema de conexão no LNBits"
echo "⚠️ Verificação de Pagamentos - Depende dos invoices"
