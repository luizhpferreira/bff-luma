#!/bin/bash

echo "🔐 Testando login..."

# Fazer login
LOGIN_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d "{
    \"email\": \"01383972281\",
    \"password\": \"#Ruiter1\"
  }")

echo "📋 Resposta do login:"
echo "$LOGIN_RESPONSE" | jq '.' 2>/dev/null || echo "$LOGIN_RESPONSE"
echo ""

# Extrair o token JWT
TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.data.token' 2>/dev/null)

if [ "$TOKEN" != "null" ] && [ "$TOKEN" != "" ]; then
    echo "✅ Login realizado com sucesso!"
    echo "🔑 Token JWT obtido"
    echo ""
    
    echo "🧾 Testando criação de invoice..."
    
    # Criar invoice
    INVOICE_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/invoices \
      -H "Content-Type: application/json" \
      -H "Authorization: Bearer $TOKEN" \
      -d "{
        \"amount\": 1000,
        \"memo\": \"Teste de invoice\"
      }")
    
    echo "📋 Resposta da criação do invoice:"
    echo "$INVOICE_RESPONSE" | jq '.' 2>/dev/null || echo "$INVOICE_RESPONSE"
    echo ""
    
    # Extrair payment hash
    PAYMENT_HASH=$(echo "$INVOICE_RESPONSE" | jq -r '.data.payment_hash' 2>/dev/null)
    
    if [ "$PAYMENT_HASH" != "null" ] && [ "$PAYMENT_HASH" != "" ]; then
        echo "✅ Invoice criado com sucesso!"
        echo "🔗 Payment Hash: $PAYMENT_HASH"
        echo ""
        
        echo "🔍 Verificando status do pagamento..."
        
        # Verificar status do pagamento
        STATUS_RESPONSE=$(curl -s -X GET "http://localhost:8080/api/v1/payments/status?payment_hash=$PAYMENT_HASH" \
          -H "Authorization: Bearer $TOKEN")
        
        echo "📋 Status do pagamento:"
        echo "$STATUS_RESPONSE" | jq '.' 2>/dev/null || echo "$STATUS_RESPONSE"
        
    else
        echo "❌ Erro ao criar invoice"
    fi
    
else
    echo "❌ Erro no login"
fi

echo ""
echo "�� Teste concluído!"
