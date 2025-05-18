package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"FMgo/config"
	"FMgo/supabase"
)

func main() {
	fmt.Println("=== FMgo - Listado de Datos en Supabase ===")

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

	// Listar los datos de las tablas
	listarDatos(ctx, client, "empresas")
	listarDatos(ctx, client, "documentos")
	listarDatos(ctx, client, "certificados")
	listarDatos(ctx, client, "cafs")
	listarDatos(ctx, client, "xml_files")
	listarDatos(ctx, client, "sobres_xml")
}

func listarDatos(ctx context.Context, client *supabase.Client, tabla string) {
	fmt.Printf("\n=== Tabla '%s' ===\n", tabla)

	// Configurar solicitud HTTP para obtener los datos
	url := fmt.Sprintf("%s/rest/v1/%s?select=*&limit=5", client.GetConfig().Supabase.URL, tabla)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		fmt.Printf("Error creando solicitud para tabla %s: %v\n", tabla, err)
		return
	}

	req.Header.Set("apikey", client.GetConfig().Supabase.AnonKey)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", client.GetConfig().Supabase.AnonKey))

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		fmt.Printf("Error obteniendo datos de la tabla %s: %v\n", tabla, err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error leyendo respuesta para tabla %s: %v\n", tabla, err)
		return
	}

	// Procesar respuesta JSON
	var data []map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		fmt.Printf("Error decodificando JSON para tabla %s: %v\n", tabla, err)
		fmt.Printf("Respuesta: %s\n", string(body))
		return
	}

	if len(data) == 0 {
		fmt.Printf("No hay datos en la tabla %s\n", tabla)
		return
	}

	// Mostrar datos (primeros 5 registros)
	for i, item := range data {
		fmt.Printf("Registro %d:\n", i+1)
		// Mostrar solo algunos campos clave en lugar de todo el registro
		switch tabla {
		case "empresas":
			fmt.Printf("  ID: %s\n", item["id"])
			fmt.Printf("  Nombre: %s\n", item["nombre"])
			fmt.Printf("  RUT: %s\n", item["rut"])
			fmt.Printf("  Email: %s\n", item["email"])
		case "documentos":
			fmt.Printf("  ID: %s\n", item["id"])
			fmt.Printf("  Tipo: %s\n", item["tipo_documento"])
			fmt.Printf("  Número: %s\n", item["numero_documento"])
			fmt.Printf("  Estado: %s\n", item["estado"])
		case "certificados":
			fmt.Printf("  ID: %s\n", item["id"])
			fmt.Printf("  Nombre: %s\n", item["nombre"])
			// No mostrar archivo (binario)
		case "cafs":
			fmt.Printf("  ID: %s\n", item["id"])
			fmt.Printf("  Tipo: %s\n", item["tipo_documento"])
			fmt.Printf("  Desde: %v\n", item["desde"])
			fmt.Printf("  Hasta: %v\n", item["hasta"])
		case "xml_files":
			fmt.Printf("  ID: %s\n", item["id"])
			fmt.Printf("  Nombre: %s\n", item["nombre_archivo"])
			// No mostrar contenido (XML largo)
		case "sobres_xml":
			fmt.Printf("  ID: %s\n", item["id"])
			fmt.Printf("  Nombre: %s\n", item["nombre_archivo"])
			// No mostrar contenido (binario)
		default:
			// Mostrar todos los campos para otras tablas
			for k, v := range item {
				// Evitar imprimir valores binarios o muy largos
				switch v.(type) {
				case string:
					if len(v.(string)) > 100 {
						fmt.Printf("  %s: [contenido largo...]\n", k)
					} else {
						fmt.Printf("  %s: %v\n", k, v)
					}
				default:
					fmt.Printf("  %s: %v\n", k, v)
				}
			}
		}
		fmt.Println()
	}
}
