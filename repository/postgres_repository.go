package repository

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq" // Driver PostgreSQL
	"FMgo/models"
)

// DBExecutor define la interfaz para ejecutar operaciones en la base de datos
type DBExecutor interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

type PostgresRepository struct {
	db DBExecutor
}

func NewPostgresRepository(dataSourceName string) (*PostgresRepository, error) {
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	log.Println("Conexión a la base de datos PostgreSQL establecida.")
	return &PostgresRepository{db: db}, nil
}

// --- Implementación de métodos de la interfaz DocumentRepository ---

func (r *PostgresRepository) SaveDocumentoTributario(doc models.DocumentoTributario) error {
	// Lógica SQL para insertar un documento
	query := `INSERT INTO documentos_tributarios (...) VALUES (...)`
	_, err := r.db.Exec(query /* ... parámetros ... */)
	return err
}

func (r *PostgresRepository) GetDocumentoTributario(tipo string, folio int) (*models.DocumentoTributario, error) {
	// Lógica SQL para obtener un documento
	// query := `SELECT ... FROM documentos_tributarios WHERE tipo_documento = $1 AND folio = $2` // Eliminado porque no se usa
	// ...
	return nil, nil // Placeholder
}

func (r *PostgresRepository) UpdateDocumentoTributario(doc models.DocumentoTributario) error {
	// Lógica SQL para actualizar un documento
	query := `UPDATE documentos_tributarios SET ... WHERE id = $1`
	_, err := r.db.Exec(query, doc.ID /* ... otros parámetros ... */)
	return err
}

func (r *PostgresRepository) GetControlFolio(tipoDocumento string) (*models.ControlFolio, error) {
	// query := `SELECT rango_inicial, rango_final, folio_actual, folios_disponibles, ultimo_uso, estado_caf, alerta_generada
	//           FROM control_folios WHERE tipo_documento = $1` // Eliminado porque no se usa
	var control models.ControlFolio
	control.TipoDocumento = tipoDocumento // Asignar el tipo de documento
	// ...
	return &control, nil // Placeholder
}

func (r *PostgresRepository) UpdateControlFolio(control models.ControlFolio) error {
	query := `UPDATE control_folios SET folio_actual = $1, folios_disponibles = $2, ultimo_uso = $3, alerta_generada = $4
              WHERE tipo_documento = $5`
	_, err := r.db.Exec(query,
		control.FolioActual,
		control.FoliosDisponibles,
		control.UltimoUso,
		control.AlertaGenerada,
		control.TipoDocumento,
	)
	return err
}
