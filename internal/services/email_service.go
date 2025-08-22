package services

import (
	"log"
)

// EmailService representa o serviço de envio de emails
type EmailService struct {
	// Em produção, aqui teríamos configurações de SMTP
	// smtpHost, smtpPort, username, password, etc.
}

// NewEmailService cria um novo serviço de email
func NewEmailService() *EmailService {
	return &EmailService{}
}

// SendPasswordResetEmail envia email de reset de senha
func (s *EmailService) SendPasswordResetEmail(email, token string) error {
	// Em produção, aqui enviaríamos um email real
	// Por enquanto, apenas logamos o token para facilitar os testes
	
	log.Printf("📧 Email de reset de senha enviado para: %s", email)
	log.Printf("🔑 Token de reset: %s", token)
	log.Printf("🌐 Link de reset: http://localhost:3000/reset-password?token=%s", token)
	
	return nil
}

// SendWelcomeEmail envia email de boas-vindas
func (s *EmailService) SendWelcomeEmail(email, walletID string) error {
	log.Printf("📧 Email de boas-vindas enviado para: %s", email)
	log.Printf("💳 Wallet ID: %s", walletID)
	
	return nil
}

// ValidateEmail valida formato de email (simples)
func (s *EmailService) ValidateEmail(email string) bool {
	// Validação básica de email
	if len(email) < 5 {
		return false
	}
	
	// Verifica se contém @
	hasAt := false
	for _, char := range email {
		if char == '@' {
			hasAt = true
			break
		}
	}
	
	return hasAt
}
