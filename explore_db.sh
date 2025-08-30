#!/bin/bash

# Script para explorar o banco de dados SQLite
DB_FILE="bff_luma.db"

echo "=== Explorador do Banco de Dados SQLite ==="
echo "Arquivo: $DB_FILE"
echo ""

# Verificar se o arquivo existe
if [ ! -f "$DB_FILE" ]; then
    echo "❌ Arquivo $DB_FILE não encontrado!"
    exit 1
fi

# Verificar tamanho do arquivo
FILE_SIZE=$(stat -c%s "$DB_FILE")
echo "📁 Tamanho do arquivo: $FILE_SIZE bytes"
echo ""

if [ $FILE_SIZE -eq 0 ]; then
    echo "⚠️  O banco de dados está vazio!"
    echo "Para inicializar o banco, execute:"
    echo "  sqlite3 $DB_FILE < migrate.sql"
    echo ""
    exit 0
fi

echo "📋 Tabelas no banco:"
sqlite3 "$DB_FILE" ".tables"
echo ""

echo "🏗️  Estrutura das tabelas:"
sqlite3 "$DB_FILE" ".schema"
echo ""

echo "📊 Contagem de registros por tabela:"
for table in $(sqlite3 "$DB_FILE" ".tables"); do
    count=$(sqlite3 "$DB_FILE" "SELECT COUNT(*) FROM $table;")
    echo "  $table: $count registros"
done
echo ""

echo "🔍 Para explorar mais, use:"
echo "  sqlite3 $DB_FILE"
echo ""
echo "Comandos úteis dentro do sqlite3:"
echo "  .tables                    - Listar tabelas"
echo "  .schema                    - Ver estrutura"
echo "  SELECT * FROM tabela;      - Ver dados"
echo "  .mode column               - Formato de colunas"
echo "  .headers on                - Mostrar cabeçalhos"
echo "  .quit                      - Sair"
