package handlers

import (
	"encoding/json"
	"net/http"

	"FMgo/models"
)

// SIIResponseHandler maneja las respuestas del SII
type SIIResponseHandler struct {
	// Aquí se pueden agregar dependencias necesarias
}

// NewSIIResponseHandler crea una nueva instancia del handler
func NewSIIResponseHandler() *SIIResponseHandler {
	return &SIIResponseHandler{}
}

// ProcesarRespuestaHandler maneja el procesamiento de respuestas del SII
func (h *SIIResponseHandler) ProcesarRespuestaHandler(w http.ResponseWriter, r *http.Request) {
	var respuesta models.RespuestaSII
	if err := json.NewDecoder(r.Body).Decode(&respuesta); err != nil {
		http.Error(w, "Error al decodificar la respuesta", http.StatusBadRequest)
		return
	}

	// Procesar la respuesta según el estado
	switch respuesta.Estado {
	case "OK":
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"mensaje": "Documento procesado correctamente",
			"trackId": respuesta.TrackID,
		})
	case "RECHAZADO":
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error":   "Documento rechazado",
			"errores": respuesta.Errores,
		})
	default:
		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(map[string]string{
			"mensaje": "Documento en proceso",
			"trackId": respuesta.TrackID,
		})
	}
}
