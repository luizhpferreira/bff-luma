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
	log.Printf("DEBUG: Request recebido - Username: %s, Password: %s", req.Username, req.Password)

	// Validação básica
	if req.Username == "" {
		respondWithError(w, http.StatusBadRequest, "username é obrigatório", "")
		return
	}

	// Validar se o username é um CPF válido
	if err := h.walletService.ValidateCPF(req.Username); err != nil {
		respondWithError(w, http.StatusBadRequest, "CPF inválido", err.Error())
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
	wallet, err := h.walletService.CreateWallet(req.Username, req.Password)
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

	// Validar se o email (que agora é CPF) é um CPF válido
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
