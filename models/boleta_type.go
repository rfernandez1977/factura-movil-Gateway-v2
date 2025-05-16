package models

import (
	"fmt"
	"time"
)

// BOLETAType represents the boleta type for SII integration
type BOLETAType struct {
	Documento DocumentoBOL `xml:"Documento" json:"documento"`
	TmstFirma string       `xml:"TmstFirma" json:"tmst_firma"`
	Signature SignatureBOL `xml:"Signature,omitempty" json:"signature,omitempty"`
	TED       string       `xml:",omitempty" json:"ted,omitempty"` // Timbre Electrónico de Documentos
}

// DocumentoBOL represents the document structure for electronic boletas
type DocumentoBOL struct {
	ID           string          `xml:"ID,attr" json:"id"`
	Encabezado   EncabezadoBOL   `xml:"Encabezado" json:"encabezado"`
	Detalle      []DetalleBOL    `xml:"Detalle" json:"detalle"`
	SubTotInfo   []SubTotInfo    `xml:"SubTotInfo,omitempty" json:"sub_tot_info,omitempty"`
	DscRcgGlobal []DscRcgGlobal  `xml:"DscRcgGlobal,omitempty" json:"dsc_rcg_global,omitempty"`
	Referencia   []ReferenciaBOL `xml:"Referencia,omitempty" json:"referencia,omitempty"`
	TED          string          `xml:"TED,omitempty" json:"ted,omitempty"`
	TmstFirma    string          `xml:"TmstFirma,omitempty" json:"tmst_firma,omitempty"`
}

// EncabezadoBOL represents the header of an electronic boleta
type EncabezadoBOL struct {
	IdDoc      IdDocBOL     `xml:"IdDoc" json:"id_doc"`
	Emisor     EmisorBOL    `xml:"Emisor" json:"emisor"`
	Receptor   ReceptorBOL  `xml:"Receptor" json:"receptor"`
	Totales    TotalesBOL   `xml:"Totales" json:"totales"`
	OtraMoneda []OtraMoneda `xml:"OtraMoneda,omitempty" json:"otra_moneda,omitempty"`
}

// IdDocBOL represents the identification of the document
type IdDocBOL struct {
	TipoDTE      string `xml:"TipoDTE" json:"tipo_dte"`
	Folio        int    `xml:"Folio" json:"folio"`
	FchEmis      string `xml:"FchEmis" json:"fch_emis"`
	IndServicio  int    `xml:"IndServicio,omitempty" json:"ind_servicio,omitempty"`
	IndMntNeto   int    `xml:"IndMntNeto,omitempty" json:"ind_mnt_neto,omitempty"`
	PeriodoDesde string `xml:"PeriodoDesde,omitempty" json:"periodo_desde,omitempty"`
	PeriodoHasta string `xml:"PeriodoHasta,omitempty" json:"periodo_hasta,omitempty"`
	FchVenc      string `xml:"FchVenc,omitempty" json:"fch_venc,omitempty"`
}

// EmisorBOL represents the emitter of the boleta
type EmisorBOL struct {
	RUTEmisor    string `xml:"RUTEmisor" json:"rut_emisor"`
	RznSocEmisor string `xml:"RznSocEmisor" json:"rzn_soc_emisor"`
	GiroEmisor   string `xml:"GiroEmisor" json:"giro_emisor"`
	DirOrigen    string `xml:"DirOrigen,omitempty" json:"dir_origen,omitempty"`
	CmnaOrigen   string `xml:"CmnaOrigen,omitempty" json:"cmna_origen,omitempty"`
	CiudadOrigen string `xml:"CiudadOrigen,omitempty" json:"ciudad_origen,omitempty"`
	CdgSIISucur  string `xml:"CdgSIISucur,omitempty" json:"cdg_sii_sucur,omitempty"`
}

// ReceptorBOL represents the receiver of the boleta
type ReceptorBOL struct {
	RUTRecep    string `xml:"RUTRecep" json:"rut_recep"`
	RznSocRecep string `xml:"RznSocRecep,omitempty" json:"rzn_soc_recep,omitempty"`
	DirRecep    string `xml:"DirRecep,omitempty" json:"dir_recep,omitempty"`
	CmnaRecep   string `xml:"CmnaRecep,omitempty" json:"cmna_recep,omitempty"`
	CiudadRecep string `xml:"CiudadRecep,omitempty" json:"ciudad_recep,omitempty"`
}

