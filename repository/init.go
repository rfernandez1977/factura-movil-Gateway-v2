package repository

import (
	"fmt"
	"log"
	"sync"

	"FMgo/supabase"
)

var (
	instance *SupabaseRepository
	once     sync.Once
)

// InitializeRepository inicializa el repositorio con la configuración desde el archivo especificado
func InitializeRepository(configPath string) (*SupabaseRepository, error) {
	var initErr error
	once.Do(func() {
		log.Printf("Inicializando repositorio con configuración desde: %s", configPath)

		// Inicializar cliente de Supabase
		client, err := supabase.InitClientWithConfig(configPath)
		if err != nil {
			initErr = fmt.Errorf("error inicializando cliente Supabase: %w", err)
			return
		}

		// Crear instancia del repositorio
		instance = NewSupabaseRepository(client)
		log.Println("Repositorio inicializado correctamente")
	})

	if initErr != nil {
		return nil, initErr
	}

	return instance, nil
}

// GetRepository devuelve la instancia actual del repositorio
// Si no se ha inicializado, retorna un error
func GetRepository() (*SupabaseRepository, error) {
	if instance == nil {
		return nil, fmt.Errorf("el repositorio no ha sido inicializado")
	}
	return instance, nil
}

// IsInitialized devuelve true si el repositorio ha sido inicializado
func IsInitialized() bool {
	return instance != nil
}
