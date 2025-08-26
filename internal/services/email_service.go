package services

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/smtp"
)

// EmailService representa o serviço de envio de emails
type EmailService struct {
	host     string
	port     string
	username string
	password string
	fromEmail string
	fromName  string
	useTLS   bool
	enabled  bool
	appDomain string
	appProtocol string
}

// NewEmailService cria um novo serviço de email
func NewEmailService(host, port, username, password, fromEmail, fromName string, useTLS bool, appDomain, appProtocol string) *EmailService {
	enabled := host != "" && username != "" && password != ""
	
	if enabled {
		log.Printf("📧 Email Service: Configurado para %s:%s", host, port)
	} else {
		log.Printf("📧 Email Service: Modo simulado (SMTP não configurado)")
	}
	
	return &EmailService{
		host:      host,
		port:      port,
		username:  username,
		password:  password,
		fromEmail: fromEmail,
		fromName:  fromName,
		useTLS:    useTLS,
		enabled:   enabled,
		appDomain: appDomain,
		appProtocol: appProtocol,
	}
}

// SendPasswordResetEmail envia email de reset de senha
func (s *EmailService) SendPasswordResetEmail(email, token string) error {
	if !s.enabled {
		// Modo simulado
		log.Printf("📧 Email de reset de senha enviado para: %s", email)
		log.Printf("🔑 Token de reset: %s", token)
		log.Printf("🌐 Link de reset: %s://%s/reset-password?token=%s", s.appProtocol, s.appDomain, token)
		return nil
	}

	subject := "Recuperação de Senha - Luma"
	body := s.buildPasswordResetEmailBody(email, token)
	
	return s.sendEmail(email, subject, body)
}

// SendWelcomeEmail envia email de boas-vindas
func (s *EmailService) SendWelcomeEmail(email, walletID string) error {
	if !s.enabled {
		// Modo simulado
		log.Printf("📧 Email de boas-vindas enviado para: %s", email)
		log.Printf("💳 Wallet ID: %s", walletID)
		return nil
	}

	subject := "Bem-vindo a Luma!"
	body := s.buildWelcomeEmailBody(email, walletID)
	
	return s.sendEmail(email, subject, body)
}

