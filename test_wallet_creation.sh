#!/bin/bash

echo "🧪 Testando criação de carteira no BFF..."

# Dados de teste (CPF válido diferente)
CPF="01383972281"
EMAIL="luizferreiralps@gmail.com"
PASSWORD="#Ruiter1"

echo "📝 Dados de teste:"
echo "  CPF: $CPF"
echo "  Email: $EMAIL"
echo "  Senha: $PASSWORD"
echo ""

# Criar carteira
echo "🚀 Criando carteira..."
RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/wallets \
  -H "Content-Type: application/json" \
  -d "{
    \"username\": \"$CPF\",
    \"email\": \"$EMAIL\",
    \"password\": \"$PASSWORD\",
    \"password_repeat\": \"$PASSWORD\"
  }")

echo "📋 Resposta da criação:"
echo "$RESPONSE" | jq '.' 2>/dev/null || echo "$RESPONSE"
echo ""

# Verificar se a carteira foi criada
if echo "$RESPONSE" | grep -q "success.*true"; then
    echo "✅ Carteira criada com sucesso!"
    
    # Extrair wallet_id da resposta
    WALLET_ID=$(echo "$RESPONSE" | jq -r '.data.wallet_id' 2>/dev/null)
    if [ "$WALLET_ID" != "null" ] && [ "$WALLET_ID" != "" ]; then
        echo "🆔 Wallet ID: $WALLET_ID"
    fi
    
    echo ""
    echo "📧 Verifique o email para confirmar a conta..."
    echo "🔗 Link de confirmação: http://localhost:8080/confirm-email?token=TOKEN_DO_EMAIL"
    
else
    echo "❌ Erro na criação da carteira"
    echo "🔍 Verifique os logs do BFF para mais detalhes"
fi

echo ""
echo "🏁 Teste concluído!"
