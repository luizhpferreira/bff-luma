package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"bff-luma/internal/middleware"
	"bff-luma/internal/models"
	"bff-luma/internal/services"
)

// WalletHandler representa o handler para operações de carteira
type WalletHandler struct {
	walletService  *services.WalletService
	cleanupService *services.CleanupService
	rateLimiter    *services.RateLimiter
}

// NewWalletHandler cria um novo handler de carteiras
func NewWalletHandler(walletService *services.WalletService, cleanupService *services.CleanupService, rateLimiter *services.RateLimiter) *WalletHandler {
	return &WalletHandler{
		walletService:  walletService,
		cleanupService: cleanupService,
		rateLimiter:    rateLimiter,
	}
}

// CreateWallet cria uma nova carteira
func (h *WalletHandler) CreateWallet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	var req models.CreateWalletRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Erro ao decodificar requisição", err.Error())
		return
	}
	
	// Debug: log do request
	log.Printf("DEBUG: Request recebido - Username: %s, Email: %s, Password: %s", req.Username, req.Email, req.Password)

	// Validação básica
	if req.Username == "" {
		respondWithError(w, http.StatusBadRequest, "CPF é obrigatório", "")
		return
	}

	// Validar se o username é um CPF válido
	if err := h.walletService.ValidateCPF(req.Username); err != nil {
		respondWithError(w, http.StatusBadRequest, "CPF inválido", err.Error())
		return
	}

	if req.Email == "" {
		respondWithError(w, http.StatusBadRequest, "email é obrigatório", "")
		return
	}

	if req.Password == "" {
		respondWithError(w, http.StatusBadRequest, "password é obrigatório", "")
		return
	}

	if req.PasswordRepeat == "" {
		respondWithError(w, http.StatusBadRequest, "password_repeat é obrigatório", "")
		return
	}

	// Valida se as senhas coincidem
	if req.Password != req.PasswordRepeat {
		respondWithError(w, http.StatusBadRequest, "as senhas não coincidem", "")
		return
	}

	// Valida se a senha é forte
	if err := h.walletService.ValidatePasswordStrength(req.Password); err != nil {
		respondWithError(w, http.StatusBadRequest, "erro na validação da senha", err.Error())
		return
	}

	// Cria a carteira (cada usuário terá sua própria wallet no LNBits)
	wallet, err := h.walletService.CreateWallet(req.Username, req.Email, req.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Erro ao criar carteira", err.Error())
		return
	}

	response := &models.CreateWalletResponse{
		WalletID: wallet.WalletID,
		Email:    wallet.Email,
		Message:  "Carteira criada com sucesso",
	}

	respondWithSuccess(w, http.StatusCreated, "Carteira criada com sucesso", response)
}

