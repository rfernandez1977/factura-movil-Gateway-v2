package utils

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims define la estructura de los claims del token JWT
type Claims struct {
	UserID string `json:"user_id"`
	Rut    string `json:"rut"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// TokenClaims contiene las claims para un token JWT
type TokenClaims struct {
	UserID string `json:"user_id"`
	Rut    string `json:"rut"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// JWTUtils contiene utilidades para JWT
type JWTUtils struct {
	secretKey    string
	issuer       string
	expiration   time.Duration
	refreshToken time.Duration
}

// NewJWTUtils crea una nueva instancia de JWTUtils
func NewJWTUtils() *JWTUtils {
	return &JWTUtils{
		secretKey:    getEnv("JWT_SECRET_KEY", "default_secret_key"),
		issuer:       getEnv("JWT_ISSUER", "fmgo"),
		expiration:   24 * time.Hour,
		refreshToken: 7 * 24 * time.Hour,
	}
}

// getEnv obtiene una variable de entorno o retorna un valor por defecto
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// GenerateToken genera un token JWT
func (j *JWTUtils) GenerateToken(userID, rut, role string) (string, error) {
	claims := TokenClaims{
		UserID: userID,
		Rut:    rut,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.expiration)),
			Issuer:    j.issuer,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(j.secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken valida un token JWT
func (j *JWTUtils) ValidateToken(tokenString string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("token inválido")
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok {
		return nil, errors.New("no se pudo obtener las claims")
	}

	return claims, nil
}

// RefreshToken refresca un token JWT
func (j *JWTUtils) RefreshToken(tokenString string) (string, error) {
	claims, err := j.ValidateToken(tokenString)
	if err != nil {
		return "", err
	}

	// Actualizar fechas en RegisteredClaims
	claims.RegisteredClaims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(j.expiration))
	claims.RegisteredClaims.IssuedAt = jwt.NewNumericDate(time.Now())

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	newTokenString, err := token.SignedString([]byte(j.secretKey))
	if err != nil {
		return "", err
	}

	return newTokenString, nil
}

// HasRole verifica si el token tiene un rol específico
func HasRole(tokenString, role string, j *JWTUtils) (bool, error) {
	claims, err := j.ValidateToken(tokenString)
	if err != nil {
		return false, err
	}
	return claims.Role == role, nil
}

// GetUserID obtiene el ID de usuario del token
func GetUserID(tokenString string, j *JWTUtils) (string, error) {
	claims, err := j.ValidateToken(tokenString)
	if err != nil {
		return "", err
	}
	return claims.UserID, nil
}

// GetRut obtiene el RUT del token
func GetRut(tokenString string, j *JWTUtils) (string, error) {
	claims, err := j.ValidateToken(tokenString)
	if err != nil {
		return "", err
	}
	return claims.Rut, nil
}
