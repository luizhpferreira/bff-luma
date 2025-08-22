# Exemplos de cURL para BFF Luma API

## ✅ Endpoints Funcionando

### 1. Criar Carteira (Cadastro)
```bash
curl -X POST http://localhost:8080/api/v1/wallets \
  -H "Content-Type: application/json" \
  -d '{
    "email": "usuario@exemplo.com",
    "password": "minhasenha123",
    "password_repeat": "minhasenha123"
  }'
```

**Resposta esperada:**
```json
{
  "success": true,
  "message": "Carteira criada com sucesso",
  "data": {
    "wallet_id": "4aef1b41b711488583cc4f5b500b02fc",
    "email": "usuario@exemplo.com",
    "message": "Carteira criada com sucesso"
  }
}
```

### 2. Login
```bash
curl -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "usuario@exemplo.com",
    "password": "minhasenha123"
  }'
```

**Resposta esperada:**
```json
{
  "success": true,
  "message": "Login realizado com sucesso",
  "data": {
    "wallet_id": "4aef1b41b711488583cc4f5b500b02fc",
    "email": "usuario@exemplo.com",
    "message": "Login realizado com sucesso"
  }
}
```

### 3. Obter Informações da Carteira
```bash
curl "http://localhost:8080/api/v1/wallets?email=usuario@exemplo.com"
```

**Resposta esperada:**
```json
{
  "success": true,
  "message": "Informações da carteira",
  "data": {
    "id": 1,
    "email": "usuario@exemplo.com",
    "wallet_id": "4aef1b41b711488583cc4f5b500b02fc",
    "created_at": "2025-08-21T22:47:51.04792317-03:00",
    "updated_at": "2025-08-21T22:47:51.04792317-03:00"
  }
}
```

## ⚠️ Endpoints com Problemas de Conexão LNBits

### 4. Criar Invoice
```bash
curl -X POST http://localhost:8080/api/v1/invoices \
  -H "Content-Type: application/json" \
  -d '{
    "email": "usuario@exemplo.com",
    "amount": 1000,
    "memo": "Teste de invoice"
  }'
```

**Erro esperado (problema de conexão LNBits):**
```json
{
  "success": false,
  "error": "erro ao criar invoice no LNBits: erro na resposta do LNBits: 520 - {\"detail\":\"[Errno 111] Connection refused\",\"status\":\"pending\"}",
  "message": "Erro ao criar invoice"
}
```

### 5. Verificar Status do Pagamento
```bash
curl "http://localhost:8080/api/v1/payments/status?email=usuario@exemplo.com&payment_hash=abc123"
```

**Erro esperado (problema de conexão LNBits):**
```json
{
  "success": false,
  "error": "erro ao verificar status do pagamento: erro na resposta do LNBits: 404 - {\"detail\":\"Payment does not exist.\"}",
  "message": "Erro ao verificar status do pagamento"
}
```

## 🔍 Testes de Validação

### 6. Teste de Senhas que Não Coincidem
```bash
curl -X POST http://localhost:8080/api/v1/wallets \
  -H "Content-Type: application/json" \
  -d '{
    "email": "teste2@exemplo.com",
    "password": "123456",
    "password_repeat": "654321"
  }'
```

**Resposta esperada:**
```json
{
  "success": false,
  "error": "as senhas não coincidem",
  "message": "Erro ao criar carteira"
}
```

### 7. Teste de Login com Senha Incorreta
```bash
curl -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "usuario@exemplo.com",
    "password": "senhaerrada"
  }'
```

**Resposta esperada:**
```json
{
  "success": false,
  "error": "senha incorreta",
  "message": "Erro no login"
}
```

## 🧪 Script de Teste Completo

```bash
#!/bin/bash

echo "🧪 Testando BFF Luma API"
echo "=========================="

# Gerar email único
EMAIL="teste_$(date +%s)@exemplo.com"
PASSWORD="123456"

echo "📧 Email de teste: $EMAIL"
echo ""

# 1. Criar carteira
echo "1️⃣ Criando carteira..."
RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/wallets \
  -H "Content-Type: application/json" \
  -d "{
    \"email\": \"$EMAIL\",
    \"password\": \"$PASSWORD\",
    \"password_repeat\": \"$PASSWORD\"
  }")

echo "$RESPONSE" | jq '.'
echo ""

# 2. Login
echo "2️⃣ Fazendo login..."
RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d "{
    \"email\": \"$EMAIL\",
    \"password\": \"$PASSWORD\"
  }")

echo "$RESPONSE" | jq '.'
echo ""

# 3. Obter informações da carteira
echo "3️⃣ Obtendo informações da carteira..."
RESPONSE=$(curl -s "http://localhost:8080/api/v1/wallets?email=$EMAIL")
echo "$RESPONSE" | jq '.'
echo ""

# 4. Tentar criar invoice (pode falhar por conexão LNBits)
echo "4️⃣ Tentando criar invoice..."
RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/invoices \
  -H "Content-Type: application/json" \
  -d "{
    \"email\": \"$EMAIL\",
    \"amount\": 1000,
    \"memo\": \"Teste de invoice\"
  }")

echo "$RESPONSE" | jq '.'
echo ""

echo "✅ Testes concluídos!"
```

## 📝 Notas

- **Cadastro**: Use `email`, `password` e `password_repeat` para criar uma nova carteira
- **Login**: Use `email` e `password` para autenticar
- **Validação**: O sistema verifica se as senhas coincidem no cadastro
- **Segurança**: As senhas são armazenadas no banco (em produção, deve ser hash)
- **LNBits**: Problemas de conexão são esperados nos endpoints de invoice e pagamento
