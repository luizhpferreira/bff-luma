# BFF Luma - API Backend for Frontend

Uma API BFF (Backend for Frontend) em Go que integra com LNBits para gerenciar carteiras Lightning de forma segura.

## 🚀 Características

- **Segurança**: As chaves das carteiras nunca são expostas para o frontend
- **Integração LNBits**: Criação automática de carteiras via Admin Key
- **SQLite**: Armazenamento local de mapeamento wallet_id ↔ app_user_id
- **API RESTful**: Endpoints padronizados para operações de carteira
- **Graceful Shutdown**: Encerramento seguro do servidor

## 📋 Pré-requisitos

- Go 1.21+
- LNBits rodando localmente ou remotamente
- Admin Key do LNBits configurada

## 🛠️ Instalação

1. Clone o repositório:
```bash
git clone <repository-url>
cd bff_luma
```

2. Configure as variáveis de ambiente no arquivo `.env`:
```env
# Configurações da aplicação
APP_PORT=8080

# Configurações JWT
JWT_SECRET=supersecreto123456789

# Configurações LNbits
LNBITS_BASE_URL=http://127.0.0.1:5000
LNBITS_ADMIN_KEY=sua_admin_key_aqui
LNBITS_WEBHOOK_SECRET=seu_webhook_secret_aqui
```

3. Execute as dependências:
```bash
go mod tidy
```

4. Execute o servidor:
```bash
go run cmd/server/main.go
```

## 📚 API Endpoints

### Health Check
```
GET /health
```
Verifica se a API está funcionando.

**Resposta:**
```json
{
  "success": true,
  "message": "API funcionando",
  "data": {
    "status": "ok",
    "service": "BFF Luma API",
    "version": "1.0.0"
  }
}
```

### Criar Carteira (Cadastro)
```
POST /api/v1/wallets
```

**Request Body:**
```json
{
  "email": "usuario@exemplo.com",
  "password": "MinhaSenha@123",
  "password_repeat": "MinhaSenha@123"
}
```

**Requisitos da Senha:**
- Mínimo 8 caracteres
- Pelo menos uma letra maiúscula
- Pelo menos uma letra minúscula
- Pelo menos um número
- Pelo menos um caractere especial
- Não pode conter sequências comuns (123, abc, qwe, etc.)
- Não pode ter mais de 2 caracteres iguais consecutivos

**Resposta:**
```json
{
  "success": true,
  "message": "Carteira criada com sucesso",
  "data": {
    "wallet_id": "abc123def456",
    "email": "usuario@exemplo.com",
    "message": "Carteira criada com sucesso"
  }
}
```

### Login
```
POST /api/v1/login
```

**Request Body:**
```json
{
  "email": "usuario@exemplo.com",
  "password": "MinhaSenha@123"
}
```

**Resposta:**
```json
{
  "success": true,
  "message": "Login realizado com sucesso",
  "data": {
    "wallet_id": "abc123def456",
    "email": "usuario@exemplo.com",
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "message": "Login realizado com sucesso"
  }
}
```

### Refresh Token
```
POST /api/v1/refresh
```

**Headers:**
```
Authorization: Bearer <token>
```

**Resposta:**
```json
{
  "success": true,
  "message": "Token renovado com sucesso",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "message": "Token renovado com sucesso"
  }
}
```

### Recuperação de Senha
```
POST /api/v1/forgot-password
```

**Request Body:**
```json
{
  "email": "usuario@exemplo.com"
}
```

**Resposta:**
```json
{
  "success": true,
  "message": "Solicitação processada",
  "data": {
    "message": "Se o email existir em nossa base, você receberá um link de recuperação"
  }
}
```

### Reset de Senha
```
POST /api/v1/reset-password
```

**Request Body:**
```json
{
  "token": "token-gerado-no-email",
  "new_password": "NovaSenha@2024!",
  "new_password_repeat": "NovaSenha@2024!"
}
```

**Resposta:**
```json
{
  "success": true,
  "message": "Senha redefinida com sucesso",
  "data": {
    "message": "Senha redefinida com sucesso"
  }
}
```

### Limpeza de Tokens (Administração)
```
POST /api/v1/admin/cleanup
```

**Resposta:**
```json
{
  "success": true,
  "message": "Limpeza executada com sucesso",
  "data": {
    "expired_tokens": 0,
    "interval": "1h0m0s",
    "is_running": true
  }
}
```

### Estatísticas de Limpeza (Administração)
```
GET /api/v1/admin/cleanup/stats
```

**Resposta:**
```json
{
  "success": true,
  "message": "Estatísticas obtidas com sucesso",
  "data": {
    "expired_tokens": 0,
    "interval": "1h0m0s",
    "is_running": true
  }
}
```

### Obter Informações da Carteira
```
GET /api/v1/wallets
```

**Headers:**
```
Authorization: Bearer <token>
```

**Resposta:**
```json
{
  "success": true,
  "message": "Informações da carteira",
  "data": {
    "id": 1,
    "email": "usuario@exemplo.com",
    "wallet_id": "abc123def456",
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  }
}
```

### Criar Invoice
```
POST /api/v1/invoices
```

**Headers:**
```
Authorization: Bearer <token>
```

**Request Body:**
```json
{
  "amount": 1000,
  "memo": "Pagamento teste"
}
```

**Resposta:**
```json
{
  "success": true,
  "message": "Invoice criado com sucesso",
  "data": {
    "payment_request": "lnbc10u1p3qkqkqpp5...",
    "payment_hash": "abc123def456...",
    "amount": 1000,
    "memo": "Pagamento teste",
    "expires_at": 1642234567
  }
}
```

