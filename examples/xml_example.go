package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"FMgo/config"
	"FMgo/models"
	"FMgo/services"
)

func main() {
	// Cargar configuración
	supabaseConfig, err := config.NewSupabaseConfig()
	if err != nil {
		log.Fatalf("Error al cargar configuración: %v", err)
	}

	// Crear servicios
	empresaService := services.NewEmpresaService(supabaseConfig)
	xmlService := services.NewXMLService(supabaseConfig)

	// Crear contexto
	ctx := context.Background()

	// Crear una empresa de ejemplo
	empresa := &models.Empresa{
		Rut:              "76.123.456-7",
		RazonSocial:      "Empresa Ejemplo SpA",
		Giro:             "Desarrollo de Software",
		Direccion:        "Av. Principal 123",
		Comuna:           "Santiago",
		Ciudad:           "Santiago",
		Telefono:         "+56 2 1234 5678",
		Email:            "contacto@empresaejemplo.cl",
		ResolucionNumero: "12345",
		ResolucionFecha:  time.Now(),
	}

	// Guardar empresa
	err = empresaService.CrearEmpresa(ctx, empresa)
	if err != nil {
		log.Fatalf("Error al crear empresa: %v", err)
	}

	// Crear un documento de ejemplo
	documento := &models.Documento{
		EmpresaID:     empresa.ID,
		TipoDocumento: "33", // Factura electrónica
		Folio:         1,
		Estado:        "PENDIENTE",
		FechaEmision:  time.Now(),
	}

	// Guardar documento
	err = empresaService.GuardarDocumento(ctx, documento)
	if err != nil {
		log.Fatalf("Error al guardar documento: %v", err)
	}

	// Generar XML para el documento
	xmlContent, err := xmlService.GenerarXMLFactura(ctx, documento, empresa)
	if err != nil {
		log.Fatalf("Error al generar XML: %v", err)
	}

	// Validar XML
	err = xmlService.ValidarXML(xmlContent)
	if err != nil {
		log.Fatalf("Error al validar XML: %v", err)
	}

	// Guardar XML
	err = xmlService.GuardarXML(ctx, documento.ID, xmlContent)
	if err != nil {
		log.Fatalf("Error al guardar XML: %v", err)
	}

	fmt.Println("XML generado y guardado exitosamente:")
	fmt.Println(xmlContent)

	// Obtener XML guardado
	xmlGuardado, err := xmlService.ObtenerXML(ctx, documento.ID)
	if err != nil {
		log.Fatalf("Error al obtener XML: %v", err)
	}

	fmt.Println("\nXML recuperado de la base de datos:")
	fmt.Println(xmlGuardado)
}
