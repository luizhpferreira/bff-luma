# Funcionalidade de Pagamento de Invoices

## Visão Geral

O BFF Luma implementa uma funcionalidade completa para pagamento de invoices Lightning Network. Os usuários autenticados podem pagar invoices BOLT11 usando suas carteiras.

## Arquitetura

### Componentes Implementados

1. **Modelos** (`internal/models/wallet.go`)
   - `PaymentRequest`: Requisição para pagar invoice
   - `PaymentResponse`: Resposta do pagamento

2. **Serviço LNBits** (`internal/services/lnbits.go`)
   - `PayInvoice()`: Integração com LNBits para pagamento

3. **Serviço Wallet** (`internal/services/wallet_service.go`)
   - `PayInvoice()`: Lógica de negócio para pagamento

4. **Handler** (`internal/handlers/wallet_handler.go`)
   - `PayInvoice()`: Endpoint HTTP para pagamento

5. **Rotas** (`cmd/server/main.go`)
   - `POST /api/v1/payments`: Endpoint de pagamento

## Fluxo de Pagamento

```
1. Usuário autenticado → JWT Token
2. Envia payment_request (BOLT11)
3. Sistema valida JWT e identifica carteira
4. Usa AdminKey da carteira para pagar
5. LNBits processa pagamento
6. Retorna status do pagamento
```

## Endpoints

### POST /api/v1/payments

**Headers obrigatórios:**
```
Authorization: Bearer <JWT_TOKEN>
Content-Type: application/json
```

**Body:**
```json
{
  "payment_request": "lnbc1..."
}
```

**Resposta de sucesso:**
```json
{
  "success": true,
  "message": "Invoice pago com sucesso",
  "data": {
    "payment_hash": "abc123...",
    "paid": true,
    "amount": 1000,
    "memo": "Pagamento teste"
  }
}
```

## Segurança

- ✅ Autenticação JWT obrigatória
- ✅ AdminKey nunca exposta no frontend
- ✅ Validação de carteira do usuário
- ✅ Rate limiting aplicado

## Testes

### Scripts de Teste

1. **test_pay_invoice.sh**: Verificação da implementação
2. **test_pay_invoice_advanced.sh**: Testes automatizados da API

### Como Testar

```bash
# 1. Iniciar servidor
cd bff_luma && make run

# 2. Executar testes
./test_pay_invoice_advanced.sh

# 3. Teste manual com invoice real
curl -X POST http://localhost:8080/api/v1/payments \
  -H 'Content-Type: application/json' \
  -H 'Authorization: Bearer <JWT_TOKEN>' \
  -d '{"payment_request": "lnbc1..."}'
```

## Próximos Passos

1. ✅ Implementação no BFF
2. 🔄 Integração com aplicativo móvel
3. 🔄 Tratamento de erros específicos
4. 🔄 Logs detalhados para debugging
5. 🔄 Monitoramento de transações

## Status

**✅ CONCLUÍDO**: Funcionalidade de pagamento implementada e testada no BFF.
