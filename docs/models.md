# Modelos del Sistema

## Tipos de Documentos Tributarios

### TipoDTE
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

### BOLETAType
```go
// Ver models/tipo_documento.go para la definición completa
type BOLETAType string

const (
    BOLETATypeElectronica BOLETAType = "ELECTRONICA"
    BOLETATypePapel      BOLETAType = "PAPEL"
)
```

## Documentos Tributarios

### DocumentoTributario
```go
// Ver models/documento_tributario.go para la definición completa
type DocumentoTributario struct {
    ID              string          `json:"id" bson:"_id,omitempty"`
    Tipo            string          `json:"tipo" bson:"tipo"`
    Folio           int             `json:"folio" bson:"folio"`
    FechaEmision    time.Time       `json:"fecha_emision" bson:"fecha_emision"`
    FechaVencimiento time.Time      `json:"fecha_vencimiento" bson:"fecha_vencimiento"`
    MontoTotal      float64         `json:"monto_total" bson:"monto_total"`
    Estado          EstadoDocumento `json:"estado" bson:"estado"`
    Emisor          Emisor          `json:"emisor" bson:"emisor"`
    Receptor        Receptor        `json:"receptor" bson:"receptor"`
    Detalles        []Detalle       `json:"detalles" bson:"detalles"`
    Impuestos       []Impuesto      `json:"impuestos" bson:"impuestos"`
    Referencias     []Referencia    `json:"referencias" bson:"referencias"`
    Timestamps      Timestamps      `json:"timestamps" bson:"timestamps"`
}
```

## Modelos XML

Todos los tipos XML están definidos en `models/xml.go`:

### DTEXMLModel
```go
type DTEXMLModel struct {
    XMLName     struct{} `xml:"DTE"`
    Version     string   `xml:"version,attr"`
    Documento   DocumentoXMLModel `xml:"Documento"`
    Signature   SignatureXMLModel `xml:"Signature,omitempty"`
}
```

### DocumentoXMLModel
```go
type DocumentoXMLModel struct {
    XMLName     struct{} `xml:"Documento"`
    Encabezado  EncabezadoXMLModel `xml:"Encabezado"`
    Detalle     []DetalleXMLModel `xml:"Detalle"`
    Referencias []ReferenciaXMLModel `xml:"Referencias>Referencia,omitempty"`
}
```

### EncabezadoXMLModel
```go
type EncabezadoXMLModel struct {
    XMLName     struct{} `xml:"Encabezado"`
    IDDocumento IDDocumentoXMLModel `xml:"IdDoc"`
    Emisor      EmisorXMLModel `xml:"Emisor"`
    Receptor    ReceptorXMLModel `xml:"Receptor"`
    Totales     TotalesXMLModel `xml:"Totales"`
}
```

### IDDocumentoXMLModel
```go
type IDDocumentoXMLModel struct {
    XMLName     struct{} `xml:"IdDoc"`
    TipoDTE     string   `xml:"TipoDTE"`
    Folio       int      `xml:"Folio"`
    FechaEmision string  `xml:"FchEmis"`
}
```

### EmisorXMLModel
```go
type EmisorXMLModel struct {
    XMLName     struct{} `xml:"Emisor"`
    RUT         string   `xml:"RUTEmisor"`
    RazonSocial string   `xml:"RznSoc"`
    Giro        string   `xml:"GiroEmis"`
    Direccion   string   `xml:"DirOrigen"`
    Comuna      string   `xml:"CmnaOrigen"`
    Ciudad      string   `xml:"CiudadOrigen"`
}
```

### ReceptorXMLModel
```go
type ReceptorXMLModel struct {
    XMLName     struct{} `xml:"Receptor"`
    RUT         string   `xml:"RUTRecep"`
    RazonSocial string   `xml:"RznSocRecep"`
    Giro        string   `xml:"GiroRecep"`
    Direccion   string   `xml:"DirRecep"`
    Comuna      string   `xml:"CmnaRecep"`
    Ciudad      string   `xml:"CiudadRecep"`
}
```

### TotalesXMLModel
```go
type TotalesXMLModel struct {
    XMLName     struct{} `xml:"Totales"`
    MontoNeto   float64  `xml:"MntNeto"`
    MontoExento float64  `xml:"MntExe"`
    IVA         float64  `xml:"IVA"`
    MontoTotal  float64  `xml:"MntTotal"`
}
```

### DetalleXMLModel
```go
type DetalleXMLModel struct {
    XMLName       struct{} `xml:"Detalle"`
    NmbItem       string   `xml:"NmbItem"`
    QtyItem       float64  `xml:"QtyItem"`
    PrcItem       float64  `xml:"PrcItem"`
    MontoItem     float64  `xml:"MontoItem"`
    DescuentoPct  float64  `xml:"DescuentoPct,omitempty"`
    DescuentoMonto float64 `xml:"DescuentoMonto,omitempty"`
    Impuestos     []ImpuestoXMLModel `xml:"Impuestos>Impuesto,omitempty"`
}
```

### ImpuestoXMLModel
```go
type ImpuestoXMLModel struct {
    XMLName     struct{} `xml:"Impuesto"`
    Tipo        string   `xml:"TipoImp"`
    Tasa        float64  `xml:"TasaImp"`
    Monto       float64  `xml:"MontoImp"`
}
```

### ReferenciaXMLModel
```go
type ReferenciaXMLModel struct {
    XMLName     struct{} `xml:"Referencia"`
    Tipo        string   `xml:"TpoDocRef"`
    Folio       int      `xml:"FolioRef"`
    Fecha       string   `xml:"FchRef"`
    Codigo      string   `xml:"CodRef,omitempty"`
    Razon       string   `xml:"RazonRef,omitempty"`
}
```