// TotalesBOL represents the totals of the boleta
type TotalesBOL struct {
	MntNeto       float64 `xml:"MntNeto,omitempty" json:"mnt_neto,omitempty"`
	MntExe        float64 `xml:"MntExe,omitempty" json:"mnt_exe,omitempty"`
	IVA           float64 `xml:"IVA,omitempty" json:"iva,omitempty"`
	MntTotal      float64 `xml:"MntTotal" json:"mnt_total"`
	MontoNF       float64 `xml:"MontoNF,omitempty" json:"monto_nf,omitempty"`
	MontoPeriodo  float64 `xml:"MontoPeriodo,omitempty" json:"monto_periodo,omitempty"`
	SaldoAnterior float64 `xml:"SaldoAnterior,omitempty" json:"saldo_anterior,omitempty"`
	VlrPagar      float64 `xml:"VlrPagar,omitempty" json:"vlr_pagar,omitempty"`
}

// DetalleBOL represents a boleta detail
type DetalleBOL struct {
	NroLinDet       int             `xml:"NroLinDet" json:"nro_lin_det"`
	CdgItem         CdgItem         `xml:"CdgItem,omitempty" json:"cdg_item,omitempty"`
	IndExe          int             `xml:"IndExe,omitempty" json:"ind_exe,omitempty"`
	NmbItem         string          `xml:"NmbItem" json:"nmb_item"`
	DscItem         string          `xml:"DscItem,omitempty" json:"dsc_item,omitempty"`
	QtyItem         float64         `xml:"QtyItem,omitempty" json:"qty_item,omitempty"`
	UnmdItem        string          `xml:"UnmdItem,omitempty" json:"unmd_item,omitempty"`
	PrcItem         float64         `xml:"PrcItem,omitempty" json:"prc_item,omitempty"`
	DescuentoMonto  float64         `xml:"DescuentoMonto,omitempty" json:"descuento_monto,omitempty"`
	DescuentoPct    float64         `xml:"DescuentoPct,omitempty" json:"descuento_pct,omitempty"`
	RecargoPct      float64         `xml:"RecargoPct,omitempty" json:"recargo_pct,omitempty"`
	RecargoMonto    float64         `xml:"RecargoMonto,omitempty" json:"recargo_monto,omitempty"`
	MontoItem       float64         `xml:"MontoItem" json:"monto_item"`
	ItemEspectaculo ItemEspectaculo `xml:"ItemEspectaculo,omitempty" json:"item_espectaculo,omitempty"`
}

// CdgItem represents an item code
type CdgItem struct {
	TpoCodigo string `xml:"TpoCodigo" json:"tpo_codigo"`
	VlrCodigo string `xml:"VlrCodigo" json:"vlr_codigo"`
}

// ItemEspectaculo represents item information for shows
type ItemEspectaculo struct {
	FolioTicket      string `xml:"FolioTicket" json:"folio_ticket"`
	NmbEvento        string `xml:"NmbEvento" json:"nmb_evento"`
	Tipo             string `xml:"Tipo" json:"tipo"`
	FchEvento        string `xml:"FchEvento" json:"fch_evento"`
	LugarEvento      string `xml:"LugarEvento" json:"lugar_evento"`
	UbicacionAsiento string `xml:"UbicacionAsiento" json:"ubicacion_asiento"`
}

// SubTotInfo represents subtotals information
type SubTotInfo struct {
	NroSTI         int     `xml:"NroSTI" json:"nro_sti"`
	GlosaSTI       string  `xml:"GlosaSTI" json:"glosa_sti"`
	SubTotMntNeto  float64 `xml:"SubTotMntNeto,omitempty" json:"sub_tot_mnt_neto,omitempty"`
	SubTotMntExe   float64 `xml:"SubTotMntExe,omitempty" json:"sub_tot_mnt_exe,omitempty"`
	SubTotMntIVA   float64 `xml:"SubTotMntIVA,omitempty" json:"sub_tot_mnt_iva,omitempty"`
	SubTotMntTotal float64 `xml:"SubTotMntTotal" json:"sub_tot_mnt_total"`
	SubTotIVANoRec float64 `xml:"SubTotIVANoRec,omitempty" json:"sub_tot_iva_no_rec,omitempty"`
}

