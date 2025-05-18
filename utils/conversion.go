package utils

import (
	"time"

	"FMgo/core/dte"
	"FMgo/domain"
	"FMgo/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// DomainToModel functions convert domain types to model types

// DomainItemToModelDetalle converts domain.Item to models.Detalle
func DomainItemToModelDetalle(item domain.Item) models.Detalle {
	return models.Detalle{
		Descripcion:    item.Descripcion,
		Cantidad:       int(item.Cantidad),
		PrecioUnitario: item.PrecioUnit,
		MontoItem:      item.MontoTotal,
	}
}

// DomainItemsToModelDetalles converts []domain.Item to []models.Detalle
func DomainItemsToModelDetalles(items []domain.Item) []models.Detalle {
	result := make([]models.Detalle, len(items))
	for i, item := range items {
		result[i] = DomainItemToModelDetalle(item)
	}
	return result
}

// DomainDocumentoToModelDocumento converts domain.DocumentoTributario to models.DocumentoTributario
func DomainDocumentoToModelDocumento(doc domain.DocumentoTributario) models.DocumentoTributario {
	return models.DocumentoTributario{
		ID:           doc.ID.Hex(),
		TipoDTE:      doc.TipoDocumento,
		Folio:        int(doc.Folio),
		FechaEmision: doc.FechaEmision,
		RUTEmisor:    doc.RutEmisor,
		RUTReceptor:  doc.RutReceptor,
		MontoTotal:   doc.MontoTotal,
		MontoNeto:    doc.MontoNeto,
		MontoExento:  doc.MontoExento,
		MontoIVA:     doc.MontoIVA,
		Estado:       models.EstadoDTE(doc.Estado),
		Timestamps: models.Timestamps{
			Creado:     doc.FechaCreacion.Format(time.RFC3339),
			Modificado: doc.FechaActualizacion.Format(time.RFC3339),
		},
	}
}

// CoreDTEToModelDTE converts core/dte.DTE to models.DocumentoTributario
func CoreDTEToModelDTE(d *dte.DTE) models.DocumentoTributario {
	doc := d.Documento
	id := primitive.NewObjectID()
	folio := int(doc.Encabezado.IDDocumento.Folio)

	return models.DocumentoTributario{
		ID:           id.Hex(),
		TipoDTE:      doc.Encabezado.IDDocumento.TipoDTE,
		Folio:        folio,
		FechaEmision: doc.Encabezado.IDDocumento.FechaEmision,
		RUTEmisor:    doc.Encabezado.Emisor.RUT,
		RUTReceptor:  doc.Encabezado.Receptor.RUT,
		MontoTotal:   doc.Encabezado.Totales.MontoTotal,
		MontoNeto:    doc.Encabezado.Totales.MontoNeto,
		MontoExento:  doc.Encabezado.Totales.MontoExento,
		MontoIVA:     doc.Encabezado.Totales.IVA,
		Estado:       models.EstadoDTEAceptado,
		XML:          d.XML,
		Timestamps: models.Timestamps{
			Creado:     d.FechaCreacion.Format(time.RFC3339),
			Modificado: time.Now().Format(time.RFC3339),
		},
	}
}

// ModelToCore functions convert model types to core types

// ModelDocumentoToCoreDTE converts models.DocumentoTributario to core/dte.DTE
func ModelDocumentoToCoreDTE(doc models.DocumentoTributario) *dte.DTE {
	emisor := dte.Emisor{
		RUT:         doc.RUTEmisor,
		RazonSocial: doc.Emisor.RazonSocial,
		Giro:        doc.Emisor.GiroComercial,
		Direccion:   doc.Emisor.Direccion,
		Comuna:      doc.Emisor.Comuna,
		Ciudad:      doc.Emisor.Ciudad,
	}

	receptor := dte.Receptor{
		RUT:         doc.RUTReceptor,
		RazonSocial: doc.Receptor.RazonSocial,
		Giro:        doc.Receptor.GiroComercial,
		Direccion:   doc.Receptor.Direccion,
		Comuna:      doc.Receptor.Comuna,
		Ciudad:      doc.Receptor.Ciudad,
	}

	idDoc := dte.IDDocumento{
		TipoDTE:      doc.TipoDTE,
		Folio:        doc.Folio,
		FechaEmision: doc.FechaEmision,
	}

	totales := dte.Totales{
		MontoNeto:   doc.MontoNeto,
		MontoExento: doc.MontoExento,
		IVA:         doc.MontoIVA,
		MontoTotal:  doc.MontoTotal,
	}

	encabezado := dte.Encabezado{
		IDDocumento: idDoc,
		Emisor:      emisor,
		Receptor:    receptor,
		Totales:     totales,
	}

	// Convert detalles to core/dte.Detalle
	detalles := make([]dte.Detalle, len(doc.Detalles))
	for i, detalle := range doc.Detalles {
		detalles[i] = dte.Detalle{
			NumeroLinea: i + 1,
			Nombre:      detalle.Descripcion,
			Descripcion: detalle.Descripcion,
			Cantidad:    float64(detalle.Cantidad),
			Precio:      detalle.PrecioUnitario,
			MontoItem:   detalle.MontoItem,
		}
	}

	// Convertir el campo Creado de string a time.Time
	var fechaCreacion time.Time
	if doc.Timestamps.Creado != "" {
		fechaCreacion, _ = time.Parse(time.RFC3339, doc.Timestamps.Creado)
	} else {
		fechaCreacion = time.Now()
	}

	return &dte.DTE{
		ID:            doc.ID,
		Documento:     dte.Documento{Encabezado: encabezado, Detalles: detalles},
		FechaCreacion: fechaCreacion,
		Estado:        string(doc.Estado),
		XML:           doc.XML,
	}
}

// Additional conversion functions can be added as needed for other types

// DomainItemToModelItem converts domain.Item to models.Item
func DomainItemToModelItem(item domain.Item, lineNumber int) models.Item {
	return models.Item{
		ID:             item.ID.Hex(),
		NumeroLinea:    lineNumber,
		Codigo:         "",
		Nombre:         item.Descripcion,
		Descripcion:    item.Descripcion,
		Cantidad:       item.Cantidad,
		PrecioUnitario: item.PrecioUnit,
		MontoItem:      item.MontoTotal,
	}
}

// DomainItemsToModelItems converts []domain.Item to []models.Item
func DomainItemsToModelItems(items []domain.Item) []models.Item {
	result := make([]models.Item, len(items))
	for i, item := range items {
		result[i] = DomainItemToModelItem(item, i+1)
	}
	return result
}
