#!/bin/bash

echo "🔧 Configuração do Cloudflare Tunnel Token"
echo "=========================================="

# Verifica se o token foi fornecido como argumento
if [ -z "$1" ]; then
    echo "❌ Token não fornecido!"
    echo ""
    echo "📋 Uso:"
    echo "  ./setup_cloudflare.sh SEU_TOKEN_AQUI"
    echo ""
    echo "🔍 Para obter o token:"
    echo "  1. Acesse: https://dash.cloudflare.com/"
    echo "  2. Vá em 'Zero Trust' > 'Access' > 'Tunnels'"
    echo "  3. Crie um novo tunnel ou use um existente"
    echo "  4. Copie o token fornecido"
    echo ""
    echo "💡 Exemplo:"
    echo "  ./setup_cloudflare.sh eyJhIjoiY2xvdWRmbGFyZS10dW5uZWwiLCJ0Ijoi... "
    exit 1
fi

TOKEN=$1

# Verifica se o arquivo .env existe
if [ ! -f .env ]; then
    echo "❌ Arquivo .env não encontrado!"
    echo "   Execute este script no diretório do projeto BFF Luma"
    exit 1
fi

# Substitui o token no arquivo .env
if grep -q "CLOUDFLARE_TUNNEL_TOKEN=" .env; then
    # Se já existe, substitui
    sed -i "s/CLOUDFLARE_TUNNEL_TOKEN=.*/CLOUDFLARE_TUNNEL_TOKEN=$TOKEN/" .env
    echo "✅ Token atualizado no arquivo .env"
else
    # Se não existe, adiciona
    echo "" >> .env
    echo "# Cloudflare Tunnel Configuration" >> .env
    echo "CLOUDFLARE_TUNNEL_TOKEN=$TOKEN" >> .env
    echo "✅ Token adicionado ao arquivo .env"
fi

echo ""
echo "🎯 Configuração concluída!"
echo "📋 Para usar o cloudflared tunnel:"
echo "   docker compose --profile tunnel up -d"
echo ""
echo "📋 Para verificar se está funcionando:"
echo "   docker compose logs cloudflared"
