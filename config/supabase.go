package config

// Este archivo se mantiene como referencia para futuras implementaciones específicas de Supabase
// La configuración principal se ha movido a config.go para evitar duplicación

// GetSiiEndpoint obtiene el endpoint del SII según el ambiente
func (c *SupabaseConfig) GetSiiEndpoint() string {
	// Usar ambiente desde la configuración
	ambiente := "certificacion" // valor por defecto

	// Comprobar ambiente de producción
	if ambiente == "produccion" {
		return "https://palena.sii.cl"
	}
	return "https://maullin.sii.cl"
}
