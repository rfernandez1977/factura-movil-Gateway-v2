package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Contact representa un contacto del cliente
type Contact struct {
	Name        string   `json:"name"`
	Position    string   `json:"position"`
	Department  string   `json:"department"`
	Email       string   `json:"email"`
	Phones      []string `json:"phones"`
	MainContact bool     `json:"mainContact"`
}

// ValidateContactHandler maneja la validación de contactos
func (h *ClientHandlers) ValidateContactHandler(c *gin.Context) {
	var contacts []Contact

	if err := c.ShouldBindJSON(&contacts); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Datos de contacto inválidos",
			"code":  "CLI_CONT_001",
		})
		return
	}

	// Validar que exista un contacto principal
	hasMainContact := false
	for _, contact := range contacts {
		if contact.MainContact {
			if hasMainContact {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "Solo puede haber un contacto principal",
					"code":  "CLI_CONT_002",
				})
				return
			}
			hasMainContact = true
		}
	}
}
