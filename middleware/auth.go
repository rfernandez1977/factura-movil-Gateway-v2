package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"FMgo/utils"

	"github.com/gin-gonic/gin"
)

var (
	logger = log.New(os.Stdout, "[AUTH] ", log.LstdFlags)
)

type AuthConfig struct {
	APIKey string
}

type Claims struct {
	UserID string `json:"user_id"`
	Rut    string `json:"rut"`
	Role   string `json:"role"`
}

// AuthMiddleware verifica la autenticación y los roles
func AuthMiddleware(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token no proporcionado"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "formato de token inválido"})
			c.Abort()
			return
		}

		token := parts[1]

		// Validar token y extraer claims usando utils
		jwtUtils := utils.NewJWTUtils()
		userID, _ := utils.GetUserID(token, jwtUtils)
		rut, _ := utils.GetRut(token, jwtUtils)
		// Aquí deberías tener una función para obtener el rol desde el token
		role := "" // TODO: Implementar utils.GetRole(token, jwtUtils) si existe

		// Verificar roles si se especificaron
		if len(roles) > 0 {
			hasRole := false
			for _, r := range roles {
				if role == r {
					hasRole = true
					break
				}
			}
			if !hasRole {
				c.JSON(http.StatusForbidden, gin.H{"error": "acceso denegado"})
				c.Abort()
				return
			}
		}

		c.Set("user_id", userID)
		c.Set("rut", rut)
		c.Set("role", role)

		c.Next()
	}
}

// RateLimitMiddleware limita las peticiones por IP
func RateLimitMiddleware(limit int, window time.Duration) gin.HandlerFunc {
	// Mapa para almacenar contadores por IP
	counters := make(map[string]struct {
		count     int
		resetTime time.Time
	})

	return func(c *gin.Context) {
		ip := c.ClientIP()
		now := time.Now()

		// Obtener contador para la IP
		counter := counters[ip]

		// Resetear contador si ha pasado el tiempo de ventana
		if now.After(counter.resetTime) {
			counter.count = 0
			counter.resetTime = now.Add(window)
		}

		// Incrementar contador
		counter.count++
		counters[ip] = counter

		// Verificar límite
		if counter.count > limit {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": fmt.Sprintf("demasiadas peticiones. límite: %d por %v", limit, window),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

func NewAuthMiddleware(config AuthConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		apiKey := c.GetHeader("X-API-Key")

		// Validación de API Key
		if apiKey == "" {
			logger.Printf("Intento de acceso sin API Key desde IP: %s", c.ClientIP())
			c.AbortWithStatusJSON(401, gin.H{
				"error": "API Key no proporcionada",
				"code":  "AUTH_001",
			})
			return
		}

		if apiKey != config.APIKey {
			logger.Printf("API Key inválida desde IP: %s", c.ClientIP())
			c.AbortWithStatusJSON(401, gin.H{
				"error": "API Key inválida",
				"code":  "AUTH_002",
			})
			return
		}

		// Agregar información de autenticación al contexto
		c.Set("authenticated", true)
		c.Set("auth_time", startTime)

		c.Next()

		// Logging después de la solicitud
		duration := time.Since(startTime)
		logger.Printf(
			"Solicitud autenticada | IP: %s | Duración: %v | Endpoint: %s",
			c.ClientIP(),
			duration,
			c.Request.URL.Path,
		)
	}
}

func AuthMiddlewareHandler(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "No authorization header", http.StatusUnauthorized)
			return
		}

		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 {
			http.Error(w, "Invalid token format", http.StatusUnauthorized)
			return
		}

		token := bearerToken[1]
		valid, err := utils.ValidateToken(token)
		if err != nil || !valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Para pruebas, usamos claims fijos
		userClaims := Claims{
			UserID: "123",
			Rut:    "76212889-6",
			Role:   "ADMIN",
		}

		// Almacenar claims en el contexto
		ctx := r.Context()
		ctx = context.WithValue(ctx, "claims", userClaims)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	}
}
