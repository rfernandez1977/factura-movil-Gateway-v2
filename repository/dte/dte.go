package dte

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/cursor/FMgo/core/dte"
)

// Repository maneja el almacenamiento de DTE
type Repository struct {
	db *sql.DB
}

// NewRepository crea un nuevo repositorio de DTE
func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// Save guarda un DTE en la base de datos
func (r *Repository) Save(dte *dte.DTE) error {
	// Iniciar transacción
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("error al iniciar transacción: %w", err)
	}
	defer tx.Rollback()

	// Insertar DTE
	query := `
		INSERT INTO dte (
			id, version, tipo_dte, folio, fecha_emision,
			rut_emisor, razon_social_emisor, giro_emisor,
			direccion_emisor, comuna_emisor, ciudad_emisor,
			correo_emisor, rut_receptor, razon_social_receptor,
			giro_receptor, direccion_receptor, comuna_receptor,
			ciudad_receptor, monto_neto, tasa_iva, monto_iva,
			monto_total, track_id, estado, fecha_creacion
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12,
			$13, $14, $15, $16, $17, $18, $19, $20, $21, $22,
			$23, $24, $25, $26
		)
	`

	_, err = tx.Exec(query,
		dte.Documento.ID,
		dte.Version,
		dte.Documento.Encabezado.ID.TipoDTE,
		dte.Documento.Encabezado.ID.Folio,
		dte.Documento.Encabezado.ID.FechaEmision,
		dte.Documento.Encabezado.Emisor.RUT,
		dte.Documento.Encabezado.Emisor.RazonSocial,
		dte.Documento.Encabezado.Emisor.Giro,
		dte.Documento.Encabezado.Emisor.Direccion,
		dte.Documento.Encabezado.Emisor.Comuna,
		dte.Documento.Encabezado.Emisor.Ciudad,
		dte.Documento.Encabezado.Emisor.Correo,
		dte.Documento.Encabezado.Receptor.RUT,
		dte.Documento.Encabezado.Receptor.RazonSocial,
		dte.Documento.Encabezado.Receptor.Giro,
		dte.Documento.Encabezado.Receptor.Direccion,
		dte.Documento.Encabezado.Receptor.Comuna,
		dte.Documento.Encabezado.Receptor.Ciudad,
		dte.Documento.Encabezado.Totales.MontoNeto,
		dte.Documento.Encabezado.Totales.IVA,
		dte.Documento.Encabezado.Totales.MontoIVA,
		dte.Documento.Encabezado.Totales.MontoTotal,
		"",          // track_id
		"PENDIENTE", // estado
		time.Now(),
	)
	if err != nil {
		return fmt.Errorf("error al insertar DTE: %w", err)
	}

	// Insertar detalles
	for _, detalle := range dte.Documento.Detalle {
		query = `
			INSERT INTO dte_detalle (
				dte_id, numero_linea, nombre_item,
				cantidad, precio_unitario, monto_item
			) VALUES ($1, $2, $3, $4, $5, $6)
		`

		_, err = tx.Exec(query,
			dte.Documento.ID,
			detalle.NroLinDet,
			detalle.NombreItem,
			detalle.Cantidad,
			detalle.PrecioUnit,
			detalle.MontoItem,
		)
		if err != nil {
			return fmt.Errorf("error al insertar detalle: %w", err)
		}
	}

	// Confirmar transacción
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("error al confirmar transacción: %w", err)
	}

	return nil
}

