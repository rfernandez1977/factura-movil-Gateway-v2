package supabaseexample

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"FMgo/config"
	"FMgo/supabase"
	"github.com/joho/godotenv"
)

// RunExample ejecuta los ejemplos de Supabase
func RunExample() {
	// Cargar archivo .env
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	// Mostrar menú de opciones
	fmt.Println("=== FMgo - Ejemplos de Supabase ===")
	fmt.Println("1. Verificar conexión a Supabase")
	fmt.Println("2. Probar Repositorio de Documentos")
	fmt.Println("0. Salir")

	var opcion int
	fmt.Print("\nSeleccione una opción: ")
	fmt.Scan(&opcion)

	switch opcion {
	case 1:
		ejecutarVerificacionConexion()
	case 2:
		fmt.Println("\nEjecutando ejemplo de Repositorio de Documentos...\n")
		RunRepositoryExample()
	case 0:
		fmt.Println("Saliendo...")
		os.Exit(0)
	default:
		fmt.Println("Opción no válida")
	}
}

func ejecutarVerificacionConexion() {
	fmt.Println("Verificando la conexión con Supabase...")

	// Imprimir variables de entorno para diagnóstico
	fmt.Println("\nVariables de entorno:")
	for _, env := range os.Environ() {
		if strings.HasPrefix(env, "SUPABASE_") {
			fmt.Println(" -", env)
		}
	}

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

	fmt.Println("\nConfiguración de Supabase:")
	fmt.Printf("- Supabase URL: %s\n", cfg.Supabase.URL)
	fmt.Printf("- Database Host: %s\n", cfg.Database.Host)
	fmt.Printf("- Database Port: %d\n", cfg.Database.Port)
	fmt.Printf("- Database Name: %s\n", cfg.Database.Name)
	fmt.Printf("- Database User: %s\n", cfg.Database.User)
	fmt.Printf("- SSL Mode: %s\n", cfg.Database.SSLMode)

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
