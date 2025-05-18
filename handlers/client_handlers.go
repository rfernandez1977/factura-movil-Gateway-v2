package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"FMgo/api"
)

type ClientHandlers struct {
	client *api.FacturaMovilClient
}

func NewClientHandlers(client *api.FacturaMovilClient) *ClientHandlers {
	return &ClientHandlers{client: client}
}

// SearchClientsHandler maneja la búsqueda de clientes
func (h *ClientHandlers) SearchClientsHandler(c *gin.Context) {
	searchTerm := c.Query("search")
	if searchTerm == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Término de búsqueda requerido",
			"code":  "CLI_001",
		})
		return
	}

	resp, err := h.client.SearchClients(searchTerm)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "Error en la búsqueda",
			"code":   "CLI_002",
			"detail": err.Error(),
		})
		return
	}

	var clients interface{}
	if err := json.Unmarshal(resp, &clients); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "Error al procesar resultados",
			"code":   "CLI_003",
			"detail": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, clients)
}
