#!/bin/bash

echo "🔍 Verificando token de confirmação..."

# Usar variáveis de ambiente para conectar ao banco
export PGPASSWORD="bff_luma"

# Buscar o token
TOKEN=$(psql -h localhost -U bff_luma -d bff_luma -t -c "SELECT email_confirmation_token FROM wallets WHERE email = 'luizferreiralps@gmail.com' AND email_confirmation_token IS NOT NULL;" 2>/dev/null | xargs)

if [ -z "$TOKEN" ]; then
    echo "❌ Token não encontrado"
    echo "📋 Verificando dados da carteira:"
    psql -h localhost -U bff_luma -d bff_luma -c "SELECT cpf, email, wallet_id, admin_key, invoice_key, email_confirmed FROM wallets WHERE email = 'luizferreiralps@gmail.com';" 2>/dev/null
else
    echo "🔑 Token encontrado: $TOKEN"
    echo ""
    echo "✅ Confirmando email..."
    
    RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/confirm-email \
      -H "Content-Type: application/json" \
      -d "{\"token\": \"$TOKEN\"}")
    
    echo "📋 Resposta:"
    echo "$RESPONSE" | jq '.' 2>/dev/null || echo "$RESPONSE"
    
    if echo "$RESPONSE" | grep -q "success.*true"; then
        echo ""
        echo "🎉 Email confirmado! Verificando dados da wallet:"
        psql -h localhost -U bff_luma -d bff_luma -c "SELECT cpf, email, wallet_id, admin_key, invoice_key, email_confirmed FROM wallets WHERE email = 'luizferreiralps@gmail.com';" 2>/dev/null
    fi
fi
