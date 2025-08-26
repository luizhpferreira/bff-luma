package services

import (
	"fmt"
	"log"
	"strings"
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



// ValidateCPF valida se o CPF é válido
func (s *WalletService) ValidateCPF(cpf string) error {
	// Remove caracteres não numéricos
	cpf = strings.ReplaceAll(cpf, ".", "")
	cpf = strings.ReplaceAll(cpf, "-", "")
	
	if len(cpf) != 11 {
		return fmt.Errorf("CPF deve ter 11 dígitos")
	}
	
	// Verifica se todos os dígitos são iguais
	if cpf == "00000000000" || cpf == "11111111111" || cpf == "22222222222" ||
		cpf == "33333333333" || cpf == "44444444444" || cpf == "55555555555" ||
		cpf == "66666666666" || cpf == "77777777777" || cpf == "88888888888" ||
		cpf == "99999999999" {
		return fmt.Errorf("CPF inválido")
	}
	
	// Validação dos dígitos verificadores
	sum := 0
	for i := 0; i < 9; i++ {
		sum += int(cpf[i]-'0') * (10 - i)
	}
	remainder := (sum * 10) % 11
	if remainder == 10 || remainder == 11 {
		remainder = 0
	}
	if remainder != int(cpf[9]-'0') {
		return fmt.Errorf("CPF inválido")
	}
	
	sum = 0
	for i := 0; i < 10; i++ {
		sum += int(cpf[i]-'0') * (11 - i)
	}
	remainder = (sum * 10) % 11
	if remainder == 10 || remainder == 11 {
		remainder = 0
	}
	if remainder != int(cpf[10]-'0') {
		return fmt.Errorf("CPF inválido")
	}
	
	return nil
}

// CreateWallet cria uma nova carteira para o usuário (apenas dados básicos, sem LNBits)
func (s *WalletService) CreateWallet(username, email, password string) (*models.Wallet, error) {
	// Verifica se a carteira já existe
	exists, err := s.db.WalletExists(username)
	if err != nil {
		return nil, fmt.Errorf("erro ao verificar existência da carteira: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("carteira já existe para este CPF")
	}

	// Gera hash da senha
	hashedPassword, err := s.password.HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("erro ao gerar hash da senha: %w", err)
	}

	// Cria objeto wallet apenas com dados básicos (sem LNBits ainda)
	wallet := &models.Wallet{
		CPF:      username,        // CPF do usuário
		Email:    email,           // Email do usuário
		Password: hashedPassword,  // Senha hasheada
		// WalletID, AdminKey, InvoiceKey serão preenchidos após confirmação
	}

	// Salva no banco de dados (apenas dados básicos)
	if err := s.db.CreateWallet(wallet); err != nil {
		return nil, fmt.Errorf("erro ao salvar dados básicos no banco: %w", err)
	}

	// Salva a senha original temporariamente para uso na confirmação
	if err := s.db.SaveOriginalPassword(email, password); err != nil {
		log.Printf("⚠️ Erro ao salvar senha original temporária para %s: %v", email, err)
	}

	// Gera token de confirmação de email
	confirmationToken := uuid.New().String()
	expiresAt := time.Now().Add(1 * time.Hour) // Token expira em 1 hora
	
	// Salva o token no banco
	if err := s.db.CreateEmailConfirmationToken(email, confirmationToken, expiresAt); err != nil {
		log.Printf("⚠️ Erro ao salvar token de confirmação para %s: %v", email, err)
	} else {
		log.Printf("🔑 Token de confirmação criado para %s: %s", email, confirmationToken)
	}

	// Envia email de confirmação
	log.Printf("🚀 Iniciando envio de email de confirmação para %s", email)
	go func() {
		log.Printf("📧 Tentando enviar email de confirmação para %s", email)
		if s.email == nil {
			log.Printf("❌ Email service é nil!")
			return
		}
		if err := s.email.SendEmailConfirmation(email, confirmationToken); err != nil {
			log.Printf("⚠️ Erro ao enviar email de confirmação para %s: %v", email, err)
		} else {
			log.Printf("📧 Email de confirmação enviado com sucesso para %s", email)
		}
	}()

	log.Printf("Dados básicos salvos para CPF %s e email %s (aguardando confirmação)", username, email)

	return wallet, nil
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

	// Cria o invoice no LNBits usando a invoice key da carteira do usuário
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
	// Remove formatação do CPF se houver
	cleanCPF := strings.ReplaceAll(req.Email, ".", "")
	cleanCPF = strings.ReplaceAll(cleanCPF, "-", "")
	cleanCPF = strings.ReplaceAll(cleanCPF, " ", "")
	
	// Busca a carteira pelo CPF
	wallet, err := s.db.GetWalletByCPF(cleanCPF)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar carteira por CPF: %w", err)
	}

	if wallet == nil {
		return nil, fmt.Errorf("carteira não encontrada para o CPF %s", req.Email)
	}

	// Verifica se o email foi confirmado
	if !wallet.EmailConfirmed {
		return nil, fmt.Errorf("email não confirmado. Verifique sua caixa de entrada e confirme seu email antes de fazer login")
	}

	// Verifica se a carteira foi criada no LNBits
	if wallet.WalletID == "" {
		return nil, fmt.Errorf("carteira não está ativa. Entre em contato com o suporte")
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
		Username: wallet.CPF, // CPF do usuário
		Token:    token,
		Message:  "Login realizado com sucesso",
	}

	return response, nil
}

