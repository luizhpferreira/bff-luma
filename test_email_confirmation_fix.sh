#!/bin/bash

# Teste de confirmação de email após correção da tabela temp_passwords
echo "🧪 Testando confirmação de email após correção..."

# URL base
BASE_URL="http://localhost:8080"

# Dados de teste
TEST_EMAIL="teste_$(date +%s)@exemplo.com"
TEST_CPF="01383972281"  # CPF válido
TEST_PASSWORD="#Ruiter1"  # Senha forte com 12 caracteres

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
    echo "📋 Últimos logs do aplicativo:"
    docker logs bff_luma_app --tail 20
    exit 1
fi

echo "🔑 Token encontrado: $TOKEN"

# 3. Confirmar email
echo "3️⃣ Confirmando email..."
CONFIRM_RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/confirm-email" \
  -H "Content-Type: application/json" \
  -d "{
    \"token\": \"$TOKEN\"
  }")

echo "Resposta da confirmação: $CONFIRM_RESPONSE"

# 4. Verificar se a confirmação foi bem-sucedida
if echo "$CONFIRM_RESPONSE" | grep -q "Email confirmado com sucesso"; then
    echo "✅ Confirmação de email funcionando corretamente!"
else
    echo "❌ Erro na confirmação de email"
    echo "📋 Logs do aplicativo:"
    docker logs bff_luma_app --tail 10
fi

# 5. Verificar se a tabela temp_passwords foi usada corretamente
echo "4️⃣ Verificando uso da tabela temp_passwords..."
TEMP_PASSWORD_COUNT=$(docker exec bff_luma_postgres psql -U postgres -d bff_luma -t -c "SELECT COUNT(*) FROM temp_passwords WHERE email = '$TEST_EMAIL';" | xargs)

echo "📊 Registros na tabela temp_passwords para $TEST_EMAIL: $TEMP_PASSWORD_COUNT"

# 6. Verificar se a carteira foi confirmada
echo "5️⃣ Verificando status da carteira..."
WALLET_STATUS=$(docker exec bff_luma_postgres psql -U postgres -d bff_luma -t -c "SELECT email_confirmed FROM wallets WHERE email = '$TEST_EMAIL';" | xargs)

echo "📊 Status de confirmação da carteira: $WALLET_STATUS"

if [ "$WALLET_STATUS" = "t" ]; then
    echo "✅ Carteira confirmada com sucesso!"
else
    echo "❌ Carteira não foi confirmada"
fi

echo "�� Teste concluído!"