// ReferenciaBOL represents a reference in a boleta
type ReferenciaBOL struct {
	NroLinRef int    `xml:"NroLinRef" json:"nro_lin_ref"`
	CodRef    string `xml:"CodRef,omitempty" json:"cod_ref,omitempty"`
	TpoDocRef string `xml:"TpoDocRef" json:"tpo_doc_ref"`
	FolioRef  string `xml:"FolioRef" json:"folio_ref"`
	FchRef    string `xml:"FchRef" json:"fch_ref"`
	RazonRef  string `xml:"RazonRef,omitempty" json:"razon_ref,omitempty"`
}

// DscRcgGlobal represents a global discount or charge
type DscRcgGlobal struct {
	NroLinDR int     `xml:"NroLinDR" json:"nro_lin_dr"`
	TpoMov   string  `xml:"TpoMov" json:"tpo_mov"`
	GlosaDR  string  `xml:"GlosaDR,omitempty" json:"glosa_dr,omitempty"`
	TpoValor string  `xml:"TpoValor" json:"tpo_valor"`
	ValorDR  float64 `xml:"ValorDR" json:"valor_dr"`
	IndExeDR int     `xml:"IndExeDR,omitempty" json:"ind_exe_dr,omitempty"`
}

// OtraMoneda represents amounts in another currency
type OtraMoneda struct {
	TpoMoneda          string  `xml:"TpoMoneda" json:"tpo_moneda"`
	TpoCambio          float64 `xml:"TpoCambio" json:"tpo_cambio"`
	MntNetoOtrMnda     float64 `xml:"MntNetoOtrMnda,omitempty" json:"mnt_neto_otr_mnda,omitempty"`
	MntExeOtrMnda      float64 `xml:"MntExeOtrMnda,omitempty" json:"mnt_exe_otr_mnda,omitempty"`
	MntFaeCarneOtrMnda float64 `xml:"MntFaeCarneOtrMnda,omitempty" json:"mnt_fae_carne_otr_mnda,omitempty"`
	MntMargComOtrMnda  float64 `xml:"MntMargComOtrMnda,omitempty" json:"mnt_marg_com_otr_mnda,omitempty"`
	IVAOtrMnda         float64 `xml:"IVAOtrMnda,omitempty" json:"iva_otr_mnda,omitempty"`
	MntTotOtrMnda      float64 `xml:"MntTotOtrMnda" json:"mnt_tot_otr_mnda"`
}

// SignatureBOL represents the digital signature of a boleta
type SignatureBOL struct {
	SignedInfo     SignedInfoBOL `xml:"SignedInfo" json:"signed_info"`
	SignatureValue string        `xml:"SignatureValue" json:"signature_value"`
	KeyInfo        KeyInfoBOL    `xml:"KeyInfo" json:"key_info"`
}

// SignedInfoBOL represents the signed information
type SignedInfoBOL struct {
	CanonicalizationMethod CanonicalizationMethodBOL `xml:"CanonicalizationMethod" json:"canonicalization_method"`
	SignatureMethod        SignatureMethodBOL        `xml:"SignatureMethod" json:"signature_method"`
	Reference              ReferenceBOL              `xml:"Reference" json:"reference"`
}

// CanonicalizationMethodBOL represents the canonicalization method
type CanonicalizationMethodBOL struct {
	Algorithm string `xml:"Algorithm,attr" json:"algorithm"`
}

// SignatureMethodBOL represents the signature method
type SignatureMethodBOL struct {
	Algorithm string `xml:"Algorithm,attr" json:"algorithm"`
}

// ReferenceBOL represents the reference to sign
type ReferenceBOL struct {
	URI          string          `xml:"URI,attr" json:"uri"`
	DigestMethod DigestMethodBOL `xml:"DigestMethod" json:"digest_method"`
	DigestValue  string          `xml:"DigestValue" json:"digest_value"`
}

