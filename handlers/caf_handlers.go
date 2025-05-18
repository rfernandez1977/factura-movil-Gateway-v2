package handlers

import (
	"time"

	"FMgo/api"
	"github.com/gin-gonic/gin"
)

// CAFHandlers maneja las operaciones de CAF
type CAFHandlers struct {
	client *api.FacturaMovilClient
}

// NewCAFHandlers crea una nueva instancia de CAFHandlers
func NewCAFHandlers(client *api.FacturaMovilClient) *CAFHandlers {
	return &CAFHandlers{
		client: client,
	}
}

// GetCAFHandler maneja la obtención de CAF
func (h *CAFHandlers) GetCAFHandler(c *gin.Context) {
	// Implementación pendiente
	c.JSON(200, gin.H{
		"message": "CAF handler",
	})
}

type CAFValidation struct {
	TipoDocumento     string    `json:"tipoDocumento"`
	RangoInicial      int       `json:"rangoInicial"`
	RangoFinal        int       `json:"rangoFinal"`
	FechaVencimiento  time.Time `json:"fechaVencimiento"`
	EstadoCAF         string    `json:"estadoCAF"`
	FoliosDisponibles int       `json:"foliosDisponibles"`
}

func (h *CAFHandlers) ValidateCAFHandler(c *gin.Context) {
	var caf CAFValidation

	// Validación de vigencia del CAF
	if time.Now().After(caf.FechaVencimiento) {
		c.JSON(400, gin.H{
			"error":   "CAF vencido",
			"codigo":  "CAF_001",
			"detalle": "El CAF ha expirado, solicite uno nuevo",
		})
		return
	}

	// Validación de folios disponibles
	if caf.FoliosDisponibles < 100 {
		c.JSON(200, gin.H{
			"advertencia":     "Folios próximos a agotarse",
			"codigo":          "CAF_002",
			"foliosRestantes": caf.FoliosDisponibles,
		})
	}
}
