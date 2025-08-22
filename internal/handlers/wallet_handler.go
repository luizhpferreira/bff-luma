package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"bff-luma/internal/middleware"
	"bff-luma/internal/models"
	"bff-luma/internal/services"
)

// WalletHandler representa o handler para operações de carteira
type WalletHandler struct {
	walletService *services.WalletService
}

// NewWalletHandler cria um novo handler de carteiras
func NewWalletHandler(walletService *services.WalletService) *WalletHandler {
	return &WalletHandler{
		walletService: walletService,
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

	// Validação básica
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

	response, err := h.walletService.CreateWallet(&req)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Erro ao criar carteira", err.Error())
		return
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
		respondWithError(w, http.StatusBadRequest, "email é obrigatório", "")
		return
	}

	if req.Password == "" {
		respondWithError(w, http.StatusBadRequest, "password é obrigatório", "")
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
