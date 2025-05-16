package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/fmgo/api"
)

type ProductHandlers struct {
	client *api.FacturaMovilClient
}

func NewProductHandlers(client *api.FacturaMovilClient) *ProductHandlers {
	return &ProductHandlers{client: client}
}

// CreateProductHandler maneja la creación de productos
func (h *ProductHandlers) CreateProductHandler(c *gin.Context) {
	var product interface{}
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "Datos de producto inválidos",
			"code":   "PROD_001",
			"detail": err.Error(),
		})
		return
	}

	resp, err := h.client.CreateProduct(product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "Error al crear producto",
			"code":   "PROD_002",
			"detail": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, json.RawMessage(resp))
}

// ListProductsHandler maneja el listado de productos
func (h *ProductHandlers) ListProductsHandler(c *gin.Context) {
	params := make(map[string]string)
	// Parámetros de filtrado comunes
	if category := c.Query("category"); category != "" {
		params["category"] = category
	}
	if search := c.Query("search"); search != "" {
		params["search"] = search
	}

	resp, err := h.client.ListProducts(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "Error al listar productos",
			"code":   "PROD_003",
			"detail": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, json.RawMessage(resp))
}
