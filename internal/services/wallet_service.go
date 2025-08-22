package services

import (
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"

	"bff-luma/internal/database"
	"bff-luma/internal/models"
)

// WalletService representa o serviço de gerenciamento de carteiras
type WalletService struct {
	db       *database.Database
	lnbits   *LNBitsService
	jwt      *JWTService
	email    *EmailService
	password *PasswordService
}

// NewWalletService cria um novo serviço de carteiras
func NewWalletService(db *database.Database, lnbits *LNBitsService, jwt *JWTService, email *EmailService, password *PasswordService) *WalletService {
	return &WalletService{
		db:       db,
		lnbits:   lnbits,
		jwt:      jwt,
		email:    email,
		password: password,
	}
}



// CreateWallet cria uma nova carteira para um usuário
func (s *WalletService) CreateWallet(req *models.CreateWalletRequest) (*models.CreateWalletResponse, error) {
	// Verifica se as senhas coincidem
	if req.Password != req.PasswordRepeat {
		return nil, fmt.Errorf("as senhas não coincidem")
	}

	// Valida se a senha é forte
	if err := s.password.ValidatePasswordStrength(req.Password); err != nil {
		return nil, fmt.Errorf("erro na validação da senha: %w", err)
	}

	// Verifica se já existe uma carteira para este email
	exists, err := s.db.WalletExists(req.Email)
	if err != nil {
		return nil, fmt.Errorf("erro ao verificar existência da carteira: %w", err)
	}

	if exists {
		return nil, fmt.Errorf("carteira já existe para o email %s", req.Email)
	}

	// Cria a carteira no LNBits usando o email como username
	wallet, err := s.lnbits.CreateWallet(req.Email, req.Password)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar carteira no LNBits: %w", err)
	}

	// Gera hash da senha
	hashedPassword, err := s.password.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("erro ao gerar hash da senha: %w", err)
	}

	// Associa o email e senha hash à carteira
	wallet.Email = req.Email
	wallet.Password = hashedPassword

	// Salva no banco de dados
	if err := s.db.CreateWallet(wallet); err != nil {
		return nil, fmt.Errorf("erro ao salvar carteira no banco: %w", err)
	}

	// Envia email de boas-vindas
	go func() {
		if err := s.email.SendWelcomeEmail(req.Email, wallet.WalletID); err != nil {
			log.Printf("⚠️ Erro ao enviar email de boas-vindas para %s: %v", req.Email, err)
		}
	}()

	log.Printf("Carteira criada com sucesso para email %s: %s", req.Email, wallet.WalletID)

	response := &models.CreateWalletResponse{
		WalletID: wallet.WalletID,
		Email:    wallet.Email,
		Message:  "Carteira criada com sucesso",
	}

	return response, nil
}

// CreateInvoice cria um invoice para um usuário
func (s *WalletService) CreateInvoice(req *models.InvoiceRequest) (*models.InvoiceResponse, error) {
	// Busca a carteira do usuário
	wallet, err := s.db.GetWalletByEmail(req.Email)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar carteira: %w", err)
	}

	if wallet == nil {
		return nil, fmt.Errorf("carteira não encontrada para o email %s", req.Email)
	}

	// Cria o invoice no LNBits usando a invoice key da carteira
	invoice, err := s.lnbits.CreateInvoice(wallet.InvoiceKey, req.Amount, req.Memo)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar invoice no LNBits: %w", err)
	}

	log.Printf("Invoice criado com sucesso para email %s: %s", req.Email, invoice.PaymentHash)

	return invoice, nil
}

// CheckPaymentStatus verifica o status de um pagamento
func (s *WalletService) CheckPaymentStatus(email, paymentHash string) (*models.PaymentStatus, error) {
	// Busca a carteira do usuário
	wallet, err := s.db.GetWalletByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar carteira: %w", err)
	}

	if wallet == nil {
		return nil, fmt.Errorf("carteira não encontrada para o email %s", email)
	}

	// Verifica o status do pagamento no LNBits
	paymentStatus, err := s.lnbits.CheckPayment(wallet.InvoiceKey, paymentHash)
	if err != nil {
		return nil, fmt.Errorf("erro ao verificar status do pagamento: %w", err)
	}

	return paymentStatus, nil
}

