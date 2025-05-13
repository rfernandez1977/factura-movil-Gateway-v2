# Documentación de Servicios

## Servicios Principales

### DTEService
El servicio principal para el manejo de Documentos Tributarios Electrónicos (DTE).

```go
// Definición en services/dte_service.go
type DTEService struct {
    config  *config.SupabaseConfig
    db      *mongo.Database
    sii     sii.SIIClientInterface
    storage *StorageService
    empresa *EmpresaService
}
```

#### Métodos principales:
- `CrearDocumento`: Crea un nuevo documento DTE
- `FirmarDocumento`: Firma un documento DTE
- `EnviarDocumento`: Envía un documento DTE al SII
- `ConsultarEstado`: Consulta el estado de un documento en el SII

### DTEGenerator
Generador de documentos DTE y XML.

```go
// Definición en services/dte_generator.go
type DTEGenerator struct {
    caf *models.CAF
}
```

#### Métodos principales:
- `GenerarDTE`: Genera un documento DTE
- `GenerarDTEXML`: Genera un DTE en formato XML
- `GenerarSobre`: Genera un sobre para enviar al SII

### SIIClient
El servicio `SIIClient` maneja la comunicación con el Servicio de Impuestos Internos (SII).

```go
// Definición en services/sii/sii.go
type SIIClientInterface interface {
    ObtenerSemilla() (string, error)
    ObtenerToken(semilla string) (string, error)
    EnviarDTE(sobre *models.Sobre, token string) error
}
```

### CAFManager
El servicio `CAFManager` maneja los Códigos de Autorización de Folios (CAF).

```go
// Definición en services/caf/caf.go
type Manager struct {
    config *config.CAFConfig
}
```

#### Métodos principales:
- `CargarCAF`: Carga un CAF desde un archivo XML
- `ValidarCAF`: Valida que el CAF esté dentro del rango de folios
- `ObtenerCAF`: Obtiene un CAF por tipo de documento

## Flujos de Trabajo

### 1. Emisión de Documento Tributario
1. Validar CAF
2. Generar DTE
3. Firmar documento
4. Enviar al SII
5. Consultar estado

### 2. Consulta de Estado
1. Obtener semilla
2. Obtener token
3. Consultar estado
4. Procesar respuesta

### 3. Validación de CAF
1. Cargar CAF
2. Validar rango de folios
3. Verificar vigencia
4. Actualizar contador

## Manejo de Errores

### SIIService
```go
type SIIService interface {
    ConsultarEstado(trackID string) (*models.EstadoSII, error)
    EnviarDTE(dte []byte) (*models.EstadoSII, error)
    ConsultarDTE(tipoDTE, folio, rutEmisor string) (*models.EstadoSII, error)
    VerificarComunicacion() error
}
```

### Constantes y Enumeraciones

#### Tipos de DTE
```go
// Ver models/tipo_documento.go para la definición completa
type TipoDTE int

const (
    FacturaElectronica TipoDTE = 33
    BoletaElectronica TipoDTE = 39
    NotaCreditoElectronica TipoDTE = 61
    NotaDebitoElectronica TipoDTE = 56
    GuiaDespachoElectronica TipoDTE = 52
)
```

#### Tipos de Boleta
```