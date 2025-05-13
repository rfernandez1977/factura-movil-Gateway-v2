package models

import "time"

// DTEXMLModel representa un DTE en formato XML
type DTEXMLModel struct {
	XMLName   struct{}          `xml:"DTE"`
	Version   string            `xml:"version,attr"`
	Documento DocumentoXMLModel `xml:"Documento"`
	Signature *FirmaXMLModel    `xml:"Signature,omitempty"`
}

// DocumentoXMLModel representa el contenido del DTE en formato XML
type DocumentoXMLModel struct {
	ID           string             `xml:"ID,attr"`
	Encabezado   EncabezadoXMLModel `xml:"Encabezado"`
	Detalle      []DetalleDTEXML    `xml:"Detalle"`
	SubTotInfo   *SubTotInfoXML     `xml:"SubTotInfo,omitempty"`
	DscRcgGlobal []DscRcgGlobalXML  `xml:"DscRcgGlobal,omitempty"`
	Referencia   []ReferenciaXML    `xml:"Referencia,omitempty"`
	TED          *TEDXML            `xml:"TED,omitempty"`
	TmstFirma    string             `xml:"TmstFirma,omitempty"`
}

// EncabezadoXMLModel contiene la información principal del DTE en formato XML
type EncabezadoXMLModel struct {
	IdDoc    IDDocumentoXMLModel `xml:"IdDoc"`
	Emisor   EmisorXMLModel      `xml:"Emisor"`
	Receptor ReceptorXMLModel    `xml:"Receptor"`
	Totales  TotalesXMLModel     `xml:"Totales"`
}

// IDDocumentoXMLModel contiene la información de identificación del documento en formato XML
type IDDocumentoXMLModel struct {
	TipoDTE       string  `xml:"TipoDTE"`
	Folio         int     `xml:"Folio"`
	FechaEmision  string  `xml:"FchEmis"`
	IndNoRebaja   *int    `xml:"IndNoRebaja,omitempty"`
	TipoDespacho  *int    `xml:"TipoDespacho,omitempty"`
	IndTraslado   *int    `xml:"IndTraslado,omitempty"`
	TpoImpresion  *string `xml:"TpoImpresion,omitempty"`
	IndServicio   *int    `xml:"IndServicio,omitempty"`
	MntBruto      *int    `xml:"MntBruto,omitempty"`
	TpoTranCompra *string `xml:"TpoTranCompra,omitempty"`
	TpoTranVenta  *string `xml:"TpoTranVenta,omitempty"`
	FmaPago       *string `xml:"FmaPago,omitempty"`
	FchVenc       *string `xml:"FchVenc,omitempty"`
}

// EmisorXMLModel contiene la información del emisor del DTE en formato XML
type EmisorXMLModel struct {
	RUT         string   `xml:"RUTEmisor"`
	RazonSocial string   `xml:"RznSoc"`
	Giro        string   `xml:"GiroEmis"`
	Acteco      []string `xml:"Acteco,omitempty"`
	Direccion   string   `xml:"DirOrigen"`
	Comuna      string   `xml:"CmnaOrigen"`
	Ciudad      string   `xml:"CiudadOrigen"`
	Telefono    *string  `xml:"Telefono,omitempty"`
	Correo      *string  `xml:"CorreoEmisor,omitempty"`
	CdgSIISucur *string  `xml:"CdgSIISucur,omitempty"`
}

// ReceptorXMLModel contiene la información del receptor del DTE en formato XML
type ReceptorXMLModel struct {
	RUT         string  `xml:"RUTRecep"`
	RazonSocial string  `xml:"RznSocRecep"`
	Giro        string  `xml:"GiroRecep"`
	Direccion   string  `xml:"DirRecep"`
	Comuna      string  `xml:"CmnaRecep"`
	Ciudad      string  `xml:"CiudadRecep"`
	Telefono    *string `xml:"Telefono,omitempty"`
	Correo      *string `xml:"CorreoRecep,omitempty"`
	Contacto    *string `xml:"Contacto,omitempty"`
}

