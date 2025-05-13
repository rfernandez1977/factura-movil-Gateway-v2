package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/cursor/FMgo/config"
	"github.com/supabase-community/postgrest-go"
)

// DatabaseService maneja las operaciones con la base de datos
type DatabaseService struct {
	config *config.Config
}

// NewDatabaseService crea una nueva instancia del servicio de base de datos
func NewDatabaseService(config *config.Config) *DatabaseService {
	return &DatabaseService{
		config: config,
	}
}

// GetClient retorna el cliente de la base de datos
func (s *DatabaseService) GetClient() *postgrest.Client {
	return s.config.Client
}

// From establece la tabla a consultar
func (s *DatabaseService) From(table string) *postgrest.QueryBuilder {
	query := s.config.Client.From(table)
	return query
}

// Query ejecuta una consulta en la base de datos
func (s *DatabaseService) Query(ctx context.Context, query string, args ...interface{}) (interface{}, error) {
	client := s.GetClient()
	if client == nil {
		return nil, fmt.Errorf("cliente de base de datos no inicializado")
	}

	// Aquí deberías implementar la lógica específica para tu base de datos
	// Por ejemplo, para PostgreSQL podrías usar:
	// return client.(*pgx.Conn).Query(ctx, query, args...)

	return nil, fmt.Errorf("método no implementado")
}

// Execute ejecuta una operación en la base de datos
func (s *DatabaseService) Execute(ctx context.Context, query string, args ...interface{}) error {
	client := s.GetClient()
	if client == nil {
		return fmt.Errorf("cliente de base de datos no inicializado")
	}

	// Aquí deberías implementar la lógica específica para tu base de datos
	// Por ejemplo, para PostgreSQL podrías usar:
	// _, err := client.(*pgx.Conn).Exec(ctx, query, args...)
	// return err

	return fmt.Errorf("método no implementado")
}

// Transaction ejecuta una transacción en la base de datos
func (s *DatabaseService) Transaction(ctx context.Context, fn func(tx interface{}) error) error {
	client := s.GetClient()
	if client == nil {
		return fmt.Errorf("cliente de base de datos no inicializado")
	}

	// Aquí deberías implementar la lógica específica para tu base de datos
	// Por ejemplo, para PostgreSQL podrías usar:
	// return client.(*pgx.Conn).BeginTx(ctx, pgx.TxOptions{}).RunInTx(ctx, fn)

	return fmt.Errorf("método no implementado")
}

// Usuario representa un usuario en la base de datos
type Usuario struct {
	ID               string                 `json:"id"`
	Email            string                 `json:"email"`
	EmailVerified    bool                   `json:"email_verified"`
	Role             string                 `json:"role"`
	CreatedAt        time.Time              `json:"created_at"`
	UpdatedAt        time.Time              `json:"updated_at"`
	LastSignInAt     time.Time              `json:"last_sign_in_at"`
	EmailConfirmedAt time.Time              `json:"email_confirmed_at"`
	IsSuperAdmin     bool                   `json:"is_super_admin"`
	IsAnonymous      bool                   `json:"is_anonymous"`
	IsSSOUser        bool                   `json:"is_sso_user"`
	Providers        []string               `json:"providers"`
	UserMetadata     map[string]interface{} `json:"user_metadata"`
	AppMetadata      map[string]interface{} `json:"app_metadata"`
}

// ObtenerUsuario obtiene un usuario por su ID
func (s *DatabaseService) ObtenerUsuario(ctx context.Context, userID string) (*Usuario, error) {
	var usuario Usuario
	resp, _, err := s.config.Client.From("auth.users").
		Select("*", "", false).
		Eq("id", userID).
		Single().
		Execute()

	if err != nil {
		return nil, fmt.Errorf("error al obtener usuario: %v", err)
	}

	if err := json.Unmarshal(resp, &usuario); err != nil {
		return nil, fmt.Errorf("error al decodificar usuario: %v", err)
	}

	return &usuario, nil
}

// ActualizarUsuario actualiza la información de un usuario
func (s *DatabaseService) ActualizarUsuario(ctx context.Context, userID string, updates map[string]interface{}) error {
	_, _, err := s.config.Client.From("auth.users").
		Update(updates, "", "").
		Eq("id", userID).
		Execute()

	if err != nil {
		return fmt.Errorf("error al actualizar usuario: %v", err)
	}

	return nil
}

// CrearUsuario crea un nuevo usuario
func (s *DatabaseService) CrearUsuario(ctx context.Context, email, password string, metadata map[string]interface{}) (*Usuario, error) {
	// Crear el usuario en auth.users
	userData := map[string]interface{}{
		"email":         email,
		"password":      password,
		"user_metadata": metadata,
	}

	resp, _, err := s.config.Client.From("auth.users").
		Insert(userData, false, "", "", "").
		Single().
		Execute()

	if err != nil {
		return nil, fmt.Errorf("error al crear usuario: %v", err)
	}

	var usuario Usuario
	if err := json.Unmarshal(resp, &usuario); err != nil {
		return nil, fmt.Errorf("error al decodificar usuario: %v", err)
	}

	return &usuario, nil
}

// EliminarUsuario elimina un usuario
func (s *DatabaseService) EliminarUsuario(ctx context.Context, userID string) error {
	_, _, err := s.config.Client.From("auth.users").
		Delete("", "").
		Eq("id", userID).
		Execute()

	if err != nil {
		return fmt.Errorf("error al eliminar usuario: %v", err)
	}

	return nil
}

// ListarUsuarios obtiene una lista de usuarios con filtros opcionales
func (s *DatabaseService) ListarUsuarios(ctx context.Context, filtros map[string]interface{}) ([]Usuario, error) {
	query := s.config.Client.From("auth.users").Select("*", "", false)

	// Aplicar filtros
	for key, value := range filtros {
		query = query.Eq(key, fmt.Sprintf("%v", value))
	}

	resp, _, err := query.Execute()
	if err != nil {
		return nil, fmt.Errorf("error al listar usuarios: %v", err)
	}

	var usuarios []Usuario
	if err := json.Unmarshal(resp, &usuarios); err != nil {
		return nil, fmt.Errorf("error al decodificar usuarios: %v", err)
	}

	return usuarios, nil
}

// VerificarEmail verifica el email de un usuario
func (s *DatabaseService) VerificarEmail(ctx context.Context, userID string) error {
	updates := map[string]interface{}{
		"email_verified":     true,
		"email_confirmed_at": time.Now(),
	}

	return s.ActualizarUsuario(ctx, userID, updates)
}

// ActualizarUltimoAcceso actualiza la fecha del último acceso de un usuario
func (s *DatabaseService) ActualizarUltimoAcceso(ctx context.Context, userID string) error {
	updates := map[string]interface{}{
		"last_sign_in_at": time.Now(),
	}

	return s.ActualizarUsuario(ctx, userID, updates)
}

// FindByID busca un registro por su ID
func (s *DatabaseService) FindByID(id string) (interface{}, error) {
	query := s.From("your_table")
	result, _, err := query.Select("*", "", false).Eq("id", id).Single().Execute()
	if err != nil {
		return nil, err
	}
	return result, nil
}
