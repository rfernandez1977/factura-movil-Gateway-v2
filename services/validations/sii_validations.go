package validations

import (
	"errors"
	"fmt"
	"time"

	"github.com/cursor/FMgo/models"
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
	case *models.DTEXMLModel:
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

	// Validar fechas
	if respuesta.FechaProceso.IsZero() {
		return errors.New("fecha de proceso no proporcionada")
	}

	// Validar errores
	if respuesta.Estado == "RECHAZADO" && len(respuesta.Errores) == 0 {
		return errors.New("documento rechazado sin errores especificados")
	}

	return nil
}

// validarDTE valida un DTE según las reglas del SII
func (v *SIIValidator) validarDTE(dte *models.DTEXMLModel) error {
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
func (v *SIIValidator) validarEncabezado(enc models.EncabezadoXMLModel) error {
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
func (v *SIIValidator) validarIdDoc(id models.IDDocumentoXMLModel) error {
	// Validar tipo de documento
	if id.TipoDTE == "" {
		v.agregarError(models.ErrorValidacionSchema, "TipoDTE", "Tipo de documento no especificado")
	}

	// Validar folio
	if id.Folio <= 0 {
		v.agregarError(models.ErrorValidacionSchema, "Folio", "Folio inválido")
	}

	// Validar fecha de emisión
	fechaEmision, err := time.Parse("2006-01-02", id.FechaEmision)
	if err != nil {
		v.agregarError(models.ErrorValidacionFormato, "FechaEmision", "Formato de fecha inválido")
	} else {
		// Validar que la fecha no sea futura
		if fechaEmision.After(time.Now()) {
			v.agregarError(models.ErrorValidacionNegocio, "FechaEmision", "Fecha de emisión futura no permitida")
		}
		// Validar antigüedad máxima
		if time.Since(fechaEmision).Hours() > float64(v.config.MaxDiasAntiguedad*24) {
			v.agregarError(models.ErrorValidacionNegocio, "FechaEmision", "Documento excede antigüedad máxima permitida")
		}
	}

	return nil
}

// validarEmisor valida los datos del emisor
func (v *SIIValidator) validarEmisor(emisor models.EmisorXMLModel) error {
	if emisor.RUT == "" {
		v.agregarError(models.ErrorValidacionSchema, "RUTEmisor", "RUT del emisor no especificado")
	} else if !v.validarFormatoRUT(emisor.RUT) {
		v.agregarError(models.ErrorValidacionFormato, "RUTEmisor", "Formato de RUT inválido")
	}

	if emisor.RazonSocial == "" {
		v.agregarError(models.ErrorValidacionSchema, "RznSoc", "Razón social no especificada")
	}

	if emisor.Giro == "" {
		v.agregarError(models.ErrorValidacionSchema, "GiroEmis", "Giro no especificado")
	}

	if emisor.Direccion == "" {
		v.agregarError(models.ErrorValidacionSchema, "DirOrigen", "Dirección no especificada")
	}

	if emisor.Comuna == "" {
		v.agregarError(models.ErrorValidacionSchema, "CmnaOrigen", "Comuna no especificada")
	}

	if emisor.Ciudad == "" {
		v.agregarError(models.ErrorValidacionSchema, "CiudadOrigen", "Ciudad no especificada")
	}

	return nil
}

// validarReceptor valida los datos del receptor
func (v *SIIValidator) validarReceptor(receptor models.ReceptorXMLModel) error {
	if receptor.RUT == "" {
		v.agregarError(models.ErrorValidacionSchema, "RUTRecep", "RUT del receptor no especificado")
	} else if !v.validarFormatoRUT(receptor.RUT) {
		v.agregarError(models.ErrorValidacionFormato, "RUTRecep", "Formato de RUT inválido")
	}

	if receptor.RazonSocial == "" {
		v.agregarError(models.ErrorValidacionSchema, "RznSocRecep", "Razón social no especificada")
	}

	if receptor.Giro == "" {
		v.agregarError(models.ErrorValidacionSchema, "GiroRecep", "Giro no especificado")
	}

	if receptor.Direccion == "" {
		v.agregarError(models.ErrorValidacionSchema, "DirRecep", "Dirección no especificada")
	}

	if receptor.Comuna == "" {
		v.agregarError(models.ErrorValidacionSchema, "CmnaRecep", "Comuna no especificada")
	}

	if receptor.Ciudad == "" {
		v.agregarError(models.ErrorValidacionSchema, "CiudadRecep", "Ciudad no especificada")
	}

	return nil
}

// validarDetalle valida el detalle del DTE
func (v *SIIValidator) validarDetalle(detalle []models.DetalleDTEXML) error {
	if len(detalle) == 0 {
		v.agregarError(models.ErrorValidacionSchema, "Detalle", "Detalle no especificado")
		return nil
	}

	if len(detalle) > v.config.MaxItems {
		v.agregarError(models.ErrorValidacionNegocio, "Detalle", fmt.Sprintf("Número de ítems excede el máximo permitido (%d)", v.config.MaxItems))
		return nil
	}

	for i, item := range detalle {
		if item.NroLinDet <= 0 {
			v.agregarError(models.ErrorValidacionSchema, fmt.Sprintf("Detalle[%d].NroLinDet", i), "Número de línea inválido")
		}

		if item.Nombre == "" {
			v.agregarError(models.ErrorValidacionSchema, fmt.Sprintf("Detalle[%d].NmbItem", i), "Nombre del ítem no especificado")
		}

		if item.MontoItem <= 0 {
			v.agregarError(models.ErrorValidacionNegocio, fmt.Sprintf("Detalle[%d].MontoItem", i), "Monto del ítem debe ser mayor a cero")
		}
	}

	return nil
}

// validarTotales valida los totales del DTE
func (v *SIIValidator) validarTotales(totales models.TotalesXMLModel) error {
	if totales.MntTotal <= 0 {
		v.agregarError(models.ErrorValidacionNegocio, "MntTotal", "Monto total debe ser mayor a cero")
	}

	if totales.MntTotal > int64(v.config.MaxMontoTotal) {
		v.agregarError(models.ErrorValidacionNegocio, "MntTotal", fmt.Sprintf("Monto total excede el máximo permitido (%f)", v.config.MaxMontoTotal))
	}

	// Validar IVA si corresponde
	if totales.MntNeto != nil && totales.IVA != nil {
		ivaCalculado := float64(*totales.MntNeto) * (v.config.TasaIVA / 100.0)
		if float64(*totales.IVA) != ivaCalculado {
			v.agregarError(models.ErrorValidacionNegocio, "IVA", "Monto de IVA no corresponde a la tasa vigente")
		}
	}

	return nil
}

// validarFirma valida la firma del DTE
func (v *SIIValidator) validarFirma(firma *models.FirmaXMLModel) error {
	if firma.SignatureValue == "" {
		v.agregarError(models.ErrorValidacionFirma, "SignatureValue", "Valor de firma no especificado")
	}

	if firma.KeyInfo.X509Data.X509Certificate == "" {
		v.agregarError(models.ErrorValidacionFirma, "X509Certificate", "Certificado no especificado")
	}

	return nil
}

// validarRespuestaSII valida una respuesta del SII
func (v *SIIValidator) validarRespuestaSII(resp *models.RespuestaSII) error {
	if resp.Estado == "" {
		v.agregarError(models.ErrorValidacionSchema, "Estado", "Estado no especificado")
	}

	if resp.TrackID == "" {
		v.agregarError(models.ErrorValidacionSchema, "TrackID", "TrackID no especificado")
	}

	if resp.FechaProceso.IsZero() {
		v.agregarError(models.ErrorValidacionSchema, "FechaProceso", "Fecha de proceso no especificada")
	}

	return nil
}

// validarEstadoSII valida un estado del SII
func (v *SIIValidator) validarEstadoSII(estado *models.EstadoSII) error {
	if estado.Estado == "" {
		v.agregarError(models.ErrorValidacionSchema, "Estado", "Estado no especificado")
	}

	if estado.TrackID == "" {
		v.agregarError(models.ErrorValidacionSchema, "TrackID", "TrackID no especificado")
	}

	return nil
}

// agregarError agrega un error de validación a la lista
func (v *SIIValidator) agregarError(codigo models.CodigoError, campo, mensaje string) {
	v.errores = append(v.errores, models.ErrorValidacion{
		Codigo:    string(codigo),
		Campo:     campo,
		Mensaje:   mensaje,
		Timestamp: time.Now(),
	})
}

// validarFormatoRUT valida el formato de un RUT chileno
func (v *SIIValidator) validarFormatoRUT(rut string) bool {
	// Implementar validación de formato RUT chileno
	// Formato: 12345678-9 o 12.345.678-9
	return true // TODO: Implementar validación real
}
