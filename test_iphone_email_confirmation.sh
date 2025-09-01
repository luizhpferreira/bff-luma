#!/bin/bash

# Teste específico para reproduzir o problema do iPhone com confirmação de email
echo "🍎 Testando problema de confirmação de email no iPhone..."

# URL base
BASE_URL="http://localhost:8080"

# Dados de teste
TEST_EMAIL="amanda_marinelli_$(date +%s)@exemplo.com"
TEST_CPF="01383972281"  # CPF de teste
TEST_PASSWORD="#Amanda123"  # Senha forte

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

# 3. Simular o comportamento do iPhone - fazer múltiplas tentativas
echo "3️⃣ Simulando comportamento do iPhone - múltiplas tentativas..."

for i in {1..5}; do
    echo "🔄 Tentativa $i..."
    
    # Primeiro, acessar a página de confirmação (como o iPhone faz)
    PAGE_RESPONSE=$(curl -s -X GET "$BASE_URL/confirm-email?token=$TOKEN")
    echo "📄 Página de confirmação acessada (status: $?)"
    
    # Depois, fazer a requisição POST (como o JavaScript da página faz)
    CONFIRM_RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/confirm-email" \
      -H "Content-Type: application/json" \
      -d "{
        \"token\": \"$TOKEN\"
      }")
    
    echo "📤 Resposta da confirmação: $CONFIRM_RESPONSE"
    
    # Verificar se foi bem-sucedida
    if echo "$CONFIRM_RESPONSE" | grep -q "success.*true"; then
        echo "✅ Confirmação bem-sucedida na tentativa $i!"
        break
    else
        echo "❌ Falha na tentativa $i"
        
        # Verificar status da carteira no banco
        WALLET_STATUS=$(docker exec bff_luma_postgres psql -U postgres -d bff_luma -t -c "SELECT email_confirmed FROM wallets WHERE email = '$TEST_EMAIL';" | xargs)
        echo "📊 Status da carteira após tentativa $i: $WALLET_STATUS"
        
        # Aguardar um pouco antes da próxima tentativa
        sleep 2
    fi
done

# 4. Verificar status final
echo "4️⃣ Verificando status final..."
FINAL_STATUS=$(docker exec bff_luma_postgres psql -U postgres -d bff_luma -t -c "SELECT email_confirmed, wallet_id FROM wallets WHERE email = '$TEST_EMAIL';" | xargs)

echo "📊 Status final da carteira: $FINAL_STATUS"

# 5. Verificar logs para entender o que aconteceu
echo "5️⃣ Últimos logs do aplicativo:"
docker logs bff_luma_app --tail 30

echo "�� Teste concluído!"
