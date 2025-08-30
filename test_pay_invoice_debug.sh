#!/bin/bash

echo "🧪 Teste de Debug - Pagamento de Invoices"
echo "========================================="

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
TEST_EMAIL="luizferreiralps@gmail.com"
TEST_CPF="01383972281"  # CPF da carteira existente
TEST_PASSWORD="#Ruiter1"  # Senha da carteira existente
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
log_info "🧪 Iniciando testes de debug..."

# Teste 1: Fazer login com carteira existente
echo ""
log_info "Teste 1: Fazendo login com carteira existente..."
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

# Teste 2: Verificar informações da carteira
echo ""
log_info "Teste 2: Verificando informações da carteira..."
response=$(make_request "GET" "/wallets" "" "Authorization: Bearer $JWT_TOKEN")
http_code=$?

if [ $http_code -eq 200 ]; then
    log_success "Informações da carteira obtidas"
    echo "Resposta: $response"
else
    log_error "Erro ao obter informações da carteira (HTTP $http_code)"
    echo "Resposta: $response"
fi

# Teste 3: Criar um invoice para testar o pagamento
echo ""
log_info "Teste 3: Criando invoice para teste..."
CREATE_INVOICE_DATA="{
    \"amount\": 1000,
    \"memo\": \"Teste de pagamento\"
}"

response=$(make_request "POST" "/invoices" "$CREATE_INVOICE_DATA" "Authorization: Bearer $JWT_TOKEN")
http_code=$?

if [ $http_code -eq 201 ]; then
    log_success "Invoice criado com sucesso"
    PAYMENT_REQUEST=$(echo "$response" | jq -r '.data.payment_request' 2>/dev/null)
    PAYMENT_HASH=$(echo "$response" | jq -r '.data.payment_hash' 2>/dev/null)
    echo "Payment Request: $PAYMENT_REQUEST"
    echo "Payment Hash: $PAYMENT_HASH"
else
    log_error "Erro ao criar invoice (HTTP $http_code)"
    echo "Resposta: $response"
    exit 1
fi

# Teste 4: Tentar pagar o invoice criado
echo ""
log_info "Teste 4: Tentando pagar o invoice criado..."
PAYMENT_DATA="{
    \"payment_request\": \"$PAYMENT_REQUEST\"
}"

response=$(make_request "POST" "/payments" "$PAYMENT_DATA" "Authorization: Bearer $JWT_TOKEN")
http_code=$?

if [ $http_code -eq 200 ]; then
    log_success "Pagamento realizado com sucesso!"
    echo "Resposta: $response"
else
    log_error "Erro ao pagar invoice (HTTP $http_code)"
    echo "Resposta: $response"
fi

# Teste 5: Verificar status do pagamento
echo ""
log_info "Teste 5: Verificando status do pagamento..."
response=$(make_request "GET" "/payments/status?payment_hash=$PAYMENT_HASH" "" "Authorization: Bearer $JWT_TOKEN")
http_code=$?

if [ $http_code -eq 200 ]; then
    log_success "Status do pagamento verificado"
    echo "Resposta: $response"
else
    log_warning "Erro ao verificar status (HTTP $http_code)"
    echo "Resposta: $response"
fi

echo ""
log_success "✅ Testes de debug concluídos!"

echo ""
log_info "📋 Resumo dos testes:"
echo "✅ Login com carteira existente"
echo "✅ Verificação de informações da carteira"
echo "✅ Criação de invoice para teste"
echo "✅ Tentativa de pagamento do invoice"
echo "✅ Verificação de status do pagamento"

echo ""
log_info "💡 Dados de teste:"
echo "CPF: $TEST_CPF"
echo "Email: $TEST_EMAIL"
echo "JWT Token: $JWT_TOKEN"
echo "Payment Request: $PAYMENT_REQUEST"
echo "Payment Hash: $PAYMENT_HASH"

echo ""
log_success "�� Debug concluído!"
