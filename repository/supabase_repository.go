package repository

import (
	"context"
	"fmt"

	"FMgo/supabase"
)

// SupabaseRepository implementa el patrón repositorio para Supabase
type SupabaseRepository struct {
	client *supabase.Client
}

// NewSupabaseRepository crea una nueva instancia del repositorio
func NewSupabaseRepository(client *supabase.Client) *SupabaseRepository {
	return &SupabaseRepository{
		client: client,
	}
}

// --------------------------------
// Operaciones para Empresas
// --------------------------------

// CreateEmpresa crea una nueva empresa
func (r *SupabaseRepository) CreateEmpresa(ctx context.Context, empresa *Empresa) (*Empresa, error) {
	// Convertir la empresa a un mapa para la inserción
	data := map[string]interface{}{
		"nombre":       empresa.Nombre,
		"rut":          empresa.RUT,
		"direccion":    empresa.Direccion,
		"telefono":     empresa.Telefono,
		"email":        empresa.Email,
		"rut_firma":    empresa.RUTFirma,
		"nombre_firma": empresa.NombreFirma,
		"clave_firma":  empresa.ClaveFirma,
	}

	// Insertar registro
	result, err := supabase.InsertRecord(ctx, r.client, "empresas", data)
	if err != nil {
		return nil, fmt.Errorf("error insertando empresa: %w", err)
	}

	// Convertir resultado a modelo Empresa
	empresaInsertada := &Empresa{
		ID:          fmt.Sprintf("%v", result["id"]),
		Nombre:      fmt.Sprintf("%v", result["nombre"]),
		RUT:         fmt.Sprintf("%v", result["rut"]),
		Direccion:   fmt.Sprintf("%v", result["direccion"]),
		Telefono:    fmt.Sprintf("%v", result["telefono"]),
		Email:       fmt.Sprintf("%v", result["email"]),
		RUTFirma:    fmt.Sprintf("%v", result["rut_firma"]),
		NombreFirma: fmt.Sprintf("%v", result["nombre_firma"]),
		ClaveFirma:  fmt.Sprintf("%v", result["clave_firma"]),
	}

	return empresaInsertada, nil
}

// GetEmpresaByID obtiene una empresa por su ID
func (r *SupabaseRepository) GetEmpresaByID(ctx context.Context, id string) (*Empresa, error) {
	// Obtener registro
	result, err := supabase.GetRecordByID(ctx, r.client, "empresas", id)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo empresa: %w", err)
	}

	// Convertir resultado a modelo Empresa
	empresa := &Empresa{
		ID:          fmt.Sprintf("%v", result["id"]),
		Nombre:      fmt.Sprintf("%v", result["nombre"]),
		RUT:         fmt.Sprintf("%v", result["rut"]),
		Direccion:   fmt.Sprintf("%v", result["direccion"]),
		Telefono:    fmt.Sprintf("%v", result["telefono"]),
		Email:       fmt.Sprintf("%v", result["email"]),
		RUTFirma:    fmt.Sprintf("%v", result["rut_firma"]),
		NombreFirma: fmt.Sprintf("%v", result["nombre_firma"]),
		ClaveFirma:  fmt.Sprintf("%v", result["clave_firma"]),
	}

	return empresa, nil
}

// GetEmpresaByRUT obtiene una empresa por su RUT
func (r *SupabaseRepository) GetEmpresaByRUT(ctx context.Context, rut string) (*Empresa, error) {
	// Definir filtros
	filters := map[string]string{
		"rut": rut,
	}

	// Buscar registros
	results, err := supabase.QueryRecords(ctx, r.client, "empresas", filters, 1)
	if err != nil {
		return nil, fmt.Errorf("error buscando empresa: %w", err)
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("empresa con RUT %s no encontrada", rut)
	}

	// Convertir resultado a modelo Empresa
	result := results[0]
	empresa := &Empresa{
		ID:          fmt.Sprintf("%v", result["id"]),
		Nombre:      fmt.Sprintf("%v", result["nombre"]),
		RUT:         fmt.Sprintf("%v", result["rut"]),
		Direccion:   fmt.Sprintf("%v", result["direccion"]),
		Telefono:    fmt.Sprintf("%v", result["telefono"]),
		Email:       fmt.Sprintf("%v", result["email"]),
		RUTFirma:    fmt.Sprintf("%v", result["rut_firma"]),
		NombreFirma: fmt.Sprintf("%v", result["nombre_firma"]),
		ClaveFirma:  fmt.Sprintf("%v", result["clave_firma"]),
	}

	return empresa, nil
}

