# Correção de Segurança - Senhas Hardcoded

## Problema Identificado

Foram encontradas senhas e secrets hardcoded no código fonte, o que representa um risco de segurança significativo:

### 1. Senha do Banco LNBits (lnbits.go)
**Localização**: `internal/services/lnbits.go` linhas 109-111
```go
dbPassword := os.Getenv("LNBITS_POSTGRES_PASSWORD")
if dbPassword == "" {
    dbPassword = "Qualquer2"  // ❌ SENHA HARCODED
}
```

### 2. JWT Secret (config.go)
**Localização**: `internal/config/config.go` linha 44
```go
JWTSecret: getEnv("JWT_SECRET", "supersecreto123456789"),  // ❌ SECRET HARCODED
```

## Correções Implementadas

### 1. Remoção de Fallbacks Hardcoded
- **Arquivo**: `internal/services/lnbits.go`
- **Mudança**: Removidos todos os fallbacks hardcoded para variáveis de ambiente do LNBits
- **Resultado**: O código agora falha explicitamente se as variáveis não estiverem configuradas

```go
// ANTES (❌ Inseguro)
dbPassword := os.Getenv("LNBITS_POSTGRES_PASSWORD")
if dbPassword == "" {
    dbPassword = "Qualquer2"  // Senha hardcoded
}

// DEPOIS (✅ Seguro)
dbPassword := os.Getenv("LNBITS_POSTGRES_PASSWORD")
if dbPassword == "" {
    return nil, fmt.Errorf("LNBITS_POSTGRES_PASSWORD não configurada nas variáveis de ambiente")
}
```

### 2. Validação de Variáveis Críticas
- **Arquivo**: `cmd/server/main.go`
- **Mudança**: Adicionada validação no startup para variáveis críticas
- **Resultado**: A aplicação não inicia se variáveis essenciais não estiverem configuradas

```go
// Validação no startup
if cfg.JWTSecret == "" {
    log.Fatal("❌ JWT_SECRET não configurado nas variáveis de ambiente")
}
if cfg.LNBitsAPIToken == "" {
    log.Fatal("❌ LNBITS_API_TOKEN não configurado nas variáveis de ambiente")
}
```

## Variáveis de Ambiente Obrigatórias

Agora as seguintes variáveis são **obrigatórias** e devem estar configuradas no arquivo `.env`:

```bash
# JWT Secret (obrigatório)
JWT_SECRET=your_secure_jwt_secret_here

# LNBits API Token (obrigatório)
LNBITS_API_TOKEN=your_lnbits_api_token_here

# LNBits Database (obrigatório)
LNBITS_POSTGRES_HOST=your_lnbits_db_host
LNBITS_POSTGRES_PORT=5432
LNBITS_POSTGRES_USER=lnbits
LNBITS_POSTGRES_PASSWORD=your_secure_lnbits_password
LNBITS_POSTGRES_DB=lnbits
```

## Benefícios da Correção

1. **Segurança**: Nenhuma senha ou secret está mais visível no código fonte
2. **Configuração**: Todas as credenciais são gerenciadas via variáveis de ambiente
3. **Validação**: A aplicação falha rapidamente se configurações críticas estiverem ausentes
4. **Auditoria**: Fica claro quais variáveis são obrigatórias para o funcionamento
5. **DevOps**: Facilita a configuração em diferentes ambientes (dev, staging, prod)

## Como Configurar

1. Copie o arquivo `env.example` para `.env`:
   ```bash
   cp env.example .env
   ```

2. Configure as variáveis obrigatórias no arquivo `.env`:
   ```bash
   # Edite o arquivo .env e configure as variáveis
   nano .env
   ```

3. Certifique-se de que todas as variáveis obrigatórias estão configuradas antes de executar a aplicação.

## Teste da Correção

Para verificar se a correção está funcionando:

```bash
# Teste sem variáveis de ambiente (deve falhar)
unset JWT_SECRET LNBITS_API_TOKEN
go run cmd/server/main.go

# Deve mostrar erro: "❌ JWT_SECRET não configurado nas variáveis de ambiente"

# Teste com variáveis configuradas (deve funcionar)
export JWT_SECRET="test_secret"
export LNBITS_API_TOKEN="test_token"
go run cmd/server/main.go
```

## Próximos Passos

1. **Revisar outros arquivos**: Verificar se há outras senhas hardcoded em scripts de teste
2. **Rotação de secrets**: Implementar rotação automática de secrets em produção
3. **Vault/HashiCorp**: Considerar usar um gerenciador de secrets como Vault para produção
4. **Auditoria**: Revisar regularmente o código em busca de credenciais expostas
