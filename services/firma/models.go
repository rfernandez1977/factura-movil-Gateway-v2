package firma

// ConfiguracionFirma contiene la configuración para los servicios de firma
type ConfiguracionFirma struct {
	// RutaCertificado es la ruta al archivo del certificado (PFX o PEM)
	RutaCertificado string

	// RutaLlave es la ruta al archivo de la llave privada (solo para PEM)
	RutaLlave string

	// Password es la contraseña del certificado (solo para PFX)
	Password string

	// RutEmpresa es el RUT de la empresa emisora
	RutEmpresa string

	// TTLCache es el tiempo de vida de los certificados en caché (opcional)
	TTLCache int64

	// MaxItemsCache es el número máximo de certificados en caché (opcional)
	MaxItemsCache int
} 