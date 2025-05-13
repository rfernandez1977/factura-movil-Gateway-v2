package constants

// Estados de respuesta SII
const (
	EstadoOK              = "0"
	EstadoErrorSchema     = "1"
	EstadoErrorSemilla    = "2"
	EstadoErrorToken      = "3"
	EstadoErrorFirma      = "4"
	EstadoErrorValidacion = "5"
	EstadoErrorEnvio      = "6"
)

// Códigos de error específicos del SII
const (
	// Errores de Schema
	ErrorSchemaInvalido = "00101"
	ErrorFirmaInvalida  = "00102"
	ErrorRUTInvalido    = "00103"
	ErrorFormatoFecha   = "00104"
	ErrorMontoInvalido  = "00105"

	// Errores de Negocio
	ErrorDTEDuplicado      = "00201"
	ErrorRUTNoAutorizado   = "00202"
	ErrorCAFVencido        = "00203"
	ErrorFolioNoDisponible = "00204"
	ErrorRangoFolio        = "00205"
)

// URLs de servicios SII (Certificación)
const (
	URLSemilla      = "https://maullin.sii.cl/DTEWS/CrSeed.jws?WSDL"
	URLToken        = "https://maullin.sii.cl/DTEWS/GetTokenFromSeed.jws?WSDL"
	URLRecepcionDTE = "https://maullin.sii.cl/DTEWS/RecepcionDTE.jws?WSDL"
)

// Mensajes de error
var ErrorMessages = map[string]string{
	ErrorSchemaInvalido:    "El esquema XML del DTE es inválido",
	ErrorFirmaInvalida:     "La firma digital del DTE es inválida",
	ErrorRUTInvalido:       "El RUT especificado no es válido",
	ErrorFormatoFecha:      "El formato de la fecha es inválido",
	ErrorMontoInvalido:     "El monto especificado es inválido",
	ErrorDTEDuplicado:      "El DTE ya fue recibido anteriormente",
	ErrorRUTNoAutorizado:   "El RUT no está autorizado para emitir DTEs",
	ErrorCAFVencido:        "El CAF está vencido",
	ErrorFolioNoDisponible: "El folio no está disponible para uso",
	ErrorRangoFolio:        "El folio está fuera del rango autorizado",
}
