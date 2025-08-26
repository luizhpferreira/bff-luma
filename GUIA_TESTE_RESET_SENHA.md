# 🎯 Guia de Teste - Reset de Senha com Deep Link

## ✅ Problema Resolvido!

O problema foi corrigido! O backend estava configurado para usar `localhost` em vez do domínio real. Agora o fluxo está funcionando corretamente.

## 🔧 Correções Aplicadas

### Backend (.env)
```bash
# Antes (problemático)
APP_DOMAIN=localhost 
APP_PROTOCOL=http

# Depois (correto)
APP_DOMAIN=luma.app.br 
APP_PROTOCOL=https
```

### Reinicialização
- Backend reiniciado para aplicar novas configurações
- Página de reset agora usa HTTPS e domínio correto

## 🧪 Como Testar o Fluxo Completo

### 1. Solicitar Reset de Senha
1. Abra o aplicativo móvel
2. Vá para "Esqueci a senha"
3. Digite um email válido
4. Clique em "Enviar Email"

### 2. Verificar Email
1. Verifique sua caixa de entrada
2. Procure pelo email de "Recuperação de Senha - Luma"
3. O email deve conter um link como: `https://luma.app.br/reset-password?token=TOKEN`

### 3. Clicar no Link do Email
1. Clique no link "🔑 Redefinir Senha" no email
2. **IMPORTANTE**: Clique no link no **dispositivo móvel** onde o app está instalado
3. O navegador deve abrir com a página de validação

### 4. Página Web de Validação
1. A página deve mostrar "Validando token..."
2. Se o token for válido, mostrará "✅ Token Válido!"
3. Clique no botão "📱 Abrir App"

### 5. Deep Link para o App
1. O app deve abrir automaticamente
2. Deve navegar para a tela "Redefinir Senha"
3. O token deve ser extraído automaticamente da URL

### 6. Definir Nova Senha
1. Digite sua nova senha
2. Confirme a nova senha
3. Clique em "Redefinir Senha"
4. Deve mostrar sucesso e redirecionar para login

## 🔗 Links de Teste

### Página de Reset (para testar)
```
https://luma.app.br/reset-password?token=test-token
```

### Deep Link Direto (para testar app)
```
bffluma://reset-password?token=test-token
```

## 📱 Teste no Dispositivo Móvel

### Android
1. Abra o Chrome no Android
2. Digite: `bffluma://reset-password?token=test-token`
3. Pressione Enter
4. O app deve abrir automaticamente

### iOS
1. Abra o Safari no iOS
2. Digite: `bffluma://reset-password?token=test-token`
3. Pressione Enter
4. O app deve abrir automaticamente

## 🚨 Problemas Comuns e Soluções

### 1. App não abre com deep link
- **Causa**: App não está instalado ou scheme não configurado
- **Solução**: Verificar se o app está instalado e se `app.json` tem `"scheme": "bffluma"`

### 2. Página web não carrega
- **Causa**: Backend não está rodando ou domínio incorreto
- **Solução**: Verificar se o backend está rodando em `https://luma.app.br`

### 3. Token inválido
- **Causa**: Token expirou (1 hora) ou já foi usado
- **Solução**: Solicitar novo reset de senha

### 4. Email não recebido
- **Causa**: Configurações SMTP ou email incorreto
- **Solução**: Verificar configurações SMTP no backend

## ✅ Checklist de Verificação

- [ ] Backend rodando em `https://luma.app.br`
- [ ] Configurações `.env` com domínio correto
- [ ] App instalado no dispositivo
- [ ] Scheme `bffluma` configurado no `app.json`
- [ ] Deep linking configurado no `AppNavigator.tsx`
- [ ] Tela `ResetPasswordScreen` recebe token via rota

## 🎉 Resultado Esperado

Após seguir este guia:
1. ✅ Email enviado com link correto
2. ✅ Página web carrega e valida token
3. ✅ Deep link abre o app automaticamente
4. ✅ App navega para tela de reset com token
5. ✅ Usuário pode definir nova senha
6. ✅ Senha é atualizada com sucesso

## 📞 Suporte

Se ainda houver problemas:
1. Verifique os logs do backend
2. Teste em dispositivo físico (não emulador)
3. Verifique se todas as configurações estão corretas
4. Consulte a documentação completa em `RECUPERACAO_SENHA.md`
