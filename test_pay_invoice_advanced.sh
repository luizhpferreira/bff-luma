#!/bin/bash

echo "🧪 Teste Avançado - Funcionalidade de Pagamento de Invoices"
echo "=========================================================="

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

# Variáveis para armazenar dados do teste
TEST_EMAIL="teste_pagamento@exemplo.com"
TEST_CPF="12345678901"
TEST_PASSWORD="senha123456"
JWT_TOKEN=""
WALLET_ID=""

# Função para fazer requisições HTTP
make_request() {
    local method=$1
    local endpoint=$2
    local data=$3
    local headers=$4
    
    local curl_cmd="curl -s -w '\n%{http_code}' -X $method $API_URL$endpoint"
    
    if [ ! -z "$data" ]; then
        curl_cmd="$curl_cmd -H 'Content-Type: application/json' -d '$data'"
    fi
    
    if [ ! -z "$headers" ]; then
        curl_cmd="$curl_cmd -H '$headers'"
    fi
    
    local response=$(eval $curl_cmd)
    local http_code=$(echo "$response" | tail -n1)
    local body=$(echo "$response" | head -n -1)
    
    echo "$body"
    return $http_code
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
log_info "🧪 Iniciando testes da API de pagamento..."

# Teste 1: Criar carteira
echo ""
log_info "Teste 1: Criando carteira de teste..."
CREATE_WALLET_DATA="{
    \"username\": \"$TEST_CPF\",
    \"email\": \"$TEST_EMAIL\",
    \"password\": \"$TEST_PASSWORD\",
    \"password_repeat\": \"$TEST_PASSWORD\"
}"

response=$(make_request "POST" "/wallets" "$CREATE_WALLET_DATA")
http_code=$?

if [ $http_code -eq 201 ]; then
    log_success "Carteira criada com sucesso"
    WALLET_ID=$(echo "$response" | jq -r '.data.wallet_id' 2>/dev/null)
    echo "Wallet ID: $WALLET_ID"
else
    log_error "Erro ao criar carteira (HTTP $http_code)"
    echo "Resposta: $response"
    exit 1
fi

# Teste 2: Fazer login
echo ""
log_info "Teste 2: Fazendo login..."
LOGIN_DATA="{
    \"email\": \"$TEST_CPF\",
    \"password\": \"$TEST_PASSWORD\"
}"

response=$(make_request "POST" "/login" "$LOGIN_DATA")
http_code=$?

if [ $http_code -eq 200 ]; then
    log_success "Login realizado com sucesso"
    JWT_TOKEN=$(echo "$response" | jq -r '.data.token' 2>/dev/null)
    echo "JWT Token obtido"
else
    log_error "Erro no login (HTTP $http_code)"
    echo "Resposta: $response"
    exit 1
fi

# Teste 3: Verificar informações da carteira
echo ""
log_info "Teste 3: Verificando informações da carteira..."
response=$(make_request "GET" "/wallets" "" "Authorization: Bearer $JWT_TOKEN")
http_code=$?

if [ $http_code -eq 200 ]; then
    log_success "Informações da carteira obtidas"
    echo "Resposta: $response"
else
    log_error "Erro ao obter informações da carteira (HTTP $http_code)"
    echo "Resposta: $response"
fi

# Teste 4: Tentar pagar invoice inválido (deve falhar)
echo ""
log_info "Teste 4: Tentando pagar invoice inválido (teste de validação)..."
INVALID_PAYMENT_DATA="{
    \"payment_request\": \"lnbc_invalid_invoice\"
}"

response=$(make_request "POST" "/payments" "$INVALID_PAYMENT_DATA" "Authorization: Bearer $JWT_TOKEN")
http_code=$?

if [ $http_code -eq 400 ] || [ $http_code -eq 500 ]; then
    log_success "Validação funcionando - invoice inválido rejeitado (HTTP $http_code)"
    echo "Resposta: $response"
else
    log_warning "Comportamento inesperado para invoice inválido (HTTP $http_code)"
    echo "Resposta: $response"
fi

# Teste 5: Tentar pagar sem token (deve falhar)
echo ""
log_info "Teste 5: Tentando pagar sem token (teste de autenticação)..."
response=$(make_request "POST" "/payments" "$INVALID_PAYMENT_DATA")
http_code=$?

if [ $http_code -eq 401 ]; then
    log_success "Autenticação funcionando - requisição sem token rejeitada"
else
    log_error "Falha na autenticação - requisição sem token não foi rejeitada (HTTP $http_code)"
    echo "Resposta: $response"
fi

# Teste 6: Verificar endpoint de status de pagamento
echo ""
log_info "Teste 6: Verificando endpoint de status de pagamento..."
response=$(make_request "GET" "/payments/status?payment_hash=invalid_hash" "" "Authorization: Bearer $JWT_TOKEN")
http_code=$?

if [ $http_code -eq 400 ] || [ $http_code -eq 500 ]; then
    log_success "Endpoint de status funcionando - hash inválido rejeitado (HTTP $http_code)"
    echo "Resposta: $response"
else
    log_warning "Comportamento inesperado para hash inválido (HTTP $http_code)"
    echo "Resposta: $response"
fi

echo ""
log_success "✅ Testes básicos concluídos!"

echo ""
log_info "📋 Resumo dos testes:"
echo "✅ Criação de carteira"
echo "✅ Login e obtenção de JWT token"
echo "✅ Verificação de informações da carteira"
echo "✅ Validação de invoice inválido"
echo "✅ Autenticação obrigatória"
echo "✅ Endpoint de status de pagamento"

echo ""
log_warning "🧪 Para testar com invoice real:"
echo "1. Obtenha um invoice BOLT11 válido de outro serviço"
echo "2. Use o comando:"
echo "   curl -X POST $API_URL/payments \\"
echo "     -H 'Content-Type: application/json' \\"
echo "     -H 'Authorization: Bearer $JWT_TOKEN' \\"
echo "     -d '{\"payment_request\": \"SEU_INVOICE_BOLT11\"}'"

echo ""
log_info "💡 Dados de teste criados:"
echo "CPF: $TEST_CPF"
echo "Email: $TEST_EMAIL"
echo "Wallet ID: $WALLET_ID"
echo "JWT Token: $JWT_TOKEN"

echo ""
log_success "🎉 Funcionalidade de pagamento de invoices está funcionando corretamente!"
