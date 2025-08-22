# BFF Luma - Docker Setup

Este documento descreve como executar o BFF Luma usando Docker Compose com PostgreSQL e Cloudflare Tunnel.

## 🚀 Quick Start

### 1. Pré-requisitos

- Docker
- Docker Compose
- Token do Cloudflare Tunnel (opcional)

### 2. Configuração

1. **Use seu arquivo `.env` existente ou copie o exemplo:**
   ```bash
   # Se já tem um .env, apenas configure as variáveis necessárias
   # Se não tem, copie o exemplo:
   cp env.example .env
   ```

2. **Edite o arquivo `.env` com suas configurações:**
   ```bash
   nano .env
   ```

   **Variáveis obrigatórias:**
   - `JWT_SECRET` - Chave secreta para JWT (mude em produção!)
   - `LNBITS_ADMIN_KEY` - Sua chave admin do LNBits
   - `LNBITS_WEBHOOK_SECRET` - Secret para webhooks do LNBits

   **Variáveis opcionais (têm valores padrão):**
   - `POSTGRES_PASSWORD` - Senha do PostgreSQL (padrão: postgres)
   - `SMTP_*` - Configurações de email (opcional)
   - `CLOUDFLARE_TUNNEL_TOKEN` - Token do Cloudflare Tunnel (se usar)

   **Exemplo de `.env` mínimo:**
   ```bash
   JWT_SECRET=sua_chave_secreta_aqui
   LNBITS_ADMIN_KEY=sua_admin_key_lnbits
   LNBITS_WEBHOOK_SECRET=seu_webhook_secret
   ```

### 3. Execução

#### Apenas aplicação + PostgreSQL:
```bash
docker-compose up -d
```

#### Com Cloudflare Tunnel:
```bash
docker-compose --profile tunnel up -d
```

### 4. Verificação

- **Health Check:** http://localhost:8080/health
- **API:** http://localhost:8080/api/v1
- **PostgreSQL:** localhost:5433

## 🔧 Variáveis de Ambiente

### Como funcionam
O Docker Compose usa variáveis do arquivo `.env` automaticamente. Se uma variável não estiver definida, usa o valor padrão após `:-`.

### Exemplos:
```bash
# No .env
JWT_SECRET=minha_chave_secreta
LNBITS_ADMIN_KEY=minha_admin_key

# No docker-compose.yml
- JWT_SECRET=${JWT_SECRET}                    # Obrigatório
- LNBITS_BASE_URL=${LNBITS_BASE_URL:-http://127.0.0.1:5000}  # Com padrão
```

### Variáveis por ambiente
- **Desenvolvimento:** Use `docker-compose.override.yml` (automático)
- **Produção:** Configure todas as variáveis no `.env`
- **Teste:** Crie um `.env.test` e use `--env-file .env.test`

### Vantagens de usar o mesmo .env
- **Simplicidade:** Um único arquivo para tudo
- **Consistência:** Mesmas configurações em todos os ambientes
- **Manutenção:** Menos arquivos para gerenciar

## 📋 Serviços

### PostgreSQL
- **Porta:** 5433 (externa) / 5432 (interna)
- **Database:** bff_luma
- **Usuário:** postgres
- **Senha:** postgres (mude em produção!)

### BFF Luma App
- **Porta:** 8080
- **Health Check:** /health
- **API:** /api/v1

### Cloudflare Tunnel (Opcional)
- **Container:** cloudflare/cloudflared
- **Perfil:** tunnel
- **Dependência:** bff_luma

## 🔧 Comandos Úteis

### Logs
```bash
# Todos os serviços
docker-compose logs -f

# Apenas aplicação
docker-compose logs -f bff_luma

# Apenas PostgreSQL
docker-compose logs -f postgres

# Apenas Cloudflare Tunnel
docker-compose logs -f cloudflared
```

### Parar serviços
```bash
# Parar todos
docker-compose down

# Parar e remover volumes (cuidado: perde dados!)
docker-compose down -v
```

### Rebuild
```bash
# Rebuild da aplicação
docker-compose build bff_luma

# Rebuild e restart
docker-compose up -d --build
```

