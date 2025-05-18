package handlers

import (
	"encoding/json"
	"net/http"

	"FMgo/api"
	"FMgo/models"
	"FMgo/services/validations"
)

// ValidationHandler maneja las rutas de validación
type ValidationHandler struct {
	validationService *validations.ValidationService
}

// NewValidationHandler crea un nuevo ValidationHandler
func NewValidationHandler(validationService *validations.ValidationService) *ValidationHandler {
	return &ValidationHandler{
		validationService: validationService,
	}
}

// RegisterRoutes registra las rutas del manejador
func (h *ValidationHandler) RegisterRoutes(router *api.Router) {
	router.Post("/api/validate/factura", h.ValidateFactura)
	router.Post("/api/validate/boleta", h.ValidateBoleta)
}

// ValidateFactura valida una factura
func (h *ValidationHandler) ValidateFactura(w http.ResponseWriter, r *http.Request) {
	var factura models.Factura
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&factura); err != nil {
		api.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	errors := h.validationService.ValidateFactura(&factura)
	if len(errors) > 0 {
		api.RespondWithJSON(w, http.StatusBadRequest, errors)
		return
	}

	api.RespondWithJSON(w, http.StatusOK, map[string]bool{"valid": true})
}

// ValidateBoleta valida una boleta
func (h *ValidationHandler) ValidateBoleta(w http.ResponseWriter, r *http.Request) {
	var boleta models.Boleta
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&boleta); err != nil {
		api.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	errors := h.validationService.ValidateBoleta(&boleta)
	if len(errors) > 0 {
		api.RespondWithJSON(w, http.StatusBadRequest, errors)
		return
	}

	api.RespondWithJSON(w, http.StatusOK, map[string]bool{"valid": true})
}
