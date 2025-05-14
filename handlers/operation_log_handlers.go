package handlers

import (
	"net/http"

	"github.com/cursor/FMgo/api"
	"github.com/cursor/FMgo/services"
)

// OperationLogHandler maneja las rutas de logs de operaciones
type OperationLogHandler struct {
	logService *services.LogService
}

// NewOperationLogHandler crea un nuevo OperationLogHandler
func NewOperationLogHandler(logService *services.LogService) *OperationLogHandler {
	return &OperationLogHandler{
		logService: logService,
	}
}

// RegisterRoutes registra las rutas del manejador
func (h *OperationLogHandler) RegisterRoutes(router *api.Router) {
	router.Get("/api/logs", h.GetLogs)
	router.Get("/api/logs/:id", h.GetLog)
}

// GetLogs devuelve todos los logs de operaciones
func (h *OperationLogHandler) GetLogs(w http.ResponseWriter, r *http.Request) {
	logs, err := h.logService.GetLogs()
	if err != nil {
		api.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	api.RespondWithJSON(w, http.StatusOK, logs)
}

// GetLog devuelve un log específico
func (h *OperationLogHandler) GetLog(w http.ResponseWriter, r *http.Request) {
	// TODO: Extraer ID de la URL con algún router que soporte parámetros
	id := r.URL.Query().Get("id")
	if id == "" {
		api.RespondWithError(w, http.StatusBadRequest, "ID de log no proporcionado")
		return
	}

	log, err := h.logService.GetLog(id)
	if err != nil {
		api.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	api.RespondWithJSON(w, http.StatusOK, log)
}
