package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/fmgo/config"
	"github.com/fmgo/services"
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

	// Crear servicio de base de datos
	dbService := services.NewDatabaseService(supabaseConfig)

	// Crear contexto
	ctx := context.Background()

	// Ejemplo: Obtener usuario existente
	userID := "3611d070-57b2-4f23-83f8-79f43b319144"
	usuario, err := dbService.ObtenerUsuario(ctx, userID)
	if err != nil {
		log.Printf("Error al obtener usuario: %v", err)
	} else {
		fmt.Printf("Usuario encontrado: %s\n", usuario.Email)
	}

	// Ejemplo: Actualizar usuario
	updates := map[string]interface{}{
		"user_metadata": map[string]interface{}{
			"last_login": time.Now(),
			"ip_address": "192.168.1.1",
		},
	}

	err = dbService.ActualizarUsuario(ctx, userID, updates)
	if err != nil {
		log.Printf("Error al actualizar usuario: %v", err)
	} else {
		fmt.Println("Usuario actualizado exitosamente")
	}

	// Ejemplo: Listar usuarios
	filtros := map[string]interface{}{
		"email_verified": true,
		"role":           "authenticated",
	}

	usuarios, err := dbService.ListarUsuarios(ctx, filtros)
	if err != nil {
		log.Printf("Error al listar usuarios: %v", err)
	} else {
		fmt.Printf("Usuarios encontrados: %d\n", len(usuarios))
		for _, u := range usuarios {
			fmt.Printf("- %s (%s)\n", u.Email, u.ID)
		}
	}

	// Ejemplo: Verificar email
	err = dbService.VerificarEmail(ctx, userID)
	if err != nil {
		log.Printf("Error al verificar email: %v", err)
	} else {
		fmt.Println("Email verificado exitosamente")
	}

	// Ejemplo: Actualizar último acceso
	err = dbService.ActualizarUltimoAcceso(ctx, userID)
	if err != nil {
		log.Printf("Error al actualizar último acceso: %v", err)
	} else {
		fmt.Println("Último acceso actualizado exitosamente")
	}
}
