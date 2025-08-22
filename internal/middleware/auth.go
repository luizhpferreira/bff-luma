package middleware

import (
	"context"
	"net/http"
	"strings"

	"bff-luma/internal/services"
)

// AuthMiddleware cria um middleware de autenticação JWT
func AuthMiddleware(jwtService *services.JWTService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extrai o token do header Authorization
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Token de autorização não fornecido", http.StatusUnauthorized)
				return
			}

			// Verifica se o header tem o formato "Bearer <token>"
			tokenParts := strings.Split(authHeader, " ")
			if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
				http.Error(w, "Formato de autorização inválido. Use: Bearer <token>", http.StatusUnauthorized)
				return
			}

			tokenString := tokenParts[1]

			// Valida o token
			claims, err := jwtService.ValidateToken(tokenString)
			if err != nil {
				http.Error(w, "Token inválido: "+err.Error(), http.StatusUnauthorized)
				return
			}

			// Adiciona as claims ao contexto da requisição
			ctx := r.Context()
			ctx = context.WithValue(ctx, "user_email", claims.Email)
			ctx = context.WithValue(ctx, "wallet_id", claims.WalletID)
			ctx = context.WithValue(ctx, "claims", claims)

			// Chama o próximo handler com o contexto atualizado
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserEmail extrai o email do usuário do contexto
func GetUserEmail(r *http.Request) string {
	if email, ok := r.Context().Value("user_email").(string); ok {
		return email
	}
	return ""
}

// GetWalletID extrai o wallet_id do contexto
func GetWalletID(r *http.Request) string {
	if walletID, ok := r.Context().Value("wallet_id").(string); ok {
		return walletID
	}
	return ""
}

// GetClaims extrai as claims completas do contexto
func GetClaims(r *http.Request) *services.Claims {
	if claims, ok := r.Context().Value("claims").(*services.Claims); ok {
		return claims
	}
	return nil
}