### Verificar Status do Pagamento
```
GET /api/v1/payments/status?payment_hash=abc123def456
```

**Headers:**
```
Authorization: Bearer <token>
```

**Resposta:**
```json
{
  "success": true,
  "message": "Status do pagamento verificado",
  "data": {
    "payment_hash": "abc123def456",
    "paid": true,
    "amount": 1000,
    "memo": "Pagamento teste",
    "email": "usuario@exemplo.com",
    "paid_at": 1642234567
  }
}
```

## 🔒 Segurança

- **Chaves Protegidas**: As chaves Admin e Invoice das carteiras são armazenadas apenas no backend
- **Senhas**: As senhas são armazenadas com hash bcrypt (custo 12) no banco de dados
- **Validação**: Todos os inputs são validados antes do processamento
- **Autenticação**: Sistema de login com email e senha
- **JWT**: Tokens JWT para sessões seguras (expira em 24 horas)
- **Middleware**: Autenticação obrigatória para endpoints protegidos
- **Recuperação de Senha**: Sistema seguro de reset com tokens únicos (expira em 1 hora)
- **Limpeza Automática**: Remoção automática de tokens expirados a cada hora
- **Logs**: Operações importantes são logadas para auditoria
- **CORS**: Configurado para permitir requisições cross-origin

## 🔐 Requisitos de Senha Forte

Para garantir a segurança das contas, todas as senhas devem atender aos seguintes requisitos:

### 📋 Regras Obrigatórias:
- **Mínimo 8 caracteres**
- **Pelo menos 1 letra maiúscula** (A-Z)
- **Pelo menos 1 letra minúscula** (a-z)
- **Pelo menos 1 número** (0-9)
- **Pelo menos 1 caractere especial** (!@#$%^&*()_+-=[]{}|;:,.<>?)

### 🚫 Restrições:
- **Não pode conter sequências comuns** como "123", "abc", "password", "senha"
- **Não pode ter mais de 2 caracteres idênticos consecutivos** (ex: "aaa", "111")

### 🔒 Armazenamento Seguro:
- **Hash bcrypt** com custo 12 (padrão de segurança)
- **Salt único** para cada senha
- **Impossível reverter** o hash para obter a senha original

### ✅ Exemplos de Senhas Válidas:
- `B@nco2024!`
- `MinhaSenha@123`
- `Secure#Pass1`
- `P@ssw0rd!`

### ❌ Exemplos de Senhas Inválidas:
- `12345678` (só números, sem maiúsculas/minúsculas/especiais)
- `password` (sequência comum, sem números/especiais)
- `Senha123` (sem caracteres especiais)
- `B@nco2024` (sem caracteres especiais)
- `aaa@Bc123` (3 'a' consecutivos)

## 🗄️ Banco de Dados

O sistema utiliza SQLite para armazenar o mapeamento entre `app_user_id` e `wallet_id`. A estrutura da tabela:

```sql
CREATE TABLE wallets (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    wallet_id TEXT NOT NULL UNIQUE,
    admin_key TEXT NOT NULL,
    invoice_key TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

## 🔄 Fluxo de Operação

1. **Cadastro de Usuário**:
   - Frontend envia `email` + `password` + `password_repeat`
   - Backend valida se as senhas coincidem
   - Backend chama LNBits com Admin Key
   - LNBits retorna `wallet_id`, `admin_key`, `invoice_key`
   - Backend salva mapeamento no SQLite
   - Frontend recebe apenas `wallet_id`

2. **Login**:
   - Frontend envia `email` + `password`
   - Backend verifica credenciais no banco
   - Backend gera token JWT e retorna `wallet_id`, `token` e confirmação

3. **Criação de Invoice**:
   - Frontend envia `Authorization: Bearer <token>` + `amount`
   - Backend valida token JWT e extrai email
   - Backend busca `invoice_key` no banco
   - Backend chama LNBits com `invoice_key`
   - Frontend recebe `payment_request` (BOLT11)

4. **Verificação de Pagamento**:
   - Frontend envia `Authorization: Bearer <token>` + `payment_hash`
   - Backend busca `invoice_key` e verifica no LNBits
   - Frontend recebe status do pagamento

## 🚀 Executando em Produção

1. Configure as variáveis de ambiente adequadas
2. Use um proxy reverso (nginx, traefik)
3. Configure SSL/TLS
4. Monitore logs e métricas
5. Configure backup do banco SQLite

## 🧪 Testando

### Criar uma carteira:
```bash
curl -X POST http://localhost:8080/api/v1/wallets \
  -H "Content-Type: application/json" \
  -d '{"app_user_id": "test_user_123"}'
```

### Criar um invoice:
```bash
curl -X POST http://localhost:8080/api/v1/invoices \
  -H "Content-Type: application/json" \
  -d '{"app_user_id": "test_user_123", "amount": 1000, "memo": "Teste"}'
```

### Verificar status:
```bash
curl "http://localhost:8080/api/v1/payments/status?app_user_id=test_user_123&payment_hash=abc123"
```

## 📝 Logs

O sistema gera logs para:
- Criação de carteiras
- Criação de invoices
- Erros de operação
- Inicialização e shutdown do servidor

## 🤝 Contribuindo

1. Fork o projeto
2. Crie uma branch para sua feature
3. Commit suas mudanças
4. Push para a branch
5. Abra um Pull Request

## 📄 Licença

Este projeto está sob a licença MIT.