### SignatureXMLModel
```go
type SignatureXMLModel struct {
    XMLName     struct{} `xml:"Signature"`
    SignedInfo  SignedInfoXMLModel `xml:"SignedInfo"`
    SignatureValue string `xml:"SignatureValue"`
    KeyInfo     KeyInfoXMLModel `xml:"KeyInfo"`
}
```

### SignedInfoXMLModel
```go
type SignedInfoXMLModel struct {
    XMLName     struct{} `xml:"SignedInfo"`
    CanonicalizationMethod CanonicalizationMethodXMLModel `xml:"CanonicalizationMethod"`
    SignatureMethod SignatureMethodXMLModel `xml:"SignatureMethod"`
    Reference    ReferenceXMLModel `xml:"Reference"`
}
```

### CanonicalizationMethodXMLModel
```go
type CanonicalizationMethodXMLModel struct {
    XMLName     struct{} `xml:"CanonicalizationMethod"`
    Algorithm   string   `xml:"Algorithm,attr"`
}
```

### SignatureMethodXMLModel
```go
type SignatureMethodXMLModel struct {
    XMLName     struct{} `xml:"SignatureMethod"`
    Algorithm   string   `xml:"Algorithm,attr"`
}
```

### ReferenceXMLModel
```go
type ReferenceXMLModel struct {
    XMLName     struct{} `xml:"Reference"`
    URI         string   `xml:"URI,attr"`
    Transforms  TransformsXMLModel `xml:"Transforms"`
    DigestMethod DigestMethodXMLModel `xml:"DigestMethod"`
    DigestValue string   `xml:"DigestValue"`
}
```

### TransformsXMLModel
```go
type TransformsXMLModel struct {
    XMLName     struct{} `xml:"Transforms"`
    Transform   []TransformXMLModel `xml:"Transform"`
}
```

### TransformXMLModel
```go
type TransformXMLModel struct {
    XMLName     struct{} `xml:"Transform"`
    Algorithm   string   `xml:"Algorithm,attr"`
}
```

### DigestMethodXMLModel
```go
type DigestMethodXMLModel struct {
    XMLName     struct{} `xml:"DigestMethod"`
    Algorithm   string   `xml:"Algorithm,attr"`
}
```

### KeyInfoXMLModel
```go
type KeyInfoXMLModel struct {
    XMLName     struct{} `xml:"KeyInfo"`
    X509Data    X509DataXMLModel `xml:"X509Data"`
}
```

### X509DataXMLModel
```go
type X509DataXMLModel struct {
    XMLName     struct{} `xml:"X509Data"`
    X509Certificate string `xml:"X509Certificate"`
}
```

### Funciones Auxiliares
```go
func NewDTEXMLModel(version string, documento DocumentoXMLModel) *DTEXMLModel
func NewDocumentoXMLModel(encabezado EncabezadoXMLModel, detalles []DetalleXMLModel) *DocumentoXMLModel
func NewEncabezadoXMLModel(idDocumento IDDocumentoXMLModel, emisor EmisorXMLModel, receptor ReceptorXMLModel, totales TotalesXMLModel) *EncabezadoXMLModel
func NewIDDocumentoXMLModel(tipoDTE string, folio int, fechaEmision time.Time) *IDDocumentoXMLModel
func NewEmisorXMLModel(rut, razonSocial, giro, direccion, comuna, ciudad string) *EmisorXMLModel
func NewReceptorXMLModel(rut, razonSocial, giro, direccion, comuna, ciudad string) *ReceptorXMLModel
func NewTotalesXMLModel(montoNeto, montoExento, iva, montoTotal float64) *TotalesXMLModel
func NewDetalleXMLModel(nmbItem string, qtyItem, prcItem float64) *DetalleXMLModel
func NewImpuestoXMLModel(tipo string, tasa, monto float64) *ImpuestoXMLModel
func NewReferenciaXMLModel(tipo string, folio int, fecha string) *ReferenciaXMLModel
func NewSignatureXMLModel(signedInfo SignedInfoXMLModel, signatureValue string, keyInfo KeyInfoXMLModel) *SignatureXMLModel
func NewSignedInfoXMLModel(canonicalizationMethod CanonicalizationMethodXMLModel, signatureMethod SignatureMethodXMLModel, reference ReferenceXMLModel) *SignedInfoXMLModel
func NewCanonicalizationMethodXMLModel(algorithm string) *CanonicalizationMethodXMLModel
func NewSignatureMethodXMLModel(algorithm string) *SignatureMethodXMLModel
func NewReferenceXMLModel(uri string, transforms TransformsXMLModel, digestMethod DigestMethodXMLModel, digestValue string) *ReferenceXMLModel
func NewTransformsXMLModel(transforms []TransformXMLModel) *TransformsXMLModel
func NewTransformXMLModel(algorithm string) *TransformXMLModel
func NewDigestMethodXMLModel(algorithm string) *DigestMethodXMLModel
func NewKeyInfoXMLModel(x509Data X509DataXMLModel) *KeyInfoXMLModel
func NewX509DataXMLModel(x509Certificate string) *X509DataXMLModel
```

## Modelos CAF

### CAF (Modelo Principal)
```go
// Ver models/caf.go para la definición completa
type CAF struct {
    ID               string    `json:"id"`
    TipoDocumento    string    `json:"tipo_documento"`
    RutEmisor        string    `json:"rut_emisor"`
    RangoInicial     int       `json:"rango_inicial"`
    RangoFinal       int       `json:"rango_final"`
    FechaVencimiento time.Time `json:"fecha_vencimiento"`
    Activo           bool      `json:"activo"`
    Estado           string    `json:"estado"`
    FolioActual      int       `json:"folio_actual"`
    UpdatedAt        time.Time `json:"updated_at"`
}
```

