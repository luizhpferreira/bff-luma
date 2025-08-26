package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"bff-luma/internal/models"
)

// LNBitsService representa o serviço de integração com LNBits
type LNBitsService struct {
	baseURL        string
	apiToken       string
	webhookSecret  string
	httpClient     *http.Client
}

// LNBitsUserResponse representa a resposta da criação de usuário no LNBits
type LNBitsUserResponse struct {
	ID              string                 `json:"id"`
	Email           *string                `json:"email"`
	Username        string                 `json:"username"`
	Password        string                 `json:"password"`
	PasswordRepeat  string                 `json:"password_repeat"`
	PubKey          *string                `json:"pubkey"`
	ExternalID      *string                `json:"external_id"`
	Extensions      *string                `json:"extensions"`
	Extra           map[string]interface{} `json:"extra"`
}

// LNBitsInvoiceRequest representa a requisição para criar invoice no LNBits
type LNBitsInvoiceRequest struct {
	Out    bool   `json:"out"`
	Amount int64  `json:"amount"`
	Memo   string `json:"memo,omitempty"`
}

// LNBitsInvoiceResponse representa a resposta da criação de invoice no LNBits
type LNBitsInvoiceResponse struct {
	PaymentRequest string `json:"payment_request"`
	PaymentHash    string `json:"payment_hash"`
	Amount         int64  `json:"amount"`
	Memo           string `json:"memo"`
	Time           int64  `json:"time"`
	ExpiresAt      int64  `json:"expires_at"`
}

// LNBitsPaymentResponse representa a resposta de verificação de pagamento
type LNBitsPaymentResponse struct {
	Paid    bool   `json:"paid"`
	Amount  int64  `json:"amount"`
	Memo    string `json:"memo"`
	Time    int64  `json:"time"`
	Bolt11  string `json:"bolt11"`
	Preimage string `json:"preimage"`
}

// LNBitsChannelRequest representa a requisição para criar canal no LNBits
type LNBitsChannelRequest struct {
	NodeURI string `json:"node_uri"`
	Amount  int64  `json:"amount"`
	Private bool   `json:"private"`
}

// LNBitsChannelResponse representa a resposta da criação de canal
type LNBitsChannelResponse struct {
	ChannelID string `json:"channel_id"`
	NodeURI   string `json:"node_uri"`
	Amount    int64  `json:"amount"`
	Private   bool   `json:"private"`
	Status    string `json:"status"`
}

