package services

import (
	"log"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// RateLimiter gerencia limites de taxa para diferentes operações
type RateLimiter struct {
	// Limitadores por IP
	ipLimiters map[string]*rate.Limiter
	// Limitadores por email (para login)
	emailLimiters map[string]*rate.Limiter
	// Mutex para thread safety
	mu sync.RWMutex
	// Configurações
	config RateLimitConfig
}

// RateLimitConfig configurações de rate limiting
type RateLimitConfig struct {
	// Login: 5 tentativas por 15 minutos por email
	LoginAttemptsPerEmail int
	LoginWindowDuration   time.Duration
	
	// IP: 100 requisições por minuto por IP
	RequestsPerIP int
	IPWindowDuration time.Duration
	
	// Recuperação de senha: 3 tentativas por hora por email
	PasswordResetPerEmail int
	PasswordResetWindowDuration time.Duration
	
	// Limpeza de limitadores antigos a cada 1 hora
	CleanupInterval time.Duration
}

// NewRateLimiter cria um novo rate limiter
func NewRateLimiter() *RateLimiter {
	config := RateLimitConfig{
		LoginAttemptsPerEmail:      5,
		LoginWindowDuration:        15 * time.Minute,
		RequestsPerIP:              100,
		IPWindowDuration:           1 * time.Minute,
		PasswordResetPerEmail:      3,
		PasswordResetWindowDuration: 1 * time.Hour,
		CleanupInterval:            1 * time.Hour,
	}

	rl := &RateLimiter{
		ipLimiters:    make(map[string]*rate.Limiter),
		emailLimiters: make(map[string]*rate.Limiter),
		config:        config,
	}

	// Inicia limpeza automática
	go rl.startCleanup()

	return rl
}

// startCleanup inicia a limpeza automática de limitadores antigos
func (rl *RateLimiter) startCleanup() {
	ticker := time.NewTicker(rl.config.CleanupInterval)
	defer ticker.Stop()

	for range ticker.C {
		rl.cleanup()
	}
}

// cleanup remove limitadores antigos para economizar memória
func (rl *RateLimiter) cleanup() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	// Por simplicidade, vamos limpar todos os limitadores a cada hora
	// Em produção, você poderia rastrear timestamps de último uso
	if len(rl.ipLimiters) > 1000 {
		rl.ipLimiters = make(map[string]*rate.Limiter)
		log.Println("🧹 Rate limiter: limitadores de IP limpos")
	}
	
	if len(rl.emailLimiters) > 1000 {
		rl.emailLimiters = make(map[string]*rate.Limiter)
		log.Println("🧹 Rate limiter: limitadores de email limpos")
	}
}

// getIPLimiter obtém ou cria um limitador para um IP
func (rl *RateLimiter) getIPLimiter(ip string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	limiter, exists := rl.ipLimiters[ip]
	if !exists {
		limiter = rate.NewLimiter(rate.Every(rl.config.IPWindowDuration/time.Duration(rl.config.RequestsPerIP)), rl.config.RequestsPerIP)
		rl.ipLimiters[ip] = limiter
	}

	return limiter
}

// getEmailLimiter obtém ou cria um limitador para um email
func (rl *RateLimiter) getEmailLimiter(email string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	limiter, exists := rl.emailLimiters[email]
	if !exists {
		limiter = rate.NewLimiter(rate.Every(rl.config.LoginWindowDuration/time.Duration(rl.config.LoginAttemptsPerEmail)), rl.config.LoginAttemptsPerEmail)
		rl.emailLimiters[email] = limiter
	}

	return limiter
}

// getPasswordResetLimiter obtém ou cria um limitador para reset de senha
func (rl *RateLimiter) getPasswordResetLimiter(email string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	key := "reset_" + email
	limiter, exists := rl.emailLimiters[key]
	if !exists {
		limiter = rate.NewLimiter(rate.Every(rl.config.PasswordResetWindowDuration/time.Duration(rl.config.PasswordResetPerEmail)), rl.config.PasswordResetPerEmail)
		rl.emailLimiters[key] = limiter
	}

	return limiter
}

// AllowIP verifica se um IP pode fazer uma requisição
func (rl *RateLimiter) AllowIP(ip string) bool {
	limiter := rl.getIPLimiter(ip)
	return limiter.Allow()
}

// AllowLogin verifica se um email pode tentar login
func (rl *RateLimiter) AllowLogin(email string) bool {
	limiter := rl.getEmailLimiter(email)
	return limiter.Allow()
}

// AllowPasswordReset verifica se um email pode solicitar reset de senha
func (rl *RateLimiter) AllowPasswordReset(email string) bool {
	limiter := rl.getPasswordResetLimiter(email)
	return limiter.Allow()
}

// GetRemainingAttempts retorna tentativas restantes para login
func (rl *RateLimiter) GetRemainingAttempts(email string) int {
	limiter := rl.getEmailLimiter(email)
	return int(limiter.TokensAt(time.Now()))
}

// GetRemainingPasswordResets retorna tentativas restantes para reset de senha
func (rl *RateLimiter) GetRemainingPasswordResets(email string) int {
	limiter := rl.getPasswordResetLimiter(email)
	return int(limiter.TokensAt(time.Now()))
}

// GetStats retorna estatísticas do rate limiter
func (rl *RateLimiter) GetStats() map[string]interface{} {
	rl.mu.RLock()
	defer rl.mu.RUnlock()

	return map[string]interface{}{
		"ip_limiters_count":     len(rl.ipLimiters),
		"email_limiters_count":  len(rl.emailLimiters),
		"login_attempts_limit":  rl.config.LoginAttemptsPerEmail,
		"login_window":          rl.config.LoginWindowDuration.String(),
		"ip_requests_limit":     rl.config.RequestsPerIP,
		"ip_window":             rl.config.IPWindowDuration.String(),
		"reset_attempts_limit":  rl.config.PasswordResetPerEmail,
		"reset_window":          rl.config.PasswordResetWindowDuration.String(),
	}
}