### CAFRequest (Solicitud Interna)
```go
// Ver models/caf.go para la definición completa
type CAFRequest struct {
    TipoDocumento string `json:"tipo_documento"`
    RutEmisor     string `json:"rut_emisor"`
    Cantidad      int    `json:"cantidad"`
}
```

### SIICAFRequest (Solicitud SII)
```go
// Ver services/caf_service.go para la definición completa
type SIICAFRequest struct {
    RUTEmisor      string
    TipoDTE        string
    FolioInicial   int
    FolioFinal     int
    FechaSolicitud time.Time
}
```

### SIICAFResponse (Respuesta SII)
```go
// Ver services/caf_service.go para la definición completa
type SIICAFResponse struct {
    Estado         string    `xml:"ESTADO"`
    Glosa          string    `xml:"GLOSA"`
    TrackID        string    `xml:"TRACKID"`
    FechaRespuesta time.Time `xml:"FECHARESPUESTA"`
    URLDescarga    string    `xml:"URLDESCARGA"`
}
```

### CAFMetadata (Metadatos Archivo)
```go
// Ver services/caf_service.go para la definición completa
type CAFMetadata struct {
    RUTEmisor        string
    TipoDTE          string
    FolioInicial     int
    FolioFinal       int
    FechaEmision     time.Time
    FechaVencimiento time.Time
    Estado           string
    Hash             string
}
```

### SIICAFXML (XML SII)
```go
// Ver services/caf_service.go para la definición completa
type SIICAFXML struct {
    XMLName           xml.Name `xml:"AUTORIZACION"`
    Version           string   `xml:"CAF>version,attr"`
    RUTEmisor         string   `xml:"CAF>DA>RE"`
    RazonSocial       string   `xml:"CAF>DA>RS"`
    TipoDTE           string   `xml:"CAF>DA>TD"`
    FolioInicial      int      `xml:"CAF>DA>RNG>D"`
    FolioFinal        int      `xml:"CAF>DA>RNG>H"`
    FechaAutorizacion string   `xml:"CAF>DA>FA"`
    Modulo            string   `xml:"CAF>DA>RSAPK>M"`
    Exponente         string   `xml:"CAF>DA>RSAPK>E"`
    IDK               string   `xml:"CAF>DA>IDK"`
    Firma             string   `xml:"CAF>FRMA"`
    PrivateKey        string   `xml:"RSASK"`
    PublicKey         string   `xml:"RSAPUBK"`
}
```

### CAFXML (XML DTE)
```go
// Ver models/xml.go para la definición completa
type CAFXML struct {
    Version string `xml:"version,attr"`
    DA      DAXML  `xml:"DA"`
    FRMA    string `xml:"FRMA"`
}
```

## Modelos SII

Los siguientes tipos están definidos en `models/sii_types.go` y son utilizados para la comunicación con el SII:

### RespuestaSII
```go
// Definición en models/sii_types.go
type RespuestaSII struct {
    XMLName         struct{}          `xml:"RespuestaDTE"`
    Version         string            `xml:"version,attr"`
    Estado          string            `xml:"Estado"`
    Glosa           string            `xml:"Glosa"`
    TrackID         string            `xml:"TrackID"`
    NumeroAtencion  string            `xml:"NumeroAtencion"`
    FechaProceso    time.Time         `xml:"FechaProceso"`
    FechaRecepcion  time.Time         `xml:"FechaRecepcion"`
    FechaAceptacion time.Time         `xml:"FechaAceptacion"`
    FechaRechazo    time.Time         `xml:"FechaRechazo"`
    Errores         []ErrorSII        `xml:"Errores>Error"`
    Detalles        map[string]string `xml:"Detalles>Detalle"`
    Estadistica     *EstadisticaSII   `xml:"Estadistica,omitempty"`
}
```

### EstadoSII
```go
// Definición en models/sii_types.go
type EstadoSII struct {
    XMLName struct{}   `xml:"Estado"`
    Version string     `xml:"version,attr"`
    Estado  string     `xml:"Estado"`
    Glosa   string     `xml:"Glosa"`
    TrackID string     `xml:"TrackID"`
    Fecha   time.Time  `xml:"Fecha,omitempty"`
    Errores []ErrorSII `xml:"Errores>Error,omitempty"`
}
```

### ErrorSII
```go
// Definición en models/sii_types.go
type ErrorSII struct {
    XMLName     struct{} `xml:"Error"`
    Codigo      string   `xml:"Codigo"`
    Descripcion string   `xml:"Descripcion"`
    Detalle     string   `xml:"Detalle,omitempty"`
}
```

### DocumentoRechazado
```go
// Definición en models/sii_types.go
type DocumentoRechazado struct {
    XMLName      struct{}  `xml:"DocumentoRechazado"`
    Folio        int       `xml:"Folio"`
    TipoDoc      string    `xml:"TipoDoc"`
    RutEmisor    string    `xml:"RutEmisor"`
    RutReceptor  string    `xml:"RutReceptor"`
    Motivo       string    `xml:"Motivo"`
    FechaRechazo time.Time `xml:"FechaRechazo"`
}
```

### InformacionContribuyente
```go
// Definición en models/sii_types.go
type InformacionContribuyente struct {
    RUT         string `xml:"RUT"`
    RazonSocial string `xml:"RazonSocial"`
    Giro        string `xml:"Giro"`
    Direccion   string `xml:"Direccion"`
    Comuna      string `xml:"Comuna"`
    Ciudad      string `xml:"Ciudad"`
}
```

