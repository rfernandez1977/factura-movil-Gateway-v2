# FMgo - Cliente SII para Facturación Electrónica

## Estado del Proyecto: Fase de Certificación 🚀

FMgo es un cliente robusto para la integración con los servicios del SII (Servicio de Impuestos Internos) de Chile, actualmente en fase de certificación.

### Características Principales

- ✅ Cliente SOAP completo para servicios del SII
- ✅ Manejo de autenticación y tokens
- ✅ Envío y consulta de DTEs
- ✅ Validación XSD de documentos
- ✅ Firma digital de documentos
- ✅ Pruebas unitarias completas
- 🔄 En proceso de certificación SII

## Estructura del Proyecto
```
FMgo/
├── core/
│   └── firma/
│       ├── models/
│       │   └── configuracion.go
│       └── interfaces/
│           └── firma_service.go
├── services/
│   ├── base_firma_service.go
│   ├── sii_firma_service.go
│   ├── cert_cache.go
│   └── test_data/
│       └── test_cert.go
├── docs/
│   └── firma_digital.md
└── README.md
```

### Componentes Principales

#### Cliente SII (`core/sii/client/`)
- Implementación oficial para la comunicación con el SII
- Características:
  - Manejo robusto de certificados
  - Sistema de reintentos configurable
  - Logging estructurado
  - Pruebas unitarias completas
  - Soporte para ambiente de certificación y producción

Para usar el cliente SII:
```go
import "FMgo/core/sii/client"

// Crear cliente
config := &models.Config{
    SII: models.SIIConfig{
        BaseURL:    "https://palena.sii.cl",
        CertPath:   "/path/to/cert.pem",
        KeyPath:    "/path/to/key.pem",
        RetryCount: 3,
        Timeout:    30,
    },
}

logger := logger.NewLogger()
client, err := client.NewHTTPClient(config, logger)
```

## Requisitos
- Go 1.21 o superior
- OpenSSL
- Certificado digital válido para el SII

## Instalación
```bash
go get github.com/usuario/FMgo
```

## Configuración
```go
config := &siimodels.ConfigSII{
    Ambiente:       siimodels.AmbienteCertificacion,
    BaseURL:        siimodels.URLBaseCertificacion,
    CertPath:       "/path/to/cert.p12",
    KeyPath:        "/path/to/key.pem",
    RutEmpresa:     "76555555-5",
    RutCertificado: "11111111-1",
    RetryCount:     3,
    Timeout:        30 * time.Second,
}
```

## Uso Básico

### Envío de DTE
```go
client, err := NewDTEClient(config)
if err != nil {
    log.Fatal(err)
}

resp, err := client.EnviarDTE(context.Background(), dte)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("TrackID: %s\n", resp.TrackID)
```

### Consulta de Estado
```go
estado, err := client.ConsultarEstadoDTE(context.Background(), rutEmisor, tipoDTE, folio)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Estado: %s\n", estado.Estado)
```

## Documentación

- [Guía de Desarrollo](docs/DESARROLLO.md)
- [Documentación de API](docs/API.md)
- [Guía de Validación](docs/VALIDACION.md)
- [Fase de Certificación](docs/CERTIFICACION.md)
- [Troubleshooting](docs/TROUBLESHOOTING.md)

## Estado de la Certificación

Actualmente en proceso de certificación con el SII:

- ✅ Implementación base completada
- ✅ Pruebas unitarias implementadas
- ✅ Validación XSD implementada
- ✅ Firma digital implementada
- 🔄 Set de pruebas de certificación en proceso
- ⏳ Certificación final pendiente

## Contribución

1. Fork el repositorio
2. Crea una rama para tu feature (`git checkout -b feature/AmazingFeature`)
3. Commit tus cambios (`git commit -m 'Add some AmazingFeature'`)
4. Push a la rama (`git push origin feature/AmazingFeature`)
5. Abre un Pull Request

## Licencia

Este proyecto está licenciado bajo la Licencia MIT - ver el archivo [LICENSE](LICENSE) para más detalles.

## Agradecimientos

- Equipo de desarrollo original
- Contribuidores
- SII por la documentación y soporte

## Integración SII

### Descripción
FMgo incluye una integración completa con el Servicio de Impuestos Internos (SII) para la gestión de documentos tributarios electrónicos (DTE).

### Características Principales
- Autenticación automática con el SII
- Envío y consulta de DTEs
- Manejo de certificados digitales
- Sistema de reintentos y recuperación de errores
- Soporte para ambientes de certificación y producción

### Configuración Rápida
1. Configurar certificados:
   ```json
   {
     "ambiente": "certificacion",
     "cert_path": "./certificados/cert.crt",
     "key_path": "./certificados/key.key"
   }
   ```

2. Inicializar cliente:
   ```go
   config := models.NewConfig()
   client, err := client.NewHTTPClient(config, logger)
   ```

3. Realizar operaciones:
   ```go
   // Obtener semilla
   semilla, err := client.ObtenerSemilla(ctx)
   
   // Obtener token
   token, err := client.ObtenerToken(ctx, semilla)
   ```

### Documentación
Para más detalles sobre la integración, consultar:
- [Documentación detallada](docs/sii_integration.md)
- [Ejemplos de uso](scripts/test_sii_connection.go)
- [Registro de cambios](CHANGELOG.md) 