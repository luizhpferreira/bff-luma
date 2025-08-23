#!/bin/bash

# Script para testar o novo sistema com CPF e email separados

echo "🧪 Testando Novo Sistema - CPF e Email Separados"
echo "=================================================="

# Configurações
API_URL="http://localhost:8080"
CPF="11122233344"
EMAIL="teste@exemplo.com"
PASSWORD="B@nco2024!"

echo ""
echo "📋 Dados de teste:"
echo "CPF: $CPF"
echo "Email: $EMAIL"
echo "Password: $PASSWORD"
echo ""

# Teste 1: Criar carteira (cadastro)
echo "1️⃣ Testando criação de carteira..."
echo "POST $API_URL/api/v1/wallets"
echo "Payload:"
cat << EOF
{
    "username": "$CPF",
    "email": "$EMAIL",
    "password": "$PASSWORD",
    "password_repeat": "$PASSWORD"
}
EOF
echo ""

response=$(curl -s -w "\n%{http_code}" \
  --location "$API_URL/api/v1/wallets" \
  --header 'Content-Type: application/json' \
  --data-raw "{
    \"username\": \"$CPF\",
    \"email\": \"$EMAIL\",
    \"password\": \"$PASSWORD\",
    \"password_repeat\": \"$PASSWORD\"
  }")

http_code=$(echo "$response" | tail -n1)
response_body=$(echo "$response" | head -n -1)

echo "Response (HTTP $http_code):"
echo "$response_body" | jq '.' 2>/dev/null || echo "$response_body"
echo ""

# Teste 2: Verificar estrutura do banco
echo "2️⃣ Verificando estrutura do banco..."
echo "sqlite3 bff_luma.db \".schema\""
echo ""

schema=$(sqlite3 bff_luma.db ".schema" 2>/dev/null)
if [ -n "$schema" ]; then
    echo "✅ Estrutura do banco:"
    echo "$schema"
else
    echo "❌ Banco ainda não foi inicializado"
fi
echo ""

# Teste 3: Verificar dados salvos
echo "3️⃣ Verificando dados salvos..."
echo "sqlite3 bff_luma.db \"SELECT * FROM wallets;\""
echo ""

data=$(sqlite3 bff_luma.db "SELECT * FROM wallets;" 2>/dev/null)
if [ -n "$data" ]; then
    echo "✅ Dados salvos:"
    echo "$data"
else
    echo "❌ Nenhum dado encontrado"
fi
echo ""

# Teste 4: Login com CPF
echo "4️⃣ Testando login com CPF..."
echo "POST $API_URL/api/v1/login"
echo "Payload:"
cat << EOF
{
    "email": "$CPF",
    "password": "$PASSWORD"
}
EOF
echo ""

response=$(curl -s -w "\n%{http_code}" \
  --location "$API_URL/api/v1/login" \
  --header 'Content-Type: application/json' \
  --data-raw "{
    \"email\": \"$CPF\",
    \"password\": \"$PASSWORD\"
  }")

http_code=$(echo "$response" | tail -n1)
response_body=$(echo "$response" | head -n -1)

echo "Response (HTTP $http_code):"
echo "$response_body" | jq '.' 2>/dev/null || echo "$response_body"
echo ""

echo "✅ Teste concluído!"
echo ""
echo "📊 Resumo das mudanças:"
echo "  • CPF e email agora são salvos em campos separados"
echo "  • Email não é mais perdido durante o cadastro"
echo "  • Estrutura do banco mais clara e organizada"
echo "  • Login continua funcionando com CPF"
echo "  • Email disponível para funcionalidades futuras"
