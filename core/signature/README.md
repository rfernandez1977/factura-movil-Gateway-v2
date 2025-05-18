# Módulo de Firma Digital

Este módulo proporciona la funcionalidad de firma digital para documentos XML según los requerimientos del SII.

## Estructura

```
core/signature/
├── services/       # Servicios de firma
├── models/        # Modelos de datos
├── utils/         # Utilidades comunes
└── tests/         # Tests unitarios
```

## Componentes

### Services
- `signature_service.go`: Servicio principal de firma
- `caf_service.go`: Gestión de CAF
- `certificate_service.go`: Gestión de certificados

### Models
- `signature.go`: Modelos para firma XML
- `caf.go`: Modelos para CAF
- `certificate.go`: Modelos para certificados

### Utils
- `xml_utils.go`: Utilidades para procesamiento XML
- `crypto_utils.go`: Utilidades criptográficas

## Características

- Firma de documentos XML (XML-DSIG)
- Gestión de CAF (Código de Autorización de Folios)
- Validación de certificados
- Caché de certificados
- Sistema de alertas para vencimiento
- Logging detallado de operaciones

## Uso

```go
// Crear servicio de firma
signatureService := signature.NewService(config)

// Firmar documento
signedDoc, err := signatureService.SignXML(xmlData)

// Validar firma
isValid, err := signatureService.ValidateSignature(signedDoc)

// Gestionar CAF
cafService := caf.NewService(config)
err := cafService.ValidateCAF(cafData)
```

## Tests

```bash
# Ejecutar tests
go test ./...

# Ejecutar tests con coverage
go test -cover ./...
``` 