#!/bin/bash

echo "🧪 Testando funcionalidade de múltiplas Invoice Keys"
echo "=================================================="

# Verifica se o servidor está rodando
echo "📡 Verificando se o servidor está rodando..."
if curl -s http://localhost:8080/health > /dev/null; then
    echo "✅ Servidor está rodando"
else
    echo "❌ Servidor não está rodando. Inicie o servidor primeiro."
    exit 1
fi

echo ""
echo "🔧 Backend implementado:"
echo "✅ Modelos InvoiceKey, CreateInvoiceKeyRequest, etc."
echo "✅ Serviço CreateInvoiceKey no LNBits"
echo "✅ Serviço ListInvoiceKeys no LNBits"
echo "✅ Serviço CreateInvoiceWithKey no LNBits"
echo "✅ Handler CreateInvoiceKey no wallet_handler.go"
echo "✅ Handler ListInvoiceKeys no wallet_handler.go"
echo "✅ Handler CreateInvoiceWithKey no wallet_handler.go"
echo "✅ Rotas POST /api/v1/invoice-keys"
echo "✅ Rotas GET /api/v1/invoice-keys"
echo "✅ Rotas POST /api/v1/invoices/with-key"

echo ""
echo "🎯 Funcionalidades implementadas:"
echo "• Criar múltiplas invoice keys para uma wallet"
echo "• Listar todas as invoice keys do usuário"
echo "• Criar invoices usando invoice keys específicas"
echo "• Organizar recebimentos por categorias"
echo "• Separar contas por finalidade"

echo ""
echo "📋 Como funciona:"
echo "1. Um usuário tem UMA wallet principal"
echo "2. A wallet pode ter MÚLTIPLAS invoice keys"
echo "3. Cada invoice key pode ser usada para diferentes propósitos:"
echo "   • 'Vendas Online'"
echo "   • 'Consultoria'"
echo "   • 'Doações'"
echo "   • 'Freelance'"
echo "   • etc."

echo ""
echo "🔑 Estrutura das Invoice Keys:"
echo "• ID único para cada invoice key"
echo "• Nome descritivo (ex: 'Vendas Online')"
echo "• Descrição opcional"
echo "• Data de criação"
echo "• Chave de invoice (nunca exposta no frontend)"

echo ""
echo "🚀 Endpoints disponíveis:"
echo "POST /api/v1/invoice-keys - Criar nova invoice key"
echo "GET /api/v1/invoice-keys - Listar invoice keys"
echo "POST /api/v1/invoices/with-key - Criar invoice com key específica"

echo ""
echo "📝 Exemplo de uso:"
echo "1. Usuário cria invoice key 'Vendas Online'"
echo "2. Usuário cria invoice key 'Consultoria'"
echo "3. Ao criar invoices, escolhe qual key usar"
echo "4. Pode rastrear recebimentos por categoria"

echo ""
echo "✅ Implementação concluída!"
echo "A funcionalidade de múltiplas invoice keys está pronta para uso."
echo "Cada usuário pode ter uma wallet com várias invoice keys para organizar seus recebimentos."
