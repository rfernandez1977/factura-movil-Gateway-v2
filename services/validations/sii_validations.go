package validations

import (
	"errors"
	"fmt"
	"time"

	"github.com/fmgo/models"
)

// CodigoError representa los códigos de error de validación
type CodigoError string

// Códigos de error de validación
const (
	ErrorValidacionSchema  CodigoError = "SCHEMA"
	ErrorValidacionFormato CodigoError = "FORMATO"
	ErrorValidacionNegocio CodigoError = "NEGOCIO"
	ErrorValidacionFirma   CodigoError = "FIRMA"
	ErrorValidacionSII     CodigoError = "SII"
)

// SIIValidator maneja las validaciones específicas del SII
// Se renombra para evitar conflictos con otras declaraciones
type SIIValidator struct {
	// Configuración de validación
	config  *ConfiguracionValidacion
	errores []models.ErrorValidacion
}

// ConfiguracionValidacion contiene la configuración para las validaciones
type ConfiguracionValidacion struct {
	MaxDiasAntiguedad int
	MaxMontoTotal     float64
	MaxItems          int
	TasaIVA           float64
	ValidarCAF        bool
	ValidarFirma      bool
	ValidarSchema     bool
}

// DTEXMLModel representa un DTE en formato XML
type DTEXMLModel struct {
	Documento *DocumentoXMLModel
	Signature *FirmaXMLModel
}

// DocumentoXMLModel representa la estructura de un documento XML
type DocumentoXMLModel struct {
	Encabezado EncabezadoXMLModel
	Detalle    []DetalleDTEXML
}

// EncabezadoXMLModel representa el encabezado de un DTE en formato XML
type EncabezadoXMLModel struct {
	IdDoc    IDDocumentoXMLModel
	Emisor   EmisorXMLModel
	Receptor ReceptorXMLModel
	Totales  TotalesXMLModel
}

// IDDocumentoXMLModel representa la identificación del documento en XML
type IDDocumentoXMLModel struct {
	TipoDTE      string
	Folio        int
	FechaEmision string
}

// EmisorXMLModel representa al emisor en formato XML
type EmisorXMLModel struct {
	RUT         string
	RazonSocial string
	Giro        string
	Direccion   string
	Comuna      string
	Ciudad      string
}

// ReceptorXMLModel representa al receptor en formato XML
type ReceptorXMLModel struct {
	RUT         string
	RazonSocial string
	Giro        string
	Direccion   string
	Comuna      string
	Ciudad      string
}

// TotalesXMLModel representa los totales en formato XML
type TotalesXMLModel struct {
	MntTotal int64
	MntNeto  *int64
	IVA      *int64
}

// DetalleDTEXML representa un detalle de DTE en XML
type DetalleDTEXML struct {
	NroLinDet  int
	Nombre     string
	Cantidad   float64
	PrecioUnit float64
	MontoItem  float64
}

// FirmaXMLModel representa la firma electrónica en formato XML
type FirmaXMLModel struct {
	SignatureValue string
	KeyInfo        struct {
		X509Data struct {
			X509Certificate string
		}
	}
}

// NewSIIValidator crea una nueva instancia del servicio de validación
func NewSIIValidator(config *ConfiguracionValidacion) *SIIValidator {
	if config == nil {
		config = &ConfiguracionValidacion{
			MaxDiasAntiguedad: 30,
			MaxMontoTotal:     1000000000, // 1.000.000.000
			MaxItems:          1000,
			TasaIVA:           19.0,
			ValidarCAF:        true,
			ValidarFirma:      true,
			ValidarSchema:     true,
		}
	}
	return &SIIValidator{
		config:  config,
		errores: make([]models.ErrorValidacion, 0),
	}
}

// ValidarDocumento valida un documento según las reglas del SII
func (v *SIIValidator) ValidarDocumento(doc interface{}) error {
	v.errores = make([]models.ErrorValidacion, 0)

	switch d := doc.(type) {
	case *DTEXMLModel:
		return v.validarDTE(d)
	case *models.RespuestaSII:
		return v.validarRespuestaSII(d)
	case *models.EstadoSII:
		return v.validarEstadoSII(d)
	default:
		return errors.New("tipo de documento no soportado")
	}
}

// ValidarRespuesta valida una respuesta del SII
func (v *SIIValidator) ValidarRespuesta(respuesta *models.RespuestaSII) error {
	if respuesta == nil {
		return errors.New("respuesta SII vacía")
	}

	// Validar track ID
	if respuesta.TrackID == "" {
		return errors.New("track ID no proporcionado")
	}

	// Validar estado
	if respuesta.Estado == "" {
		return errors.New("estado no proporcionado")
	}

	// No validamos FechaProceso ya que no existe en el tipo actual
	// Tampoco validamos Errores por la misma razón

	return nil
}