// DigestMethodBOL represents the digest method
type DigestMethodBOL struct {
	Algorithm string `xml:"Algorithm,attr" json:"algorithm"`
}

// KeyInfoBOL represents the key information
type KeyInfoBOL struct {
	KeyValue KeyValueBOL `xml:"KeyValue" json:"key_value"`
	X509Data X509DataBOL `xml:"X509Data" json:"x509_data"`
}

// KeyValueBOL represents the key value
type KeyValueBOL struct {
	RSAKeyValue RSAKeyValueBOL `xml:"RSAKeyValue" json:"rsa_key_value"`
}

// RSAKeyValueBOL represents the RSA key value
type RSAKeyValueBOL struct {
	Modulus  string `xml:"Modulus" json:"modulus"`
	Exponent string `xml:"Exponent" json:"exponent"`
}

// X509DataBOL represents X509 data
type X509DataBOL struct {
	X509Certificate string `xml:"X509Certificate" json:"x509_certificate"`
}

// ConvertirBoleta converts a Boleta to BOLETAType for XML
func ConvertirBoleta(boleta *Boleta) BOLETAType {
	fechaEmision := boleta.FechaEmision.Format("2006-01-02")

	boletaType := BOLETAType{
		Documento: DocumentoBOL{
			ID: "BOLETA_" + boleta.RUTEmisor + "_" + GetFormattedFolio(boleta.Folio),
			Encabezado: EncabezadoBOL{
				IdDoc: IdDocBOL{
					TipoDTE: "39", // 39 for electronic boleta
					Folio:   boleta.Folio,
					FchEmis: fechaEmision,
				},
				Emisor: EmisorBOL{
					RUTEmisor:    boleta.RUTEmisor,
					RznSocEmisor: boleta.RazonSocialEmisor,
					GiroEmisor:   boleta.GiroEmisor,
					DirOrigen:    boleta.DireccionEmisor,
					CmnaOrigen:   boleta.ComunaEmisor,
				},
				Receptor: ReceptorBOL{
					RUTRecep:    boleta.RUTReceptor,
					RznSocRecep: boleta.RazonSocialReceptor,
					DirRecep:    boleta.DireccionReceptor,
				},
				Totales: TotalesBOL{
					MntNeto:  boleta.MontoNeto,
					MntExe:   boleta.MontoExento,
					IVA:      boleta.MontoIVA,
					MntTotal: boleta.MontoTotal,
				},
			},
		},
		TmstFirma: time.Now().Format("2006-01-02T15:04:05"),
	}

	// Convertir detalles
	if boleta.Detalles != nil && len(boleta.Detalles) > 0 {
		detalles := boleta.Detalles
		for i, detalle := range detalles {
			detalleXML := DetalleBOL{
				NroLinDet: i + 1,
				NmbItem:   detalle.Descripcion,
				QtyItem:   float64(detalle.Cantidad),
				PrcItem:   detalle.Precio,
				MontoItem: detalle.Total,
			}

			boletaType.Documento.Detalle = append(boletaType.Documento.Detalle, detalleXML)
		}
	} else if boleta.Items != nil && len(boleta.Items) > 0 {
		// Usar Items si Detalles está vacío
		items := boleta.Items
		for i, item := range items {
			detalleXML := DetalleBOL{
				NroLinDet: i + 1,
				NmbItem:   item.Descripcion,
				QtyItem:   float64(item.Cantidad),
				PrcItem:   item.Precio,
				MontoItem: item.Total,
			}

			boletaType.Documento.Detalle = append(boletaType.Documento.Detalle, detalleXML)
		}
	}

	// Convertir referencias
	for i, ref := range boleta.Referencias {
		refXML := ReferenciaBOL{
			NroLinRef: i + 1,
			TpoDocRef: string(ref.TipoDocumento),
			FolioRef:  GetFormattedFolio(ref.Folio),
			FchRef:    ref.FechaReferencia.Format("2006-01-02"),
			RazonRef:  ref.RazonReferencia,
		}

		boletaType.Documento.Referencia = append(boletaType.Documento.Referencia, refXML)
	}

	return boletaType
}

// GetFormattedFolio formats the folio in string with padding of zeros
func GetFormattedFolio(folio int) string {
	return fmt.Sprintf("%010d", folio)
}
