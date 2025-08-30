package models

import (
	"time"
)

// Wallet representa uma carteira Lightning
type Wallet struct {
	ID                        int       `json:"id" db:"id"`
	CPF                       string    `json:"cpf" db:"cpf"`
	Email                     string    `json:"email" db:"email"`
	Password                  string    `json:"-" db:"password"` // Hash da senha, nunca exposto
	WalletID                  string    `json:"wallet_id" db:"wallet_id"`
	AdminKey                  string    `json:"-" db:"admin_key"` // Nunca exposto para o frontend
	InvoiceKey                string    `json:"-" db:"invoice_key"` // Nunca exposto para o frontend
	EmailConfirmed            bool      `json:"email_confirmed" db:"email_confirmed"`
	EmailConfirmationToken    *string   `json:"-" db:"email_confirmation_token"`
	EmailConfirmationExpiresAt *time.Time `json:"-" db:"email_confirmation_expires_at"`
	CreatedAt                 time.Time `json:"created_at" db:"created_at"`
	UpdatedAt                 time.Time `json:"updated_at" db:"updated_at"`
}

// CreateWalletRequest representa a requisição para criar uma carteira
type CreateWalletRequest struct {
	Username        string `json:"username" validate:"required"` // CPF do usuário
	Email           string `json:"email" validate:"required,email"` // Email do usuário
	Password        string `json:"password" validate:"required,min=8"`
	PasswordRepeat  string `json:"password_repeat" validate:"required"`
}

// LoginRequest representa a requisição de login
type LoginRequest struct {
	Email    string `json:"email" validate:"required"` // Campo usado para CPF (não email)
	Password string `json:"password" validate:"required"`
}

// CreateWalletResponse representa a resposta da criação de carteira
type CreateWalletResponse struct {
	WalletID string `json:"wallet_id"`
	Email    string `json:"email"`
	Message  string `json:"message"`
}

// LoginResponse representa a resposta do login
type LoginResponse struct {
	WalletID string `json:"wallet_id"`
	Email    string `json:"email"`
	Username string `json:"username"` // CPF do usuário
	Token    string `json:"token"`
	Message  string `json:"message"`
}

// InvoiceRequest representa a requisição para criar um invoice
type InvoiceRequest struct {
	Amount int64  `json:"amount" validate:"required,min=1"`
	Memo   string `json:"memo,omitempty"`
	Email  string `json:"email" validate:"required,email"`
}

// InvoiceResponse representa a resposta da criação de invoice
type InvoiceResponse struct {
	PaymentRequest string `json:"payment_request"`
	PaymentHash    string `json:"payment_hash"`
	Amount         int64  `json:"amount"`
	Memo           string `json:"memo,omitempty"`
	ExpiresAt      string `json:"expires_at"`
}

// PaymentStatus representa o status de um pagamento
type PaymentStatus struct {
	PaymentHash string `json:"payment_hash"`
	Paid        bool   `json:"paid"`
	Amount      int64  `json:"amount"`
	Memo        string `json:"memo,omitempty"`
	Email       string `json:"email"`
	PaidAt      *int64 `json:"paid_at,omitempty"`
}

// PaymentRequest representa a requisição para pagar um invoice
type PaymentRequest struct {
	PaymentRequest string `json:"payment_request" validate:"required"`
}

// PaymentResponse representa a resposta do pagamento de um invoice
type PaymentResponse struct {
	PaymentHash string `json:"payment_hash"`
	Paid        bool   `json:"paid"`
	Amount      int64  `json:"amount"`
	Memo        string `json:"memo,omitempty"`
}

// InvoiceKey representa uma chave de invoice individual
type InvoiceKey struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	InvoiceKey  string `json:"-"` // Nunca exposto para o frontend
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
}

// CreateInvoiceKeyRequest representa a requisição para criar uma nova invoice key
type CreateInvoiceKeyRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description,omitempty"`
}

// CreateInvoiceKeyResponse representa a resposta da criação de uma nova invoice key
type CreateInvoiceKeyResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	Message     string `json:"message"`
}

// ListInvoiceKeysResponse representa a resposta da listagem de invoice keys
type ListInvoiceKeysResponse struct {
	InvoiceKeys []InvoiceKey `json:"invoice_keys"`
	Total       int          `json:"total"`
}

// ForgotPasswordRequest representa a requisição de recuperação de senha
type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// ResetPasswordRequest representa a requisição de reset de senha
type ResetPasswordRequest struct {
	Token           string `json:"token" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=8"`
	NewPasswordRepeat string `json:"new_password_repeat" validate:"required,min=8"`
}

// ForgotPasswordResponse representa a resposta de recuperação de senha
type ForgotPasswordResponse struct {
	Message string `json:"message"`
}

// ResetPasswordResponse representa a resposta de reset de senha
type ResetPasswordResponse struct {
	Message string `json:"message"`
}

// ConfirmEmailRequest representa a requisição de confirmação de email
type ConfirmEmailRequest struct {
	Token string `json:"token" validate:"required"`
}

// ConfirmEmailResponse representa a resposta de confirmação de email
type ConfirmEmailResponse struct {
	Message string `json:"message"`
	Email   string `json:"email"`
}

// ValidateResetTokenRequest representa a requisição de validação de token de reset
type ValidateResetTokenRequest struct {
	Token string `json:"token" validate:"required"`
}

// ValidateResetTokenResponse representa a resposta de validação de token de reset
type ValidateResetTokenResponse struct {
	Valid   bool   `json:"valid"`
	Message string `json:"message"`
}