### EstadoContribuyente
```go
// Definición en models/sii_types.go
type EstadoContribuyente struct {
    Estado string `xml:"Estado"`
    Glosa  string `xml:"Glosa"`
}
```

### EstadoDocumentoSII
```go
// Definición en models/sii_types.go
type EstadoDocumentoSII struct {
    TipoDocto    string `xml:"TIPO_DOCTO"`
    Folio        int    `xml:"FOLIO"`
    FechaEmision string `xml:"FECHA_EMISION"`
    Estado       string `xml:"ESTADO"`
    Glosa        string `xml:"GLOSA"`
}
```

### RespuestaEnvioDTE
```go
// Definición en models/sii_types.go
type RespuestaEnvioDTE struct {
    XMLName        struct{} `xml:"RespuestaDTE"`
    Version        string   `xml:"version,attr"`
    Identificacion struct {
        RUTEmisor     string `xml:"RUTEMISOR"`
        RUTEnvia      string `xml:"RUTENVIA"`
        TrackID       int    `xml:"TRACKID"`
        TMSTRecepcion string `xml:"TMSTRECEPCION"`
        Estado        string `xml:"ESTADO"`
    } `xml:"IDENTIFICACION"`
    ErrorEnvio *struct {
        DetErrorEnvio []string `xml:"DETERRENVIO"`
    } `xml:"ERRORENVIO,omitempty"`
}
```

### ResumenContribuyente
```go
// Definición en models/sii_types.go
type ResumenContribuyente struct {
    RUT         string `json:"rut" xml:"RUT"`
    RazonSocial string `json:"razon_social" xml:"RazonSocial"`
    Estado      string `json:"estado" xml:"Estado"`
    Glosa       string `json:"glosa" xml:"Glosa"`
}
```

### SobreDTE
```go
// Definición en models/sii_types.go
type SobreDTE struct {
    XMLName    xml.Name              `xml:"SobreDTE" json:"-"`
    Version    string                `xml:"version,attr" json:"version"`
    Documentos []DocumentoTributario `xml:"Documentos>Documento" json:"documentos"`
    TrackID    string                `xml:"TrackID" json:"trackId"`
    FechaEnvio time.Time             `xml:"FechaEnvio" json:"fechaEnvio"`
}
```

## Modelos de Entidades

### Emisor
```go
// Ver models/emisor.go para la definición completa
type Emisor struct {
    RUT         string `json:"rut" xml:"RUT"`
    RazonSocial string `json:"razon_social" xml:"RazonSocial"`
    Giro        string `json:"giro" xml:"Giro"`
    Direccion   string `json:"direccion" xml:"Direccion"`
    Comuna      string `json:"comuna" xml:"Comuna"`
    Ciudad      string `json:"ciudad" xml:"Ciudad"`
}
```

### Receptor
```go
// Ver models/receptor.go para la definición completa
type Receptor struct {
    RUT         string `json:"rut" xml:"RUT"`
    RazonSocial string `json:"razon_social" xml:"RazonSocial"`
    Giro        string `json:"giro" xml:"Giro"`
    Direccion   string `json:"direccion" xml:"Direccion"`
    Comuna      string `json:"comuna" xml:"Comuna"`
    Ciudad      string `json:"ciudad" xml:"Ciudad"`
}
```

## Tipos de Validación

Todos los tipos de validación están definidos en `models/validation.go`:

### ValidationRule
```go
type ValidationRule struct {
    ID          string    `json:"id"`
    Nombre      string    `json:"nombre"`
    Descripcion string    `json:"descripcion"`
    Tipo        string    `json:"tipo"`
    Expresion   string    `json:"expresion"`
    Mensaje     string    `json:"mensaje"`
    Activo      bool      `json:"activo"`
    CreatedAt   time.Time `json:"createdAt"`
    UpdatedAt   time.Time `json:"updatedAt"`
}
```

### ValidationResult
```go
type ValidationResult struct {
    ID          string    `json:"id"`
    DocumentoID string    `json:"documentoId"`
    ReglaID     string    `json:"reglaId"`
    Exitoso     bool      `json:"exitoso"`
    Mensaje     string    `json:"mensaje"`
    Detalles    string    `json:"detalles"`
    CreatedAt   time.Time `json:"createdAt"`
}
```

### ValidationError
```go
type ValidationError struct {
    ID          string    `json:"id"`
    DocumentoID string    `json:"documentoId"`
    ReglaID     string    `json:"reglaId"`
    Codigo      string    `json:"codigo"`
    Mensaje     string    `json:"mensaje"`
    Detalles    string    `json:"detalles"`
    CreatedAt   time.Time `json:"createdAt"`
}
```

### ValidationStatus
```go
type ValidationStatus struct {
    ID          string    `json:"id"`
    DocumentoID string    `json:"documentoId"`
    Estado      string    `json:"estado"`
    TotalReglas int       `json:"totalReglas"`
    Exitosas    int       `json:"exitosas"`
    Fallidas    int       `json:"fallidas"`
    CreatedAt   time.Time `json:"createdAt"`
    UpdatedAt   time.Time `json:"updatedAt"`
}
```

### ValidationType
```go
type ValidationType struct {
    ID          string    `json:"id"`
    Nombre      string    `json:"nombre"`
    Descripcion string    `json:"descripcion"`
    Reglas      []string  `json:"reglas"`
    Activo      bool      `json:"activo"`
    CreatedAt   time.Time `json:"createdAt"`
    UpdatedAt   time.Time `json:"updatedAt"`
}
```

