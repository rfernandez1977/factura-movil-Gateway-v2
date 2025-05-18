package utils

import (
	"FMgo/config"
)

// GetSiiEndpoint obtiene el endpoint del SII según la configuración
func GetSiiEndpoint(config *config.SupabaseConfig) string {
	return config.GetSiiEndpoint()
}
