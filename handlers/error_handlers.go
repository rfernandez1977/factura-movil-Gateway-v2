package handlers

import (
	"github.com/gin-gonic/gin"
	"FMgo/api"
)

type ErrorHandlers struct {
	client *api.FacturaMovilClient
}

type ErrorLog struct {
	Code       string                 `json:"code"`
	Message    string                 `json:"message"`
	Stack      string                 `json:"stack"`
	Context    map[string]interface{} `json:"context"`
	Resolution string                 `json:"resolution"`
}

func (h *ErrorHandlers) HandleErrorHandler(c *gin.Context) {
	// Manejo centralizado de errores
	// Registro detallado
	// Sugerencias de resolución
}
