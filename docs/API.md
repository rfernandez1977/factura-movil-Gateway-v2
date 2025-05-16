# API Reference - FMgo

## Servicios

### ValidatorService

Servicio principal para la validación de documentos XML contra esquemas XSD.

#### Métodos

##### `NewValidatorService(basePath string) *ValidatorService`

Crea una nueva instancia del servicio de validación.

- **Parámetros**:
  - `basePath`: Ruta base donde se encuentran los esquemas XSD
- **Retorna**: Nueva instancia de `ValidatorService`

##### `CargarEsquema(nombre string) error`

Carga un esquema XSD específico.

- **Parámetros**:
  - `nombre`: Nombre del archivo XSD a cargar
- **Retorna**: Error si hay problemas al cargar el esquema

##### `ValidarXML(xml []byte, nombreEsquema string) error`

Valida un documento XML contra un esquema específico.

- **Parámetros**:
  - `xml`: Bytes del documento XML
  - `nombreEsquema`: Nombre del esquema contra el cual validar
- **Retorna**: Error si la validación falla

##### `ValidarDocumento(doc *models.DocumentoTributarioBasico) error`

Valida un documento tributario.

- **Parámetros**:
  - `doc`: Documento tributario a validar
- **Retorna**: Error si la validación falla

##### `ValidarEnvio(envio interface{}) error`

Valida un envío de documentos.

- **Parámetros**:
  - `envio`: Envío a validar
- **Retorna**: Error si la validación falla

##### `LimpiarEsquemas()`

Libera los esquemas cargados.

### ValidatorMiddleware

Middleware para la validación automática de documentos en handlers HTTP.

#### Métodos

##### `NewValidatorMiddleware(schemaPath string) *ValidatorMiddleware`

Crea una nueva instancia del middleware de validación.

- **Parámetros**:
  - `schemaPath`: Ruta a los esquemas XSD
- **Retorna**: Nueva instancia de `ValidatorMiddleware`

##### `WithDocumentoValidation(next func(ctx context.Context, doc *models.DocumentoTributarioBasico) error) func(ctx context.Context, doc *models.DocumentoTributarioBasico) error`

Envuelve un handler con validación de documento.

- **Parámetros**:
  - `next`: Handler a envolver
- **Retorna**: Handler envuelto con validación

##### `WithEnvioValidation(next func(ctx context.Context, envio interface{}) error) func(ctx context.Context, envio interface{}) error`

Envuelve un handler con validación de envío.

- **Parámetros**:
  - `next`: Handler a envolver
- **Retorna**: Handler envuelto con validación

## Modelos

### DocumentoTributarioBasico

```go
type DocumentoTributarioBasico struct {
    ID           string    
    TipoDTE      string    
    Folio        int       
    FechaEmision time.Time 
    RutEmisor    string    
    RutReceptor  string    
    MontoTotal   int       
    MontoNeto    int       
    MontoIVA     int       
    Estado       string    
}
```

### DTE

```go
type DTE struct {
    XMLName   xml.Name
    Version   string   
    Documento Documento
}
```

### Documento

```go
type Documento struct {
    XMLName    xml.Name   
    ID         string     
    Encabezado Encabezado 
    Detalle    []Detalle  
}
```

## Códigos de Error

| Código | Descripción |
|--------|-------------|
| E001 | Error al cargar esquema |
| E002 | Esquema no encontrado |
| E003 | Error al parsear XML |
| E004 | Error de validación XML |
| E005 | Tipo de documento no soportado |

## Ejemplos de Uso

### Validación Simple

```go
validator := services.NewValidatorService("schema_dte")
doc := &models.DocumentoTributarioBasico{
    ID:           "T33F1",
    TipoDTE:      "33",
    Folio:        1,
    FechaEmision: time.Now(),
    RutEmisor:    "76.123.456-7",
    RutReceptor:  "77.888.999-0",
    MontoTotal:   119000,
    MontoNeto:    100000,
    MontoIVA:     19000,
    Estado:       "PENDIENTE",
}

err := validator.ValidarDocumento(doc)
if err != nil {
    log.Fatalf("Error validando documento: %v", err)
}
```

### Uso con Middleware

```go
func configurarAPI(router *mux.Router) {
    middleware := middleware.NewValidatorMiddleware("schema_dte")
    
    router.HandleFunc("/documento", middleware.WithDocumentoValidation(handleDocumento))
    router.HandleFunc("/envio", middleware.WithEnvioValidation(handleEnvio))
}

func handleDocumento(ctx context.Context, doc *models.DocumentoTributarioBasico) error {
    // Procesar documento ya validado
    return nil
}
``` 