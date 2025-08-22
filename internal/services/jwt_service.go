package services

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTService gerencia tokens JWT
type JWTService struct {
	secretKey []byte
}

// Claims representa as claims do JWT
type Claims struct {
	Email    string `json:"email"`
	WalletID string `json:"wallet_id"`
	jwt.RegisteredClaims
}

// NewJWTService cria um novo serviço JWT
func NewJWTService(secretKey string) *JWTService {
	return &JWTService{
		secretKey: []byte(secretKey),
	}
}

// GenerateToken gera um novo token JWT
func (s *JWTService) GenerateToken(email, walletID string) (string, error) {
	// Token expira em 24 horas
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &Claims{
		Email:    email,
		WalletID: walletID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "bff-luma",
			Subject:   email,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.secretKey)
	if err != nil {
		return "", fmt.Errorf("erro ao gerar token: %w", err)
	}

	return tokenString, nil
}

// ValidateToken valida um token JWT
func (s *JWTService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Verifica se o método de assinatura é correto
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("método de assinatura inesperado: %v", token.Header["alg"])
		}
		return s.secretKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("erro ao validar token: %w", err)
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("token inválido")
}

// RefreshToken gera um novo token com base em um token válido
func (s *JWTService) RefreshToken(tokenString string) (string, error) {
	claims, err := s.ValidateToken(tokenString)
	if err != nil {
		return "", fmt.Errorf("erro ao validar token para refresh: %w", err)
	}

	// Gera um novo token com as mesmas claims
	return s.GenerateToken(claims.Email, claims.WalletID)
}