// TotalesXMLModel contiene los montos del DTE en formato XML
type TotalesXMLModel struct {
	MntNeto       *int64                 `xml:"MntNeto,omitempty"`
	MntExe        *int64                 `xml:"MntExe,omitempty"`
	TasaIVA       *float64               `xml:"TasaIVA,omitempty"`
	IVA           *int64                 `xml:"IVA,omitempty"`
	IVAProp       *int64                 `xml:"IVAProp,omitempty"`
	IVATerc       *int64                 `xml:"IVATerc,omitempty"`
	ImptoReten    []ImpuestoRetencionXML `xml:"ImptoReten,omitempty"`
	IVANoRet      *int64                 `xml:"IVANoRet,omitempty"`
	MntTotal      int64                  `xml:"MntTotal"`
	MontoNF       *int64                 `xml:"MontoNF,omitempty"`
	MontoPeriodo  *int64                 `xml:"MontoPeriodo,omitempty"`
	SaldoAnterior *int64                 `xml:"SaldoAnterior,omitempty"`
	VlrPagar      *int64                 `xml:"VlrPagar,omitempty"`
}

// DetalleDTEXML representa un ítem del detalle del DTE
type DetalleDTEXML struct {
	NroLinDet      int          `xml:"NroLinDet"`
	CdgItem        []CodItemXML `xml:"CdgItem,omitempty"`
	IndExe         *int         `xml:"IndExe,omitempty"`
	Nombre         string       `xml:"NmbItem"`
	Descripcion    *string      `xml:"DscItem,omitempty"`
	Cantidad       *float64     `xml:"QtyItem,omitempty"`
	Unidad         *string      `xml:"UnmdItem,omitempty"`
	Precio         *float64     `xml:"PrcItem,omitempty"`
	DescuentoPct   *float64     `xml:"DescuentoPct,omitempty"`
	DescuentoMonto *int64       `xml:"DescuentoMonto,omitempty"`
	RecargoPct     *float64     `xml:"RecargoPct,omitempty"`
	RecargoMonto   *int64       `xml:"RecargoMonto,omitempty"`
	MontoItem      int64        `xml:"MontoItem"`
}

// CodItemXML representa un código de ítem
type CodItemXML struct {
	TipoCodigo string `xml:"TpoCodigo"`
	Codigo     string `xml:"VlrCodigo"`
}

// DscRcgGlobalXML representa un descuento o recargo global
type DscRcgGlobalXML struct {
	NroLinDR  int     `xml:"NroLinDR"`
	TipoMov   string  `xml:"TpoMov"`
	GlosaDR   *string `xml:"GlosaDR,omitempty"`
	TipoValor string  `xml:"TpoValor"`
	ValorDR   float64 `xml:"ValorDR"`
	IndExeDR  *int    `xml:"IndExeDR,omitempty"`
}

// ReferenciaXML representa una referencia a otro documento
type ReferenciaXML struct {
	NroLinRef  int     `xml:"NroLinRef"`
	TipoDocRef string  `xml:"TpoDocRef"`
	FolioRef   string  `xml:"FolioRef"`
	FechaRef   string  `xml:"FchRef"`
	CodRef     *string `xml:"CodRef,omitempty"`
	RazonRef   *string `xml:"RazonRef,omitempty"`
}

// ImpuestoRetencionXML representa un impuesto de retención
type ImpuestoRetencionXML struct {
	TipoImp  string  `xml:"TipoImp"`
	TasaImp  float64 `xml:"TasaImp"`
	MontoImp int64   `xml:"MontoImp"`
}

// TEDXML representa el Timbre Electrónico del DTE
type TEDXML struct {
	Version string   `xml:"version,attr"`
	DD      DDXML    `xml:"DD"`
	FRMT    *FRMTXML `xml:"FRMT,omitempty"`
}

// DDXML representa los datos del DTE
type DDXML struct {
	RE    string `xml:"RE"`    // RUT Emisor
	TD    string `xml:"TD"`    // Tipo DTE
	F     int    `xml:"F"`     // Folio
	FE    string `xml:"FE"`    // Fecha Emisión
	RR    string `xml:"RR"`    // RUT Receptor
	RSR   string `xml:"RSR"`   // Razón Social Receptor
	MNT   int64  `xml:"MNT"`   // Monto Total
	IT1   string `xml:"IT1"`   // Primer Item
	CAF   CAFXML `xml:"CAF"`   // Código Autorización de Folios
	TSTED string `xml:"TSTED"` // Timbre SII
}

// CAFXML representa el Código de Autorización de Folios
type CAFXML struct {
	Version string `xml:"version,attr"`
	DA      DAXML  `xml:"DA"`
	FRMA    string `xml:"FRMA"`
}

