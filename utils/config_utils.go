package utils

import (
	"github.com/cursor/FMgo/config"
)

// GetSiiEndpoint obtiene el endpoint del SII según la configuración
func GetSiiEndpoint(config *config.SupabaseConfig) string {
	if config.SIIBaseURL != "" {
		return config.SIIBaseURL
	}

	if config.Ambiente == "produccion" {
		return "https://palena.sii.cl"
	}
	return "https://maullin.sii.cl"
}
