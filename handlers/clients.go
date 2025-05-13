package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/cursor/FMgo/gateway/api"
	"github.com/cursor/FMgo/gateway/metrics"
	"github.com/cursor/FMgo/gateway/models"
)

// SearchClientsHandler maneja la búsqueda de clientes
func SearchClientsHandler(fmClient *api.FacturaMovilClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		searchTerm := c.Query("search")
		if searchTerm == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Término de búsqueda requerido",
				"code":  "CLI_001",
			})
			return
		}

		// Realizar la búsqueda usando el cliente de Factura Móvil
		respBody, err := fmClient.SearchClients(searchTerm)
		if err != nil {
			log.Printf("Error al buscar clientes: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error al buscar clientes",
				"code":  "CLI_002",
			})
			metrics.RequestCounter.WithLabelValues("search_clients", "error").Inc()
			return
		}

		// Decodificar la respuesta
		var clientResp models.ClientResponse
		if err := json.Unmarshal(respBody, &clientResp); err != nil {
			log.Printf("Error al decodificar respuesta de clientes: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error al procesar la respuesta",
				"code":  "CLI_003",
			})
			return
		}

		// Devolver los resultados
		c.JSON(http.StatusOK, clientResp)
		metrics.RequestCounter.WithLabelValues("search_clients", "success").Inc()
	}
}
