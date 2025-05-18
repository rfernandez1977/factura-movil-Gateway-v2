package handlers

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"FMgo/api"
)

type NameHandlers struct {
	client *api.FacturaMovilClient
}

func NewNameHandlers(client *api.FacturaMovilClient) *NameHandlers {
	return &NameHandlers{client: client}
}

// ValidateNameHandler maneja la validación de nombres
func (h *NameHandlers) ValidateNameHandler(c *gin.Context) {
	var request struct {
		Name      string `json:"name"`
		LastName  string `json:"lastName,omitempty"`
		IsCompany bool   `json:"isCompany"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "Datos de nombre inválidos",
			"code":   "NAME_001",
			"detail": err.Error(),
		})
		return
	}

	// Validar longitud del nombre
	if len(request.Name) < 2 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Nombre demasiado corto",
			"code":  "NAME_002",
		})
		return
	}

	// Validar caracteres permitidos
	nameRegex := regexp.MustCompile(`^[a-zA-ZáéíóúÁÉÍÓÚñÑ\s]+$`)
	if !nameRegex.MatchString(request.Name) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Nombre contiene caracteres no permitidos",
			"code":  "NAME_003",
		})
		return
	}

	// Formatear nombre
	formattedName := formatName(request.Name, request.LastName, request.IsCompany)

	c.JSON(http.StatusOK, gin.H{
		"formattedName": formattedName,
		"isValid":       true,
	})
}

// FormatNameHandler maneja el formateo de nombres
func (h *NameHandlers) FormatNameHandler(c *gin.Context) {
	name := strings.TrimSpace(c.Query("name"))
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Nombre requerido",
			"code":  "NAME_004",
		})
		return
	}

	// Formatear nombre según reglas de negocio
	formattedName := standardizeName(name)

	c.JSON(http.StatusOK, gin.H{
		"original":  name,
		"formatted": formattedName,
	})
}

// Función auxiliar para formatear nombres
func formatName(name, lastName string, isCompany bool) string {
	if isCompany {
		return strings.ToUpper(strings.TrimSpace(name))
	}

	// Formatear nombre de persona
	nameParts := strings.Fields(name)
	for i, part := range nameParts {
		nameParts[i] = strings.Title(strings.ToLower(part))
	}

	if lastName != "" {
		lastNameParts := strings.Fields(lastName)
		for i, part := range lastNameParts {
			lastNameParts[i] = strings.Title(strings.ToLower(part))
		}
		return strings.Join(nameParts, " ") + " " + strings.Join(lastNameParts, " ")
	}

	return strings.Join(nameParts, " ")
}

// Función auxiliar para estandarizar nombres
func standardizeName(name string) string {
	// Remover múltiples espacios
	name = strings.Join(strings.Fields(name), " ")

	// Convertir primera letra de cada palabra a mayúscula
	words := strings.Split(name, " ")
	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(word[:1]) + strings.ToLower(word[1:])
		}
	}

	return strings.Join(words, " ")
}
