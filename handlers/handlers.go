package handlers

import (
	"github.com/cursor/FMgo/api"
)

// BaseHandler contiene la configuración común para todos los handlers
type BaseHandler struct {
	client *api.FacturaMovilClient
}

// NewBaseHandler crea una nueva instancia del handler base
func NewBaseHandler(client *api.FacturaMovilClient) *BaseHandler {
	return &BaseHandler{
		client: client,
	}
}
