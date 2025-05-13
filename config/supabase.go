package config

// Este archivo se mantiene como referencia para futuras implementaciones específicas de Supabase
// La configuración principal se ha movido a config.go para evitar duplicación

// SupabaseConfig contiene la configuración para conectarse a Supabase
type SupabaseConfig struct {
	Client interface{} `json:"-" bson:"-"`
	URL        string `json:"url" bson:"url"`
	Key        string `json:"key" bson:"key"`
	Token      string `json:"token" bson:"token"`
	Ambiente   string `json:"ambiente" bson:"ambiente"`
	BaseURL    string `json:"base_url" bson:"base_url"`
	JWTSecret  string `json:"jwt_secret" bson:"jwt_secret"`
	SIIBaseURL string `json:"sii_base_url" bson:"sii_base_url"`
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

// GetSiiEndpoint obtiene el endpoint del SII según el ambiente
func (c *SupabaseConfig) GetSiiEndpoint() string {
	if c.SIIBaseURL != "" {
		return c.SIIBaseURL
	}

	if c.Ambiente == "produccion" {
		return "https://palena.sii.cl"
	}
	return "https://maullin.sii.cl"
}
