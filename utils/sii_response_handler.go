package utils

import (
	"fmt"

	"github.com/fmgo/models"
)

// ProcesarRespuesta procesa la respuesta del SII y retorna un error si hay problemas
func ProcesarRespuestaSII(r *models.SIIResponseHTTP) error {
	if r.Header.Get("X-Estado") == "OK" {
		return nil
	}

	return fmt.Errorf("Error en la respuesta del SII: Estado=%s, Glosa=%s",
		r.Header.Get("X-Estado"), r.Header.Get("X-Glosa"))
}

// ObtenerTrackID retorna el ID de seguimiento de la transacción
func ObtenerTrackID(r *models.SIIResponseHTTP) string {
	return r.Header.Get("X-Glosa") // O ajusta según el campo correcto, si existe TrackID en Header
}

// EstadoEnvioOK verifica si el estado del envío es correcto
func EstadoEnvioOK(r *models.SIIResponseHTTP) bool {
	return r.Header.Get("X-Estado") == "OK"
}

// ObtenerDetallesEnvio retorna los detalles del envío
func ObtenerDetallesEnvio(r *models.SIIResponseHTTP) string {
	// Si Body tiene EstadoDocumento, ajusta aquí según la estructura real
	return fmt.Sprintf("Estado: %s\nGlosa: %s",
		r.Header.Get("X-Estado"),
		r.Header.Get("X-Glosa"))
}
