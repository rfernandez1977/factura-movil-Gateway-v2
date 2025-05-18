package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"FMgo/api"
)

type MunicipalityHandlers struct {
	client *api.FacturaMovilClient
}

func NewMunicipalityHandlers(client *api.FacturaMovilClient) *MunicipalityHandlers {
	return &MunicipalityHandlers{client: client}
}

// ListMunicipalitiesHandler maneja el listado de municipalidades
func (h *MunicipalityHandlers) ListMunicipalitiesHandler(c *gin.Context) {
	region := c.Query("region")
	if region == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Región requerida",
			"code":  "MUN_001",
		})
		return
	}

	resp, err := h.client.ListMunicipalities(region)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "Error al obtener municipalidades",
			"code":   "MUN_002",
			"detail": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, json.RawMessage(resp))
}