// ConfirmEmailPage exibe uma página HTML para confirmação de email
func (h *WalletHandler) ConfirmEmailPage(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "Token não fornecido", http.StatusBadRequest)
		return
	}

	// Página HTML que processa a confirmação e redireciona para o app
	html := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Confirmando Email - BFF Luma</title>
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <style>
        body { 
            font-family: Arial, sans-serif; 
            line-height: 1.6; 
            color: #333; 
            margin: 0; 
            padding: 20px; 
            background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%);
            min-height: 100vh;
            display: flex;
            align-items: center;
            justify-content: center;
        }
        .container { 
            max-width: 500px; 
            background: white; 
            padding: 40px; 
            border-radius: 15px; 
            box-shadow: 0 10px 30px rgba(0,0,0,0.2);
            text-align: center;
        }
        .header { 
            background: #28a745; 
            color: white; 
            padding: 20px; 
            border-radius: 10px; 
            margin-bottom: 30px;
        }
        .content { 
            padding: 20px 0; 
        }
        .button { 
            display: inline-block; 
            padding: 15px 30px; 
            background: #007bff; 
            color: white; 
            text-decoration: none; 
            border-radius: 8px; 
            font-weight: bold;
            margin: 10px;
            transition: background 0.3s;
        }
        .button:hover { 
            background: #0056b3; 
        }
        .success { 
            background: #28a745; 
        }
        .success:hover { 
            background: #1e7e34; 
        }
        .error { 
            background: #dc3545; 
        }
        .error:hover { 
            background: #c82333; 
        }
        .loading { 
            display: none; 
        }
        .result { 
            display: none; 
            margin-top: 20px; 
        }
        .spinner {
            border: 4px solid #f3f3f3;
            border-top: 4px solid #007bff;
            border-radius: 50%%;
            width: 40px;
            height: 40px;
            animation: spin 1s linear infinite;
            margin: 20px auto;
        }
        @keyframes spin {
            0%% { transform: rotate(0deg); }
            100%% { transform: rotate(360deg); }
        }
        .app-link {
            background: #6c5ce7;
            margin-top: 20px;
        }
        .app-link:hover {
            background: #5a4fcf;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>📧 Confirmando Email</h1>
            <p>BFF Luma</p>
        </div>
        
        <div class="content">
            <div id="loading" class="loading">
                <div class="spinner"></div>
                <p>Confirmando seu email...</p>
            </div>
            
            <div id="success" class="result">
                <h2>✅ Email Confirmado!</h2>
                <p>Seu email foi confirmado com sucesso!</p>
                <p>Agora você pode fazer login no aplicativo.</p>
                <a href="bffluma://login" class="button success">📱 Abrir App</a>
                <a href="https://play.google.com/store/apps/details?id=com.anonymous.BFFLumaMobile" class="button app-link">📥 Baixar App</a>
            </div>
            
            <div id="error" class="result">
                <h2>❌ Erro na Confirmação</h2>
                <p id="error-message">Ocorreu um erro ao confirmar seu email.</p>
                <button onclick="retryConfirmation()" class="button">🔄 Tentar Novamente</button>
                <a href="mailto:support@bffluma.com" class="button error">📧 Suporte</a>
            </div>
        </div>
    </div>

    <script>
        const token = '%s';
        
        async function confirmEmail() {
            document.getElementById('loading').style.display = 'block';
            
            try {
                const response = await fetch('/api/v1/confirm-email', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({ token: token })
                });
                
                const data = await response.json();
                
                if (response.ok && data.success) {
                    document.getElementById('loading').style.display = 'none';
                    document.getElementById('success').style.display = 'block';
                    
                    // Tenta abrir o app automaticamente após 2 segundos
                    setTimeout(() => {
                        window.location.href = 'bffluma://login';
                    }, 2000);
                } else {
                    throw new Error(data.message || 'Erro na confirmação');
                }
            } catch (error) {
                document.getElementById('loading').style.display = 'none';
                document.getElementById('error').style.display = 'block';
                document.getElementById('error-message').textContent = error.message;
            }
        }
        
        function retryConfirmation() {
            document.getElementById('success').style.display = 'none';
            document.getElementById('error').style.display = 'none';
            confirmEmail();
        }
        
        // Inicia a confirmação automaticamente
        window.onload = confirmEmail;
    </script>
