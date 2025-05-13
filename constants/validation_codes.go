package constants

// Códigos de error para validación de notas de venta
const (
    NV001 = "NV001" // Folio obligatorio
    NV002 = "NV002" // Fecha de emisión obligatoria
    NV003 = "NV003" // Formato de folio inválido
    NV004 = "NV004" // Fecha futura
    NV005 = "NV005" // Folio no correlativo
    NV006 = "NV006" // RUT emisor inválido
    NV007 = "NV007" // RUT receptor inválido
    NV008 = "NV008" // Monto neto inválido
    NV009 = "NV009" // Monto IVA inválido
    NV010 = "NV010" // Monto total inválido
)

// Mensajes de error para validación de notas de venta
const (
    MsgFolioObligatorio      = "Folio es obligatorio"
    MsgFechaEmisionObligatoria = "Fecha de emisión es obligatoria"
    MsgFormatoFolioInvalido  = "Formato de folio inválido"
    MsgFechaFutura          = "Fecha de emisión no puede ser futura"
    MsgFolioNoCorrelativo   = "Folio no correlativo"
    MsgRutEmisorInvalido    = "RUT emisor inválido"
    MsgRutReceptorInvalido  = "RUT receptor inválido"
    MsgMontoNetoInvalido    = "Monto neto debe ser mayor a 0"
    MsgMontoIvaInvalido     = "Monto IVA no puede ser negativo"
    MsgMontoTotalInvalido   = "Monto total debe ser mayor a 0"
)

// Descripciones detalladas de los códigos de error
var ValidationCodeDescriptions = map[string]string{
    NV001: "El folio es un campo obligatorio para la nota de venta",
    NV002: "La fecha de emisión es un campo obligatorio para la nota de venta",
    NV003: "El formato del folio no cumple con los requisitos establecidos",
    NV004: "La fecha de emisión no puede ser posterior a la fecha actual",
    NV005: "El folio no sigue la secuencia correlativa establecida",
    NV006: "El RUT del emisor no es válido o no está registrado",
    NV007: "El RUT del receptor no es válido o no está registrado",
    NV008: "El monto neto debe ser un valor positivo",
    NV009: "El monto de IVA no puede ser negativo",
    NV010: "El monto total debe ser un valor positivo",
} 