// NewLNBitsService cria um novo serviço LNBits
func NewLNBitsService(baseURL, apiToken, webhookSecret string) *LNBitsService {
	return &LNBitsService{
		baseURL:       baseURL,
		apiToken:      apiToken,
		webhookSecret: webhookSecret,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// CreateWallet cria um novo usuário no LNBits
func (s *LNBitsService) CreateWallet(username, email, password string) (*models.Wallet, error) {
	url := fmt.Sprintf("%s/users/api/v1/user", s.baseURL)
	
	payload := map[string]interface{}{
		"username":        username,
		"email":           email,
		"password":        password,
		"password_repeat": password,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("erro ao serializar payload: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("erro ao criar requisição: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.apiToken))

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("erro ao fazer requisição: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("erro na resposta do LNBits: %d - %s", resp.StatusCode, string(body))
	}

	var lnbitsResp LNBitsUserResponse
	if err := json.NewDecoder(resp.Body).Decode(&lnbitsResp); err != nil {
		return nil, fmt.Errorf("erro ao decodificar resposta: %w", err)
	}

	// Com a criação do usuário, precisamos criar uma wallet para ele
	// O ID do usuário será usado como wallet_id
	wallet := &models.Wallet{
		WalletID:   lnbitsResp.ID,
		AdminKey:   lnbitsResp.ID, // Temporário - precisaremos obter a wallet real
		InvoiceKey: lnbitsResp.ID, // Temporário - precisaremos obter a wallet real
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	return wallet, nil
}

// CreateInvoice cria um invoice na carteira especificada
func (s *LNBitsService) CreateInvoice(invoiceKey string, amount int64, memo string) (*models.InvoiceResponse, error) {
	url := fmt.Sprintf("%s/api/v1/payments", s.baseURL)
	
	payload := LNBitsInvoiceRequest{
		Out:    false,
		Amount: amount,
		Memo:   memo,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("erro ao serializar payload: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("erro ao criar requisição: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Api-Key", invoiceKey)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("erro ao fazer requisição: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("erro na resposta do LNBits: %d - %s", resp.StatusCode, string(body))
	}

	var lnbitsResp LNBitsInvoiceResponse
	if err := json.NewDecoder(resp.Body).Decode(&lnbitsResp); err != nil {
		return nil, fmt.Errorf("erro ao decodificar resposta: %w", err)
	}

	invoice := &models.InvoiceResponse{
		PaymentRequest: lnbitsResp.PaymentRequest,
		PaymentHash:    lnbitsResp.PaymentHash,
		Amount:         lnbitsResp.Amount,
		Memo:           lnbitsResp.Memo,
		ExpiresAt:      lnbitsResp.ExpiresAt,
	}

	return invoice, nil
}

// CheckPayment verifica o status de um pagamento
func (s *LNBitsService) CheckPayment(invoiceKey, paymentHash string) (*models.PaymentStatus, error) {
	url := fmt.Sprintf("%s/api/v1/payments/%s", s.baseURL, paymentHash)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar requisição: %w", err)
	}

	req.Header.Set("X-Api-Key", invoiceKey)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("erro ao fazer requisição: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("erro na resposta do LNBits: %d - %s", resp.StatusCode, string(body))
	}

	var lnbitsResp LNBitsPaymentResponse
	if err := json.NewDecoder(resp.Body).Decode(&lnbitsResp); err != nil {
		return nil, fmt.Errorf("erro ao decodificar resposta: %w", err)
	}

	paymentStatus := &models.PaymentStatus{
		PaymentHash: paymentHash,
		Paid:        lnbitsResp.Paid,
		Amount:      lnbitsResp.Amount,
		Memo:        lnbitsResp.Memo,
	}

	if lnbitsResp.Paid && lnbitsResp.Time > 0 {
		paymentStatus.PaidAt = &lnbitsResp.Time
	}

	return paymentStatus, nil
}

// PayInvoice paga um invoice usando a carteira do usuário
func (s *LNBitsService) PayInvoice(adminKey, paymentRequest string) (*models.PaymentResponse, error) {
	url := fmt.Sprintf("%s/api/v1/payments", s.baseURL)
	
	payload := map[string]interface{}{
		"out":             true,
		"payment_request": paymentRequest,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("erro ao serializar payload: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("erro ao criar requisição: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Api-Key", adminKey)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("erro ao fazer requisição: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("erro na resposta do LNBits: %d - %s", resp.StatusCode, string(body))
	}

	var lnbitsResp LNBitsPaymentResponse
	if err := json.NewDecoder(resp.Body).Decode(&lnbitsResp); err != nil {
		return nil, fmt.Errorf("erro ao decodificar resposta: %w", err)
	}

	payment := &models.PaymentResponse{
		PaymentHash: lnbitsResp.Preimage, // Preimage é o hash do pagamento
		Paid:        lnbitsResp.Paid,
		Amount:      lnbitsResp.Amount,
		Memo:        lnbitsResp.Memo,
	}

	return payment, nil
}

// CreateInvoiceKey cria uma nova invoice key para uma wallet
func (s *LNBitsService) CreateInvoiceKey(adminKey, name, description string) (*models.InvoiceKey, error) {
	// No LNBits, podemos criar uma nova "extension" ou usar uma abordagem diferente
	// Por enquanto, vamos gerar uma nova invoice key baseada na admin key
	// Em uma implementação real, você pode usar diferentes extensions do LNBits
	
	// Gerar um ID único para a invoice key
	keyID := fmt.Sprintf("inkey_%s_%d", adminKey[:8], time.Now().Unix())
	
	// Para simplificar, vamos usar uma variação da admin key
	// Em produção, você pode implementar uma lógica mais robusta
	newInvoiceKey := fmt.Sprintf("%s_%s", adminKey, keyID)
	
	invoiceKey := &models.InvoiceKey{
		ID:          keyID,
		Name:        name,
		InvoiceKey:  newInvoiceKey,
		Description: description,
		CreatedAt:   time.Now().Format(time.RFC3339),
	}
	
	return invoiceKey, nil
}

// ListInvoiceKeys lista todas as invoice keys de uma wallet
func (s *LNBitsService) ListInvoiceKeys(adminKey string) ([]models.InvoiceKey, error) {
	// Por enquanto, vamos retornar uma lista vazia
	// Em uma implementação real, você pode armazenar as invoice keys no banco de dados
	// ou usar as extensions do LNBits
	
	var invoiceKeys []models.InvoiceKey
	
	// Adicionar a invoice key principal (admin key)
	mainKey := models.InvoiceKey{
		ID:          "main",
		Name:        "Chave Principal",
		InvoiceKey:  adminKey, // Usar a admin key como invoice key principal
		Description: "Chave principal da wallet",
		CreatedAt:   time.Now().Format(time.RFC3339),
	}
	
	invoiceKeys = append(invoiceKeys, mainKey)
	
	return invoiceKeys, nil
}

// CreateInvoiceWithKey cria um invoice usando uma invoice key específica
func (s *LNBitsService) CreateInvoiceWithKey(invoiceKey string, amount int64, memo string) (*models.InvoiceResponse, error) {
	// Se a invoice key for a admin key, usar o método normal
	if strings.HasPrefix(invoiceKey, "lnbits_") {
		return s.CreateInvoice(invoiceKey, amount, memo)
	}
	
	// Para outras invoice keys, você pode implementar lógica específica
	// Por enquanto, vamos usar a admin key base
	adminKey := strings.Split(invoiceKey, "_")[0]
	return s.CreateInvoice(adminKey, amount, memo)
}



// CreateChannel cria um canal Lightning para um usuário específico
func (s *LNBitsService) CreateChannel(adminKey, nodeURI string, amount int64, private bool) (*LNBitsChannelResponse, error) {
	url := fmt.Sprintf("%s/api/v1/channels", s.baseURL)
	
	payload := LNBitsChannelRequest{
		NodeURI: nodeURI,
		Amount:  amount,
		Private: private,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("erro ao serializar payload: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("erro ao criar requisição: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Api-Key", adminKey) // Usa admin key da wallet do usuário

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("erro ao fazer requisição: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("erro na resposta do LNBits: %d - %s", resp.StatusCode, string(body))
	}

	var lnbitsResp LNBitsChannelResponse
	if err := json.NewDecoder(resp.Body).Decode(&lnbitsResp); err != nil {
		return nil, fmt.Errorf("erro ao decodificar resposta: %w", err)
	}

	return &lnbitsResp, nil
}
