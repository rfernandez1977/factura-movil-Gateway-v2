package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/fmgo/config"
	"github.com/fmgo/supabase"
	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("=== FMgo - Verificación de Datos en Supabase ===")
	// Cargar variables de entorno
	if err := godotenv.Load(); err != nil {
		fmt.Println("No se pudo cargar el archivo .env, usando valores de configuración por defecto")
	}

	// Cargar la configuración
	cfg, err := config.Load("config.json")
	if err != nil {
		log.Fatalf("Error al cargar la configuración: %v", err)
	}

	// Crear cliente de Supabase
	fmt.Println("\nCreando cliente de Supabase...")
	client, err := supabase.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creando cliente de Supabase: %v", err)
	}

	// Verificar conexión
	fmt.Println("Verificando conexión a Supabase...")
	ctx := context.Background()
	if err := client.Ping(ctx); err != nil {
		log.Fatalf("Error verificando conexión a Supabase: %v", err)
	}

	fmt.Println("\n✅ Conexión a Supabase establecida correctamente")

	// Verificar tablas
	verificarTabla(ctx, client, "empresas")
	verificarTabla(ctx, client, "documentos")
	verificarTabla(ctx, client, "certificados")
	verificarTabla(ctx, client, "sesiones")
	verificarTabla(ctx, client, "cafs")
	verificarTabla(ctx, client, "xml_files")
	verificarTabla(ctx, client, "sobres_xml")
}

func verificarTabla(ctx context.Context, client *supabase.Client, tabla string) {
	fmt.Printf("\nVerificando tabla '%s'...\n", tabla)

	// Configurar solicitud HTTP para obtener el conteo
	url := fmt.Sprintf("%s/rest/v1/%s?select=count", client.GetConfig().Supabase.URL, tabla)

	req, err := http.NewRequestWithContext(ctx, "HEAD", url, nil)
	if err != nil {
		fmt.Printf("Error creando solicitud para tabla %s: %v\n", tabla, err)
		return
	}

	req.Header.Set("apikey", client.GetConfig().Supabase.AnonKey)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", client.GetConfig().Supabase.AnonKey))
	req.Header.Set("Prefer", "count=exact")

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		fmt.Printf("Error verificando tabla %s: %v\n", tabla, err)
		return
	}
	defer resp.Body.Close()

	contentRange := resp.Header.Get("Content-Range")
	if contentRange == "" {
		fmt.Printf("No se pudo obtener el conteo para la tabla %s\n", tabla)
		return
	}

	fmt.Printf("Tabla '%s': %s registros\n", tabla, contentRange)
}
