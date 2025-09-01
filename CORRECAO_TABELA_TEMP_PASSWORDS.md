# Correção do Erro: Tabela `temp_passwords` não existe

## Problema Identificado

O erro ocorria durante a confirmação de email:
```
❌ Erro ao obter senha original: erro ao buscar senha original: pq: relation "temp_passwords" does not exist
```

## Causa Raiz

A tabela `temp_passwords` não estava sendo criada no método `createTables()` do arquivo `database.go`. Esta tabela é essencial para o fluxo de confirmação de email, pois armazena temporariamente a senha original do usuário durante o processo de confirmação.

## Solução Implementada

### 1. Adição da Tabela no Código

Adicionada a criação da tabela `temp_passwords` no método `createTables()` em `bff_luma/internal/database/database.go`:

```sql
CREATE TABLE IF NOT EXISTS temp_passwords (
    id SERIAL PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_temp_passwords_email ON temp_passwords(email);
```

### 2. Arquivo de Inicialização SQL

Criado arquivo `bff_luma/init.sql/init.sql` para garantir que a tabela seja criada mesmo em bancos existentes:

```sql
-- Criação da tabela temp_passwords se não existir
CREATE TABLE IF NOT EXISTS temp_passwords (
    id SERIAL PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Criação do índice se não existir
CREATE INDEX IF NOT EXISTS idx_temp_passwords_email ON temp_passwords(email);
```

### 3. Correção do Docker Compose

Atualizado o `docker-compose.yml` para apontar para o arquivo correto:

```yaml
volumes:
  - ./init.sql/init.sql:/docker-entrypoint-initdb.d/init.sql
```

### 4. Criação Manual da Tabela

Para bancos já existentes, a tabela foi criada manualmente:

```sql
CREATE TABLE IF NOT EXISTS temp_passwords (
    id SERIAL PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_temp_passwords_email ON temp_passwords(email);
```

## Funcionalidade da Tabela

A tabela `temp_passwords` é usada para:

1. **Armazenar senha original temporariamente**: Durante a criação da carteira, a senha original é salva nesta tabela
2. **Recuperar senha para LNBits**: Durante a confirmação de email, a senha é recuperada para criar a carteira no LNBits
3. **Limpeza automática**: A senha é removida após a confirmação bem-sucedida

## Fluxo Corrigido

1. ✅ Usuário cria conta → Senha salva em `temp_passwords`
2. ✅ Token de confirmação gerado
3. ✅ Usuário confirma email → Senha recuperada de `temp_passwords`
4. ✅ Carteira criada no LNBits
5. ✅ Senha removida de `temp_passwords`
6. ✅ Email de boas-vindas enviado

## Teste de Validação

O teste foi executado com sucesso:

```bash
# Criação da carteira
{"success":true,"message":"Carteira criada com sucesso"}

# Confirmação do email
{"success":true,"message":"Email confirmado com sucesso"}

# Verificação do status
email_confirmed: t (true)
```

## Resultado

✅ **Problema resolvido**: A confirmação de email agora funciona corretamente
✅ **Tabela criada**: `temp_passwords` existe e está funcionando
✅ **Fluxo completo**: Criação → Confirmação → LNBits → Limpeza
✅ **Segurança mantida**: Senhas temporárias são removidas após uso

## Arquivos Modificados

- `bff_luma/internal/database/database.go` - Adicionada criação da tabela
- `bff_luma/init.sql/init.sql` - Criado arquivo de inicialização
- `bff_luma/docker-compose.yml` - Corrigido caminho do arquivo SQL
- `bff_luma/test_email_confirmation_fix.sh` - Script de teste criado