// validarDTE valida un DTE según las reglas del SII
func (v *SIIValidator) validarDTE(dte *DTEXMLModel) error {
	if err := v.validarEncabezado(dte.Documento.Encabezado); err != nil {
		return err
	}

	if err := v.validarDetalle(dte.Documento.Detalle); err != nil {
		return err
	}

	if err := v.validarTotales(dte.Documento.Encabezado.Totales); err != nil {
		return err
	}

	if v.config.ValidarFirma && dte.Signature != nil {
		if err := v.validarFirma(dte.Signature); err != nil {
			return err
		}
	}

	if len(v.errores) > 0 {
		return fmt.Errorf("errores de validación: %v", v.errores)
	}

	return nil
}

// validarEncabezado valida el encabezado del DTE
func (v *SIIValidator) validarEncabezado(enc EncabezadoXMLModel) error {
	if err := v.validarIdDoc(enc.IdDoc); err != nil {
		return err
	}

	if err := v.validarEmisor(enc.Emisor); err != nil {
		return err
	}

	if err := v.validarReceptor(enc.Receptor); err != nil {
		return err
	}

	return nil
}

// validarIdDoc valida la identificación del documento
func (v *SIIValidator) validarIdDoc(id IDDocumentoXMLModel) error {
	// Validar tipo de documento
	if id.TipoDTE == "" {
		v.agregarError(ErrorValidacionSchema, "TipoDTE", "Tipo de documento no especificado")
	}

	// Validar folio
	if id.Folio <= 0 {
		v.agregarError(ErrorValidacionSchema, "Folio", "Folio inválido")
	}

	// Validar fecha de emisión
	fechaEmision, err := time.Parse("2006-01-02", id.FechaEmision)
	if err != nil {
		v.agregarError(ErrorValidacionFormato, "FechaEmision", "Formato de fecha inválido")
	} else {
		// Validar que la fecha no sea futura
		if fechaEmision.After(time.Now()) {
			v.agregarError(ErrorValidacionNegocio, "FechaEmision", "Fecha de emisión futura no permitida")
		}
		// Validar antigüedad máxima
		if time.Since(fechaEmision).Hours() > float64(v.config.MaxDiasAntiguedad*24) {
			v.agregarError(ErrorValidacionNegocio, "FechaEmision", "Documento excede antigüedad máxima permitida")
		}
	}

	return nil
}

// validarEmisor valida los datos del emisor
func (v *SIIValidator) validarEmisor(emisor EmisorXMLModel) error {
	if emisor.RUT == "" {
		v.agregarError(ErrorValidacionSchema, "RUTEmisor", "RUT del emisor no especificado")
	} else if !v.validarFormatoRUT(emisor.RUT) {
		v.agregarError(ErrorValidacionFormato, "RUTEmisor", "Formato de RUT inválido")
	}

	if emisor.RazonSocial == "" {
		v.agregarError(ErrorValidacionSchema, "RznSoc", "Razón social no especificada")
	}

	if emisor.Giro == "" {
		v.agregarError(ErrorValidacionSchema, "GiroEmis", "Giro no especificado")
	}

	if emisor.Direccion == "" {
		v.agregarError(ErrorValidacionSchema, "DirOrigen", "Dirección no especificada")
	}

	if emisor.Comuna == "" {
		v.agregarError(ErrorValidacionSchema, "CmnaOrigen", "Comuna no especificada")
	}

	if emisor.Ciudad == "" {
		v.agregarError(ErrorValidacionSchema, "CiudadOrigen", "Ciudad no especificada")
	}

	return nil
}

// validarReceptor valida los datos del receptor
func (v *SIIValidator) validarReceptor(receptor ReceptorXMLModel) error {
	if receptor.RUT == "" {
		v.agregarError(ErrorValidacionSchema, "RUTRecep", "RUT del receptor no especificado")
	} else if !v.validarFormatoRUT(receptor.RUT) {
		v.agregarError(ErrorValidacionFormato, "RUTRecep", "Formato de RUT inválido")
	}

	if receptor.RazonSocial == "" {
		v.agregarError(ErrorValidacionSchema, "RznSocRecep", "Razón social no especificada")
	}

	if receptor.Giro == "" {
		v.agregarError(ErrorValidacionSchema, "GiroRecep", "Giro no especificado")
	}

	if receptor.Direccion == "" {
		v.agregarError(ErrorValidacionSchema, "DirRecep", "Dirección no especificada")
	}

	if receptor.Comuna == "" {
		v.agregarError(ErrorValidacionSchema, "CmnaRecep", "Comuna no especificada")
	}

	if receptor.Ciudad == "" {
		v.agregarError(ErrorValidacionSchema, "CiudadRecep", "Ciudad no especificada")
	}

	return nil
}