// Get obtiene un DTE por su ID
func (r *Repository) Get(id string) (*dte.DTE, error) {
	// Obtener DTE
	query := `
		SELECT
			id, version, tipo_dte, folio, fecha_emision,
			rut_emisor, razon_social_emisor, giro_emisor,
			direccion_emisor, comuna_emisor, ciudad_emisor,
			correo_emisor, rut_receptor, razon_social_receptor,
			giro_receptor, direccion_receptor, comuna_receptor,
			ciudad_receptor, monto_neto, tasa_iva, monto_iva,
			monto_total, track_id, estado
		FROM dte
		WHERE id = $1
	`

	var d dte.DTE
	var fechaEmision string
	err := r.db.QueryRow(query, id).Scan(
		&d.Documento.ID,
		&d.Version,
		&d.Documento.Encabezado.ID.TipoDTE,
		&d.Documento.Encabezado.ID.Folio,
		&fechaEmision,
		&d.Documento.Encabezado.Emisor.RUT,
		&d.Documento.Encabezado.Emisor.RazonSocial,
		&d.Documento.Encabezado.Emisor.Giro,
		&d.Documento.Encabezado.Emisor.Direccion,
		&d.Documento.Encabezado.Emisor.Comuna,
		&d.Documento.Encabezado.Emisor.Ciudad,
		&d.Documento.Encabezado.Emisor.Correo,
		&d.Documento.Encabezado.Receptor.RUT,
		&d.Documento.Encabezado.Receptor.RazonSocial,
		&d.Documento.Encabezado.Receptor.Giro,
		&d.Documento.Encabezado.Receptor.Direccion,
		&d.Documento.Encabezado.Receptor.Comuna,
		&d.Documento.Encabezado.Receptor.Ciudad,
		&d.Documento.Encabezado.Totales.MontoNeto,
		&d.Documento.Encabezado.Totales.IVA,
		&d.Documento.Encabezado.Totales.MontoIVA,
		&d.Documento.Encabezado.Totales.MontoTotal,
		&d.TrackID,
		&d.Estado,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("DTE no encontrado")
		}
		return nil, fmt.Errorf("error al obtener DTE: %w", err)
	}

	// Obtener detalles
	query = `
		SELECT
			numero_linea, nombre_item, cantidad,
			precio_unitario, monto_item
		FROM dte_detalle
		WHERE dte_id = $1
		ORDER BY numero_linea
	`

	rows, err := r.db.Query(query, id)
	if err != nil {
		return nil, fmt.Errorf("error al obtener detalles: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var detalle dte.Detalle
		err := rows.Scan(
			&detalle.NroLinDet,
			&detalle.NombreItem,
			&detalle.Cantidad,
			&detalle.PrecioUnit,
			&detalle.MontoItem,
		)
		if err != nil {
			return nil, fmt.Errorf("error al escanear detalle: %w", err)
		}
		d.Documento.Detalle = append(d.Documento.Detalle, detalle)
	}

	return &d, nil
}

// UpdateEstado actualiza el estado de un DTE
func (r *Repository) UpdateEstado(id string, estado string, trackID string) error {
	query := `
		UPDATE dte
		SET estado = $1, track_id = $2, fecha_actualizacion = $3
		WHERE id = $4
	`

	_, err := r.db.Exec(query, estado, trackID, time.Now(), id)
	if err != nil {
		return fmt.Errorf("error al actualizar estado: %w", err)
	}

	return nil
}

// GetByFolio obtiene un DTE por su folio y tipo
func (r *Repository) GetByFolio(tipoDTE int, folio int) (*dte.DTE, error) {
	query := `
		SELECT id
		FROM dte
		WHERE tipo_dte = $1 AND folio = $2
	`

	var id string
	err := r.db.QueryRow(query, tipoDTE, folio).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("DTE no encontrado")
		}
		return nil, fmt.Errorf("error al obtener DTE: %w", err)
	}

	return r.Get(id)
}

// GetByRUT obtiene los DTE de un emisor
func (r *Repository) GetByRUT(rut string) ([]*dte.DTE, error) {
	query := `
		SELECT id
		FROM dte
		WHERE rut_emisor = $1
		ORDER BY fecha_emision DESC
	`

	rows, err := r.db.Query(query, rut)
	if err != nil {
		return nil, fmt.Errorf("error al obtener DTE: %w", err)
	}
	defer rows.Close()

	var dtes []*dte.DTE
	for rows.Next() {
		var id string
		err := rows.Scan(&id)
		if err != nil {
			return nil, fmt.Errorf("error al escanear ID: %w", err)
		}

		d, err := r.Get(id)
		if err != nil {
			return nil, fmt.Errorf("error al obtener DTE %s: %w", id, err)
		}

		dtes = append(dtes, d)
	}

	return dtes, nil
}

// GetByFecha obtiene los DTE emitidos en un rango de fechas
func (r *Repository) GetByFecha(desde, hasta time.Time) ([]*dte.DTE, error) {
	query := `
		SELECT id
		FROM dte
		WHERE fecha_emision BETWEEN $1 AND $2
		ORDER BY fecha_emision DESC
	`

	rows, err := r.db.Query(query, desde, hasta)
	if err != nil {
		return nil, fmt.Errorf("error al obtener DTE: %w", err)
	}
	defer rows.Close()

	var dtes []*dte.DTE
	for rows.Next() {
		var id string
		err := rows.Scan(&id)
		if err != nil {
			return nil, fmt.Errorf("error al escanear ID: %w", err)
		}

		d, err := r.Get(id)
		if err != nil {
			return nil, fmt.Errorf("error al obtener DTE %s: %w", id, err)
		}

		dtes = append(dtes, d)
	}

	return dtes, nil
}
