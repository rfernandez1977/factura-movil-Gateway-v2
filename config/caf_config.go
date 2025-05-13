package config

// TipoDocumento representa un tipo de documento tributario
type TipoDocumento struct {
	ID     int    `json:"id" yaml:"id"`
	Nombre string `json:"nombre" yaml:"nombre"`
}

// CAFConfig representa la configuración para el manejo de CAF
type CAFConfig struct {
	RutaCAF        string          `json:"ruta_caf" yaml:"ruta_caf"`
	TiposDocumento []TipoDocumento `json:"tipos_documento" yaml:"tipos_documento"`
}
