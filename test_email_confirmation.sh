#!/bin/bash

echo "🧪 Testando confirmação de email..."

# Buscar o token de confirmação no banco de dados
echo "🔍 Buscando token de confirmação no banco..."

# Conectar ao banco e buscar o token
TOKEN=$(psql -h localhost -U bff_luma -d bff_luma -t -c "SELECT email_confirmation_token FROM wallets WHERE email = 'luizferreiralps@gmail.com' AND email_confirmation_token IS NOT NULL;" 2>/dev/null | xargs)

if [ -z "$TOKEN" ]; then
    echo "❌ Token não encontrado no banco de dados"
    echo "🔍 Verificando se a carteira existe..."
    
    # Verificar se a carteira existe
    WALLET_EXISTS=$(psql -h localhost -U bff_luma -d bff_luma -t -c "SELECT COUNT(*) FROM wallets WHERE email = 'luizferreiralps@gmail.com';" 2>/dev/null | xargs)
    
    if [ "$WALLET_EXISTS" = "1" ]; then
        echo "✅ Carteira encontrada no banco"
        echo "📋 Dados da carteira:"
        psql -h localhost -U bff_luma -d bff_luma -c "SELECT cpf, email, wallet_id, admin_key, invoice_key, email_confirmed FROM wallets WHERE email = 'luizferreiralps@gmail.com';" 2>/dev/null
    else
        echo "❌ Carteira não encontrada no banco"
    fi
    exit 1
fi

echo "🔑 Token encontrado: $TOKEN"
echo ""

# Confirmar email usando o token
echo "✅ Confirmando email..."
RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/confirm-email \
  -H "Content-Type: application/json" \
  -d "{\"token\": \"$TOKEN\"}")

echo "📋 Resposta da confirmação:"
echo "$RESPONSE" | jq '.' 2>/dev/null || echo "$RESPONSE"
echo ""

# Verificar se a confirmação foi bem-sucedida
if echo "$RESPONSE" | grep -q "success.*true"; then
    echo "✅ Email confirmado com sucesso!"
    
    # Verificar se a wallet foi criada no LNBits
    echo ""
    echo "🔍 Verificando dados da wallet no banco..."
    psql -h localhost -U bff_luma -d bff_luma -c "SELECT cpf, email, wallet_id, admin_key, invoice_key, email_confirmed FROM wallets WHERE email = 'luizferreiralps@gmail.com';" 2>/dev/null
    
    # Testar login
    echo ""
    echo "🔐 Testando login..."
    LOGIN_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/login \
      -H "Content-Type: application/json" \
      -d "{
        \"email\": \"01383972281\",
        \"password\": \"#Ruiter1\"
      }")
    
    echo "📋 Resposta do login:"
    echo "$LOGIN_RESPONSE" | jq '.' 2>/dev/null || echo "$LOGIN_RESPONSE"
    
else
    echo "❌ Erro na confirmação do email"
    echo "🔍 Verifique os logs do BFF para mais detalhes"
fi

echo ""
echo "�� Teste concluído!"