### ValidationConfig
```go
type ValidationConfig struct {
    ID            string    `json:"id"`
    Tipo          string    `json:"tipo"`
    Reglas        []string  `json:"reglas"`
    MaxErrores    int       `json:"maxErrores"`
    StopOnError   bool      `json:"stopOnError"`
    CreatedAt     time.Time `json:"createdAt"`
    UpdatedAt     time.Time `json:"updatedAt"`
}
```

### ValidationRequest
```go
type ValidationRequest struct {
    ID          string    `json:"id"`
    DocumentoID string    `json:"documentoId"`
    Tipo        string    `json:"tipo"`
    Config      ValidationConfig `json:"config"`
    CreatedAt   time.Time `json:"createdAt"`
}
```

### ValidationResponse
```go
type ValidationResponse struct {
    ID          string    `json:"id"`
    RequestID   string    `json:"requestId"`
    DocumentoID string    `json:"documentoId"`
    Estado      string    `json:"estado"`
    Resultados  []ValidationResult `json:"resultados"`
    Errores     []ValidationError `json:"errores"`
    CreatedAt   time.Time `json:"createdAt"`
}
```

### ValidationMetadata
```go
type ValidationMetadata struct {
    ID          string    `json:"id"`
    DocumentoID string    `json:"documentoId"`
    Version     string    `json:"version"`
    Timestamp   time.Time `json:"timestamp"`
    Usuario     string    `json:"usuario"`
    Sistema     string    `json:"sistema"`
}
```

### Funciones Auxiliares
```go
func NewValidationRule(nombre, descripcion, tipo, expresion, mensaje string) *ValidationRule
func NewValidationResult(documentoID, reglaID string, exitoso bool, mensaje, detalles string) *ValidationResult
func NewValidationError(documentoID, reglaID, codigo, mensaje, detalles string) *ValidationError
func NewValidationStatus(documentoID string, totalReglas int) *ValidationStatus
func NewValidationType(nombre, descripcion string, reglas []string) *ValidationType
func NewValidationConfig(tipo string, reglas []string) *ValidationConfig
func NewValidationRequest(documentoID, tipo string, config ValidationConfig) *ValidationRequest
func NewValidationResponse(requestID, documentoID string) *ValidationResponse
func NewValidationMetadata(documentoID, version, usuario, sistema string) *ValidationMetadata
```

## Modelos de Estado

Todos los tipos de estado están definidos en `models/estados.go`:

### EstadoSII
```go
type EstadoSII struct {
    Codigo      int       `json:"codigo"`
    Descripcion string    `json:"descripcion"`
    Timestamp   time.Time `json:"timestamp"`
}
```

### EstadoDocumento
```go
type EstadoDocumento string

const (
    EstadoDocumentoPendiente    EstadoDocumento = "PENDIENTE"
    EstadoDocumentoProcesando   EstadoDocumento = "PROCESANDO"
    EstadoDocumentoCompletado   EstadoDocumento = "COMPLETADO"
    EstadoDocumentoError        EstadoDocumento = "ERROR"
    EstadoDocumentoRechazado    EstadoDocumento = "RECHAZADO"
    EstadoDocumentoAnulado      EstadoDocumento = "ANULADO"
    EstadoDocumentoEnviado      EstadoDocumento = "ENVIADO"
    EstadoDocumentoAceptado     EstadoDocumento = "ACEPTADO"
    EstadoDocumentoRechazadoSII EstadoDocumento = "RECHAZADO_SII"
)
```

### EstadoFlujo
```go
type EstadoFlujo string

const (
    EstadoFlujoPendiente   EstadoFlujo = "PENDIENTE"
    EstadoFlujoProcesando  EstadoFlujo = "PROCESANDO"
    EstadoFlujoCompletado  EstadoFlujo = "COMPLETADO"
    EstadoFlujoError       EstadoFlujo = "ERROR"
    EstadoFlujoCancelado   EstadoFlujo = "CANCELADO"
    EstadoFlujoPausado     EstadoFlujo = "PAUSADO"
    EstadoFlujoReanudado   EstadoFlujo = "REANUDADO"
)
```

### EstadoPaso
```go
type EstadoPaso string

const (
    EstadoPasoPendiente   EstadoPaso = "PENDIENTE"
    EstadoPasoProcesando  EstadoPaso = "PROCESANDO"
    EstadoPasoCompletado  EstadoPaso = "COMPLETADO"
    EstadoPasoError       EstadoPaso = "ERROR"
    EstadoPasoCancelado   EstadoPaso = "CANCELADO"
    EstadoPasoPausado     EstadoPaso = "PAUSADO"
    EstadoPasoReanudado   EstadoPaso = "REANUDADO"
)
```

### EstadoNotificacion
```go
type EstadoNotificacion string

const (
    EstadoNotificacionPendiente   EstadoNotificacion = "PENDIENTE"
    EstadoNotificacionEnviada     EstadoNotificacion = "ENVIADA"
    EstadoNotificacionEntregada   EstadoNotificacion = "ENTREGADA"
    EstadoNotificacionLeida       EstadoNotificacion = "LEIDA"
    EstadoNotificacionError       EstadoNotificacion = "ERROR"
    EstadoNotificacionCancelada   EstadoNotificacion = "CANCELADA"
)
```

### EstadoIntegracionERP
```go
type EstadoIntegracionERP string

const (
    EstadoIntegracionERPPendiente   EstadoIntegracionERP = "PENDIENTE"
    EstadoIntegracionERPProcesando  EstadoIntegracionERP = "PROCESANDO"
    EstadoIntegracionERPCompletado  EstadoIntegracionERP = "COMPLETADO"
    EstadoIntegracionERPError       EstadoIntegracionERP = "ERROR"
    EstadoIntegracionERPCancelado   EstadoIntegracionERP = "CANCELADO"
)
```

