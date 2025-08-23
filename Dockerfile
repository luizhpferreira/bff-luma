# Build stage
FROM golang:1.24 AS builder

# Instala dependências necessárias
RUN apt-get update && apt-get install -y git ca-certificates tzdata gcc libc6-dev && rm -rf /var/lib/apt/lists/*

# Define o diretório de trabalho
WORKDIR /app

# Copia os arquivos de dependências
COPY go.mod go.sum ./

# Baixa as dependências
RUN go mod download

# Copia o código fonte
COPY . .

# Compila o aplicativo (CGO habilitado para PostgreSQL)
RUN CGO_ENABLED=1 GOOS=linux go build -ldflags "-linkmode external -extldflags -static" -a -installsuffix cgo -o main ./cmd/server

# Final stage
FROM ubuntu:22.04

# Instala ca-certificates para HTTPS
RUN apt-get update && apt-get install -y ca-certificates tzdata && rm -rf /var/lib/apt/lists/*

# Cria usuário não-root
RUN groupadd -g 1001 appgroup && \
    useradd -u 1001 -g appgroup -s /bin/bash -m appuser

# Define o diretório de trabalho
WORKDIR /app

# Copia o binário compilado
COPY --from=builder /app/main .

# Muda a propriedade para o usuário não-root
RUN chown -R appuser:appgroup /app

# Muda para o usuário não-root
USER appuser

# Expõe a porta
EXPOSE 8080

# Comando para executar o aplicativo
CMD ["./main"]
