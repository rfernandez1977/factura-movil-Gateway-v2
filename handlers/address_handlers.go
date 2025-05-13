package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/cursor/FMgo/api"
)

type AddressHandlers struct {
	client *api.FacturaMovilClient
}

func NewAddressHandlers(client *api.FacturaMovilClient) *AddressHandlers {
	return &AddressHandlers{client: client}
}

// CreateAddressHandler maneja la creación de direcciones adicionales
func (h *AddressHandlers) CreateAddressHandler(c *gin.Context) {
	clientID := c.Param("clientId")
	if clientID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID de cliente requerido",
			"code":  "ADDR_001",
		})
		return
	}

	var address interface{}
	if err := c.ShouldBindJSON(&address); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "Datos de dirección inválidos",
			"code":   "ADDR_002",
			"detail": err.Error(),
		})
		return
	}

	resp, err := h.client.CreateAddress(clientID, address)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "Error al crear dirección",
			"code":   "ADDR_003",
			"detail": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, json.RawMessage(resp))
}

// UpdateAddressHandler maneja la actualización de direcciones
func (h *AddressHandlers) UpdateAddressHandler(c *gin.Context) {
	addressID := c.Param("addressId")
	if addressID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID de dirección requerido",
			"code":  "ADDR_004",
		})
		return
	}

	var address interface{}
	if err := c.ShouldBindJSON(&address); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "Datos de dirección inválidos",
			"code":   "ADDR_005",
			"detail": err.Error(),
		})
		return
	}

	resp, err := h.client.UpdateAddress(addressID, address)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "Error al actualizar dirección",
			"code":   "ADDR_006",
			"detail": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, json.RawMessage(resp))
}

// DeleteAddressHandler maneja la eliminación de direcciones
func (h *AddressHandlers) DeleteAddressHandler(c *gin.Context) {
	addressID := c.Param("addressId")
	if addressID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID de dirección requerido",
			"code":  "ADDR_007",
		})
		return
	}

	err := h.client.DeleteAddress(addressID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "Error al eliminar dirección",
			"code":   "ADDR_008",
			"detail": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Dirección eliminada exitosamente",
	})
}

// ListAddressesHandler maneja el listado de direcciones de un cliente
func (h *AddressHandlers) ListAddressesHandler(c *gin.Context) {
	clientID := c.Param("clientId")
	if clientID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID de cliente requerido",
			"code":  "ADDR_009",
		})
		return
	}

	resp, err := h.client.ListAddresses(clientID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "Error al listar direcciones",
			"code":   "ADDR_010",
			"detail": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, json.RawMessage(resp))
}