</body>
</html>`, token)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
}

// ResetPasswordPage exibe uma página HTML para reset de senha
func (h *WalletHandler) ResetPasswordPage(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "Token não fornecido", http.StatusBadRequest)
		return
	}

	// Página HTML que processa o token e redireciona para o app
	html := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Redefinir Senha - BFF Luma</title>
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <style>
        body { 
            font-family: Arial, sans-serif; 
            line-height: 1.6; 
            color: #333; 
            margin: 0; 
            padding: 20px; 
            background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%);
            min-height: 100vh;
            display: flex;
            align-items: center;
            justify-content: center;
        }
        .container { 
            max-width: 500px; 
            background: white; 
            padding: 40px; 
            border-radius: 15px; 
            box-shadow: 0 10px 30px rgba(0,0,0,0.2);
            text-align: center;
        }
        .header { 
            background: #007bff; 
            color: white; 
            padding: 20px; 
            border-radius: 10px; 
            margin-bottom: 30px;
        }
        .content { 
            padding: 20px 0; 
        }
        .button { 
            display: inline-block; 
            padding: 15px 30px; 
            background: #007bff; 
            color: white; 
            text-decoration: none; 
            border-radius: 8px; 
            font-weight: bold;
            margin: 10px;
            transition: background 0.3s;
        }
        .button:hover { 
            background: #0056b3; 
        }
        .success { 
            background: #28a745; 
        }
        .success:hover { 
            background: #1e7e34; 
        }
        .error { 
            background: #dc3545; 
        }
        .error:hover { 
            background: #c82333; 
        }
        .loading { 
            display: none; 
        }
        .result { 
            display: none; 
            margin-top: 20px; 
        }
        .spinner {
            border: 4px solid #f3f3f3;
            border-top: 4px solid #007bff;
            border-radius: 50%%;
            width: 40px;
            height: 40px;
            animation: spin 1s linear infinite;
            margin: 20px auto;
        }
        @keyframes spin {
            0%% { transform: rotate(0deg); }
            100%% { transform: rotate(360deg); }
        }
        .app-link {
            background: #6c5ce7;
            margin-top: 20px;
        }
        .app-link:hover {
            background: #5a4fcf;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🔐 Redefinir Senha</h1>
            <p>BFF Luma</p>
        </div>
        
        <div class="content">
            <div id="loading" class="loading">
                <div class="spinner"></div>
                <p>Validando token...</p>
            </div>
            
            <div id="success" class="result">
                <h2>✅ Token Válido!</h2>
                <p>Seu token de recuperação é válido!</p>
                <p>Clique no botão abaixo para abrir o aplicativo e redefinir sua senha.</p>
                <a href="bffluma://reset-password?token=%s" class="button success">📱 Abrir App</a>
                <a href="https://play.google.com/store/apps/details?id=com.anonymous.BFFLumaMobile" class="button app-link">📥 Baixar App</a>
            </div>
            
            <div id="error" class="result">
                <h2>❌ Token Inválido</h2>
                <p id="error-message">O token de recuperação é inválido ou expirou.</p>
                <button onclick="retryValidation()" class="button">🔄 Tentar Novamente</button>
                <a href="mailto:support@bffluma.com" class="button error">📧 Suporte</a>
            </div>
        </div>
    </div>

    <script>
        const token = '%s';
        
        async function validateToken() {
            document.getElementById('loading').style.display = 'block';
            
            try {
                const response = await fetch('/api/v1/validate-reset-token', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({ token: token })
                });
                
                const data = await response.json();
                
                if (response.ok && data.success) {
                    document.getElementById('loading').style.display = 'none';
                    document.getElementById('success').style.display = 'block';
                    
                    // Tenta abrir o app automaticamente após 2 segundos
                    setTimeout(() => {
                        window.location.href = 'bffluma://reset-password?token=' + token;
                    }, 2000);
                } else {
                    throw new Error(data.message || 'Token inválido');
                }
            } catch (error) {
                document.getElementById('loading').style.display = 'none';
                document.getElementById('error').style.display = 'block';
                document.getElementById('error-message').textContent = error.message;
            }
        }
        
        function retryValidation() {
            document.getElementById('success').style.display = 'none';
            document.getElementById('error').style.display = 'none';
            validateToken();
        }
        
        // Inicia a validação automaticamente
        window.onload = validateToken;
    </script>
</body>
</html>`, token, token)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
}

// ConfirmEmail confirma o email do usuário
func (h *WalletHandler) ConfirmEmail(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	var req models.ConfirmEmailRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Erro ao decodificar requisição", err.Error())
		return
	}

	if req.Token == "" {
		respondWithError(w, http.StatusBadRequest, "Token é obrigatório", "")
		return
	}

	// Busca a carteira antes de confirmar para obter o email
	wallet, err := h.walletService.GetWalletByEmailConfirmationToken(req.Token)
	if err != nil || wallet == nil {
		respondWithError(w, http.StatusBadRequest, "Token inválido ou expirado", "")
		return
	}

	// Confirma o email
	if err := h.walletService.ConfirmEmail(req.Token); err != nil {
		respondWithError(w, http.StatusBadRequest, "Erro ao confirmar email", err.Error())
		return
	}

	response := &models.ConfirmEmailResponse{
		Message: "Email confirmado com sucesso!",
		Email:   wallet.Email,
	}

	respondWithSuccess(w, http.StatusOK, "Email confirmado com sucesso", response)
}



