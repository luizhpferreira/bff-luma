#!/bin/bash

echo "🎯 Teste Final da Correção do Problema da Amanda..."

# URL base
BASE_URL="http://localhost:8080"

# Dados de teste
TEST_EMAIL="amanda_final_$(date +%s)@exemplo.com"
TEST_CPF="39130037122"  # CPF válido
TEST_PASSWORD="#Amanda123!"

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

# 3. Testar primeira confirmação
echo "3️⃣ Primeira confirmação (deve funcionar)..."
FIRST_RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/confirm-email" \
  -H "Content-Type: application/json" \
  -d "{
    \"token\": \"$TOKEN\"
  }")

echo "📤 Resposta da primeira confirmação: $FIRST_RESPONSE"

# 4. Verificar se a primeira confirmação foi bem-sucedida
if echo "$FIRST_RESPONSE" | grep -q "success.*true"; then
    echo "✅ Primeira confirmação bem-sucedida!"
else
    echo "❌ Primeira confirmação falhou"
    exit 1
fi

# 5. Verificar status no banco
echo "4️⃣ Verificando status no banco..."
WALLET_STATUS=$(docker exec bff_luma_postgres psql -U postgres -d bff_luma -t -c "SELECT email_confirmed, wallet_id FROM wallets WHERE email = '$TEST_EMAIL';" | xargs)

echo "📊 Status da carteira: $WALLET_STATUS"

# 6. Testar segunda confirmação (deve retornar token inválido, mas isso é esperado)
echo "5️⃣ Segunda confirmação (deve retornar token inválido - isso é esperado)..."
SECOND_RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/confirm-email" \
  -H "Content-Type: application/json" \
  -d "{
    \"token\": \"$TOKEN\"
  }")

echo "📤 Resposta da segunda confirmação: $SECOND_RESPONSE"

# 7. Verificar se a carteira ainda está confirmada
echo "6️⃣ Verificando se a carteira permanece confirmada..."
FINAL_STATUS=$(docker exec bff_luma_postgres psql -U postgres -d bff_luma -t -c "SELECT email_confirmed FROM wallets WHERE email = '$TEST_EMAIL';" | xargs)

if [ "$FINAL_STATUS" = "t" ]; then
    echo "✅ Carteira permanece confirmada após segunda tentativa!"
    echo "🎉 Correção funcionando corretamente!"
else
    echo "❌ Carteira não está mais confirmada"
fi

echo "🎯 Teste final concluído!"
