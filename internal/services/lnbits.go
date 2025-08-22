package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"bff-luma/internal/models"
)

// LNBitsService representa o serviço de integração com LNBits
type LNBitsService struct {
	baseURL        string
	adminKey       string
	webhookSecret  string
	httpClient     *http.Client
}

// LNBitsWalletResponse representa a resposta da criação de carteira no LNBits
type LNBitsWalletResponse struct {
	ID        string `json:"id"`
	User      string `json:"user"`
	Name      string `json:"name"`
	AdminKey  string `json:"adminkey"`
	InKey     string `json:"inkey"`
	Deleted   bool   `json:"deleted"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Currency  string `json:"currency"`
	BalanceMsat int64 `json:"balance_msat"`
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

// NewLNBitsService cria um novo serviço LNBits
func NewLNBitsService(baseURL, adminKey, webhookSecret string) *LNBitsService {
	return &LNBitsService{
		baseURL:       baseURL,
		adminKey:      adminKey,
		webhookSecret: webhookSecret,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// CreateWallet cria uma nova carteira no LNBits
func (s *LNBitsService) CreateWallet(username, password string) (*models.Wallet, error) {
	url := fmt.Sprintf("%s/api/v1/wallet", s.baseURL)
	
	payload := map[string]interface{}{
		"username": username,
		"password": password,
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
	req.Header.Set("X-Api-Key", s.adminKey)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("erro ao fazer requisição: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("erro na resposta do LNBits: %d - %s", resp.StatusCode, string(body))
	}

	var lnbitsResp LNBitsWalletResponse
	if err := json.NewDecoder(resp.Body).Decode(&lnbitsResp); err != nil {
		return nil, fmt.Errorf("erro ao decodificar resposta: %w", err)
	}

	wallet := &models.Wallet{
		WalletID:   lnbitsResp.ID,
		AdminKey:   lnbitsResp.AdminKey,
		InvoiceKey: lnbitsResp.InKey,
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
