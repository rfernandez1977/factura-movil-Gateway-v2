package utils

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/xml"
	"fmt"

	"FMgo/models"
)

// SIIUtils contiene utilidades específicas para el SII
type SIIUtils struct{}

// NewSIIUtils crea una nueva instancia de SIIUtils
func NewSIIUtils() *SIIUtils {
	return &SIIUtils{}
}

// GenerateSIIXML genera el XML específico para el SII
func (u *SIIUtils) GenerateSIIXML(doc *models.DocumentoTributario) ([]byte, error) {
	// Generar el hash del documento - actualmente no se utiliza, se calculará en otra etapa
	// hash := u.GenerateDocumentHash(doc)

	// Crear la estructura XML del SII
	siiDoc := models.DTEXMLModel{
		Version: "1.0",
		Documento: models.DocumentoXMLModel{
			ID: fmt.Sprintf("DTE_%s_%d", doc.RUTEmisor, doc.Folio),
			Encabezado: models.EncabezadoXMLModel{
				IdDoc: models.IDDocumentoXMLModel{
					TipoDTE:      doc.TipoDTE,
					Folio:        doc.Folio,
					FechaEmision: doc.FechaEmision.Format("2006-01-02"),
				},
				Emisor: models.EmisorXMLModel{
					RUT:         doc.RUTEmisor,
					RazonSocial: doc.RUTEmisor, // TODO: Obtener razón social del emisor
					Giro:        "",            // TODO: Obtener giro del emisor
					Direccion:   "",            // TODO: Obtener dirección del emisor
					Comuna:      "",            // TODO: Obtener comuna del emisor
					Ciudad:      "",            // TODO: Obtener ciudad del emisor
				},
				Receptor: models.ReceptorXMLModel{
					RUT:         doc.RUTReceptor,
					RazonSocial: doc.RUTReceptor, // TODO: Obtener razón social del receptor
					Giro:        "",              // TODO: Obtener giro del receptor
					Direccion:   "",              // TODO: Obtener dirección del receptor
					Comuna:      "",              // TODO: Obtener comuna del receptor
					Ciudad:      "",              // TODO: Obtener ciudad del receptor
				},
				Totales: models.TotalesXMLModel{
					MntNeto:  intPtr(int64(doc.MontoNeto)),
					TasaIVA:  float64Ptr(19),
					IVA:      intPtr(int64(doc.MontoIVA)),
					MntTotal: int64(doc.MontoTotal),
				},
			},
			Detalle: make([]models.DetalleDTEXML, len(doc.Detalles)),
		},
	}

	// Agregar detalles en lugar de items
	for i, detalle := range doc.Detalles {
		cantidad := float64(detalle.Cantidad)
		precio := float64(detalle.PrecioUnitario)

		siiDoc.Documento.Detalle[i] = models.DetalleDTEXML{
			NroLinDet: i + 1,
			Nombre:    detalle.Descripcion,
			Cantidad:  &cantidad,
			Precio:    &precio,
			MontoItem: int64(detalle.MontoItem),
		}
	}

	// Generar el XML
	xmlData, err := xml.MarshalIndent(siiDoc, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("error al generar XML del SII: %v", err)
	}

	return xmlData, nil
}

// GenerateDocumentHash genera el hash de un documento para el SII
func (u *SIIUtils) GenerateDocumentHash(doc *models.DocumentoTributario) string {
	// Concatenar los campos relevantes para el hash
	data := fmt.Sprintf("%s|%d|%s|%s|%.2f",
		doc.TipoDTE,
		doc.Folio,
		doc.FechaEmision.Format("2006-01-02"),
		doc.RUTEmisor,
		doc.MontoTotal)

	// Calcular el hash SHA-1
	hash := sha1.Sum([]byte(data))
	return base64.StdEncoding.EncodeToString(hash[:])
}

// Helper functions for creating pointers to primitives
func intPtr(i int64) *int64 {
	return &i
}

func float64Ptr(f float64) *float64 {
	return &f
}

// ValidateTimbreElectronico valida el timbre electrónico del SII
func (u *SIIUtils) ValidateTimbreElectronico(timbre string) error {
	if timbre == "" {
		return fmt.Errorf("el timbre electrónico es requerido")
	}

	// Decodificar el timbre
	_, err := base64.StdEncoding.DecodeString(timbre)
	if err != nil {
		return fmt.Errorf("el timbre electrónico no es válido: %v", err)
	}

	return nil
}

// GenerateDocumentID genera un ID único para un documento del SII
func (u *SIIUtils) GenerateDocumentID(rutEmisor string, folio int64) string {
	return fmt.Sprintf("DTE_%s_%d", rutEmisor, folio)
}