// DAXML representa los datos de autorización
type DAXML struct {
	RE    string `xml:"RE"`    // RUT Empresa
	RS    string `xml:"RS"`    // Razón Social
	TD    string `xml:"TD"`    // Tipo DTE
	RNG   RNGXML `xml:"RNG"`   // Rango
	FA    string `xml:"FA"`    // Fecha Autorización
	RSAPK RSAPK  `xml:"RSAPK"` // Llave Pública
	IDK   int    `xml:"IDK"`   // ID Llave
}

// RNGXML representa el rango de folios
type RNGXML struct {
	D int `xml:"D"` // Desde
	H int `xml:"H"` // Hasta
}

// RSAPK representa la llave pública RSA
type RSAPK struct {
	M string `xml:"M"` // Módulo
	E string `xml:"E"` // Exponente
}

// FRMTXML representa el formato del timbre
type FRMTXML struct {
	Version string `xml:"version,attr"`
	URI     string `xml:",chardata"`
}

// FirmaXMLModel representa la firma digital del DTE
type FirmaXMLModel struct {
	XMLName        struct{}      `xml:"Signature"`
	SignedInfo     SignedInfoXML `xml:"SignedInfo"`
	SignatureValue string        `xml:"SignatureValue"`
	KeyInfo        KeyInfoXML    `xml:"KeyInfo"`
}

// SignedInfoXML representa la información firmada
type SignedInfoXML struct {
	CanonicalizationMethod CanonicalizationMethodXML `xml:"CanonicalizationMethod"`
	SignatureMethod        SignatureMethodXML        `xml:"SignatureMethod"`
	Reference              ReferenceSignatureXML     `xml:"Reference"`
}

// CanonicalizationMethodXML representa el método de canonicalización
type CanonicalizationMethodXML struct {
	Algorithm string `xml:"Algorithm,attr"`
}

// SignatureMethodXML representa el método de firma
type SignatureMethodXML struct {
	Algorithm string `xml:"Algorithm,attr"`
}

// ReferenceSignatureXML representa la referencia de la firma
type ReferenceSignatureXML struct {
	URI          string          `xml:"URI,attr"`
	Transforms   TransformsXML   `xml:"Transforms"`
	DigestMethod DigestMethodXML `xml:"DigestMethod"`
	DigestValue  string          `xml:"DigestValue"`
}

// TransformsXML representa las transformaciones aplicadas
type TransformsXML struct {
	Transform []TransformXML `xml:"Transform"`
}

// TransformXML representa una transformación
type TransformXML struct {
	Algorithm string `xml:"Algorithm,attr"`
}

// DigestMethodXML representa el método de digest
type DigestMethodXML struct {
	Algorithm string `xml:"Algorithm,attr"`
}

// KeyInfoXML representa la información de la llave
type KeyInfoXML struct {
	KeyValue KeyValueXML `xml:"KeyValue"`
	X509Data X509DataXML `xml:"X509Data"`
}

// KeyValueXML representa el valor de la llave
type KeyValueXML struct {
	RSAKeyValue RSAKeyValueXML `xml:"RSAKeyValue"`
}

// RSAKeyValueXML representa el valor de la llave RSA
type RSAKeyValueXML struct {
	Modulus  string `xml:"Modulus"`
	Exponent string `xml:"Exponent"`
}

// X509DataXML representa los datos del certificado X509
type X509DataXML struct {
	X509Certificate string `xml:"X509Certificate"`
}

// SubTotInfoXML representa la información de subtotales
type SubTotInfoXML struct {
	NroLinea    int     `xml:"NroLinea"`
	GlosaLinea  string  `xml:"GlosaLinea,omitempty"`
	Valor       float64 `xml:"Valor"`
	LineaQuiebr *int    `xml:"LineaQuiebr,omitempty"`
}

// Funciones auxiliares para crear modelos XML
func NewDTEXMLModel(version string, documento DocumentoXMLModel) *DTEXMLModel {
	return &DTEXMLModel{
		Version:   version,
		Documento: documento,
	}
}

func NewDocumentoXMLModel(encabezado EncabezadoXMLModel, detalles []DetalleDTEXML) *DocumentoXMLModel {
	return &DocumentoXMLModel{
		Encabezado: encabezado,
		Detalle:    detalles,
	}
}

