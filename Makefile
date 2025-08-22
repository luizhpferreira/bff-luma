# Makefile para BFF Luma

.PHONY: build run test clean help

# Variáveis
BINARY_NAME=bff-luma
MAIN_PATH=cmd/server/main.go

# Comandos principais
build:
	@echo "🔨 Compilando BFF Luma..."
	go build -o $(BINARY_NAME) $(MAIN_PATH)
	@echo "✅ Compilação concluída!"

run:
	@echo "🚀 Executando BFF Luma..."
	go run $(MAIN_PATH)

dev:
	@echo "🔄 Executando em modo desenvolvimento..."
	@echo "📝 Logs detalhados ativados..."
	go run $(MAIN_PATH)

test:
	@echo "🧪 Executando testes..."
	go test ./...

test-api:
	@echo "🧪 Testando API..."
	@if [ -f "./test_api.sh" ]; then \
		./test_api.sh; \
	else \
		echo "❌ Script de teste não encontrado"; \
	fi

clean:
	@echo "🧹 Limpando arquivos..."
	rm -f $(BINARY_NAME)
	rm -f *.db
	@echo "✅ Limpeza concluída!"

deps:
	@echo "📦 Baixando dependências..."
	go mod tidy
	go mod vendor
	@echo "✅ Dependências atualizadas!"

install:
	@echo "📥 Instalando dependências..."
	go mod download
	@echo "✅ Instalação concluída!"

lint:
	@echo "🔍 Executando linter..."
	golangci-lint run

format:
	@echo "🎨 Formatando código..."
	go fmt ./...

# Comando de ajuda
help:
	@echo "📋 Comandos disponíveis:"
	@echo ""
	@echo "  build     - Compila o projeto"
	@echo "  run       - Executa o servidor"
	@echo "  dev       - Executa em modo desenvolvimento"
	@echo "  test      - Executa testes"
	@echo "  test-api  - Testa a API com script"
	@echo "  clean     - Remove arquivos compilados"
	@echo "  deps      - Atualiza dependências"
	@echo "  install   - Instala dependências"
	@echo "  lint      - Executa linter"
	@echo "  format    - Formata código"
	@echo "  help      - Mostra esta ajuda"
	@echo ""
	@echo "💡 Exemplo: make run"