### EstadoSesion
```go
type EstadoSesion string

const (
    EstadoSesionActiva     EstadoSesion = "ACTIVA"
    EstadoSesionExpirada   EstadoSesion = "EXPIRADA"
    EstadoSesionCerrada    EstadoSesion = "CERRADA"
    EstadoSesionBloqueada  EstadoSesion = "BLOQUEADA"
)
```

### EstadoCAF
```go
type EstadoCAF string

const (
    EstadoCAFActivo     EstadoCAF = "ACTIVO"
    EstadoCAFAgotado    EstadoCAF = "AGOTADO"
    EstadoCAFExpirado   EstadoCAF = "EXPIRADO"
    EstadoCAFAnulado    EstadoCAF = "ANULADO"
    EstadoCAFPendiente  EstadoCAF = "PENDIENTE"
)
```

### EstadoDTE
```go
type EstadoDTE string

const (
    EstadoDTEPendiente    EstadoDTE = "PENDIENTE"
    EstadoDTEProcesando   EstadoDTE = "PROCESANDO"
    EstadoDTECompletado   EstadoDTE = "COMPLETADO"
    EstadoDTEError        EstadoDTE = "ERROR"
    EstadoDTERechazado    EstadoDTE = "RECHAZADO"
    EstadoDTEAnulado      EstadoDTE = "ANULADO"
    EstadoDTEEnviado      EstadoDTE = "ENVIADO"
    EstadoDTEAceptado     EstadoDTE = "ACEPTADO"
    EstadoDTERechazadoSII EstadoDTE = "RECHAZADO_SII"
)
```

### Funciones Auxiliares
```go
func NewEstadoSII(codigo int, descripcion string) *EstadoSII
```

## Modelos de Error

Todos los tipos de error están definidos en `models/errors.go`:

### ErrorSII
```go
type ErrorSII struct {
    Codigo    int       `json:"codigo"`
    Mensaje   string    `json:"mensaje"`
    Timestamp time.Time `json:"timestamp"`
}
```

### ErrorValidacion
```go
type ErrorValidacion struct {
    Campo   string `json:"campo"`
    Mensaje string `json:"mensaje"`
}
```

### ErrorSistema
```go
type ErrorSistema struct {
    Codigo    int       `json:"codigo"`
    Mensaje   string    `json:"mensaje"`
    Timestamp time.Time `json:"timestamp"`
}
```

### ErrorIntegracion
```go
type ErrorIntegracion struct {
    Sistema  string    `json:"sistema"`
    Codigo   int       `json:"codigo"`
    Mensaje  string    `json:"mensaje"`
    Timestamp time.Time `json:"timestamp"`
}

### ErrorResponse
```go
type ErrorResponse struct {
    Error   string `json:"error"`
    Message string `json:"message"`
    Code    int    `json:"code"`
}

### ErrorLog
```go
type ErrorLog struct {
    ID        string    `json:"id" bson:"_id,omitempty"`
    Tipo      string    `json:"tipo" bson:"tipo"`
    Mensaje   string    `json:"mensaje" bson:"mensaje"`
    Detalles  string    `json:"detalles" bson:"detalles"`
    Timestamp time.Time `json:"timestamp" bson:"timestamp"`
    Usuario   string    `json:"usuario" bson:"usuario"`
}

### Constantes
```go
const (
    ErrorTipoSII         = "SII"
    ErrorTipoValidacion  = "VALIDACION"
    ErrorTipoSistema     = "SISTEMA"
    ErrorTipoIntegracion = "INTEGRACION"
)
```

### Funciones Auxiliares
```go
func NewErrorSII(codigo int, mensaje string) *ErrorSII
func NewErrorValidacion(campo, mensaje string) *ErrorValidacion
func NewErrorSistema(codigo int, mensaje string) *ErrorSistema
func NewErrorIntegracion(sistema string, codigo int, mensaje string) *ErrorIntegracion
func NewErrorResponse(error string, message string, code int) *ErrorResponse
func NewErrorLog(tipo, mensaje, detalles, usuario string) *ErrorLog
```

## Modelos de Documento

Los siguientes tipos de documento están definidos en sus respectivos archivos:

### DocumentoTributario
```go
// Definición en models/documento_tributario.go
type DocumentoTributario struct {
    TipoDTE         TipoDTE         `json:"tipo_dte" xml:"TipoDTE"`
    Folio           int             `json:"folio" xml:"Folio"`
    FechaEmision    time.Time       `json:"fecha_emision" xml:"FechaEmision"`
    Emisor          Emisor          `json:"emisor" xml:"Emisor"`
    Receptor        Receptor        `json:"receptor" xml:"Receptor"`
    MontosImpuestos MontosImpuestos `json:"montos_impuestos" xml:"MontosImpuestos"`
    Items           []Item          `json:"items" xml:"Items>Item"`
    Referencias     []Referencia    `json:"referencias" xml:"Referencias>Referencia,omitempty"`
}
```

### Documento
```go
// Definición en models/documento.go
type Documento struct {
    ID              string    `json:"id" bson:"_id"`
    Tipo            string    `json:"tipo" bson:"tipo"`
    Estado          string    `json:"estado" bson:"estado"`
    FechaCreacion   time.Time `json:"fecha_creacion" bson:"fecha_creacion"`
    FechaActualizacion time.Time `json:"fecha_actualizacion" bson:"fecha_actualizacion"`
    Metadata        map[string]interface{} `json:"metadata" bson:"metadata"`
}
```

### DTEDocument
```go
// Definición en models/dte_document.go
type DTEDocument struct {
    ID              string    `json:"id" bson:"_id"`
    TipoDTE         string    `json:"tipo_dte" bson:"tipo_dte"`
    Folio           int       `json:"folio" bson:"folio"`
    FechaEmision    time.Time `json:"fecha_emision" bson:"fecha_emision"`
    Estado          string    `json:"estado" bson:"estado"`
    XML             string    `json:"xml" bson:"xml"`
    PDF             string    `json:"pdf" bson:"pdf"`
    Firma           string    `json:"firma" bson:"firma"`
    TrackID         string    `json:"track_id" bson:"track_id"`
    FechaCreacion   time.Time `json:"fecha_creacion" bson:"fecha_creacion"`
    FechaActualizacion time.Time `json:"fecha_actualizacion" bson:"fecha_actualizacion"`
}
```

### DocumentoAlmacenado
```go
// Definición en models/documento_almacenado.go
type DocumentoAlmacenado struct {
    ID              string    `json:"id" bson:"_id"`
    Tipo            string    `json:"tipo" bson:"tipo"`
    Ruta            string    `json:"ruta" bson:"ruta"`
    NombreArchivo   string    `json:"nombre_archivo" bson:"nombre_archivo"`
    Tamaño          int64     `json:"tamaño" bson:"tamaño"`
    Hash            string    `json:"hash" bson:"hash"`
    FechaCreacion   time.Time `json:"fecha_creacion" bson:"fecha_creacion"`
    FechaActualizacion time.Time `json:"fecha_actualizacion" bson:"fecha_actualizacion"`
}

