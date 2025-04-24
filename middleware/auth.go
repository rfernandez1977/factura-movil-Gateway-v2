package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// APIKeyMiddleware validates API key for authentication
func APIKeyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")
		if apiKey != "your-secure-api-key" { // Replace with your API key
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid API key"})
			c.Abort()
			return
		}
		c.Next()
	}
}