// Login autentica um usuário
func (s *WalletService) Login(req *models.LoginRequest) (*models.LoginResponse, error) {
	wallet, err := s.db.GetWalletByEmail(req.Email)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar carteira: %w", err)
	}

	if wallet == nil {
		return nil, fmt.Errorf("carteira não encontrada para o email %s", req.Email)
	}

	// Verifica a senha usando bcrypt
	if err := s.password.CheckPassword(req.Password, wallet.Password); err != nil {
		return nil, fmt.Errorf("senha incorreta")
	}

	// Gera token JWT
	token, err := s.jwt.GenerateToken(wallet.Email, wallet.WalletID)
	if err != nil {
		return nil, fmt.Errorf("erro ao gerar token: %w", err)
	}

	response := &models.LoginResponse{
		WalletID: wallet.WalletID,
		Email:    wallet.Email,
		Token:    token,
		Message:  "Login realizado com sucesso",
	}

	return response, nil
}

// GetWalletInfo retorna informações da carteira (sem as chaves sensíveis)
func (s *WalletService) GetWalletInfo(email string) (*models.Wallet, error) {
	wallet, err := s.db.GetWalletByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar carteira: %w", err)
	}

	if wallet == nil {
		return nil, fmt.Errorf("carteira não encontrada para o email %s", email)
	}

	// Remove as chaves sensíveis antes de retornar
	wallet.AdminKey = ""
	wallet.InvoiceKey = ""
	wallet.Password = ""

	return wallet, nil
}

// RefreshToken renova um token JWT
func (s *WalletService) RefreshToken(tokenString string) (string, error) {
	return s.jwt.RefreshToken(tokenString)
}

// ForgotPassword inicia o processo de recuperação de senha
func (s *WalletService) ForgotPassword(req *models.ForgotPasswordRequest) (*models.ForgotPasswordResponse, error) {
	// Verifica se o email existe
	exists, err := s.db.WalletExists(req.Email)
	if err != nil {
		return nil, fmt.Errorf("erro ao verificar existência do usuário: %w", err)
	}

	if !exists {
		// Por segurança, não revelamos se o email existe ou não
		log.Printf("Tentativa de reset de senha para email inexistente: %s", req.Email)
		return &models.ForgotPasswordResponse{
			Message: "Se o email existir em nossa base, você receberá um link de recuperação",
		}, nil
	}

	// Gera token único
	token := uuid.New().String()
	
	// Token expira em 1 hora
	expiresAt := time.Now().Add(1 * time.Hour)

	// Salva token no banco
	err = s.db.CreateResetToken(req.Email, token, expiresAt)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar token de reset: %w", err)
	}

	// Envia email (simulado)
	err = s.email.SendPasswordResetEmail(req.Email, token)
	if err != nil {
		return nil, fmt.Errorf("erro ao enviar email: %w", err)
	}

	log.Printf("Token de reset criado para %s: %s", req.Email, token)

	return &models.ForgotPasswordResponse{
		Message: "Se o email existir em nossa base, você receberá um link de recuperação",
	}, nil
}

// ResetPassword redefine a senha usando um token válido
func (s *WalletService) ResetPassword(req *models.ResetPasswordRequest) (*models.ResetPasswordResponse, error) {
	// Valida se as senhas coincidem
	if req.NewPassword != req.NewPasswordRepeat {
		return nil, fmt.Errorf("as senhas não coincidem")
	}

	// Valida força da senha
	if err := s.password.ValidatePasswordStrength(req.NewPassword); err != nil {
		return nil, fmt.Errorf("senha não atende aos requisitos de segurança: %w", err)
	}

	// Busca token válido
	email, expiresAt, used, err := s.db.GetResetToken(req.Token)
	if err != nil {
		return nil, fmt.Errorf("erro ao validar token: %w", err)
	}

	if email == "" {
		return nil, fmt.Errorf("token inválido ou expirado")
	}

	if used {
		return nil, fmt.Errorf("token já foi utilizado")
	}

	if time.Now().After(expiresAt) {
		return nil, fmt.Errorf("token expirado")
	}

	// Gera hash da nova senha
	hashedPassword, err := s.password.HashPassword(req.NewPassword)
	if err != nil {
		return nil, fmt.Errorf("erro ao gerar hash da senha: %w", err)
	}

	// Atualiza senha
	err = s.db.UpdatePassword(email, hashedPassword)
	if err != nil {
		return nil, fmt.Errorf("erro ao atualizar senha: %w", err)
	}

	// Marca token como usado
	err = s.db.MarkResetTokenAsUsed(req.Token)
	if err != nil {
		log.Printf("Erro ao marcar token como usado: %v", err)
		// Não falha a operação por causa disso
	}

	log.Printf("Senha redefinida com sucesso para: %s", email)

	return &models.ResetPasswordResponse{
		Message: "Senha redefinida com sucesso",
	}, nil
}