// UpdateEmpresa actualiza una empresa existente
func (r *SupabaseRepository) UpdateEmpresa(ctx context.Context, empresa *Empresa) (*Empresa, error) {
	// Verificar que la empresa tiene un ID
	if empresa.ID == "" {
		return nil, fmt.Errorf("no se puede actualizar una empresa sin ID")
	}

	// Convertir la empresa a un mapa para la actualización
	data := map[string]interface{}{
		"nombre":       empresa.Nombre,
		"rut":          empresa.RUT,
		"direccion":    empresa.Direccion,
		"telefono":     empresa.Telefono,
		"email":        empresa.Email,
		"rut_firma":    empresa.RUTFirma,
		"nombre_firma": empresa.NombreFirma,
		"clave_firma":  empresa.ClaveFirma,
	}

	// Actualizar registro
	result, err := supabase.UpdateRecord(ctx, r.client, "empresas", empresa.ID, data)
	if err != nil {
		return nil, fmt.Errorf("error actualizando empresa: %w", err)
	}

	// Convertir resultado a modelo Empresa
	empresaActualizada := &Empresa{
		ID:          fmt.Sprintf("%v", result["id"]),
		Nombre:      fmt.Sprintf("%v", result["nombre"]),
		RUT:         fmt.Sprintf("%v", result["rut"]),
		Direccion:   fmt.Sprintf("%v", result["direccion"]),
		Telefono:    fmt.Sprintf("%v", result["telefono"]),
		Email:       fmt.Sprintf("%v", result["email"]),
		RUTFirma:    fmt.Sprintf("%v", result["rut_firma"]),
		NombreFirma: fmt.Sprintf("%v", result["nombre_firma"]),
		ClaveFirma:  fmt.Sprintf("%v", result["clave_firma"]),
	}

	return empresaActualizada, nil
}

// DeleteEmpresa elimina una empresa por su ID
func (r *SupabaseRepository) DeleteEmpresa(ctx context.Context, id string) error {
	// Eliminar registro
	err := supabase.DeleteRecord(ctx, r.client, "empresas", id)
	if err != nil {
		return fmt.Errorf("error eliminando empresa: %w", err)
	}

	return nil
}

// ListEmpresas obtiene una lista de empresas
func (r *SupabaseRepository) ListEmpresas(ctx context.Context, limit int) ([]*Empresa, error) {
	// Obtener registros
	results, err := supabase.ListTableData(ctx, r.client, "empresas", limit)
	if err != nil {
		return nil, fmt.Errorf("error listando empresas: %w", err)
	}

	// Convertir resultados a modelo Empresa
	empresas := make([]*Empresa, 0, len(results))
	for _, result := range results {
		empresa := &Empresa{
			ID:          fmt.Sprintf("%v", result["id"]),
			Nombre:      fmt.Sprintf("%v", result["nombre"]),
			RUT:         fmt.Sprintf("%v", result["rut"]),
			Direccion:   fmt.Sprintf("%v", result["direccion"]),
			Telefono:    fmt.Sprintf("%v", result["telefono"]),
			Email:       fmt.Sprintf("%v", result["email"]),
			RUTFirma:    fmt.Sprintf("%v", result["rut_firma"]),
			NombreFirma: fmt.Sprintf("%v", result["nombre_firma"]),
			ClaveFirma:  fmt.Sprintf("%v", result["clave_firma"]),
		}
		empresas = append(empresas, empresa)
	}

	return empresas, nil
}

// --------------------------------
// Operaciones para Documentos
// --------------------------------

// CreateDocumento crea un nuevo documento
func (r *SupabaseRepository) CreateDocumento(ctx context.Context, documento *Documento) (*Documento, error) {
	// Convertir el documento a un mapa para la inserción
	data := map[string]interface{}{
		"empresa_id":       documento.EmpresaID,
		"tipo_documento":   documento.TipoDocumento,
		"numero_documento": documento.NumeroDocumento,
		"fecha_emision":    documento.FechaEmision,
		"monto":            documento.Monto,
		"estado":           documento.Estado,
	}

	// Insertar registro
	result, err := supabase.InsertRecord(ctx, r.client, "documentos", data)
	if err != nil {
		return nil, fmt.Errorf("error insertando documento: %w", err)
	}

	// Convertir resultado a modelo Documento
	documentoInsertado := &Documento{
		ID:              fmt.Sprintf("%v", result["id"]),
		EmpresaID:       fmt.Sprintf("%v", result["empresa_id"]),
		TipoDocumento:   fmt.Sprintf("%v", result["tipo_documento"]),
		NumeroDocumento: fmt.Sprintf("%v", result["numero_documento"]),
		FechaEmision:    fmt.Sprintf("%v", result["fecha_emision"]),
		Monto:           result["monto"].(float64),
		Estado:          fmt.Sprintf("%v", result["estado"]),
	}

	return documentoInsertado, nil
}

