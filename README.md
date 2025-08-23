<<<<<<< HEAD
# BFF Luma

Backend for Frontend (BFF) para integração com LNBits, fornecendo uma API REST para gerenciamento de carteiras Lightning.

## Funcionalidades

- Criação de carteiras Lightning via LNBits
- Autenticação JWT
- Gerenciamento de usuários
- Criação de faturas
- Verificação de status de pagamentos
- Reset de senha via email
- Rate limiting
- Limpeza automática de tokens expirados
- Suporte a PostgreSQL e SQLite

## Tecnologias

- **Go 1.23** - Linguagem principal
- **Chi Router** - Roteamento HTTP
- **PostgreSQL/SQLite** - Banco de dados
- **JWT** - Autenticação
- **LNBits** - Backend Lightning
- **SMTP** - Envio de emails
- **Docker** - Containerização

## Estrutura do Projeto

```
bff_luma/
├── cmd/
│   └── server/
│       └── main.go          # Ponto de entrada da aplicação
├── internal/
│   ├── config/
│   │   └── config.go        # Configurações da aplicação
│   ├── database/
│   │   └── database.go      # Camada de acesso a dados
│   ├── handlers/
│   │   └── wallet_handler.go # Handlers HTTP
│   ├── middleware/
│   │   ├── auth.go          # Middleware de autenticação
│   │   └── rate_limit.go    # Middleware de rate limiting
│   ├── models/
│   │   ├── response.go      # Modelos de resposta
│   │   └── wallet.go        # Modelos de dados
│   └── services/
│       ├── cleanup_service.go    # Serviço de limpeza
│       ├── email_service.go      # Serviço de email
│       ├── jwt_service.go        # Serviço JWT
│       ├── lnbits.go             # Integração com LNBits
│       ├── password_service.go   # Serviço de senhas
│       ├── rate_limiter.go       # Rate limiter
│       └── wallet_service.go     # Serviço de carteiras
├── docker-compose.yml       # Orquestração de containers
├── Dockerfile              # Imagem Docker da aplicação
├── .dockerignore           # Arquivos ignorados pelo Docker
├── go.mod                  # Dependências Go
├── go.sum                  # Checksums das dependências
├── env.example             # Exemplo de variáveis de ambiente
├── Makefile                # Comandos de automação
└── README.md               # Este arquivo
```

## Instalação

### Opção 1: Docker Compose (Recomendado)

1. Clone o repositório:
```bash
git clone <url-do-repositorio>
cd bff_luma
```

2. Configure as variáveis de ambiente:
```bash
cp env.example .env
# Edite o arquivo .env com suas configurações
```

3. Execute com Docker Compose:
```bash
# Iniciar todos os serviços
make docker-up

# Ou usando docker-compose diretamente
docker-compose up -d
```

4. Verifique o status:
```bash
make status
```

### Opção 2: Desenvolvimento Local

#### Pré-requisitos

- Go 1.23 ou superior
- PostgreSQL ou SQLite3
- Acesso a um servidor LNBits

#### Configuração

1. Clone o repositório:
```bash
git clone <url-do-repositorio>
cd bff_luma
```

2. Instale as dependências:
```bash
make dev-deps
```

3. Configure as variáveis de ambiente:
```bash
make dev-setup
# Edite o arquivo .env com suas configurações
```

4. Execute a aplicação:
```bash
make run
```

## Configuração das Variáveis de Ambiente

Copie o arquivo `env.example` para `.env` e configure as seguintes variáveis:

### Configurações da Aplicação
- `APP_PORT`: Porta onde a aplicação será executada (padrão: 8080)
- `JWT_SECRET`: Chave secreta para assinatura de JWT

### Configurações LNBits
- `LNBITS_BASE_URL`: URL base do servidor LNBits
- `LNBITS_ADMIN_KEY`: Chave de administrador do LNBits
- `LNBITS_WEBHOOK_SECRET`: Segredo para webhooks
- `LNBITS_PORT`: Porta do LNBits (padrão: 5000)
- `LNBITS_SECRET`: Segredo do LNBits
- `LNBITS_ADMIN_UI`: Habilitar UI de admin (true/false)
- `LNBITS_ACTIVE_EXTENSIONS`: Extensões ativas (ex: usermanager)

### Configurações do Banco de Dados
- `DATABASE_TYPE`: Tipo de banco (postgres/sqlite)
- `DATABASE_URL`: URL de conexão PostgreSQL
- `DATABASE_PATH`: Caminho para arquivo SQLite (fallback)
- `POSTGRES_USER`: Usuário PostgreSQL
- `POSTGRES_PASSWORD`: Senha PostgreSQL
- `POSTGRES_DB`: Nome do banco PostgreSQL

### Configurações Bitcoin/CLN
- `BITCOIN_RPCUSER`: Usuário RPC Bitcoin
- `BITCOIN_RPCPASSWORD`: Senha RPC Bitcoin
- `CLN_ALIAS`: Alias do nó CLN

### Configurações Cloudflare
- `CLOUDFLARED_TOKEN`: Token do Cloudflare Tunnel
- `PUBLIC_HOST`: Host público

