package main

import (
	"fmt"
	"log"

	"FMgo/config"
)

func main() {
	// Cargar configuración
	cfg, err := config.Load("config.json")
	if err != nil {
		log.Fatalf("Error cargando configuración: %v", err)
	}

	// Mostrar configuración (sin mostrar valores sensibles)
	fmt.Printf("Environment: %s\n", cfg.Env)
	fmt.Printf("Supabase URL: %s\n", cfg.Supabase.URL)
	fmt.Printf("Database Host: %s\n", cfg.Database.Host)
	fmt.Printf("Database Port: %d\n", cfg.Database.Port)
	fmt.Printf("Database Name: %s\n", cfg.Database.Name)
	fmt.Printf("Database User: %s\n", cfg.Database.User)
	fmt.Printf("SSL Mode: %s\n", cfg.Database.SSLMode)
}
