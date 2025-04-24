package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/usuario/gateway/api"
	"github.com/usuario/gateway/metrics"
	"github.com/gin-gonic/gin"
)

// validateEntityData validates required fields for clients and products
func validateEntityData(entityData map[string]interface{}, entityType string) error {
	var requiredFields []string
	switch entityType {
	case "cliente":
		requiredFields = []string{"code", "name", "address"}
	case "producto":
		requiredFields = []string{"code", "name", "price"}
	}
	for _, field := range requiredFields {
		if _, exists := entityData[field]; !exists {
			return fmt.Errorf("missing required field: %s", field)
		}
	}
	return nil
}

// CreateEntityHandler handles creation of entities (clients, products)
func CreateEntityHandler(entityType string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var entityData map[string]interface{}
		if err := c.BindJSON(&entityData); err != nil {
			log.Printf("Invalid request payload for %s: %v", entityType, err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}

		// Validate entity data
		if err := validateEntityData(entityData, entityType); err != nil {
			log.Printf("Validation failed for %s: %v", entityType, err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var endpoint string
		switch entityType {
		case "cliente":
			endpoint = "/services/common/client"
		case "producto":
			endpoint = "/services/common/product"
		}

		resp, err := api.CallFacturaMovil("POST", endpoint, entityData, 2)
		if err != nil {
			log.Printf("Failed to communicate with Factura Móvil for %s: %v", entityType, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to communicate with Factura Móvil: " + err.Error()})
			metrics.RequestCounter.WithLabelValues(entityType, "error").Inc()
			return
		}
		defer resp.Body.Close()

		log.Printf("%s created successfully", entityType)
		c.JSON(http.StatusOK, gin.H{"status": fmt.Sprintf("%s created successfully", entityType)})
		metrics.RequestCounter.WithLabelValues(entityType, "success").Inc()
	}
}