package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"FMgo/api"
	"FMgo/utils"
)

type EmailHandlers struct {
	client *api.FacturaMovilClient
}

func NewEmailHandlers(client *api.FacturaMovilClient) *EmailHandlers {
	return &EmailHandlers{client: client}
}

// ValidateEmailHandler maneja la validación de correos electrónicos
func (h *EmailHandlers) ValidateEmailHandler(c *gin.Context) {
	var request struct {
		Email string `json:"email" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "Datos de correo inválidos",
			"code":   "EMAIL_001",
			"detail": err.Error(),
		})
		return
	}

	// Limpiar y normalizar el correo
	email := utils.CleanEmail(request.Email)

	// Validar el correo usando el validador unificado
	if err := utils.ValidateEmail(email); err != nil {
		var code string
		switch {
		case err.Error() == "el correo electrónico no puede estar vacío":
			code = "EMAIL_002"
		case err.Error() == "el correo electrónico excede la longitud máxima permitida (254 caracteres)":
			code = "EMAIL_003"
		case err.Error() == "formato de correo electrónico inválido":
			code = "EMAIL_004"
		case err.Error() == "dominio de correo no válido":
			code = "EMAIL_005"
		default:
			code = "EMAIL_006"
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"code":  code,
			"email": email,
			"valid": false,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"email":   email,
		"valid":   true,
		"message": "Correo electrónico válido",
	})
}
