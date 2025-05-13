package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/cursor/FMgo/api"

	"github.com/gin-gonic/gin"
)

type PriceHandlers struct {
	client *api.FacturaMovilClient
}

func NewPriceHandlers(client *api.FacturaMovilClient) *PriceHandlers {
	return &PriceHandlers{client: client}
}

// ValidatePriceHandler maneja la validación de precios
func (h *PriceHandlers) ValidatePriceHandler(c *gin.Context) {
	var request struct {
		Price       float64 `json:"price"`
		Currency    string  `json:"currency"`
		ProductType string  `json:"productType"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "Datos de precio inválidos",
			"code":   "PRICE_001",
			"detail": err.Error(),
		})
		return
	}

	// Validar precio negativo
	if request.Price < 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "El precio no puede ser negativo",
			"code":  "PRICE_002",
		})
		return
	}

	// Validar moneda
	if !isValidCurrency(request.Currency) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Moneda inválida",
			"code":  "PRICE_003",
		})
		return
	}

	// Formatear precio según reglas de negocio
	formattedPrice := formatPrice(request.Price, request.Currency)

	c.JSON(http.StatusOK, gin.H{
		"originalPrice":  request.Price,
		"formattedPrice": formattedPrice,
		"isValid":        true,
	})
}

// CalculateTaxesHandler maneja el cálculo de impuestos
func (h *PriceHandlers) CalculateTaxesHandler(c *gin.Context) {
	var request struct {
		NetPrice float64 `json:"netPrice"`
		TaxRate  float64 `json:"taxRate"`
		IsExempt bool    `json:"isExempt"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "Datos inválidos para cálculo de impuestos",
			"code":   "PRICE_004",
			"detail": err.Error(),
		})
		return
	}

	// Calcular impuestos
	var taxAmount float64
	var totalPrice float64

	if !request.IsExempt {
		taxAmount = request.NetPrice * (request.TaxRate / 100)
		totalPrice = request.NetPrice + taxAmount
	} else {
		totalPrice = request.NetPrice
	}

	c.JSON(http.StatusOK, gin.H{
		"netPrice":   request.NetPrice,
		"taxAmount":  taxAmount,
		"totalPrice": totalPrice,
		"isExempt":   request.IsExempt,
	})
}

// UpdatePriceListHandler maneja la actualización de listas de precios
func (h *PriceHandlers) UpdatePriceListHandler(c *gin.Context) {
	var request struct {
		ProductID string  `json:"productId"`
		NewPrice  float64 `json:"newPrice"`
		ListType  string  `json:"listType"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "Datos inválidos para actualización de precio",
			"code":   "PRICE_005",
			"detail": err.Error(),
		})
		return
	}

	// Validar y actualizar precio en la lista
	resp, err := h.client.UpdateProductPrice(request.ProductID, request.NewPrice, request.ListType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "Error al actualizar precio",
			"code":   "PRICE_006",
			"detail": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, json.RawMessage(resp))
}

// Función auxiliar para validar moneda
func isValidCurrency(currency string) bool {
	validCurrencies := map[string]bool{
		"CLP": true,
		"USD": true,
		"EUR": true,
		"UF":  true,
	}
	return validCurrencies[currency]
}

// Función auxiliar para formatear precio
func formatPrice(price float64, currency string) string {
	switch currency {
	case "CLP":
		return "$ " + strconv.FormatFloat(price, 'f', 0, 64)
	case "USD":
		return "USD " + strconv.FormatFloat(price, 'f', 2, 64)
	case "EUR":
		return "€ " + strconv.FormatFloat(price, 'f', 2, 64)
	case "UF":
		return "UF " + strconv.FormatFloat(price, 'f', 4, 64)
	default:
		return strconv.FormatFloat(price, 'f', 2, 64)
	}
}