// Login autentica um usuário
func (h *WalletHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Erro ao decodificar requisição", err.Error())
		return
	}

	// Validação básica
	if req.Email == "" {
		respondWithError(w, http.StatusBadRequest, "CPF é obrigatório", "")
		return
	}

	// Validar se o CPF é válido
	if err := h.walletService.ValidateCPF(req.Email); err != nil {
		respondWithError(w, http.StatusBadRequest, "CPF inválido", err.Error())
		return
	}

	if req.Password == "" {
		respondWithError(w, http.StatusBadRequest, "password é obrigatório", "")
		return
	}

	// Verifica rate limit por CPF
	if !h.rateLimiter.AllowLogin(req.Email) {
		remaining := h.rateLimiter.GetRemainingAttempts(req.Email)
		respondWithError(w, http.StatusTooManyRequests, "Muitas tentativas de login. Tente novamente em 15 minutos.", fmt.Sprintf("Tentativas restantes: %d", remaining))
		return
	}

	response, err := h.walletService.Login(&req)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Erro no login", err.Error())
		return
	}

	respondWithSuccess(w, http.StatusOK, "Login realizado com sucesso", response)
}

// RefreshToken renova um token JWT
func (h *WalletHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	// Extrai o token do header Authorization
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		respondWithError(w, http.StatusUnauthorized, "Token de autorização não fornecido", "")
		return
	}

	// Verifica se o header tem o formato "Bearer <token>"
	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		respondWithError(w, http.StatusUnauthorized, "Formato de autorização inválido. Use: Bearer <token>", "")
		return
	}

	tokenString := tokenParts[1]

	// Renova o token
	newToken, err := h.walletService.RefreshToken(tokenString)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Erro ao renovar token", err.Error())
		return
	}

	response := map[string]interface{}{
		"token":   newToken,
		"message": "Token renovado com sucesso",
	}

	respondWithSuccess(w, http.StatusOK, "Token renovado com sucesso", response)
}

// ForgotPassword inicia o processo de recuperação de senha
func (h *WalletHandler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	var req models.ForgotPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Erro ao decodificar requisição", err.Error())
		return
	}

	// Validação básica
	if req.Email == "" {
		respondWithError(w, http.StatusBadRequest, "email é obrigatório", "")
		return
	}

	// Verifica rate limit para reset de senha
	if !h.rateLimiter.AllowPasswordReset(req.Email) {
		remaining := h.rateLimiter.GetRemainingPasswordResets(req.Email)
		respondWithError(w, http.StatusTooManyRequests, "Muitas tentativas de recuperação de senha. Tente novamente em 1 hora.", fmt.Sprintf("Tentativas restantes: %d", remaining))
		return
	}

	response, err := h.walletService.ForgotPassword(&req)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Erro ao processar recuperação de senha", err.Error())
		return
	}

	respondWithSuccess(w, http.StatusOK, "Solicitação processada", response)
}

// ValidateResetToken valida um token de reset de senha
func (h *WalletHandler) ValidateResetToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	var req models.ValidateResetTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Erro ao decodificar requisição", err.Error())
		return
	}

	if req.Token == "" {
		respondWithError(w, http.StatusBadRequest, "token é obrigatório", "")
		return
	}

	// Valida o token
	valid, err := h.walletService.ValidateResetToken(req.Token)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Erro ao validar token", err.Error())
		return
	}

	if !valid {
		respondWithError(w, http.StatusBadRequest, "Token inválido ou expirado", "")
		return
	}

	response := &models.ValidateResetTokenResponse{
		Valid:   true,
		Message: "Token válido",
	}

	respondWithSuccess(w, http.StatusOK, "Token válido", response)
}

