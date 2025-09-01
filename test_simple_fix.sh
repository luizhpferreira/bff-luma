#!/bin/bash

echo "🧪 Teste simples da correção do problema do iPhone..."

# Criar carteira com CPF válido
echo "1️⃣ Criando carteira..."
RESPONSE=$(curl -s -X POST "http://localhost:8080/api/v1/wallets" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "11023887933",
    "email": "test_fix@exemplo.com",
    "password": "#Test123!",
    "password_repeat": "#Test123!"
  }')

echo "Resposta: $RESPONSE"

# Buscar token
echo "2️⃣ Buscando token..."
TOKEN=$(docker logs bff_luma_app --tail 30 | grep "Token de confirmação criado para test_fix@exemplo.com" | tail -1 | sed 's/.*Token de confirmação criado para.*: //')

if [ -z "$TOKEN" ]; then
    echo "❌ Token não encontrado"
    exit 1
fi

echo "🔑 Token: $TOKEN"

# Testar primeira confirmação
echo "3️⃣ Primeira confirmação..."
FIRST=$(curl -s -X POST "http://localhost:8080/api/v1/confirm-email" \
  -H "Content-Type: application/json" \
  -d "{\"token\": \"$TOKEN\"}")

echo "Primeira: $FIRST"

# Testar segunda confirmação
echo "4️⃣ Segunda confirmação..."
SECOND=$(curl -s -X POST "http://localhost:8080/api/v1/confirm-email" \
  -H "Content-Type: application/json" \
  -d "{\"token\": \"$TOKEN\"}")

echo "Segunda: $SECOND"

# Verificar se a segunda retorna "já confirmado"
if echo "$SECOND" | grep -q "Email já foi confirmado"; then
    echo "✅ Correção funcionando!"
else
    echo "❌ Correção não funcionou"
fi
