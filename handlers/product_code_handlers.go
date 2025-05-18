package handlers

import (
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"FMgo/api"
)

type ProductCodeHandlers struct {
	client *api.FacturaMovilClient
}

func NewProductCodeHandlers(client *api.FacturaMovilClient) *ProductCodeHandlers {
	return &ProductCodeHandlers{client: client}
}

// ValidateProductCodeHandler maneja la validación de códigos de producto
func (h *ProductCodeHandlers) ValidateProductCodeHandler(c *gin.Context) {
	var request struct {
		Code string `json:"code"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "Datos de código inválidos",
			"code":   "PCODE_001",
			"detail": err.Error(),
		})
		return
	}

	// Validar longitud y caracteres
	codeRegex := regexp.MustCompile(`^[A-Z0-9-]{3,20}$`)
	if !codeRegex.MatchString(request.Code) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Formato de código inválido (debe tener entre 3 y 20 caracteres alfanuméricos)",
			"code":  "PCODE_002",
		})
		return
	}

	// Verificar si el código ya existe
	exists, err := h.client.CheckProductCodeExists(request.Code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "Error al verificar código",
			"code":   "PCODE_003",
			"detail": err.Error(),
		})
		return
	}

	if exists {
		c.JSON(http.StatusConflict, gin.H{
			"error": "El código ya existe",
			"code":  "PCODE_004",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    request.Code,
		"isValid": true,
	})
}
