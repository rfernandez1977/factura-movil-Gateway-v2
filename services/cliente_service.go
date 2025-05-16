package services

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/fmgo/config"
	"github.com/fmgo/models"
	"github.com/supabase-community/postgrest-go"
)

// ClienteService maneja la lógica de negocio de clientes
type ClienteService struct {
	client *postgrest.Client
}

// NewClienteService crea una nueva instancia del servicio de cliente
func NewClienteService(config *config.SupabaseConfig) *ClienteService {
	return &ClienteService{
		client: config.Client.(*postgrest.Client),
	}
}

// GetClienteByCode obtiene un cliente por su código
func (s *ClienteService) GetClienteByCode(code string) (*models.Client, error) {
	var cliente models.Client
	resp, _, err := s.client.From("clientes").
		Select("*", "", false).
		Eq("code", code).
		Single().
		Execute()

	if err != nil {
		return nil, fmt.Errorf("error al obtener cliente: %v", err)
	}

	// Decodificar el resultado en el struct Client
	if err := json.Unmarshal(resp, &cliente); err != nil {
		return nil, fmt.Errorf("error al decodificar cliente: %v", err)
	}

	return &cliente, nil
}

// GetClienteByID obtiene un cliente por su ID
func (s *ClienteService) GetClienteByID(id int) (*models.Client, error) {
	var cliente models.Client
	resp, _, err := s.client.From("clientes").
		Select("*", "", false).
		Eq("id", strconv.Itoa(id)).
		Single().
		Execute()

	if err != nil {
		return nil, fmt.Errorf("error al obtener cliente: %v", err)
	}

	// Decodificar el resultado en el struct Client
	if err := json.Unmarshal(resp, &cliente); err != nil {
		return nil, fmt.Errorf("error al decodificar cliente: %v", err)
	}

	return &cliente, nil
}

// CrearCliente crea un nuevo cliente
func (s *ClienteService) CrearCliente(cliente *models.Client) error {
	// Validar cliente
	if err := s.validarCliente(cliente); err != nil {
		return err
	}

	// Convertir el cliente a JSON para la inserción
	clienteJSON, err := json.Marshal(cliente)
	if err != nil {
		return fmt.Errorf("error al serializar cliente: %v", err)
	}

	// Guardar cliente en Supabase
	_, _, err = s.client.From("clientes").
		Insert(string(clienteJSON), false, "", "", "").
		Execute()

	if err != nil {
		return fmt.Errorf("error al guardar cliente: %v", err)
	}

	return nil
}

// ActualizarCliente actualiza un cliente existente
func (s *ClienteService) ActualizarCliente(cliente *models.Client) error {
	// Validar cliente
	if err := s.validarCliente(cliente); err != nil {
		return err
	}

	// Convertir el cliente a JSON para la actualización
	clienteJSON, err := json.Marshal(cliente)
	if err != nil {
		return fmt.Errorf("error al serializar cliente: %v", err)
	}

	// Actualizar cliente en Supabase
	_, _, err = s.client.From("clientes").
		Update(string(clienteJSON), "", "").
		Eq("id", strconv.Itoa(cliente.ID)).
		Execute()

	if err != nil {
		return fmt.Errorf("error al actualizar cliente: %v", err)
	}

	return nil
}

// EliminarCliente elimina un cliente
func (s *ClienteService) EliminarCliente(id int) error {
	// Eliminar cliente de Supabase
	_, _, err := s.client.From("clientes").
		Delete("", "").
		Eq("id", strconv.Itoa(id)).
		Execute()

	if err != nil {
		return fmt.Errorf("error al eliminar cliente: %v", err)
	}

	return nil
}

// validarCliente valida un cliente antes de crearlo o actualizarlo
func (s *ClienteService) validarCliente(cliente *models.Client) error {
	if cliente.Code == "" {
		return fmt.Errorf("código requerido")
	}
	if cliente.Name == "" {
		return fmt.Errorf("nombre requerido")
	}
	if cliente.Line == "" {
		return fmt.Errorf("giro requerido")
	}
	if cliente.Address == "" {
		return fmt.Errorf("dirección requerida")
	}
	if cliente.Municipality.Name == "" {
		return fmt.Errorf("comuna requerida")
	}
	return nil
}