// SendEmailConfirmation envia email de confirmação
func (s *EmailService) SendEmailConfirmation(email, token string) error {
	if !s.enabled {
		// Modo simulado
		log.Printf("📧 Email de confirmação enviado para: %s", email)
		log.Printf("🔑 Token de confirmação: %s", token)
		log.Printf("🌐 Link de confirmação: %s://%s/confirm-email?token=%s", s.appProtocol, s.appDomain, token)
		return nil
	}

	subject := "Confirme seu email - BFF Luma"
	body := s.buildEmailConfirmationBody(email, token)
	
	return s.sendEmail(email, subject, body)
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

// sendEmail envia um email via SMTP
func (s *EmailService) sendEmail(to, subject, body string) error {
	// Configuração do cabeçalho do email
	headers := make(map[string]string)
	headers["From"] = fmt.Sprintf("%s <%s>", s.fromName, s.fromEmail)
	headers["To"] = to
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=UTF-8"

	// Constrói a mensagem
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	// Configuração da autenticação
	auth := smtp.PlainAuth("", s.username, s.password, s.host)

	// Endereço do servidor SMTP
	addr := fmt.Sprintf("%s:%s", s.host, s.port)

	// Envia o email
	if s.useTLS {
		// Detecta automaticamente se deve usar TLS direto ou STARTTLS baseado na porta
		if s.port == "465" {
			// Porta 465: TLS direto (SSL)
			tlsConfig := &tls.Config{
				ServerName: s.host,
			}
			
			conn, err := tls.Dial("tcp", addr, tlsConfig)
			if err != nil {
				return fmt.Errorf("erro ao conectar com TLS direto: %w", err)
			}
			defer conn.Close()

			client, err := smtp.NewClient(conn, s.host)
			if err != nil {
				return fmt.Errorf("erro ao criar cliente SMTP: %w", err)
			}
			defer client.Close()

			if err = client.Auth(auth); err != nil {
				return fmt.Errorf("erro na autenticação: %w", err)
			}

			if err = client.Mail(s.fromEmail); err != nil {
				return fmt.Errorf("erro ao definir remetente: %w", err)
			}

			if err = client.Rcpt(to); err != nil {
				return fmt.Errorf("erro ao definir destinatário: %w", err)
			}

			w, err := client.Data()
			if err != nil {
				return fmt.Errorf("erro ao iniciar dados: %w", err)
			}
			defer w.Close()

			_, err = w.Write([]byte(message))
			if err != nil {
				return fmt.Errorf("erro ao escrever mensagem: %w", err)
			}
		} else {
			// Porta 587 (ou outras): STARTTLS
			client, err := smtp.Dial(addr)
			if err != nil {
				return fmt.Errorf("erro ao conectar com servidor SMTP: %w", err)
			}
			defer client.Close()

			// Inicia STARTTLS
			if err = client.StartTLS(&tls.Config{ServerName: s.host}); err != nil {
				return fmt.Errorf("erro ao iniciar STARTTLS: %w", err)
			}

			// Autentica após STARTTLS
			if err = client.Auth(auth); err != nil {
				return fmt.Errorf("erro na autenticação: %w", err)
			}

			// Define remetente
			if err = client.Mail(s.fromEmail); err != nil {
				return fmt.Errorf("erro ao definir remetente: %w", err)
			}

			// Define destinatário
			if err = client.Rcpt(to); err != nil {
				return fmt.Errorf("erro ao definir destinatário: %w", err)
			}

			// Envia dados
			w, err := client.Data()
			if err != nil {
				return fmt.Errorf("erro ao iniciar dados: %w", err)
			}
			defer w.Close()

			_, err = w.Write([]byte(message))
			if err != nil {
				return fmt.Errorf("erro ao escrever mensagem: %w", err)
			}
		}
	} else {
		// Sem TLS (não recomendado para produção)
		err := smtp.SendMail(addr, auth, s.fromEmail, []string{to}, []byte(message))
		if err != nil {
			return fmt.Errorf("erro ao enviar email: %w", err)
		}
	}

	log.Printf("📧 Email enviado com sucesso para: %s", to)
	return nil
}

// buildPasswordResetEmailBody constrói o corpo do email de reset de senha
func (s *EmailService) buildPasswordResetEmailBody(email, token string) string {
	resetLink := fmt.Sprintf("%s://%s/reset-password?token=%s", s.appProtocol, s.appDomain, token)
	
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Recuperação de Senha</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background: #007bff; color: white; padding: 20px; text-align: center; }
        .content { padding: 20px; background: #f8f9fa; }
        .button { display: inline-block; padding: 12px 24px; background: #007bff; color: white; text-decoration: none; border-radius: 5px; }
        .footer { text-align: center; padding: 20px; color: #666; font-size: 12px; }
        .warning { background: #fff3cd; border: 1px solid #ffeaa7; padding: 15px; border-radius: 5px; margin: 20px 0; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🔐 Recuperação de Senha</h1>
        </div>
        <div class="content">
            <p>Olá!</p>
            <p>Recebemos uma solicitação para redefinir a senha da sua conta na <strong>Luma</strong>.</p>
            <p>Se você não fez essa solicitação, pode ignorar este email.</p>
            
            <div style="text-align: center; margin: 30px 0;">
                <a href="%s" class="button">🔑 Redefinir Senha</a>
            </div>
            
            <div class="warning">
                <strong>⚠️ Importante:</strong>
                <ul>
                    <li>Este link expira em <strong>1 hora</strong></li>
                    <li>Use apenas em dispositivos confiáveis</li>
                    <li>Não compartilhe este link com ninguém</li>
                </ul>
            </div>
            
            <p>Se o botão não funcionar, copie e cole este link no seu navegador:</p>
            <p style="word-break: break-all; background: #f1f1f1; padding: 10px; border-radius: 3px;">%s</p>
        </div>
        <div class="footer">
            <p>Este email foi enviado automaticamente. Não responda a este email.</p>
            <p>&copy; 2024 BFF Luma. Todos os direitos reservados.</p>
        </div>
    </div>
</body>
</html>`, resetLink, resetLink)
}

// buildWelcomeEmailBody constrói o corpo do email de boas-vindas
func (s *EmailService) buildWelcomeEmailBody(email, walletID string) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Bem-vindo a Luma!</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background: #28a745; color: white; padding: 20px; text-align: center; }
        .content { padding: 20px; background: #f8f9fa; }
        .wallet-info { background: #e9ecef; padding: 15px; border-radius: 5px; margin: 20px 0; }
        .footer { text-align: center; padding: 20px; color: #666; font-size: 12px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🎉 Bem-vindo a Luma!</h1>
        </div>
        <div class="content">
            <p>Olá!</p>
            <p>🎉 <strong>Parabéns!</strong> Sua conta foi confirmada com sucesso no <strong>BFF Luma</strong>!</p>
            <p>✅ Seu email foi verificado e sua conta está ativa.</p>
            <p>Agora você pode:</p>
            <ul>
                <li>✅ Fazer login na sua conta</li>
                <li>💳 Gerenciar sua carteira Lightning</li>
                <li>💰 Criar invoices para receber pagamentos</li>
                <li>📊 Acompanhar o status dos seus pagamentos</li>
            </ul>
            
            <div class="wallet-info">
                <h3>📋 Informações da sua Carteira:</h3>
                <p><strong>Email:</strong> %s</p>
                <p><strong>Wallet ID:</strong> <code>%s</code></p>
            </div>
            
            <p><strong>🔐 Dica de Segurança:</strong></p>
            <ul>
                <li>Mantenha sua senha segura</li>
                <li>Não compartilhe suas credenciais</li>
                <li>Use sempre dispositivos confiáveis</li>
            </ul>
            
            <p>Se você tiver alguma dúvida, entre em contato conosco.</p>
        </div>
        <div class="footer">
            <p>Obrigado por escolher o BFF Luma!</p>
            <p>&copy; 2024 BFF Luma. Todos os direitos reservados.</p>
        </div>
    </div>
</body>
</html>`, email, walletID)
}

// buildEmailConfirmationBody constrói o corpo do email de confirmação
func (s *EmailService) buildEmailConfirmationBody(email, token string) string {
	confirmationLink := fmt.Sprintf("%s://%s/confirm-email?token=%s", s.appProtocol, s.appDomain, token)
	
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Confirme seu Email</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background: #007bff; color: white; padding: 20px; text-align: center; }
        .content { padding: 20px; background: #f8f9fa; }
        .button { display: inline-block; padding: 12px 24px; background: #007bff; color: white; text-decoration: none; border-radius: 5px; }
        .footer { text-align: center; padding: 20px; color: #666; font-size: 12px; }
        .warning { background: #fff3cd; border: 1px solid #ffeaa7; padding: 15px; border-radius: 5px; margin: 20px 0; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>📧 Confirme seu Email</h1>
        </div>
        <div class="content">
            <p>Olá!</p>
            <p>Para completar o cadastro da sua conta no <strong>BFF Luma</strong>, por favor confirme seu email clicando no botão abaixo:</p>
            <p><strong>Após a confirmação, você receberá um email de boas-vindas com mais informações sobre sua conta.</strong></p>
            
            <div style="text-align: center; margin: 30px 0;">
                <a href="%s" class="button">✅ Confirmar Email</a>
            </div>
            
            <div class="warning">
                <strong>⚠️ Importante:</strong>
                <ul>
                    <li>Este link expira em <strong>24 horas</strong></li>
                    <li>Use apenas em dispositivos confiáveis</li>
                    <li>Não compartilhe este link com ninguém</li>
                </ul>
            </div>
            
            <p>Se o botão não funcionar, copie e cole este link no seu navegador:</p>
            <p style="word-break: break-all; background: #f1f1f1; padding: 10px; border-radius: 3px;">%s</p>
            
            <p>Se você não criou uma conta no BFF Luma, pode ignorar este email.</p>
        </div>
        <div class="footer">
            <p>Este email foi enviado automaticamente. Não responda a este email.</p>
            <p>&copy; 2024 BFF Luma. Todos os direitos reservados.</p>
        </div>
    </div>
</body>
</html>`, confirmationLink, confirmationLink)
}