## Modelos de Configuración

Todos los tipos de configuración están definidos en `models/config.go`:

### Config
```go
type Config struct {
    // Configuración del servidor
    Server struct {
        Port         int    `json:"port"`
        Host         string `json:"host"`
        ReadTimeout  int    `json:"readTimeout"`
        WriteTimeout int    `json:"writeTimeout"`
    } `json:"server"`

    // Configuración de la base de datos
    Database struct {
        Host     string `json:"host"`
        Port     int    `json:"port"`
        User     string `json:"user"`
        Password string `json:"password"`
        Name     string `json:"name"`
        SSLMode  string `json:"sslMode"`
    } `json:"database"`

    // Configuración del SII
    SII struct {
        BaseURL      string `json:"baseUrl"`
        Timeout      int    `json:"timeout"`
        RetryCount   int    `json:"retryCount"`
        RetryDelay   int    `json:"retryDelay"`
        CertPath     string `json:"certPath"`
        KeyPath      string `json:"keyPath"`
        CertPassword string `json:"certPassword"`
    } `json:"sii"`

    // Configuración de logs
    Logging struct {
        Level      string `json:"level"`
        FilePath   string `json:"filePath"`
        MaxSize    int    `json:"maxSize"`
        MaxBackups int    `json:"maxBackups"`
        MaxAge     int    `json:"maxAge"`
        Compress   bool   `json:"compress"`
    } `json:"logging"`

    // Configuración de caché
    Cache struct {
        Enabled  bool   `json:"enabled"`
        Type     string `json:"type"`
        Address  string `json:"address"`
        Password string `json:"password"`
        DB       int    `json:"db"`
    } `json:"cache"`

    // Configuración de seguridad
    Security struct {
        JWTSecret     string `json:"jwtSecret"`
        JWTExpiration int    `json:"jwtExpiration"`
        CORSEnabled   bool   `json:"corsEnabled"`
        CORSOrigins   []string `json:"corsOrigins"`
    } `json:"security"`

    // Configuración de notificaciones
    Notifications struct {
        Email struct {
            Enabled  bool   `json:"enabled"`
            Host     string `json:"host"`
            Port     int    `json:"port"`
            Username string `json:"username"`
            Password string `json:"password"`
            From     string `json:"from"`
        } `json:"email"`
        SMS struct {
            Enabled  bool   `json:"enabled"`
            Provider string `json:"provider"`
            APIKey   string `json:"apiKey"`
            From     string `json:"from"`
        } `json:"sms"`
    } `json:"notifications"`

    // Configuración de almacenamiento
    Storage struct {
        Type     string `json:"type"`
        Local    struct {
            BasePath string `json:"basePath"`
        } `json:"local"`
        S3 struct {
            Bucket    string `json:"bucket"`
            Region    string `json:"region"`
            AccessKey string `json:"accessKey"`
            SecretKey string `json:"secretKey"`
        } `json:"s3"`
    } `json:"storage"`

    // Configuración de integración con ERP
    ERP struct {
        Enabled  bool   `json:"enabled"`
        Type     string `json:"type"`
        BaseURL  string `json:"baseUrl"`
        APIKey   string `json:"apiKey"`
        Timeout  int    `json:"timeout"`
        RetryCount int  `json:"retryCount"`
    } `json:"erp"`

    // Configuración de validación
    Validation struct {
        Enabled     bool     `json:"enabled"`
        Rules       []string `json:"rules"`
        MaxErrors   int      `json:"maxErrors"`
        StopOnError bool     `json:"stopOnError"`
    } `json:"validation"`

    // Configuración de monitoreo
    Monitoring struct {
        Enabled    bool   `json:"enabled"`
        Type       string `json:"type"`
        Endpoint   string `json:"endpoint"`
        APIKey     string `json:"apiKey"`
        Environment string `json:"environment"`
    } `json:"monitoring"`
}
```

### ConfiguracionSII
```go
type ConfiguracionSII struct {
    ID              string    `json:"id"`
    RUTEmisor       string    `json:"rutEmisor"`
    CertificadoPath string    `json:"certificadoPath"`
    ClavePrivada    string    `json:"clavePrivada"`
    Ambiente        string    `json:"ambiente"`
    Timeout         int       `json:"timeout"`
    RetryCount      int       `json:"retryCount"`
    RetryDelay      int       `json:"retryDelay"`
    CreatedAt       time.Time `json:"createdAt"`
    UpdatedAt       time.Time `json:"updatedAt"`
}
```