// ResetPassword redefine a senha usando um token
func (h *WalletHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	var req models.ResetPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Erro ao decodificar requisição", err.Error())
		return
	}

	// Validação básica
	if req.Token == "" {
		respondWithError(w, http.StatusBadRequest, "token é obrigatório", "")
		return
	}

	if req.NewPassword == "" {
		respondWithError(w, http.StatusBadRequest, "new_password é obrigatório", "")
		return
	}

	if req.NewPasswordRepeat == "" {
		respondWithError(w, http.StatusBadRequest, "new_password_repeat é obrigatório", "")
		return
	}

	response, err := h.walletService.ResetPassword(&req)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Erro ao redefinir senha", err.Error())
		return
	}

	respondWithSuccess(w, http.StatusOK, "Senha redefinida com sucesso", response)
}

// CleanupTokens executa limpeza manual de tokens expirados
func (h *WalletHandler) CleanupTokens(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	// Executa limpeza manual
	err := h.cleanupService.CleanupNow()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Erro ao executar limpeza", err.Error())
		return
	}

	// Obtém estatísticas
	stats, err := h.cleanupService.GetStats()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Erro ao obter estatísticas", err.Error())
		return
	}

	respondWithSuccess(w, http.StatusOK, "Limpeza executada com sucesso", stats)
}

// GetCleanupStats retorna estatísticas de limpeza
func (h *WalletHandler) GetCleanupStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	stats, err := h.cleanupService.GetStats()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Erro ao obter estatísticas", err.Error())
		return
	}

	respondWithSuccess(w, http.StatusOK, "Estatísticas obtidas com sucesso", stats)
}

// GetRateLimitStats retorna estatísticas do rate limiter
func (h *WalletHandler) GetRateLimitStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	stats := h.rateLimiter.GetStats()
	respondWithSuccess(w, http.StatusOK, "Estatísticas do rate limiter obtidas com sucesso", stats)
}

// CreateInvoice cria um invoice para receber pagamento
func (h *WalletHandler) CreateInvoice(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	var req models.InvoiceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Erro ao decodificar requisição", err.Error())
		return
	}

	// Obtém o email do contexto JWT
	email := middleware.GetUserEmail(r)
	if email == "" {
		respondWithError(w, http.StatusUnauthorized, "Usuário não autenticado", "")
		return
	}

	// Usa o email do JWT em vez do email da requisição
	req.Email = email

	if req.Amount <= 0 {
		respondWithError(w, http.StatusBadRequest, "amount deve ser maior que zero", "")
		return
	}

	response, err := h.walletService.CreateInvoice(&req)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Erro ao criar invoice", err.Error())
		return
	}

	respondWithSuccess(w, http.StatusCreated, "Invoice criado com sucesso", response)
}

// CheckPaymentStatus verifica o status de um pagamento
func (h *WalletHandler) CheckPaymentStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	// Obtém o email do contexto JWT
	email := middleware.GetUserEmail(r)
	if email == "" {
		respondWithError(w, http.StatusUnauthorized, "Usuário não autenticado", "")
		return
	}

	// Extrai payment_hash da query string
	paymentHash := r.URL.Query().Get("payment_hash")
	if paymentHash == "" {
		respondWithError(w, http.StatusBadRequest, "payment_hash é obrigatório", "")
		return
	}

	response, err := h.walletService.CheckPaymentStatus(email, paymentHash)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Erro ao verificar status do pagamento", err.Error())
		return
	}

	respondWithSuccess(w, http.StatusOK, "Status do pagamento verificado", response)
}

// PayInvoice paga um invoice
func (h *WalletHandler) PayInvoice(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	var req models.PaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Erro ao decodificar requisição", err.Error())
		return
	}

	// Obtém o email do contexto JWT
	email := middleware.GetUserEmail(r)
	if email == "" {
		respondWithError(w, http.StatusUnauthorized, "Usuário não autenticado", "")
		return
	}

	if req.PaymentRequest == "" {
		respondWithError(w, http.StatusBadRequest, "payment_request é obrigatório", "")
		return
	}

	response, err := h.walletService.PayInvoice(email, req.PaymentRequest)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Erro ao pagar invoice", err.Error())
		return
	}

	respondWithSuccess(w, http.StatusOK, "Invoice pago com sucesso", response)
}

