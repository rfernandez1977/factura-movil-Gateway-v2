package handlers

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/fmgo/api"
)

type BusinessActivityHandlers struct {
	client *api.FacturaMovilClient
}

func NewBusinessActivityHandlers(client *api.FacturaMovilClient) *BusinessActivityHandlers {
	return &BusinessActivityHandlers{client: client}
}

// ValidateBusinessActivityHandler maneja la validación del giro
func (h *BusinessActivityHandlers) ValidateBusinessActivityHandler(c *gin.Context) {
	var request struct {
		BusinessActivity string `json:"businessActivity"`
		RUT              string `json:"rut,omitempty"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "Datos de giro inválidos",
			"code":   "BACT_001",
			"detail": err.Error(),
		})
		return
	}

	// Validar longitud del giro
	if len(request.BusinessActivity) < 5 || len(request.BusinessActivity) > 80 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Longitud de giro inválida (debe tener entre 5 y 80 caracteres)",
			"code":  "BACT_002",
		})
		return
	}

	// Validar caracteres permitidos
	activityRegex := regexp.MustCompile(`^[a-zA-Z0-9áéíóúÁÉÍÓÚñÑ\s\.,()-]+$`)
	if !activityRegex.MatchString(request.BusinessActivity) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Giro contiene caracteres no permitidos",
			"code":  "BACT_003",
		})
		return
	}

	// Formatear giro
	formattedActivity := formatBusinessActivity(request.BusinessActivity)

	// Si se proporciona RUT, validar contra SII
	if request.RUT != "" {
		c.JSON(http.StatusNotImplemented, gin.H{"message": "TODO: Implementar ValidateBusinessActivity en FacturaMovilClient"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"originalActivity":  request.BusinessActivity,
		"formattedActivity": formattedActivity,
		"isValid":           true,
	})
}

// ListCommonActivitiesHandler lista giros comunes
func (h *BusinessActivityHandlers) ListCommonActivitiesHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "TODO: Implementar ListCommonActivities en FacturaMovilClient"})
	return
}

// Función auxiliar para formatear giro
func formatBusinessActivity(activity string) string {
	// Remover espacios múltiples
	activity = strings.Join(strings.Fields(activity), " ")

	// Convertir a mayúsculas
	activity = strings.ToUpper(activity)

	// Remover caracteres especiales duplicados
	activity = removeDuplicateSpecialChars(activity)

	return activity
}

// Función auxiliar para remover caracteres especiales duplicados
func removeDuplicateSpecialChars(text string) string {
	// Remover puntos duplicados
	for strings.Contains(text, "..") {
		text = strings.ReplaceAll(text, "..", ".")
	}

	// Remover comas duplicadas
	for strings.Contains(text, ",,") {
		text = strings.ReplaceAll(text, ",,", ",")
	}

	// Remover espacios alrededor de puntuación
	text = strings.ReplaceAll(text, " .", ".")
	text = strings.ReplaceAll(text, ". ", ".")
	text = strings.ReplaceAll(text, " ,", ",")
	text = strings.ReplaceAll(text, ", ", ",")

	return text
}