// ConfirmEmail confirma o email do usuário e cria a carteira no LNBits
func (s *WalletService) ConfirmEmail(token string) error {
	// Busca a carteira pelo token para obter o email
	wallet, err := s.db.GetWalletByEmailConfirmationToken(token)
	if err != nil {
		return fmt.Errorf("erro ao buscar carteira por token: %w", err)
	}
	
	if wallet == nil {
		return fmt.Errorf("token inválido ou expirado")
	}

	// Confirma o email no banco de dados
	if err := s.db.ConfirmEmail(token); err != nil {
		return fmt.Errorf("erro ao confirmar email: %w", err)
	}

	log.Printf("✅ Email confirmado com sucesso para token: %s", token)

	// Obtém a senha original para criar a carteira no LNBits
	originalPassword, err := s.db.GetOriginalPassword(wallet.Email)
	if err != nil {
		log.Printf("❌ Erro ao obter senha original: %v", err)
		if unconfirmErr := s.db.UnconfirmEmail(wallet.Email); unconfirmErr != nil {
			log.Printf("⚠️ Erro ao desconfirmar email após falha: %v", unconfirmErr)
		}
		return fmt.Errorf("erro ao obter senha original: %w", err)
	}

	// Cria a carteira no LNBits após confirmação
	log.Printf("🚀 Criando carteira no LNBits para %s", wallet.Email)
	lnbitsWallet, err := s.lnbits.CreateWallet(wallet.CPF, wallet.Email, originalPassword)
	if err != nil {
		// Se falhar ao criar no LNBits, marca como não confirmado novamente
		log.Printf("❌ Erro ao criar carteira no LNBits: %v", err)
		if unconfirmErr := s.db.UnconfirmEmail(wallet.Email); unconfirmErr != nil {
			log.Printf("⚠️ Erro ao desconfirmar email após falha no LNBits: %v", unconfirmErr)
		}
		return fmt.Errorf("erro ao criar carteira no LNBits: %w", err)
	}

	// Atualiza a carteira no banco com os dados do LNBits
	wallet.WalletID = lnbitsWallet.WalletID
	wallet.AdminKey = lnbitsWallet.AdminKey
	wallet.InvoiceKey = lnbitsWallet.InvoiceKey
	
	if err := s.db.UpdateWalletLNBitsData(wallet.Email, wallet.WalletID, wallet.AdminKey, wallet.InvoiceKey); err != nil {
		log.Printf("⚠️ Erro ao atualizar dados do LNBits no banco: %v", err)
		// Não falha a operação, mas loga o erro
	}

	log.Printf("✅ Carteira criada no LNBits com sucesso: %s", wallet.WalletID)

	// Remove a senha original temporária após sucesso
	if err := s.db.RemoveOriginalPassword(wallet.Email); err != nil {
		log.Printf("⚠️ Erro ao remover senha original temporária: %v", err)
	}

	// Envia email de boas-vindas após a confirmação
	log.Printf("🚀 Iniciando envio de email de boas-vindas para %s", wallet.Email)
	go func() {
		log.Printf("📧 Tentando enviar email de boas-vindas para %s", wallet.Email)
		if s.email == nil {
			log.Printf("❌ Email service é nil!")
			return
		}
		if err := s.email.SendWelcomeEmail(wallet.Email, wallet.WalletID); err != nil {
			log.Printf("⚠️ Erro ao enviar email de boas-vindas para %s: %v", wallet.Email, err)
		} else {
			log.Printf("📧 Email de boas-vindas enviado com sucesso para %s", wallet.Email)
		}
	}()

	return nil
}

