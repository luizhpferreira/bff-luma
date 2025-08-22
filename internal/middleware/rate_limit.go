package middleware

import (
	"net/http"
	"strings"

	"bff-luma/internal/services"
)

// RateLimitMiddleware cria um middleware de rate limiting
func RateLimitMiddleware(rateLimiter *services.RateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extrai IP do cliente
			ip := getClientIP(r)

			// Verifica rate limit por IP
			if !rateLimiter.AllowIP(ip) {
				http.Error(w, "Rate limit excedido. Tente novamente em alguns minutos.", http.StatusTooManyRequests)
				return
			}

			// Chama o próximo handler
			next.ServeHTTP(w, r)
		})
	}
}

// LoginRateLimitMiddleware cria um middleware específico para login
func LoginRateLimitMiddleware(rateLimiter *services.RateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extrai IP do cliente
			ip := getClientIP(r)

			// Verifica rate limit por IP
			if !rateLimiter.AllowIP(ip) {
				http.Error(w, "Rate limit excedido. Tente novamente em alguns minutos.", http.StatusTooManyRequests)
				return
			}

			// Para login, também verificamos o email (será feito no handler)
			next.ServeHTTP(w, r)
		})
	}
}

// PasswordResetRateLimitMiddleware cria um middleware específico para reset de senha
func PasswordResetRateLimitMiddleware(rateLimiter *services.RateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extrai IP do cliente
			ip := getClientIP(r)

			// Verifica rate limit por IP
			if !rateLimiter.AllowIP(ip) {
				http.Error(w, "Rate limit excedido. Tente novamente em alguns minutos.", http.StatusTooManyRequests)
				return
			}

			// Para reset de senha, também verificamos o email (será feito no handler)
			next.ServeHTTP(w, r)
		})
	}
}

// getClientIP extrai o IP real do cliente
func getClientIP(r *http.Request) string {
	// Verifica headers de proxy
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		// Pega o primeiro IP da lista
		if commaIndex := strings.Index(ip, ","); commaIndex != -1 {
			return strings.TrimSpace(ip[:commaIndex])
		}
		return ip
	}

	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}

	if ip := r.Header.Get("X-Client-IP"); ip != "" {
		return ip
	}

	// Remove a porta do endereço remoto
	remoteAddr := r.RemoteAddr
	if colonIndex := strings.LastIndex(remoteAddr, ":"); colonIndex != -1 {
		return remoteAddr[:colonIndex]
	}

	return remoteAddr
}
