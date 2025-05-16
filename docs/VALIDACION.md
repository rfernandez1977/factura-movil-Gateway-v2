# Validación XSD en FMgo

## Descripción General

El sistema de validación XSD en FMgo proporciona una capa robusta de validación para documentos tributarios electrónicos según los esquemas oficiales del SII.

## Componentes

### ValidatorService

```go
type ValidatorService struct {
    mu       sync.RWMutex
    schemas  map[string]*xsd.Schema
    basePath string
}
```

El `ValidatorService` es thread-safe y proporciona:
- Carga lazy de esquemas XSD
- Validación de documentos individuales
- Validación de envíos completos
- Manejo de múltiples tipos de documentos

### ValidatorMiddleware

```go
type ValidatorMiddleware struct {
    validator *ValidatorService
    once     sync.Once
}
```

El middleware proporciona:
- Validación automática en handlers
- Inicialización lazy de esquemas
- Manejo consistente de errores

## Tipos de Documentos Soportados

| Tipo | Código | Esquema |
|------|--------|---------|
| Factura Electrónica | 33 | DTE_v10.xsd |
| Factura Exenta | 34 | DTE_v10.xsd |
| Nota de Débito | 56 | DTE_v10.xsd |
| Nota de Crédito | 61 | DTE_v10.xsd |
| Boleta | 39 | EnvioBOLETA_v11.xsd |
| Boleta Exenta | 41 | EnvioBOLETA_v11.xsd |

## Uso

### Validación Simple

```go
validator := services.NewValidatorService("schema_dte")
err := validator.ValidarDocumento(documento)
```

### Uso con Middleware

```go
middleware := middleware.NewValidatorMiddleware("schema_dte")
handler := middleware.WithDocumentoValidation(miHandler)
```

## Manejo de Errores

El sistema proporciona mensajes de error detallados para:
- Errores de carga de esquemas
- Errores de validación XML
- Errores de estructura de documentos

## Buenas Prácticas

1. **Carga de Esquemas**:
   - Cargar esquemas al inicio de la aplicación
   - Utilizar el middleware para carga lazy
   - Liberar esquemas cuando no se necesiten

2. **Validación**:
   - Validar documentos antes de firmarlos
   - Validar antes de enviar al SII
   - Usar el middleware en endpoints públicos

3. **Performance**:
   - Reutilizar instancias del validador
   - Implementar caché de esquemas
   - Liberar recursos apropiadamente

## Logging y Debugging

El sistema incluye logging detallado para:
- Carga de esquemas
- Errores de validación
- Uso de recursos

## Esquemas XSD

Los esquemas utilizados son:
- `DTE_v10.xsd`: Esquema para documentos tributarios
- `EnvioDTE_v10.xsd`: Esquema para sobres de envío
- `EnvioBOLETA_v11.xsd`: Esquema para boletas electrónicas
- `SiiTypes_v10.xsd`: Tipos de datos comunes
- `xmldsignature_v10.xsd`: Esquema para firmas XML

## Integración con el Sistema

### En Servicios

```go
type MiServicio struct {
    validator *services.ValidatorService
}

func (s *MiServicio) ProcesarDocumento(doc *models.DocumentoTributarioBasico) error {
    if err := s.validator.ValidarDocumento(doc); err != nil {
        return err
    }
    // Continuar con el procesamiento
    return nil
}
```

### En Handlers

```go
func handleDocumento(w http.ResponseWriter, r *http.Request) {
    doc := &models.DocumentoTributarioBasico{}
    if err := json.NewDecoder(r.Body).Decode(doc); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    validator := services.NewValidatorService("schema_dte")
    if err := validator.ValidarDocumento(doc); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Procesar documento válido
}
```

## Referencias

- [Documentación SII - Esquemas XML](https://www.sii.cl/factura_electronica/factura_mercado/schema.html)
- [Especificación XML Schema 1.0](https://www.w3.org/TR/xmlschema-1/)
- [Documentación libxml2](http://xmlsoft.org/) 