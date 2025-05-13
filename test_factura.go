package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/cursor/FMgo/config"
	"github.com/cursor/FMgo/models"
	"github.com/cursor/FMgo/services"
)

func main() {
	// Configurar Supabase
	config := &config.SupabaseConfig{
		URL:    "https://tu-proyecto.supabase.co",
		APIKey: "tu-api-key",
	}

	// Inicializar servicios
	empresaService := services.NewEmpresaService(config)
	cafService := services.NewCAFService(config)
	firmaService := services.NewFirmaService(config)
	xmlService := services.NewXMLService(config)
	siiService := services.NewSIIService("https://api.sii.cl")

	ctx := context.Background()

	// 1. Obtener empresa
	empresa, err := empresaService.GetEmpresaByRut("76212889-6")
	if err != nil {
		log.Fatal("Error obteniendo empresa:", err)
	}
	fmt.Printf("Empresa encontrada: %s\n", empresa.RazonSocial)

	// 2. Obtener CAF disponible
	caf, err := cafService.GetCAFDisponible(empresa.ID, "33") // 33 = Factura
	if err != nil {
		log.Fatal("Error obteniendo CAF:", err)
	}
	fmt.Printf("CAF encontrado: Folios %d-%d\n", caf.FolioInicial, caf.FolioFinal)

	// 3. Crear factura
	factura := &models.Factura{
		DocumentoTributario: models.DocumentoTributario{
			TipoDTE:             "33",
			Folio:               caf.FolioActual,
			FechaEmision:        time.Now(),
			RutEmisor:           empresa.Rut,
			RazonSocialEmisor:   empresa.RazonSocial,
			RutReceptor:         "55555555-5",
			RazonSocialReceptor: "Cliente de Prueba",
			MontoTotal:          100000,
			MontoNeto:           84034,
			MontoIVA:            15966,
			Estado:              models.EstadoDocumentoPendiente,
			Items: []models.Item{
				{
					Descripcion: "Producto de prueba",
					Cantidad:    1,
					Precio:      100000,
					Total:       100000,
				},
			},
		},
	}

	// 4. Generar XML
	xmlData, err := xmlService.GenerarXML(factura)
	if err != nil {
		log.Fatal("Error generando XML:", err)
	}
	fmt.Println("XML generado correctamente")

	// 5. Firmar documento
	xmlFirmado, err := firmaService.FirmarXML(xmlData, factura.ID)
	if err != nil {
		log.Fatal("Error firmando documento:", err)
	}
	fmt.Println("Documento firmado correctamente")

	// 6. Enviar al SII
	trackID, err := siiService.EnviarDocumento(ctx, factura, xmlFirmado)
	if err != nil {
		log.Fatal("Error enviando al SII:", err)
	}
	fmt.Printf("Documento enviado al SII. TrackID: %s\n", trackID)

	// 7. Actualizar factura con TrackID
	factura.TrackID = trackID
	factura.Estado = models.EstadoDocumentoEnviado
	if err := empresaService.GuardarDocumento(ctx, factura); err != nil {
		log.Fatal("Error actualizando factura:", err)
	}

	// 8. Consultar estado
	time.Sleep(5 * time.Second) // Esperar a que el SII procese
	estado, err := siiService.ConsultarEstado(ctx, trackID)
	if err != nil {
		log.Fatal("Error consultando estado:", err)
	}
	fmt.Printf("Estado del documento: %s\n", estado.Estado)
	fmt.Printf("Glosa: %s\n", estado.Glosa)

	// 9. Actualizar estado final
	factura.Estado = models.EstadoDocumento(estado.Estado)
	if err := empresaService.ActualizarDocumento(ctx, factura); err != nil {
		log.Fatal("Error actualizando estado final:", err)
	}

	fmt.Println("Prueba completada exitosamente")
}