// CreateInvoiceKey cria uma nova invoice key para o usuário
func (h *WalletHandler) CreateInvoiceKey(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	var req models.CreateInvoiceKeyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Erro ao decodificar requisição", err.Error())
		return
	}

	// Obtém o email do contexto JWT
	email := middleware.GetUserEmail(r)
	if email == "" {
		respondWithError(w, http.StatusUnauthorized, "Usuário não autenticado", "")
		return
	}

	if req.Name == "" {
		respondWithError(w, http.StatusBadRequest, "name é obrigatório", "")
		return
	}

	response, err := h.walletService.CreateInvoiceKey(email, req.Name, req.Description)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Erro ao criar invoice key", err.Error())
		return
	}

	respondWithSuccess(w, http.StatusCreated, "Invoice key criada com sucesso", response)
}

// ListInvoiceKeys lista todas as invoice keys do usuário
func (h *WalletHandler) ListInvoiceKeys(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	// Obtém o email do contexto JWT
	email := middleware.GetUserEmail(r)
	if email == "" {
		respondWithError(w, http.StatusUnauthorized, "Usuário não autenticado", "")
		return
	}

	response, err := h.walletService.ListInvoiceKeys(email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Erro ao listar invoice keys", err.Error())
		return
	}

	respondWithSuccess(w, http.StatusOK, "Invoice keys listadas com sucesso", response)
}

// CreateInvoiceWithKey cria um invoice usando uma invoice key específica
func (h *WalletHandler) CreateInvoiceWithKey(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		InvoiceKeyID string `json:"invoice_key_id" validate:"required"`
		Amount       int64  `json:"amount" validate:"required,min=1"`
		Memo         string `json:"memo,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Erro ao decodificar requisição", err.Error())
		return
	}

	// Obtém o email do contexto JWT
	email := middleware.GetUserEmail(r)
	if email == "" {
		respondWithError(w, http.StatusUnauthorized, "Usuário não autenticado", "")
		return
	}

	if req.InvoiceKeyID == "" {
		respondWithError(w, http.StatusBadRequest, "invoice_key_id é obrigatório", "")
		return
	}

	if req.Amount <= 0 {
		respondWithError(w, http.StatusBadRequest, "amount deve ser maior que zero", "")
		return
	}

	response, err := h.walletService.CreateInvoiceWithKey(email, req.InvoiceKeyID, req.Amount, req.Memo)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Erro ao criar invoice", err.Error())
		return
	}

	respondWithSuccess(w, http.StatusCreated, "Invoice criado com sucesso", response)
}

// GetWalletInfo retorna informações da carteira
func (h *WalletHandler) GetWalletInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	// Obtém o email do contexto JWT
	email := middleware.GetUserEmail(r)
	if email == "" {
		respondWithError(w, http.StatusUnauthorized, "Usuário não autenticado", "")
		return
	}

	response, err := h.walletService.GetWalletInfo(email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Erro ao buscar informações da carteira", err.Error())
		return
	}

	respondWithSuccess(w, http.StatusOK, "Informações da carteira", response)
}

// GetWalletBalance retorna o saldo da carteira
func (h *WalletHandler) GetWalletBalance(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	// Obtém o email do contexto JWT
	email := middleware.GetUserEmail(r)
	if email == "" {
		respondWithError(w, http.StatusUnauthorized, "Usuário não autenticado", "")
		return
	}

	response, err := h.walletService.GetWalletBalance(email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Erro ao buscar saldo da carteira", err.Error())
		return
	}

	respondWithSuccess(w, http.StatusOK, "Saldo da carteira obtido com sucesso", response)
}

// HealthCheck verifica se a API está funcionando
func (h *WalletHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	response := map[string]interface{}{
		"status":  "ok",
		"service": "BFF Luma API",
		"version": "1.0.0",
	}

	respondWithSuccess(w, http.StatusOK, "API funcionando", response)
}

// respondWithSuccess envia uma resposta de sucesso
func respondWithSuccess(w http.ResponseWriter, statusCode int, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := models.APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	}

	json.NewEncoder(w).Encode(response)
}

// respondWithError envia uma resposta de erro
func respondWithError(w http.ResponseWriter, statusCode int, message, error string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := models.ErrorResponse{
		Success: false,
		Error:   error,
		Message: message,
	}

	json.NewEncoder(w).Encode(response)
}
