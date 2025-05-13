package config

// Este archivo se mantiene como referencia para futuras implementaciones específicas de Supabase
// La configuración principal se ha movido a config.go para evitar duplicación

// SupabaseConfig contiene la configuración para conectarse a Supabase
type SupabaseConfig struct {
	URL      string `json:"url" bson:"url"`
	Key      string `json:"key" bson:"key"`
	Token    string `json:"token" bson:"token"`
	Ambiente string `json:"ambiente" bson:"ambiente"`
	BaseURL  string `json:"base_url" bson:"base_url"`
}

// NewSupabaseConfig crea una nueva configuración de Supabase
func NewSupabaseConfig(url, key, token, ambiente, baseURL string) *SupabaseConfig {
	return &SupabaseConfig{
		URL:      url,
		Key:      key,
		Token:    token,
		Ambiente: ambiente,
		BaseURL:  baseURL,
	}
}
