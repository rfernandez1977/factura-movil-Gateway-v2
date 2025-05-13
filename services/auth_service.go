package services

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/cursor/FMgo/config"
)

// AuthService maneja la autenticación y autorización
type AuthService struct {
	config *config.SupabaseConfig
}

// NewAuthService crea una nueva instancia del servicio de autenticación
func NewAuthService(config *config.SupabaseConfig) *AuthService {
	return &AuthService{
		config: config,
	}
}

// Claims representa los claims del JWT
type Claims struct {
	Aud      string `json:"aud"`
	Role     string `json:"role"`
	Email    string `json:"email"`
	Sub      string `json:"sub"`
	UserData struct {
		EmailVerified bool `json:"email_verified"`
	} `json:"user_metadata"`
	jwt.RegisteredClaims
}

// ValidarToken valida un token JWT
func (s *AuthService) ValidarToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("método de firma inesperado: %v", token.Header["alg"])
		}
		return []byte(s.config.JWTSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("error al validar token: %v", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("token inválido")
	}

	return claims, nil
}

// GenerarToken genera un nuevo token JWT
func (s *AuthService) GenerarToken(email string, userID string) (string, error) {
	claims := &Claims{
		Aud:   "authenticated",
		Role:  "authenticated",
		Email: email,
		Sub:   userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.config.JWTSecret))
}

// VerificarSesion verifica si una sesión es válida
func (s *AuthService) VerificarSesion(ctx context.Context, token string) (bool, error) {
	claims, err := s.ValidarToken(token)
	if err != nil {
		return false, err
	}

	// Verificar que el token no haya expirado
	if claims.ExpiresAt.Before(time.Now()) {
		return false, fmt.Errorf("token expirado")
	}

	// Verificar que el usuario esté autenticado
	if claims.Role != "authenticated" {
		return false, fmt.Errorf("usuario no autenticado")
	}

	return true, nil
}

// ObtenerUsuario obtiene la información del usuario desde el token
func (s *AuthService) ObtenerUsuario(token string) (*Claims, error) {
	return s.ValidarToken(token)
}