// GetWalletByEmailConfirmationToken obtém uma carteira pelo token de confirmação
func (s *WalletService) GetWalletByEmailConfirmationToken(token string) (*models.Wallet, error) {
	return s.db.GetWalletByEmailConfirmationToken(token)
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
	exists, err := s.db.WalletExistsByEmail(req.Email)
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

// ValidateResetToken valida se um token de reset é válido
func (s *WalletService) ValidateResetToken(token string) (bool, error) {
	// Busca token válido
	email, expiresAt, used, err := s.db.GetResetToken(token)
	if err != nil {
		return false, fmt.Errorf("erro ao validar token: %w", err)
	}

	if email == "" {
		return false, nil // Token não existe
	}

	if used {
		return false, nil // Token já foi usado
	}

	if time.Now().After(expiresAt) {
		return false, nil // Token expirado
	}

	return true, nil // Token válido
}

// ValidatePasswordStrength valida se a senha é forte
func (s *WalletService) ValidatePasswordStrength(password string) error {
	return s.password.ValidatePasswordStrength(password)
}

// generateUsernameFromEmail gera um apelido inteligente baseado no email
// Exemplos:
// luiz.fernando.silva@gmail.com -> lufe_gmail
// joao.pedro@hotmail.com -> jope_hotmail
// maria@exemplo.com -> maria_exemplo
func (s *WalletService) generateUsernameFromEmail(email string) string {
	// Separa o email em partes
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		// Se não conseguir separar, usa um hash simples
		return fmt.Sprintf("user_%x", strings.ToLower(email)[:8])
	}
	
	localPart := parts[0]  // parte antes do @
	domainPart := parts[1] // parte depois do @
	
	// Remove extensão do domínio (gmail.com -> gmail)
	domainName := strings.Split(domainPart, ".")[0]
	
	// Gera apelido da parte local
	nickname := generateNickname(localPart)
	
	// Combina nickname + domínio
	username := fmt.Sprintf("%s_%s", nickname, domainName)
	
	// Garante que não passe de 20 caracteres
	if len(username) > 20 {
		// Se ainda for muito grande, trunca mantendo a proporção
		maxNickname := 20 - len(domainName) - 1 // -1 para o underscore
		if maxNickname < 2 {
			maxNickname = 2
		}
		nickname = nickname[:maxNickname]
		username = fmt.Sprintf("%s_%s", nickname, domainName)
		
		// Se ainda for muito grande, trunca tudo
		if len(username) > 20 {
			username = username[:20]
		}
	}
	
	return strings.ToLower(username)
}

// generateNickname gera um apelido inteligente da parte local do email
func generateNickname(localPart string) string {
	// Remove caracteres especiais e substitui por underscore
	cleaned := strings.ReplaceAll(strings.ReplaceAll(localPart, ".", ""), "-", "")
	
	// Se for curto, usa como está
	if len(cleaned) <= 4 {
		return cleaned
	}
	
	// Tenta gerar um apelido inteligente
	// Pega as primeiras letras de cada "palavra" separada por pontos ou hífens
	words := strings.FieldsFunc(localPart, func(c rune) bool {
		return c == '.' || c == '-' || c == '_'
	})
	
	if len(words) > 1 {
		// Pega as primeiras 2 letras de cada palavra
		var nickname strings.Builder
		for _, word := range words {
			if len(word) >= 2 {
				nickname.WriteString(word[:2])
			} else if len(word) == 1 {
				nickname.WriteString(word)
			}
			
			// Limita a 8 caracteres para o apelido
			if nickname.Len() >= 8 {
				break
			}
		}
		
		if nickname.Len() > 0 {
			return nickname.String()
		}
	}
	
	// Se não conseguir gerar apelido inteligente, usa as primeiras 6 letras
	if len(cleaned) > 6 {
		return cleaned[:6]
	}
	
	return cleaned
}


