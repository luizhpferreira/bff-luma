#!/bin/bash

# Script para testar com CPFs únicos

echo "🧪 Testando com CPF único..."

# Gerar CPF único baseado no timestamp
TIMESTAMP=$(date +%s)
CPF="123456789${TIMESTAMP: -2}"  # Últimos 2 dígitos do timestamp
EMAIL="teste${TIMESTAMP}@exemplo.com"

echo "CPF gerado: $CPF"
echo "Email gerado: $EMAIL"
echo ""

# Testar cadastro
echo "Testando cadastro..."
response=$(curl -s -w "\n%{http_code}" \
  --location "http://localhost:8080/api/v1/wallets" \
  --header 'Content-Type: application/json' \
  --data-raw "{
    \"username\": \"$CPF\",
    \"email\": \"$EMAIL\",
    \"password\": \"B@nco2024!\",
    \"password_repeat\": \"B@nco2024!\"
  }")

http_code=$(echo "$response" | tail -n1)
response_body=$(echo "$response" | head -n -1)

echo "Response (HTTP $http_code):"
echo "$response_body" | jq '.' 2>/dev/null || echo "$response_body"

if [ "$http_code" = "201" ]; then
    echo ""
    echo "✅ Sucesso! Usuário criado com CPF: $CPF"
    echo "Email: $EMAIL"
    
    # Verificar banco
    echo ""
    echo "📊 Verificando banco de dados..."
    sqlite3 bff_luma.db "SELECT id, cpf, email, wallet_id, created_at FROM wallets WHERE cpf = '$CPF';"
else
    echo ""
    echo "❌ Erro no cadastro"
fi
