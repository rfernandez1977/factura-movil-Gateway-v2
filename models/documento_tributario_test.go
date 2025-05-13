package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestDocumentoTributario(t *testing.T) {
	// Crear un documento tributario de prueba
	tiempo := time.Now()
	objID, _ := primitive.ObjectIDFromHex("5f50cf13c56e0a1d9b4fbe5a")
	doc := DocumentoTributario{
		ID:           objID,
		TipoDTE:      TipoFactura,
		Folio:        1,
		FechaEmision: tiempo,
		MontoTotal:   10000,
		Estado:       EstadoDocumentoEnviado,
	}

	// Verificar que los campos se hayan asignado correctamente
	assert.Equal(t, objID, doc.ID)
	assert.Equal(t, TipoFactura, doc.TipoDTE)
	assert.Equal(t, 1, doc.Folio)
	assert.Equal(t, tiempo, doc.FechaEmision)
	assert.Equal(t, float64(10000), doc.MontoTotal)
	assert.Equal(t, EstadoDocumentoEnviado, doc.Estado)
}

func TestControlFolio(t *testing.T) {
	// Crear un control de folio de prueba
	tiempo := time.Now()
	control := ControlFolio{
		TipoDocumento:     "33",
		RangoInicial:      1,
		RangoFinal:        100,
		FolioActual:       5,
		FoliosDisponibles: 95,
		UltimoUso:         tiempo,
		EstadoCAF:         "ACTIVO",
		AlertaGenerada:    false,
	}

	// Verificar que los campos se hayan asignado correctamente
	assert.Equal(t, "33", control.TipoDocumento)
	assert.Equal(t, 1, control.RangoInicial)
	assert.Equal(t, 100, control.RangoFinal)
	assert.Equal(t, 5, control.FolioActual)
	assert.Equal(t, 95, control.FoliosDisponibles)
	assert.Equal(t, tiempo, control.UltimoUso)
	assert.Equal(t, "ACTIVO", control.EstadoCAF)
	assert.Equal(t, false, control.AlertaGenerada)
}