func NewEncabezadoXMLModel(idDocumento IDDocumentoXMLModel, emisor EmisorXMLModel, receptor ReceptorXMLModel, totales TotalesXMLModel) *EncabezadoXMLModel {
	return &EncabezadoXMLModel{
		IdDoc:    idDocumento,
		Emisor:   emisor,
		Receptor: receptor,
		Totales:  totales,
	}
}

func NewIDDocumentoXMLModel(tipoDTE string, folio int, fechaEmision time.Time) *IDDocumentoXMLModel {
	return &IDDocumentoXMLModel{
		TipoDTE:      tipoDTE,
		Folio:        folio,
		FechaEmision: fechaEmision.Format("2006-01-02"),
	}
}

func NewEmisorXMLModel(rut, razonSocial, giro, direccion, comuna, ciudad string) *EmisorXMLModel {
	return &EmisorXMLModel{
		RUT:         rut,
		RazonSocial: razonSocial,
		Giro:        giro,
		Direccion:   direccion,
		Comuna:      comuna,
		Ciudad:      ciudad,
	}
}

func NewReceptorXMLModel(rut, razonSocial, giro, direccion, comuna, ciudad string) *ReceptorXMLModel {
	return &ReceptorXMLModel{
		RUT:         rut,
		RazonSocial: razonSocial,
		Giro:        giro,
		Direccion:   direccion,
		Comuna:      comuna,
		Ciudad:      ciudad,
	}
}

func NewTotalesXMLModel(montoNeto, montoExento, iva, montoTotal int64) *TotalesXMLModel {
	tasaIVA := float64(iva)
	return &TotalesXMLModel{
		MntNeto:       &montoNeto,
		MntExe:        &montoExento,
		TasaIVA:       &tasaIVA,
		IVA:           &iva,
		IVAProp:       &iva,
		IVATerc:       &iva,
		ImptoReten:    []ImpuestoRetencionXML{},
		IVANoRet:      &iva,
		MntTotal:      montoTotal,
		MontoNF:       &montoNeto,
		MontoPeriodo:  &montoNeto,
		SaldoAnterior: &montoNeto,
		VlrPagar:      &montoNeto,
	}
}

func NewDetalleDTEXML(nroLinDet int, cdgItem []CodItemXML, indExe *int, nombre string, descripcion *string, cantidad *float64, unidad *string, precio *float64, descuentoPct *float64, descuentoMonto *int64, recargoPct *float64, recargoMonto *int64, montoItem int64) *DetalleDTEXML {
	return &DetalleDTEXML{
		NroLinDet:      nroLinDet,
		CdgItem:        cdgItem,
		IndExe:         indExe,
		Nombre:         nombre,
		Descripcion:    descripcion,
		Cantidad:       cantidad,
		Unidad:         unidad,
		Precio:         precio,
		DescuentoPct:   descuentoPct,
		DescuentoMonto: descuentoMonto,
		RecargoPct:     recargoPct,
		RecargoMonto:   recargoMonto,
		MontoItem:      montoItem,
	}
}

func NewCodItemXML(tipoCodigo string, codigo string) *CodItemXML {
	return &CodItemXML{
		TipoCodigo: tipoCodigo,
		Codigo:     codigo,
	}
}

func NewDscRcgGlobalXML(nroLinDR int, tipoMov string, glosaDR *string, tipoValor string, valorDR float64, indExeDR *int) *DscRcgGlobalXML {
	return &DscRcgGlobalXML{
		NroLinDR:  nroLinDR,
		TipoMov:   tipoMov,
		GlosaDR:   glosaDR,
		TipoValor: tipoValor,
		ValorDR:   valorDR,
		IndExeDR:  indExeDR,
	}
}

func NewReferenciaXML(nroLinRef int, tipoDocRef string, folioRef string, fechaRef string, codRef *string, razonRef *string) *ReferenciaXML {
	return &ReferenciaXML{
		NroLinRef:  nroLinRef,
		TipoDocRef: tipoDocRef,
		FolioRef:   folioRef,
		FechaRef:   fechaRef,
		CodRef:     codRef,
		RazonRef:   razonRef,
	}
}

