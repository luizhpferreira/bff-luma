package services

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// PasswordService gerencia operações relacionadas a senhas
type PasswordService struct {
	// Custo do bcrypt (padrão: 12)
	cost int
}

// NewPasswordService cria um novo serviço de senhas
func NewPasswordService() *PasswordService {
	return &PasswordService{
		cost: bcrypt.DefaultCost, // 12
	}
}

// HashPassword gera um hash bcrypt da senha
func (s *PasswordService) HashPassword(password string) (string, error) {
	// Gera o hash da senha
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), s.cost)
	if err != nil {
		return "", fmt.Errorf("erro ao gerar hash da senha: %w", err)
	}

	return string(hashedBytes), nil
}

// CheckPassword verifica se uma senha corresponde ao hash
func (s *PasswordService) CheckPassword(password, hashedPassword string) error {
	// Compara a senha com o hash
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return fmt.Errorf("senha incorreta")
	}

	return nil
}

// ValidatePasswordStrength valida se a senha atende aos requisitos de segurança
func (s *PasswordService) ValidatePasswordStrength(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("senha deve ter pelo menos 8 caracteres")
	}

	// Verifica se contém pelo menos uma letra maiúscula
	hasUpper := false
	hasLower := false
	hasDigit := false
	hasSpecial := false

	for _, char := range password {
		switch {
		case char >= 'A' && char <= 'Z':
			hasUpper = true
		case char >= 'a' && char <= 'z':
			hasLower = true
		case char >= '0' && char <= '9':
			hasDigit = true
		case char == '!' || char == '@' || char == '#' || char == '$' || char == '%' || char == '^' || char == '&' || char == '*' || char == '(' || char == ')' || char == '-' || char == '_' || char == '+' || char == '=' || char == '[' || char == ']' || char == '{' || char == '}' || char == '|' || char == '\\' || char == ':' || char == ';' || char == '"' || char == '\'' || char == '<' || char == '>' || char == ',' || char == '.' || char == '?' || char == '/':
			hasSpecial = true
		}
	}

	if !hasUpper {
		return fmt.Errorf("senha deve conter pelo menos uma letra maiúscula")
	}

	if !hasLower {
		return fmt.Errorf("senha deve conter pelo menos uma letra minúscula")
	}

	if !hasDigit {
		return fmt.Errorf("senha deve conter pelo menos um número")
	}

	if !hasSpecial {
		return fmt.Errorf("senha deve conter pelo menos um caractere especial")
	}

	// Verifica sequências comuns
	commonSequences := []string{"123", "abc", "qwe", "asd", "zxc", "password", "senha", "admin", "user"}
	for _, seq := range commonSequences {
		if len(password) >= len(seq) {
			for i := 0; i <= len(password)-len(seq); i++ {
				if password[i:i+len(seq)] == seq {
					return fmt.Errorf("senha fraca: não pode conter sequências comuns")
				}
			}
		}
	}

	// Verifica caracteres repetidos (mais de 2 consecutivos)
	for i := 0; i < len(password)-2; i++ {
		if password[i] == password[i+1] && password[i] == password[i+2] {
			return fmt.Errorf("senha fraca: não pode conter 3 ou mais caracteres idênticos consecutivos")
		}
	}

	return nil
}
