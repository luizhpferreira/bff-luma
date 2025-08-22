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
	LNBitsAdminKey    string
	LNBitsWebhookSecret string
	DatabasePath      string
}

// LoadConfig carrega as configurações do arquivo .env
func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Arquivo .env não encontrado, usando variáveis de ambiente do sistema")
	}

	return &Config{
		AppPort:           getEnv("APP_PORT", "8080"),
		JWTSecret:         getEnv("JWT_SECRET", "supersecreto123456789"),
		LNBitsBaseURL:     getEnv("LNBITS_BASE_URL", "http://127.0.0.1:5000"),
		LNBitsAdminKey:    getEnv("LNBITS_ADMIN_KEY", ""),
		LNBitsWebhookSecret: getEnv("LNBITS_WEBHOOK_SECRET", ""),
		DatabasePath:      getEnv("DATABASE_PATH", "./bff_luma.db"),
	}
}

// getEnv obtém uma variável de ambiente ou retorna um valor padrão
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
