package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/fmgo/config"
	"github.com/fmgo/supabase"
)

func main() {
	fmt.Println("=== FMgo - Verificación de Conexión a Supabase ===")
	verificarConexionSupabase()
}

func verificarConexionSupabase() {
	fmt.Println("Verificando la conexión con Supabase...")

	// Imprimir variables de entorno para diagnóstico
	fmt.Println("\nVariables de entorno:")
	for _, env := range os.Environ() {
		if strings.HasPrefix(env, "SUPABASE_") {
			fmt.Println(" -", env)
		}
	}

	// Cargar la configuración
	configPath := "config.json"
	absPath, _ := filepath.Abs(configPath)
	fmt.Printf("Intentando cargar configuración desde: %s\n", absPath)

	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatalf("Error al cargar la configuración: %v", err)
	}

	fmt.Printf("Configuración cargada correctamente desde: %s\n", absPath)

	fmt.Println("\nConfiguración de Supabase:")
	fmt.Printf("- Supabase URL: %s\n", cfg.Supabase.URL)
	fmt.Printf("- API Key: %s...\n", cfg.Supabase.APIKey[:20])
	fmt.Printf("- Anon Key: %s...\n", cfg.Supabase.AnonKey[:20])

	// Establecer variables de entorno manualmente antes de crear el cliente
	os.Setenv("SUPABASE_URL", cfg.Supabase.URL)
	os.Setenv("SUPABASE_ANON_KEY", cfg.Supabase.AnonKey)
	os.Setenv("SUPABASE_SERVICE_KEY", cfg.Supabase.ServiceKey)
	os.Setenv("SUPABASE_JWT_SECRET", cfg.Security.JWTSecret)

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
	fmt.Println("\nPróximos pasos:")
	fmt.Println("1. Configurar la conexión a la base de datos")
	fmt.Println("2. Crear las tablas necesarias")
	fmt.Println("3. Configurar las políticas de seguridad (RLS)")
}
