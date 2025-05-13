package handlers

import (
	"time"

	"github.com/gin-gonic/gin"
)

// ClientClassification maneja la clasificación de clientes
type ClientClassification struct {
	Category       string    `json:"category"` // A, B, C
	PurchaseVolume float64   `json:"volume"`
	Frequency      int       `json:"frequency"` // Compras por mes
	LastPurchase   time.Time `json:"lastPurchase"`
	LoyaltyPoints  int       `json:"points"`
}

// ClassifyClientHandler maneja la clasificación automática
func (h *ClientHandlers) ClassifyClientHandler(c *gin.Context) {
	var clientData ClientClassification

	// Calcular categoría basada en volumen y frecuencia
	if clientData.PurchaseVolume > 1000000 && clientData.Frequency > 5 {
		clientData.Category = "A"
	} else if clientData.PurchaseVolume > 500000 || clientData.Frequency > 3 {
		clientData.Category = "B"
	} else {
		clientData.Category = "C"
	}

	// Verificar actividad reciente
	if time.Since(clientData.LastPurchase) > 90*24*time.Hour {
		clientData.Category = "INACTIVO"
	}

	c.JSON(200, clientData)
}