// GetDocumentoByID obtiene un documento por su ID
func (r *SupabaseRepository) GetDocumentoByID(ctx context.Context, id string) (*Documento, error) {
	// Obtener registro
	result, err := supabase.GetRecordByID(ctx, r.client, "documentos", id)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo documento: %w", err)
	}

	// Convertir resultado a modelo Documento
	documento := &Documento{
		ID:              fmt.Sprintf("%v", result["id"]),
		EmpresaID:       fmt.Sprintf("%v", result["empresa_id"]),
		TipoDocumento:   fmt.Sprintf("%v", result["tipo_documento"]),
		NumeroDocumento: fmt.Sprintf("%v", result["numero_documento"]),
		FechaEmision:    fmt.Sprintf("%v", result["fecha_emision"]),
		Monto:           result["monto"].(float64),
		Estado:          fmt.Sprintf("%v", result["estado"]),
	}

	return documento, nil
}

// UpdateDocumentoEstado actualiza el estado de un documento
func (r *SupabaseRepository) UpdateDocumentoEstado(ctx context.Context, id string, estado string) error {
	// Actualizar registro
	data := map[string]interface{}{
		"estado": estado,
	}

	_, err := supabase.UpdateRecord(ctx, r.client, "documentos", id, data)
	if err != nil {
		return fmt.Errorf("error actualizando estado del documento: %w", err)
	}

	return nil
}

// ListDocumentosByEmpresa obtiene una lista de documentos de una empresa
func (r *SupabaseRepository) ListDocumentosByEmpresa(ctx context.Context, empresaID string, limit int) ([]*Documento, error) {
	// Definir filtros
	filters := map[string]string{
		"empresa_id": empresaID,
	}

	// Buscar registros
	results, err := supabase.QueryRecords(ctx, r.client, "documentos", filters, limit)
	if err != nil {
		return nil, fmt.Errorf("error listando documentos: %w", err)
	}

	// Convertir resultados a modelo Documento
	documentos := make([]*Documento, 0, len(results))
	for _, result := range results {
		documento := &Documento{
			ID:              fmt.Sprintf("%v", result["id"]),
			EmpresaID:       fmt.Sprintf("%v", result["empresa_id"]),
			TipoDocumento:   fmt.Sprintf("%v", result["tipo_documento"]),
			NumeroDocumento: fmt.Sprintf("%v", result["numero_documento"]),
			FechaEmision:    fmt.Sprintf("%v", result["fecha_emision"]),
			Monto:           result["monto"].(float64),
			Estado:          fmt.Sprintf("%v", result["estado"]),
		}
		documentos = append(documentos, documento)
	}

	return documentos, nil
}

// --------------------------------
// Operaciones para CAFs
// --------------------------------

// GetCAFByTipoDocumento obtiene un CAF por tipo de documento
func (r *SupabaseRepository) GetCAFByTipoDocumento(ctx context.Context, empresaID string, tipoDocumento string) (*CAF, error) {
	// Definir filtros
	filters := map[string]string{
		"empresa_id":     empresaID,
		"tipo_documento": tipoDocumento,
	}

	// Buscar registros
	results, err := supabase.QueryRecords(ctx, r.client, "cafs", filters, 1)
	if err != nil {
		return nil, fmt.Errorf("error buscando CAF: %w", err)
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("CAF para tipo de documento %s no encontrado", tipoDocumento)
	}

	// Convertir resultado a modelo CAF
	result := results[0]
	caf := &CAF{
		ID:               fmt.Sprintf("%v", result["id"]),
		EmpresaID:        fmt.Sprintf("%v", result["empresa_id"]),
		TipoDocumento:    fmt.Sprintf("%v", result["tipo_documento"]),
		Desde:            int(result["desde"].(float64)),
		Hasta:            int(result["hasta"].(float64)),
		FechaVencimiento: fmt.Sprintf("%v", result["fecha_vencimiento"]),
	}

	// Si el archivo está presente, convertirlo a []byte
	if resultado, ok := result["archivo"]; ok && resultado != nil {
		caf.Archivo = []byte(fmt.Sprintf("%v", result["archivo"]))
	}

	return caf, nil
}