### Configurações SMTP
- `SMTP_HOST`: Servidor SMTP
- `SMTP_PORT`: Porta SMTP
- `SMTP_USERNAME`: Usuário SMTP
- `SMTP_PASSWORD`: Senha SMTP
- `SMTP_FROM_EMAIL`: Email remetente
- `SMTP_FROM_NAME`: Nome remetente
- `SMTP_USE_TLS`: Usar TLS (true/false)

## Comandos Úteis

### Docker Compose
```bash
# Iniciar todos os serviços
make docker-up

# Parar todos os serviços
make docker-down

# Ver logs de todos os serviços
make docker-logs

# Ver logs apenas da aplicação
make docker-logs-app

# Reiniciar apenas a aplicação
make docker-restart

# Ver status dos containers
make status

# Verificar saúde da aplicação
make health
```

### Desenvolvimento
```bash
# Compilar localmente
make build

# Executar localmente
make run

# Executar testes
make test

# Limpar arquivos de build
make clean

# Instalar dependências
make dev-deps
```

### Banco de Dados
```bash
# Resetar banco de dados (CUIDADO!)
make db-reset

# Fazer backup
make backup
```

### Logs Específicos
```bash
# Logs do Bitcoin
make logs-bitcoin

# Logs do CLN
make logs-cln

# Logs do LNBits
make logs-lnbits

# Logs do PostgreSQL
make logs-postgres
```

## Uso

### Endpoints da API

#### Autenticação
- `POST /api/v1/wallets` - Criar nova carteira
- `POST /api/v1/login` - Fazer login
- `POST /api/v1/refresh` - Renovar token
- `POST /api/v1/forgot-password` - Solicitar reset de senha
- `POST /api/v1/reset-password` - Resetar senha

#### Carteiras (Autenticado)
- `GET /api/v1/wallets` - Obter informações da carteira
- `POST /api/v1/invoices` - Criar fatura
- `GET /api/v1/payments/status` - Verificar status de pagamento

#### Administração
- `POST /api/v1/admin/cleanup` - Executar limpeza manual
- `GET /api/v1/admin/cleanup/stats` - Estatísticas de limpeza
- `GET /api/v1/admin/rate-limit/stats` - Estatísticas de rate limiting

#### Monitoramento
- `GET /health` - Health check

### Exemplos de Uso

#### Criar uma carteira
```bash
curl -X POST http://localhost:8080/api/v1/wallets \
  -H "Content-Type: application/json" \
  -d '{
    "email": "usuario@exemplo.com",
    "password": "senha123"
  }'
```

#### Fazer login
```bash
curl -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "usuario@exemplo.com",
    "password": "senha123"
  }'
```

#### Criar uma fatura (autenticado)
```bash
curl -X POST http://localhost:8080/api/v1/invoices \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer SEU_TOKEN_JWT" \
  -d '{
    "amount": 1000,
    "description": "Pagamento de teste"
  }'
```

## Arquitetura Docker

O projeto inclui uma stack completa com:

- **bitcoind**: Nó Bitcoin principal
- **cln**: Core Lightning Network
- **postgres**: Banco de dados PostgreSQL
- **lnbits**: Interface Lightning
- **bff-luma**: Esta aplicação BFF
- **cloudflared**: Tunnel Cloudflare
- **tor**: Proxy Tor

### Rede Docker

Todos os serviços estão na rede `luma` e podem se comunicar usando os nomes dos containers.

### Volumes

- `./data/bitcoin`: Dados do Bitcoin
- `./data/cln`: Dados do CLN
- `./data/postgres`: Dados do PostgreSQL
- `./data/tor`: Dados do Tor

## Estrutura da API

### Respostas

Todas as respostas seguem o padrão:
```json
{
  "success": true,
  "message": "Mensagem de sucesso",
  "data": {
    // Dados específicos da resposta
  }
}
```

### Códigos de Status

- `200` - Sucesso
- `201` - Criado
- `400` - Requisição inválida
- `401` - Não autorizado
- `403` - Proibido
- `404` - Não encontrado
- `429` - Muitas requisições
- `500` - Erro interno do servidor

## Segurança

- Autenticação JWT
- Rate limiting por IP
- Validação de entrada
- Sanitização de dados
- Headers de segurança CORS
- Timeout de requisições
- Containers não-root

## Monitoramento

- Health check endpoint
- Logs estruturados
- Métricas de rate limiting
- Estatísticas de limpeza automática
- Monitoramento de containers

## Troubleshooting

### Problemas Comuns

1. **Aplicação não inicia**
   ```bash
   make docker-logs-app
   ```

2. **Banco de dados não conecta**
   ```bash
   make logs-postgres
   ```

3. **LNBits não responde**
   ```bash
   make logs-lnbits
   ```

4. **Bitcoin não sincroniza**
   ```bash
   make logs-bitcoin
   ```

### Reset Completo

Para resetar completamente o ambiente:
```bash
make docker-down
docker system prune -f
docker volume prune -f
make docker-up
```

## Contribuição

1. Fork o projeto
2. Crie uma branch para sua feature (`git checkout -b feature/AmazingFeature`)
3. Commit suas mudanças (`git commit -m 'Add some AmazingFeature'`)
4. Push para a branch (`git push origin feature/AmazingFeature`)
5. Abra um Pull Request

## Licença

Este projeto está sob a licença MIT. Veja o arquivo `LICENSE` para mais detalhes.
=======
# lm-app
>>>>>>> origin/main