func NewImpuestoRetencionXML(tipoImp string, tasaImp float64, montoImp int64) *ImpuestoRetencionXML {
	return &ImpuestoRetencionXML{
		TipoImp:  tipoImp,
		TasaImp:  tasaImp,
		MontoImp: montoImp,
	}
}

func NewTEDXML(version string, dd DDXML, frmt *FRMTXML) *TEDXML {
	return &TEDXML{
		Version: version,
		DD:      dd,
		FRMT:    frmt,
	}
}

func NewDDXML(re, td, fe, rr, rsr string, f int, mnt int64, it1 string, caf CAFXML, tsted string) *DDXML {
	return &DDXML{
		RE:    re,
		TD:    td,
		F:     f,
		FE:    fe,
		RR:    rr,
		RSR:   rsr,
		MNT:   mnt,
		IT1:   it1,
		CAF:   caf,
		TSTED: tsted,
	}
}

func NewCAFXML(version string, da DAXML, frma string) *CAFXML {
	return &CAFXML{
		Version: version,
		DA:      da,
		FRMA:    frma,
	}
}

func NewDAXML(re, rs, td string, rng RNGXML, fa string, rsapk RSAPK, idk int) *DAXML {
	return &DAXML{
		RE:    re,
		RS:    rs,
		TD:    td,
		RNG:   rng,
		FA:    fa,
		RSAPK: rsapk,
		IDK:   idk,
	}
}

func NewRNGXML(d, h int) *RNGXML {
	return &RNGXML{
		D: d,
		H: h,
	}
}

func NewRSAPK(m, e string) *RSAPK {
	return &RSAPK{
		M: m,
		E: e,
	}
}

func NewFRMTXML(version string, uri string) *FRMTXML {
	return &FRMTXML{
		Version: version,
		URI:     uri,
	}
}

func NewFirmaXMLModel(signedInfo SignedInfoXML, signatureValue string, keyInfo KeyInfoXML) *FirmaXMLModel {
	return &FirmaXMLModel{
		SignedInfo:     signedInfo,
		SignatureValue: signatureValue,
		KeyInfo:        keyInfo,
	}
}

func NewSignedInfoXML(canonicalizationMethod CanonicalizationMethodXML, signatureMethod SignatureMethodXML, reference ReferenceSignatureXML) *SignedInfoXML {
	return &SignedInfoXML{
		CanonicalizationMethod: canonicalizationMethod,
		SignatureMethod:        signatureMethod,
		Reference:              reference,
	}
}

func NewCanonicalizationMethodXML(algorithm string) *CanonicalizationMethodXML {
	return &CanonicalizationMethodXML{
		Algorithm: algorithm,
	}
}

func NewSignatureMethodXML(algorithm string) *SignatureMethodXML {
	return &SignatureMethodXML{
		Algorithm: algorithm,
	}
}

func NewReferenceSignatureXML(uri string, transforms TransformsXML, digestMethod DigestMethodXML, digestValue string) *ReferenceSignatureXML {
	return &ReferenceSignatureXML{
		URI:          uri,
		Transforms:   transforms,
		DigestMethod: digestMethod,
		DigestValue:  digestValue,
	}
}

func NewTransformsXML(transforms []TransformXML) *TransformsXML {
	return &TransformsXML{
		Transform: transforms,
	}
}

func NewTransformXML(algorithm string) *TransformXML {
	return &TransformXML{
		Algorithm: algorithm,
	}
}

func NewDigestMethodXML(algorithm string) *DigestMethodXML {
	return &DigestMethodXML{
		Algorithm: algorithm,
	}
}

func NewKeyValueXML(rsaKeyValue RSAKeyValueXML) *KeyValueXML {
	return &KeyValueXML{
		RSAKeyValue: rsaKeyValue,
	}
}

func NewRSAKeyValueXML(modulus, exponent string) *RSAKeyValueXML {
	return &RSAKeyValueXML{
		Modulus:  modulus,
		Exponent: exponent,
	}
}

func NewX509DataXML(x509Certificate string) *X509DataXML {
	return &X509DataXML{
		X509Certificate: x509Certificate,
	}
}

func NewSubTotInfoXML(nroLinea int, glosaLinea string, valor float64, lineaQuiebr *int) *SubTotInfoXML {
	return &SubTotInfoXML{
		NroLinea:    nroLinea,
		GlosaLinea:  glosaLinea,
		Valor:       valor,
		LineaQuiebr: lineaQuiebr,
	}
}
