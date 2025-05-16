package main

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"time"

	"github.com/fmgo/config"
	"github.com/fmgo/models"
	"github.com/fmgo/repository/supabase"
	supaClient "github.com/fmgo/supabase"
)

// RunRepositoryExample ejecuta un ejemplo de uso del repositorio de Supabase para documentos
func RunRepositoryExample() {
	// Cargar la configuración
	// Intentar cargar desde diferentes ubicaciones relativas
	configPaths := []string{
		"../../config.json", // Desde el directorio del ejecutable (examples/supabase_example)
		"../config.json",    // Desde el directorio examples
		"config.json",       // Desde el directorio actual
	}

	var cfg *config.Config
	var err error

	for _, path := range configPaths {
		absPath, _ := filepath.Abs(path)
		fmt.Printf("Intentando cargar configuración desde: %s\n", absPath)
		cfg, err = config.Load(path)
		if err == nil {
			fmt.Printf("Configuración cargada correctamente desde: %s\n", absPath)
			break
		}
	}

	if err != nil {
		log.Fatalf("Error al cargar la configuración: %v", err)
	}

	// Crear cliente de Supabase
	client, err := supaClient.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error al crear cliente de Supabase: %v", err)
	}

	// Verificar la conexión
	if err := client.Ping(context.Background()); err != nil {
		log.Fatalf("Error al verificar la conexión con Supabase: %v", err)
	}
	fmt.Println("Conexión con Supabase establecida correctamente")

	// Crear factory de repositorios
	repoFactory := supabase.NewRepositoryFactory(client)

	// Obtener repositorio de documentos
	docRepo := repoFactory.NewDocumentoRepository()

	// Crear un documento de ejemplo
	doc := &models.DocumentoTributario{
		TipoDTE:      "33",
		Folio:        12345,
		RutEmisor:    "76555444-3",
		RutReceptor:  "55666777-8",
		MontoNeto:    1000000,
		MontoIVA:     190000,
		MontoTotal:   1190000,
		Estado:       models.EstadoDTEPendiente,
		TrackID:      "",
		FechaEmision: time.Now(),
	}

	// Guardar el documento
	ctx := context.Background()
	if err := docRepo.Create(ctx, doc); err != nil {
		log.Fatalf("Error al guardar documento: %v", err)
	}
	fmt.Printf("Documento guardado correctamente con ID: %s\n", doc.ID)

	// Obtener el documento por ID
	docRecuperado, err := docRepo.GetByID(ctx, doc.ID)
	if err != nil {
		log.Fatalf("Error al obtener documento por ID: %v", err)
	}
	fmt.Printf("Documento recuperado: Folio=%d, Monto Total=%.2f\n",
		docRecuperado.Folio, docRecuperado.MontoTotal)

	// Actualizar el estado del documento
	if err := docRepo.UpdateEstado(ctx, doc.ID, models.EstadoDTEEnviado); err != nil {
		log.Fatalf("Error al actualizar estado: %v", err)
	}
	fmt.Println("Estado actualizado correctamente")

	// Actualizar el track ID
	nuevoTrackID := "SII123456789"
	if err := docRepo.UpdateTrackID(ctx, doc.ID, nuevoTrackID); err != nil {
		log.Fatalf("Error al actualizar track ID: %v", err)
	}
	fmt.Printf("TrackID actualizado correctamente: %s\n", nuevoTrackID)

	// Listar documentos con filtros
	filtro := map[string]interface{}{
		"rut_emisor": doc.RutEmisor,
		"estado":     string(models.EstadoDTEEnviado),
	}

	documentos, err := docRepo.List(ctx, filtro, 10, 0)
	if err != nil {
		log.Fatalf("Error al listar documentos: %v", err)
	}
	fmt.Printf("Se encontraron %d documentos\n", len(documentos))

	// Contar documentos
	total, err := docRepo.Count(ctx, filtro)
	if err != nil {
		log.Fatalf("Error al contar documentos: %v", err)
	}
	fmt.Printf("Total de documentos: %d\n", total)

	fmt.Println("Ejemplo completado correctamente")
}
