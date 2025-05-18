package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"FMgo/api"
)

type PaymentHandlers struct {
	client *api.FacturaMovilClient
}

func NewPaymentHandlers(client *api.FacturaMovilClient) *PaymentHandlers {
	return &PaymentHandlers{client: client}
}

// ListPaymentMethodsHandler maneja el listado de formas de pago
func (h *PaymentHandlers) ListPaymentMethodsHandler(c *gin.Context) {
	resp, err := h.client.ListPaymentMethods()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "Error al obtener formas de pago",
			"code":   "PAY_001",
			"detail": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, json.RawMessage(resp))
}

// ListPaymentTermsHandler maneja el listado de condiciones de pago
func (h *PaymentHandlers) ListPaymentTermsHandler(c *gin.Context) {
	resp, err := h.client.ListPaymentTerms()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "Error al obtener condiciones de pago",
			"code":   "PAY_002",
			"detail": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, json.RawMessage(resp))
}
