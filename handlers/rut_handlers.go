package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/fmgo/api"

	"github.com/gin-gonic/gin"
	"github.com/fmgo/utils"
)

type RutHandlers struct {
	client *api.FacturaMovilClient
}

func NewRutHandlers(client *api.FacturaMovilClient) *RutHandlers {
	return &RutHandlers{client: client}
}

// ValidateRutHandler maneja la validación de RUT
func (h *RutHandlers) ValidateRutHandler(c *gin.Context) {
	rut := c.Query("rut")
	if rut == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "RUT requerido",
			"code":  "RUT_001",
		})
		return
	}

	// Validar RUT
	if err := utils.ValidateRUT(rut); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"code":  "RUT_002",
		})
		return
	}

	// Consultar RUT en Factura Móvil
	resp, err := h.client.ValidateRut(rut)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "Error al validar RUT",
			"code":   "RUT_003",
			"detail": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, json.RawMessage(resp))
}

// FormatRutHandler maneja el formateo de RUT
func (h *RutHandlers) FormatRutHandler(c *gin.Context) {
	rut := c.Query("rut")
	if rut == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "RUT requerido",
			"code":  "RUT_004",
		})
		return
	}

	// Formatear RUT
	formattedRut := utils.FormatRUT(rut)

	c.JSON(http.StatusOK, gin.H{
		"rut": formattedRut,
	})
}
