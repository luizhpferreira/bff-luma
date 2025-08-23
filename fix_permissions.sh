#!/bin/bash

echo "🔧 Corrigindo permissões do banco de dados..."

# Parar o backend
echo "1. Parando o backend..."
pkill bff-luma 2>/dev/null
sleep 2

# Remover banco antigo
echo "2. Removendo banco antigo..."
rm -f bff_luma.db

# Criar novo banco com permissões corretas
echo "3. Criando novo banco..."
sqlite3 bff_luma.db "SELECT 1;" > /dev/null 2>&1

# Definir permissões
echo "4. Definindo permissões..."
chmod 666 bff_luma.db
chown luiz:luiz bff_luma.db 2>/dev/null

# Verificar permissões
echo "5. Verificando permissões..."
ls -la bff_luma.db

# Testar escrita
echo "6. Testando escrita..."
sqlite3 bff_luma.db "CREATE TABLE test (id INTEGER); DROP TABLE test;" 2>/dev/null
if [ $? -eq 0 ]; then
    echo "✅ Banco de dados está funcionando corretamente"
else
    echo "❌ Problema com permissões do banco"
    exit 1
fi

echo "7. Iniciando backend..."
./bff-luma &
sleep 5

echo "8. Testando API..."
curl -s http://localhost:8080/health > /dev/null
if [ $? -eq 0 ]; then
    echo "✅ Backend iniciado com sucesso"
else
    echo "❌ Problema ao iniciar backend"
    exit 1
fi

echo ""
echo "🎉 Permissões corrigidas! O sistema está pronto para uso."
