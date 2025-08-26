# 🔍 Diagnóstico de Deep Link - Reset de Senha

## ❌ Problema: App não abre com deep link

O token está sendo validado corretamente, mas quando você clica em "Abrir App", nada acontece.

## 🔍 Possíveis Causas

### 1. App não está instalado
- **Sintoma**: Nada acontece ao clicar no deep link
- **Solução**: Instalar o app no dispositivo

### 2. Scheme não configurado corretamente
- **Sintoma**: Deep link não é reconhecido pelo sistema
- **Solução**: Verificar configuração do scheme

### 3. App não compilado com scheme
- **Sintoma**: App instalado mas deep link não funciona
- **Solução**: Recompilar app com configurações corretas

### 4. Testando em emulador
- **Sintoma**: Deep links podem não funcionar em emuladores
- **Solução**: Testar em dispositivo físico

## 🧪 Testes de Diagnóstico

### Teste 1: Verificar se o app está instalado
```bash
# No dispositivo Android
adb shell pm list packages | grep bffluma

# No dispositivo iOS (via Xcode)
# Verificar se o app aparece na lista de apps instalados
```

### Teste 2: Testar deep link manualmente
```bash
# No terminal do dispositivo
adb shell am start -W -a android.intent.action.VIEW -d "bffluma://reset-password?token=test" com.anonymous.BFFLumaMobile

# Ou abrir no navegador do dispositivo
# Digite: bffluma://reset-password?token=test
```

### Teste 3: Verificar configuração do scheme
Verifique se o arquivo `app.json` tem:
```json
{
  "expo": {
    "scheme": "bffluma"
  }
}
```

### Teste 4: Verificar se o app foi compilado corretamente
```bash
# Recompilar o app
cd mobile_luma
npx expo build:android
# ou
npx expo build:ios
```

## 🔧 Soluções

### Solução 1: Recompilar o App
```bash
cd mobile_luma

# Limpar cache
npx expo start --clear

# Recompilar
npx expo build:android
# ou
npx expo build:ios
```

### Solução 2: Verificar Configurações
1. Verificar `app.json`:
```json
{
  "expo": {
    "scheme": "bffluma",
    "android": {
      "package": "com.anonymous.BFFLumaMobile"
    },
    "ios": {
      "bundleIdentifier": "com.anonymous.BFFLumaMobile"
    }
  }
}
```

2. Verificar `AppNavigator.tsx`:
```typescript
const linking = {
  prefixes: ['bffluma://'],
  config: {
    screens: {
      ResetPassword: {
        path: 'reset-password',
        parse: {
          token: (token: string) => token,
        },
      },
    },
  },
};
```

### Solução 3: Teste Alternativo
Se o deep link não funcionar, você pode:

1. **Copiar o token** da página web
2. **Abrir o app manualmente**
3. **Ir para "Esqueci a senha"**
4. **Colar o token** na tela de reset

## 📱 Teste em Dispositivo Físico

### Android
1. Instalar o app via APK ou Google Play
2. Abrir Chrome
3. Digitar: `bffluma://reset-password?token=test`
4. Pressionar Enter
5. Verificar se o app abre

### iOS
1. Instalar o app via App Store ou TestFlight
2. Abrir Safari
3. Digitar: `bffluma://reset-password?token=test`
4. Pressionar Enter
5. Verificar se o app abre

## 🚨 Troubleshooting Avançado

### Logs do App
```bash
# Android
adb logcat | grep -i bffluma

# iOS
# Verificar logs no Xcode Console
```

### Verificar Intent Filters (Android)
```xml
<!-- No AndroidManifest.xml -->
<intent-filter>
    <action android:name="android.intent.action.VIEW" />
    <category android:name="android.intent.category.DEFAULT" />
    <category android:name="android.intent.category.BROWSABLE" />
    <data android:scheme="bffluma" />
</intent-filter>
```

### Verificar URL Schemes (iOS)
```xml
<!-- No Info.plist -->
<key>CFBundleURLTypes</key>
<array>
    <dict>
        <key>CFBundleURLName</key>
        <string>bffluma</string>
        <key>CFBundleURLSchemes</key>
        <array>
            <string>bffluma</string>
        </array>
    </dict>
</array>
```

## ✅ Checklist de Verificação

- [ ] App instalado no dispositivo
- [ ] Scheme configurado no `app.json`
- [ ] App compilado com configurações corretas
- [ ] Testando em dispositivo físico (não emulador)
- [ ] Deep link testado manualmente
- [ ] Logs verificados para erros

## 🎯 Próximos Passos

1. **Verificar se o app está instalado**
2. **Testar deep link manualmente**
3. **Recompilar app se necessário**
4. **Testar em dispositivo físico**
5. **Verificar logs para erros**

## 📞 Suporte

Se o problema persistir:
1. Verifique os logs do dispositivo
2. Teste em dispositivo físico diferente
3. Verifique se o app foi compilado corretamente
4. Considere usar fallback (copiar token manualmente)