### ConfiguracionERP
```go
type ConfiguracionERP struct {
    ID          string    `json:"id"`
    Tipo        string    `json:"tipo"`
    BaseURL     string    `json:"baseUrl"`
    APIKey      string    `json:"apiKey"`
    Timeout     int       `json:"timeout"`
    RetryCount  int       `json:"retryCount"`
    CreatedAt   time.Time `json:"createdAt"`
    UpdatedAt   time.Time `json:"updatedAt"`
}
```

### ConfiguracionValidacion
```go
type ConfiguracionValidacion struct {
    ID            string    `json:"id"`
    Tipo          string    `json:"tipo"`
    Reglas        []string  `json:"reglas"`
    MaxErrores    int       `json:"maxErrores"`
    StopOnError   bool      `json:"stopOnError"`
    CreatedAt     time.Time `json:"createdAt"`
    UpdatedAt     time.Time `json:"updatedAt"`
}
```

### ConfiguracionNotificacion
```go
type ConfiguracionNotificacion struct {
    ID          string    `json:"id"`
    Tipo        string    `json:"tipo"`
    Destinatario string   `json:"destinatario"`
    Template    string    `json:"template"`
    Activo      bool      `json:"activo"`
    CreatedAt   time.Time `json:"createdAt"`
    UpdatedAt   time.Time `json:"updatedAt"`
}
```

### Funciones Auxiliares
```go
func NewConfig() *Config
func NewConfiguracionSII(rutEmisor, certificadoPath, clavePrivada, ambiente string) *ConfiguracionSII
func NewConfiguracionERP(tipo, baseURL, apiKey string) *ConfiguracionERP
func NewConfiguracionValidacion(tipo string, reglas []string) *ConfiguracionValidacion
func NewConfiguracionNotificacion(tipo, destinatario, template string) *ConfiguracionNotificacion
```

## Tipos de Servicio

Todos los tipos de servicio están definidos en `models/services.go`:

### DTEService
```go
type DTEService struct {
    ID          string    `json:"id"`
    Nombre      string    `json:"nombre"`
    Descripcion string    `json:"descripcion"`
    Version     string    `json:"version"`
    Estado      string    `json:"estado"`
    Config      Config    `json:"config"`
    CreatedAt   time.Time `json:"createdAt"`
    UpdatedAt   time.Time `json:"updatedAt"`
}
```

### DTEGenerator
```go
type DTEGenerator struct {
    ID          string    `json:"id"`
    Nombre      string    `json:"nombre"`
    Descripcion string    `json:"descripcion"`
    Version     string    `json:"version"`
    Estado      string    `json:"estado"`
    Config      Config    `json:"config"`
    CreatedAt   time.Time `json:"createdAt"`
    UpdatedAt   time.Time `json:"updatedAt"`
}
```

### ValidationService
```go
type ValidationService struct {
    ID          string    `json:"id"`
    Nombre      string    `json:"nombre"`
    Descripcion string    `json:"descripcion"`
    Version     string    `json:"version"`
    Estado      string    `json:"estado"`
    Config      Config    `json:"config"`
    CreatedAt   time.Time `json:"createdAt"`
    UpdatedAt   time.Time `json:"updatedAt"`
}
```

### NotificationService
```go
type NotificationService struct {
    ID          string    `json:"id"`
    Nombre      string    `json:"nombre"`
    Descripcion string    `json:"descripcion"`
    Version     string    `json:"version"`
    Estado      string    `json:"estado"`
    Config      Config    `json:"config"`
    CreatedAt   time.Time `json:"createdAt"`
    UpdatedAt   time.Time `json:"updatedAt"`
}
```

### StorageService
```go
type StorageService struct {
    ID          string    `json:"id"`
    Nombre      string    `json:"nombre"`
    Descripcion string    `json:"descripcion"`
    Version     string    `json:"version"`
    Estado      string    `json:"estado"`
    Config      Config    `json:"config"`
    CreatedAt   time.Time `json:"createdAt"`
    UpdatedAt   time.Time `json:"updatedAt"`
}
```

### ERPService
```go
type ERPService struct {
    ID          string    `json:"id"`
    Nombre      string    `json:"nombre"`
    Descripcion string    `json:"descripcion"`
    Version     string    `json:"version"`
    Estado      string    `json:"estado"`
    Config      Config    `json:"config"`
    CreatedAt   time.Time `json:"createdAt"`
    UpdatedAt   time.Time `json:"updatedAt"`
}
```

### MonitoringService
```go
type MonitoringService struct {
    ID          string    `json:"id"`
    Nombre      string    `json:"nombre"`
    Descripcion string    `json:"descripcion"`
    Version     string    `json:"version"`
    Estado      string    `json:"estado"`
    Config      Config    `json:"config"`
    CreatedAt   time.Time `json:"createdAt"`
    UpdatedAt   time.Time `json:"updatedAt"`
}
```

### Funciones Auxiliares
```go
func NewDTEService(nombre, descripcion, version string, config Config) *DTEService
func NewDTEGenerator(nombre, descripcion, version string, config Config) *DTEGenerator
func NewValidationService(nombre, descripcion, version string, config Config) *ValidationService
func NewNotificationService(nombre, descripcion, version string, config Config) *NotificationService
func NewStorageService(nombre, descripcion, version string, config Config) *StorageService
func NewERPService(nombre, descripcion, version string, config Config) *ERPService
func NewMonitoringService(nombre, descripcion, version string, config Config) *MonitoringService
``` 