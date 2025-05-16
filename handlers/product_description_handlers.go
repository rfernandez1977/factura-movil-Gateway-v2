package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/fmgo/api"
)

type ProductDescriptionHandlers struct {
	client *api.FacturaMovilClient
}

func NewProductDescriptionHandlers(client *api.FacturaMovilClient) *ProductDescriptionHandlers {
	return &ProductDescriptionHandlers{client: client}
}

// ValidateDescriptionHandler maneja la validación de descripciones de producto
func (h *ProductDescriptionHandlers) ValidateDescriptionHandler(c *gin.Context) {
	var request struct {
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "Datos de descripción inválidos",
			"code":   "PDESC_001",
			"detail": err.Error(),
		})
		return
	}

	// Validar longitud
	if len(request.Description) < 3 || len(request.Description) > 1000 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Longitud de descripción inválida (debe tener entre 3 y 1000 caracteres)",
			"code":  "PDESC_002",
		})
		return
	}

	// Formatear descripción
	formattedDesc := formatProductDescription(request.Description)

	c.JSON(http.StatusOK, gin.H{
		"description": formattedDesc,
		"isValid":     true,
	})
}

// Función auxiliar para formatear descripción
func formatProductDescription(desc string) string {
	// Remover espacios múltiples
	desc = strings.Join(strings.Fields(desc), " ")

	// Primera letra en mayúscula
	if len(desc) > 0 {
		desc = strings.ToUpper(desc[:1]) + strings.ToLower(desc[1:])
	}

	// Remover caracteres especiales duplicados
	desc = removeDuplicateSpecialChars(desc)

	return desc
}
