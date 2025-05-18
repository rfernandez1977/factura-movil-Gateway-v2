package supabase

import (
	"FMgo/supabase"
)

// RepositoryFactory es una fábrica para crear repositorios de Supabase
type RepositoryFactory struct {
	client *supabase.Client
}

// NewRepositoryFactory crea una nueva fábrica de repositorios
func NewRepositoryFactory(client *supabase.Client) *RepositoryFactory {
	return &RepositoryFactory{
		client: client,
	}
}

// NewDocumentoRepository crea un nuevo repositorio de documentos tributarios
func (f *RepositoryFactory) NewDocumentoRepository() DocumentoTributarioRepository {
	return NewSupabaseDocumentoRepository(f.client)
}

// Aquí se pueden agregar más métodos para crear otros tipos de repositorios
// como empresas, usuarios, etc.
