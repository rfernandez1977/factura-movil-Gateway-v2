package services

import (
	"fmt"

	"github.com/fmgo/config"
	"github.com/fmgo/models"
)

// ItemService maneja la lógica de negocio de items
type ItemService struct {
	config *config.SupabaseConfig
}

// NewItemService crea una nueva instancia del servicio de item
func NewItemService(config *config.SupabaseConfig) *ItemService {
	return &ItemService{
		config: config,
	}
}

// GetItemByID obtiene un item por su ID
func (s *ItemService) GetItemByID(id string) (*models.Item, error) {
	var item models.Item
	err := s.config.Client.DB.From("items").
		Select("*").
		Eq("id", id).
		Single().
		Execute(&item)

	if err != nil {
		return nil, fmt.Errorf("error al obtener item: %v", err)
	}

	return &item, nil
}

// CrearItem crea un nuevo item
func (s *ItemService) CrearItem(item *models.Item) (*models.Item, error) {
	// Validar item
	if err := s.validarItem(item); err != nil {
		return nil, err
	}

	// Guardar item en Supabase
	_, err := s.config.Client.DB.From("items").
		Insert(item).
		Execute()

	if err != nil {
		return nil, fmt.Errorf("error al guardar item: %v", err)
	}

	return item, nil
}

// ActualizarItem actualiza un item existente
func (s *ItemService) ActualizarItem(item *models.Item) error {
	// Validar item
	if err := s.validarItem(item); err != nil {
		return err
	}

	// Actualizar item en Supabase
	_, err := s.config.Client.DB.From("items").
		Update(item).
		Eq("id", item.ID).
		Execute()

	if err != nil {
		return fmt.Errorf("error al actualizar item: %v", err)
	}

	return nil
}

// EliminarItem elimina un item
func (s *ItemService) EliminarItem(id string) error {
	// Eliminar item de Supabase
	_, err := s.config.Client.DB.From("items").
		Delete().
		Eq("id", id).
		Execute()

	if err != nil {
		return fmt.Errorf("error al eliminar item: %v", err)
	}

	return nil
}

// validarItem valida un item antes de crearlo o actualizarlo
func (s *ItemService) validarItem(item *models.Item) error {
	if item.Codigo == "" {
		return fmt.Errorf("código requerido")
	}
	if item.Descripcion == "" {
		return fmt.Errorf("descripción requerida")
	}
	if item.PrecioUnitario <= 0 {
		return fmt.Errorf("precio unitario debe ser mayor a cero")
	}
	return nil
}
