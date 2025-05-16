package main

import (
	"context"
	"fmt"
	"log"

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

	// Crear servicio de autenticación
	authService := services.NewAuthService(supabaseConfig)

	// Ejemplo: Generar token
	email := "rfernandez@facturamovil.cl"
	userID := "3611d070-57b2-4f23-83f8-79f43b319144"

	token, err := authService.GenerarToken(email, userID)
	if err != nil {
		log.Printf("Error al generar token: %v", err)
	} else {
		fmt.Printf("Token generado: %s\n", token)
	}

	// Ejemplo: Validar token
	claims, err := authService.ValidarToken(token)
	if err != nil {
		log.Printf("Error al validar token: %v", err)
	} else {
		fmt.Printf("Token válido para usuario: %s\n", claims.Email)
	}

	// Ejemplo: Verificar sesión
	ctx := context.Background()
	valido, err := authService.VerificarSesion(ctx, token)
	if err != nil {
		log.Printf("Error al verificar sesión: %v", err)
	} else {
		fmt.Printf("Sesión válida: %v\n", valido)
	}

	// Ejemplo: Obtener información del usuario
	usuario, err := authService.ObtenerUsuario(token)
	if err != nil {
		log.Printf("Error al obtener usuario: %v", err)
	} else {
		fmt.Printf("Usuario: %s\n", usuario.Email)
		fmt.Printf("Rol: %s\n", usuario.Role)
		fmt.Printf("ID: %s\n", usuario.Sub)
	}
}