// validarDetalle valida el detalle del DTE
func (v *SIIValidator) validarDetalle(detalle []DetalleDTEXML) error {
	if len(detalle) == 0 {
		v.agregarError(ErrorValidacionSchema, "Detalle", "Detalle no especificado")
		return nil
	}

	if len(detalle) > v.config.MaxItems {
		v.agregarError(ErrorValidacionNegocio, "Detalle", fmt.Sprintf("Número de ítems excede el máximo permitido (%d)", v.config.MaxItems))
		return nil
	}

	for i, item := range detalle {
		if item.NroLinDet <= 0 {
			v.agregarError(ErrorValidacionSchema, fmt.Sprintf("Detalle[%d].NroLinDet", i), "Número de línea inválido")
		}

		if item.Nombre == "" {
			v.agregarError(ErrorValidacionSchema, fmt.Sprintf("Detalle[%d].NmbItem", i), "Nombre del ítem no especificado")
		}

		if item.MontoItem <= 0 {
			v.agregarError(ErrorValidacionNegocio, fmt.Sprintf("Detalle[%d].MontoItem", i), "Monto del ítem debe ser mayor a cero")
		}
	}

	return nil
}

// validarTotales valida los totales del DTE
func (v *SIIValidator) validarTotales(totales TotalesXMLModel) error {
	if totales.MntTotal <= 0 {
		v.agregarError(ErrorValidacionNegocio, "MntTotal", "Monto total debe ser mayor a cero")
	}

	if totales.MntTotal > int64(v.config.MaxMontoTotal) {
		v.agregarError(ErrorValidacionNegocio, "MntTotal", fmt.Sprintf("Monto total excede el máximo permitido (%f)", v.config.MaxMontoTotal))
	}

	// Validar IVA si corresponde
	if totales.MntNeto != nil && totales.IVA != nil {
		ivaCalculado := float64(*totales.MntNeto) * (v.config.TasaIVA / 100.0)
		if float64(*totales.IVA) != ivaCalculado {
			v.agregarError(ErrorValidacionNegocio, "IVA", "Monto de IVA no corresponde a la tasa vigente")
		}
	}

	return nil
}

// validarFirma valida la firma del DTE
func (v *SIIValidator) validarFirma(firma *FirmaXMLModel) error {
	if firma.SignatureValue == "" {
		v.agregarError(ErrorValidacionFirma, "SignatureValue", "Valor de firma no especificado")
	}

	if firma.KeyInfo.X509Data.X509Certificate == "" {
		v.agregarError(ErrorValidacionFirma, "X509Certificate", "Certificado no especificado")
	}

	return nil
}

// validarRespuestaSII valida una respuesta del SII
func (v *SIIValidator) validarRespuestaSII(resp *models.RespuestaSII) error {
	if resp.Estado == "" {
		v.agregarError(ErrorValidacionSchema, "Estado", "Estado no especificado")
	}

	if resp.TrackID == "" {
		v.agregarError(ErrorValidacionSchema, "TrackID", "TrackID no especificado")
	}

	// No validamos FechaProceso ya que no existe en el modelo actual

	return nil
}

// validarEstadoSII valida un estado del SII
func (v *SIIValidator) validarEstadoSII(estado *models.EstadoSII) error {
	if estado.Estado == "" {
		v.agregarError(ErrorValidacionSchema, "Estado", "Estado no especificado")
	}

	// No validamos TrackID ya que no existe en el modelo actual

	return nil
}

// agregarError agrega un error de validación a la lista
func (v *SIIValidator) agregarError(codigo CodigoError, campo, mensaje string) {
	v.errores = append(v.errores, models.ErrorValidacion{
		Codigo:    string(codigo),
		Campo:     campo,
		Mensaje:   mensaje,
		Timestamp: time.Now().Format(time.RFC3339),
	})
}

// validarFormatoRUT valida el formato de un RUT chileno
func (v *SIIValidator) validarFormatoRUT(rut string) bool {
	// Implementar validación de formato RUT chileno
	// Formato: 12345678-9 o 12.345.678-9
	return true // TODO: Implementar validación real
}

// ValidateSign valida una firma electrónica
func (v *SIIValidator) ValidateSign(sign []byte) *models.ErrorValidacion {
	if len(sign) == 0 {
		return &models.ErrorValidacion{
			Codigo:    "SII-001",
			Mensaje:   "La firma electrónica es inválida",
			Campo:     "firma",
			Valor:     "",
			Timestamp: time.Now().Format(time.RFC3339),
		}
	}
	return nil
}
