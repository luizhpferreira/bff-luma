package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config representa as configurações da aplicação
type Config struct {
	AppPort           string
	JWTSecret         string
	LNBitsBaseURL     string
	LNBitsAPIToken    string
	LNBitsWebhookSecret string
	DatabasePath      string
	// Configurações SMTP
	SMTPHost          string
	SMTPPort          string
	SMTPUsername      string
	SMTPPassword      string
	SMTPFromEmail     string
	SMTPFromName      string
	SMTPUseTLS        bool
	// Configurações de Domínio
	AppDomain         string
	AppProtocol       string
}

// LoadConfig carrega as configurações do arquivo .env
func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Arquivo .env não encontrado, usando variáveis de ambiente do sistema")
	}
	
	// Debug: verificar se as variáveis LNBITS_POSTGRES estão carregadas
	log.Printf("🔍 Debug: LNBITS_POSTGRES_HOST=%s", os.Getenv("LNBITS_POSTGRES_HOST"))
	log.Printf("🔍 Debug: LNBITS_POSTGRES_USER=%s", os.Getenv("LNBITS_POSTGRES_USER"))
	log.Printf("🔍 Debug: LNBITS_POSTGRES_PASSWORD=%s", os.Getenv("LNBITS_POSTGRES_PASSWORD"))
	log.Printf("🔍 Debug: LNBITS_POSTGRES_DB=%s", os.Getenv("LNBITS_POSTGRES_DB"))
	log.Printf("🔍 Debug: LNBITS_POSTGRES_PORT=%s", os.Getenv("LNBITS_POSTGRES_PORT"))

	return &Config{
		AppPort:           getEnv("APP_PORT", "8080"),
		JWTSecret:         getEnv("JWT_SECRET", ""),
		LNBitsBaseURL:     getEnv("LNBITS_BASE_URL", "http://127.0.0.1:5000"),
		LNBitsAPIToken:    getEnv("LNBITS_API_TOKEN", ""),
		LNBitsWebhookSecret: getEnv("LNBITS_WEBHOOK_SECRET", ""),
		DatabasePath:      getEnv("DATABASE_URL", "postgresql://bff_luma:bff_luma@localhost:5432/bff_luma?sslmode=disable"),
		// Configurações SMTP
		SMTPHost:          getEnv("SMTP_HOST", ""),
		SMTPPort:          getEnv("SMTP_PORT", "587"),
		SMTPUsername:      getEnv("SMTP_USERNAME", ""),
		SMTPPassword:      getEnv("SMTP_PASSWORD", ""),
		SMTPFromEmail:     getEnv("SMTP_FROM_EMAIL", "noreply@bff-luma.com"),
		SMTPFromName:      getEnv("SMTP_FROM_NAME", "BFF Luma"),
		SMTPUseTLS:        getEnv("SMTP_USE_TLS", "true") == "true",
		// Configurações de Domínio
		AppDomain:         getEnv("APP_DOMAIN", "localhost"),
		AppProtocol:       getEnv("APP_PROTOCOL", "http"),
	}
}

// getEnv obtém uma variável de ambiente ou retorna um valor padrão
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
