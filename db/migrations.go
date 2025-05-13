package db

import (
	"log"

	"github.com/cursor/FMgo/models"

	"gorm.io/gorm"
)

// Migrate ejecuta las migraciones de la base de datos
func Migrate(db *gorm.DB) error {
	log.Println("Iniciando migraciones...")

	// Migrar modelos
	if err := db.AutoMigrate(
		&models.Factura{},
		&models.DetalleFactura{},
		&models.AuditLog{},
	); err != nil {
		return err
	}

	// Crear índices para facturas
	if err := db.Exec(`CREATE INDEX IF NOT EXISTS idx_facturas_rut_emisor ON facturas(rut_emisor)`).Error; err != nil {
		return err
	}
	if err := db.Exec(`CREATE INDEX IF NOT EXISTS idx_facturas_fecha_emision ON facturas(fecha_emision)`).Error; err != nil {
		return err
	}
	if err := db.Exec(`CREATE INDEX IF NOT EXISTS idx_facturas_track_id ON facturas(track_id)`).Error; err != nil {
		return err
	}
	if err := db.Exec(`CREATE INDEX IF NOT EXISTS idx_facturas_estado ON facturas(estado)`).Error; err != nil {
		return err
	}

	// Crear índices para detalles de facturas
	if err := db.Exec(`CREATE INDEX IF NOT EXISTS idx_detalles_factura_id ON detalles_facturas(factura_id)`).Error; err != nil {
		return err
	}

	log.Println("Migraciones completadas exitosamente")
	return nil
}

// createIndexes crea los índices necesarios
func createIndexes() error {
	// Índice para búsqueda por RUT emisor
	if err := DB.Exec(`
        CREATE INDEX IF NOT EXISTS idx_boletas_rut_emisor 
        ON boletas(rut_emisor);
    `).Error; err != nil {
		return err
	}

	// Índice para búsqueda por fecha de emisión
	if err := DB.Exec(`
        CREATE INDEX IF NOT EXISTS idx_boletas_fecha_emision 
        ON boletas(fecha_emision);
    `).Error; err != nil {
		return err
	}

	// Índice para búsqueda por TrackID
	if err := DB.Exec(`
        CREATE INDEX IF NOT EXISTS idx_boletas_track_id 
        ON boletas(track_id);
    `).Error; err != nil {
		return err
	}

	// Índice para búsqueda por estado
	if err := DB.Exec(`
        CREATE INDEX IF NOT EXISTS idx_boletas_estado 
        ON boletas(estado);
    `).Error; err != nil {
		return err
	}

	return nil
}
