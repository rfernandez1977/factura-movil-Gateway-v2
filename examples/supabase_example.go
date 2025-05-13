package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/cursor/FMgo/config"
	"github.com/cursor/FMgo/services"
)

func main() {
	// Cargar variables de entorno
	if err := config.LoadEnv(); err != nil {
		log.Fatalf("Error al cargar variables de entorno: %v", err)
	}

	// Inicializar configuración de Supabase
	supabaseConfig, err := config.NewSupabaseConfig()
	if err != nil {
		log.Fatalf("Error al configurar Supabase: %v", err)
	}

	// Crear servicio de Supabase
	supabaseService := services.NewSupabaseService(supabaseConfig)

	// Crear contexto
	ctx := context.Background()

	// Ejemplo: Guardar un documento
	doc := &services.Documento{
		Tipo:        "DTE",
		RutEmisor:   "11.111.111-1",
		RutReceptor: "22.222.222-2",
		Folio:       1,
		MontoTotal:  1000.0,
		Estado:      "PENDIENTE",
		XML:         "<xml>...</xml>",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err = supabaseService.GuardarDocumento(ctx, doc)
	if err != nil {
		log.Printf("Error al guardar documento: %v", err)
	} else {
		fmt.Println("Documento guardado exitosamente")
	}

	// Ejemplo: Obtener documentos
	filtros := map[string]interface{}{
		"rut_emisor": "11.111.111-1",
		"estado":     "PENDIENTE",
	}

	docs, err := supabaseService.ListarDocumentos(ctx, filtros)
	if err != nil {
		log.Printf("Error al listar documentos: %v", err)
	} else {
		fmt.Printf("Documentos encontrados: %d\n", len(docs))
		for _, d := range docs {
			fmt.Printf("- ID: %s, Tipo: %s, Monto: %.2f\n", d.ID, d.Tipo, d.MontoTotal)
		}
	}

	// Ejemplo: Guardar certificado
	certificado := []byte("certificado de prueba")
	llavePrivada := []byte("llave privada de prueba")

	err = supabaseService.GuardarCertificado(ctx, "11.111.111-1", certificado, llavePrivada)
	if err != nil {
		log.Printf("Error al guardar certificado: %v", err)
	} else {
		fmt.Println("Certificado guardado exitosamente")
	}

	// Ejemplo: Guardar sesión
	token := "token-de-prueba"
	expiracion := time.Now().Add(24 * time.Hour)

	err = supabaseService.GuardarSesion(ctx, "11.111.111-1", token, expiracion)
	if err != nil {
		log.Printf("Error al guardar sesión: %v", err)
	} else {
		fmt.Println("Sesión guardada exitosamente")
	}

	// Ejemplo: Verificar sesión
	valido, err := supabaseService.VerificarSesion(ctx, token)
	if err != nil {
		log.Printf("Error al verificar sesión: %v", err)
	} else {
		fmt.Printf("Sesión válida: %v\n", valido)
	}
}
