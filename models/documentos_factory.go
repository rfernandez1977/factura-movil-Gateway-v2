package models

import (
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/google/uuid"
)

// Funciones auxiliares
func GenerateID() string {
	return uuid.New().String()
}

func GenerateHash(data []byte) string {
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}

// Funciones auxiliares para crear documentos
func NewDocumentoTributario(tipoDTE string, folio int, fechaEmision time.Time) *DocumentoTributario {
	return &DocumentoTributario{
		ID:           GenerateID(),
		TipoDTE:      tipoDTE,
		Folio:        folio,
		FechaEmision: fechaEmision,
		Estado:       EstadoDTEPendiente,
		Timestamps: Timestamps{
			Creado:     time.Now(),
			Modificado: time.Now(),
		},
	}
}

func NewDocumento(tipo, nombre, descripcion string, contenido []byte, mimeType string) *Documento {
	return &Documento{
		ID:          GenerateID(),
		Tipo:        tipo,
		Nombre:      nombre,
		Descripcion: descripcion,
		Contenido:   contenido,
		MimeType:    mimeType,
		Size:        int64(len(contenido)),
		Hash:        GenerateHash(contenido),
		Timestamps: Timestamps{
			Creado:     time.Now(),
			Modificado: time.Now(),
		},
	}
}

func NewDocumentoAlmacenado(tipo, nombre, descripcion, ruta string, mimeType string, size int64) *DocumentoAlmacenado {
	return &DocumentoAlmacenado{
		ID:          GenerateID(),
		Tipo:        tipo,
		Nombre:      nombre,
		Descripcion: descripcion,
		Ruta:        ruta,
		MimeType:    mimeType,
		Size:        size,
		Timestamps: Timestamps{
			Creado:     time.Now(),
			Modificado: time.Now(),
		},
	}
}

func NewDocumentoSeguro(tipo, nombre, descripcion string, contenido []byte, mimeType string) *DocumentoSeguro {
	return &DocumentoSeguro{
		ID:          GenerateID(),
		Tipo:        tipo,
		Nombre:      nombre,
		Descripcion: descripcion,
		Contenido:   contenido,
		MimeType:    mimeType,
		Size:        int64(len(contenido)),
		Hash:        GenerateHash(contenido),
		Timestamps: Timestamps{
			Creado:     time.Now(),
			Modificado: time.Now(),
		},
	}
}
