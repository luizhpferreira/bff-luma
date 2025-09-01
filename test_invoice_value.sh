#!/bin/bash

echo "🧪 Testando criação de invoice com valor 1000 sats..."

# Primeiro, vamos fazer login para obter um token
echo "📝 Fazendo login..."
LOGIN_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "test123"
  }')

echo "Login response: $LOGIN_RESPONSE"

# Extrair o token da resposta
TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"token":"[^"]*"' | cut -d'"' -f4)

if [ -z "$TOKEN" ]; then
    echo "❌ Falha ao obter token de autenticação"
    exit 1
fi

echo "✅ Token obtido: ${TOKEN:0:20}..."

# Agora vamos criar um invoice com valor 1000
echo "💰 Criando invoice com valor 1000 sats..."
INVOICE_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/invoices \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "amount": 1000,
    "memo": "Teste valor invoice"
  }')

echo "Invoice response: $INVOICE_RESPONSE"

# Extrair o valor da resposta
AMOUNT=$(echo $INVOICE_RESPONSE | grep -o '"amount":[0-9]*' | cut -d':' -f2)

echo "📊 Resultado:"
echo "   Valor enviado: 1000 sats"
echo "   Valor retornado: $AMOUNT sats"

if [ "$AMOUNT" = "1000" ]; then
    echo "✅ Teste PASSOU - Valor correto!"
else
    echo "❌ Teste FALHOU - Valor incorreto!"
    echo "   Esperado: 1000, Recebido: $AMOUNT"
fi
