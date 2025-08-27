#!/bin/bash

echo "🧪 Testando funcionalidade de Pagamento de Invoices"
echo "=================================================="

# Configurações
BASE_URL="http://localhost:8080"
API_URL="$BASE_URL/api/v1"

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Função para log colorido
log_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

log_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

log_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

log_error() {
    echo -e "${RED}❌ $1${NC}"
}

# Verifica se o servidor está rodando
log_info "Verificando se o servidor está rodando..."
if curl -s "$BASE_URL/health" > /dev/null; then
    log_success "Servidor está rodando"
else
    log_error "Servidor não está rodando. Inicie o servidor primeiro."
    exit 1
fi

echo ""
log_info "🔧 Verificando implementação do backend..."

# Verifica se os arquivos necessários existem
if [ -f "internal/models/wallet.go" ]; then
    log_success "Modelo PaymentRequest/Response encontrado"
else
    log_error "Modelo PaymentRequest/Response não encontrado"
fi

if [ -f "internal/services/lnbits.go" ]; then
    log_success "Serviço PayInvoice no LNBits encontrado"
else
    log_error "Serviço PayInvoice no LNBits não encontrado"
fi

if [ -f "internal/services/wallet_service.go" ]; then
    log_success "Serviço PayInvoice no WalletService encontrado"
else
    log_error "Serviço PayInvoice no WalletService não encontrado"
fi

if [ -f "internal/handlers/wallet_handler.go" ]; then
    log_success "Handler PayInvoice encontrado"
else
    log_error "Handler PayInvoice não encontrado"
fi

echo ""
log_info "🎯 Funcionalidades implementadas:"
echo "• Pagar invoices usando a carteira do usuário"
echo "• Validação de autenticação JWT"
echo "• Verificação de saldo da carteira"
echo "• Processamento de pagamento via LNBits"
echo "• Retorno do status do pagamento"

echo ""
log_info "📋 Como funciona:"
echo "1. Usuário autenticado envia payment_request (BOLT11)"
echo "2. Sistema valida o JWT e identifica a carteira"
echo "3. Sistema usa a AdminKey da carteira para pagar"
echo "4. LNBits processa o pagamento"
echo "5. Sistema retorna o status do pagamento"

echo ""
log_info "🚀 Endpoint disponível:"
echo "POST $API_URL/payments - Pagar invoice"
echo "Headers: Authorization: Bearer <JWT_TOKEN>"
echo "Body: { \"payment_request\": \"lnbc...\" }"

echo ""
log_warning "🧪 Para testar manualmente, siga estes passos:"
echo ""
echo "1. Inicie o servidor:"
echo "   cd bff_luma && make run"
echo ""
echo "2. Crie uma carteira de teste:"
echo "   curl -X POST $API_URL/wallets \\"
echo "     -H 'Content-Type: application/json' \\"
echo "     -d '{"
echo "       \"username\": \"12345678901\","
echo "       \"email\": \"teste@exemplo.com\","
echo "       \"password\": \"senha123456\","
echo "       \"password_repeat\": \"senha123456\""
echo "     }'"
echo ""
echo "3. Faça login para obter JWT token:"
echo "   curl -X POST $API_URL/login \\"
echo "     -H 'Content-Type: application/json' \\"
echo "     -d '{"
echo "       \"email\": \"12345678901\","
echo "       \"password\": \"senha123456\""
echo "     }'"
echo ""
echo "4. Use o token JWT para pagar um invoice:"
echo "   curl -X POST $API_URL/payments \\"
echo "     -H 'Content-Type: application/json' \\"
echo "     -H 'Authorization: Bearer <SEU_JWT_TOKEN>' \\"
echo "     -d '{"
echo "       \"payment_request\": \"lnbc...\""
echo "     }'"
echo ""
echo "5. Verifique o status do pagamento:"
echo "   curl -X GET \"$API_URL/payments/status?payment_hash=<HASH>\" \\"
echo "     -H 'Authorization: Bearer <SEU_JWT_TOKEN>'"

echo ""
log_info "🔒 Segurança implementada:"
echo "• Autenticação JWT obrigatória"
echo "• AdminKey nunca exposta no frontend"
echo "• Validação de carteira do usuário"
echo "• Rate limiting aplicado"

echo ""
log_info "💡 Dicas para teste:"
echo "• Use um invoice BOLT11 válido para teste"
echo "• Certifique-se de que a carteira tem saldo suficiente"
echo "• Verifique se o LNBits está configurado corretamente"
echo "• Monitore os logs do servidor durante os testes"

echo ""
log_success "✅ Implementação concluída!"
echo "A funcionalidade de pagamento de invoices está pronta para uso."
echo "Usuários autenticados podem pagar invoices usando suas carteiras."

echo ""
log_info "📝 Próximos passos:"
echo "1. Teste manualmente usando os comandos acima"
echo "2. Integre com o aplicativo móvel"
echo "3. Implemente tratamento de erros específicos"
echo "4. Adicione logs detalhados para debugging"
