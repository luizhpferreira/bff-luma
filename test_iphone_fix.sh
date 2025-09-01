#!/bin/bash

# Teste da correção do problema do iPhone com confirmação de email
echo "🍎 Testando correção do problema do iPhone..."

# URL base
BASE_URL="http://localhost:8080"

# Dados de teste
TEST_EMAIL="amanda_fix_$(date +%s)@exemplo.com"
TEST_CPF="01383972282"  # CPF válido diferente
TEST_PASSWORD="#Amanda1#"

echo "📧 Email de teste: $TEST_EMAIL"
echo "🆔 CPF de teste: $TEST_CPF"

# 1. Criar carteira
echo "1️⃣ Criando carteira..."
CREATE_RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/wallets" \
  -H "Content-Type: application/json" \
  -d "{
    \"username\": \"$TEST_CPF\",
    \"email\": \"$TEST_EMAIL\",
    \"password\": \"$TEST_PASSWORD\",
    \"password_repeat\": \"$TEST_PASSWORD\"
  }")

echo "Resposta da criação: $CREATE_RESPONSE"

# Extrair token de confirmação dos logs
echo "2️⃣ Buscando token de confirmação nos logs..."
TOKEN=$(docker logs bff_luma_app --tail 50 | grep "Token de confirmação criado para $TEST_EMAIL" | tail -1 | sed 's/.*Token de confirmação criado para.*: //')

if [ -z "$TOKEN" ]; then
    echo "❌ Token de confirmação não encontrado nos logs"
    exit 1
fi

echo "🔑 Token encontrado: $TOKEN"

# 3. Testar primeira confirmação (deve funcionar)
echo "3️⃣ Primeira confirmação (deve funcionar)..."
FIRST_RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/confirm-email" \
  -H "Content-Type: application/json" \
  -d "{
    \"token\": \"$TOKEN\"
  }")

echo "📤 Resposta da primeira confirmação: $FIRST_RESPONSE"

# 4. Testar segunda confirmação (deve retornar "já confirmado")
echo "4️⃣ Segunda confirmação (deve retornar 'já confirmado')..."
SECOND_RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/confirm-email" \
  -H "Content-Type: application/json" \
  -d "{
    \"token\": \"$TOKEN\"
  }")

echo "📤 Resposta da segunda confirmação: $SECOND_RESPONSE"

# 5. Verificar se a segunda resposta indica que já foi confirmado
if echo "$SECOND_RESPONSE" | grep -q "Email já foi confirmado"; then
    echo "✅ Correção funcionando! Segunda tentativa retorna 'já confirmado'"
else
    echo "❌ Correção não funcionou. Segunda tentativa ainda retorna erro"
fi

# 6. Verificar status final no banco
echo "5️⃣ Verificando status final no banco..."
FINAL_STATUS=$(docker exec bff_luma_postgres psql -U postgres -d bff_luma -t -c "SELECT email_confirmed, wallet_id FROM wallets WHERE email = '$TEST_EMAIL';" | xargs)

echo "📊 Status final da carteira: $FINAL_STATUS"

echo "🍎 Teste da correção concluído!"
