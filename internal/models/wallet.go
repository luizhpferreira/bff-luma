package models

import (
	"time"
)

// Wallet representa uma carteira Lightning
type Wallet struct {
	ID         int       `json:"id" db:"id"`
	Email      string    `json:"email" db:"email"`
	Password   string    `json:"-" db:"password"` // Hash da senha, nunca exposto
	WalletID   string    `json:"wallet_id" db:"wallet_id"`
	AdminKey   string    `json:"-" db:"admin_key"` // Nunca exposto para o frontend
	InvoiceKey string    `json:"-" db:"invoice_key"` // Nunca exposto para o frontend
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

// CreateWalletRequest representa a requisição para criar uma carteira
type CreateWalletRequest struct {
	Username        string `json:"username" validate:"required"` // CPF do usuário
	Password        string `json:"password" validate:"required,min=8"`
	PasswordRepeat  string `json:"password_repeat" validate:"required"`
}

// LoginRequest representa a requisição de login
type LoginRequest struct {
	Email    string `json:"email" validate:"required"` // Campo usado para CPF
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
	ExpiresAt      int64  `json:"expires_at"`
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