### Backup do banco
```bash
# Backup
docker-compose exec postgres pg_dump -U postgres bff_luma > backup.sql

# Restore
docker-compose exec -T postgres psql -U postgres bff_luma < backup.sql

# Conectar diretamente ao PostgreSQL (se necessário)
docker-compose exec postgres psql -U postgres -d bff_luma
```

## 🌐 Cloudflare Tunnel

### Configuração

1. **Crie um tunnel no Cloudflare:**
   - Acesse: https://dash.cloudflare.com/
   - Vá em "Zero Trust" > "Access" > "Tunnels"
   - Clique em "Create a tunnel"
   - Escolha "Cloudflared"
   - Copie o token

2. **Configure o token no `.env`:**
   ```bash
   CLOUDFLARE_TUNNEL_TOKEN=seu_token_aqui
   ```

3. **Execute com o perfil tunnel:**
   ```bash
   docker-compose --profile tunnel up -d
   ```

### Configuração do Tunnel

No dashboard do Cloudflare, configure o tunnel para apontar para:
- **Service:** http://bff_luma:8080
- **Hostname:** seu-dominio.com

## 🔒 Produção

### Segurança

1. **Mude todas as senhas padrão:**
   - `JWT_SECRET`
   - `POSTGRES_PASSWORD`
   - `LNBITS_ADMIN_KEY`

2. **Use secrets do Docker:**
   ```yaml
   secrets:
     jwt_secret:
       file: ./secrets/jwt_secret.txt
     postgres_password:
       file: ./secrets/postgres_password.txt
   ```

3. **Configure HTTPS:**
   - Use Cloudflare Tunnel ou
   - Configure um proxy reverso (nginx/traefik)

### Monitoramento

```bash
# Status dos containers
docker-compose ps

# Uso de recursos
docker stats

# Health checks
curl http://localhost:8080/health
```

## 🐛 Troubleshooting

### Problemas comuns

1. **Conflito de porta PostgreSQL:**
   ```bash
   # Se a porta 5433 também estiver em uso, mude no docker-compose.yml:
   ports:
     - "5434:5432"  # Mude para uma porta livre
   
   # E atualize o DATABASE_URL no .env:
   DATABASE_URL=postgres://postgres:postgres@localhost:5434/bff_luma?sslmode=disable
   ```

2. **Conflito de porta da aplicação:**
   ```bash
   # Se a porta 8080 estiver em uso, mude no docker-compose.yml:
   ports:
     - "8081:8080"  # Mude para uma porta livre
   ```

1. **Aplicação não conecta ao banco:**
   ```bash
   docker-compose logs bff_luma
   # Verifique se o PostgreSQL está rodando
   docker-compose ps postgres
   ```

2. **Porta já em uso:**
   ```bash
   # Mude a porta no docker-compose.yml
   ports:
     - "8081:8080"  # Mude 8080 para 8081
   ```

3. **Erro de permissão:**
   ```bash
   # Remova volumes e recrie
   docker-compose down -v
   docker-compose up -d
   ```

### Logs detalhados

```bash
# Logs da aplicação com timestamps
docker-compose logs -f --timestamps bff_luma

# Logs do PostgreSQL
docker-compose logs -f postgres

# Logs do Cloudflare Tunnel
docker-compose logs -f cloudflared
```

## 📁 Estrutura de Arquivos

```
bff_luma/
├── docker-compose.yml          # Configuração principal
├── docker-compose.override.yml # Override para desenvolvimento
├── Dockerfile                  # Build da aplicação
├── .dockerignore              # Arquivos ignorados no build
├── init.sql                   # Script de inicialização do banco
├── env.example                # Exemplo de variáveis de ambiente
└── .env                       # Suas variáveis de ambiente (usar existente)
```

## 🔄 Migração do SQLite

Se você estava usando SQLite anteriormente:

1. **Faça backup dos dados:**
   ```bash
   # Se ainda tem o arquivo SQLite
   sqlite3 bff_luma.db ".dump" > backup.sql
   ```

2. **Converta para PostgreSQL:**
   ```bash
   # Ajuste o backup.sql se necessário
   # Execute no PostgreSQL
   docker-compose exec -T postgres psql -U postgres bff_luma < backup.sql
   ```

## 📞 Suporte

Para problemas ou dúvidas:
- Verifique os logs: `docker-compose logs`
- Consulte a documentação da API
- Abra uma issue no repositório
