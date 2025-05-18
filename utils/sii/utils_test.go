package sii_test

import (
	"bytes"
	"testing"
	"time"

	"FMgo/models"
	"FMgo/utils/sii"
)

func TestGenerarXMLDTE(t *testing.T) {
	// Crear documento de prueba
	doc := &models.DocumentoTributario{
		ID:           "DTE_1",
		TipoDTE:      "33",
		Folio:        1,
		RutEmisor:    "76212889-6",
		RutReceptor:  "76555555-5",
		MontoTotal:   119000,
		Estado:       "PENDIENTE",
		FechaEmision: time.Now(),
	}

	// Crear empresa de prueba
	empresa := &models.Empresa{
		RUT:         "76212889-6",
		RazonSocial: "Empresa de Prueba",
		Giro:        "Servicios",
		Direccion:   "Calle Principal 123",
		Comuna:      "Santiago",
		Ciudad:      "Santiago",
	}

	// Generar XML
	xmlData, err := sii.GenerarXMLDTE(doc, empresa)
	if err != nil {
		t.Fatalf("Error generando XML: %v", err)
	}

	// Validar que el XML no esté vacío
	if len(xmlData) == 0 {
		t.Error("XML generado está vacío")
	}

	// Validar que el XML comience con la declaración XML
	if !bytes.HasPrefix(xmlData, []byte("<?xml")) {
		t.Error("XML generado no comienza con la declaración XML")
	}
}

func TestValidarRespuestaSII(t *testing.T) {
	tests := []struct {
		name    string
		resp    *models.RespuestaSII
		wantErr bool
	}{
		{
			name: "Respuesta válida",
			resp: &models.RespuestaSII{
				Estado:       "OK",
				Glosa:        "Documento procesado correctamente",
				TrackID:      "123",
				FechaProceso: time.Now(),
			},
			wantErr: false,
		},
		{
			name: "Respuesta con error",
			resp: &models.RespuestaSII{
				Estado:       "ERROR",
				Glosa:        "Error en el documento",
				TrackID:      "123",
				FechaProceso: time.Now(),
				Errores: []models.ErrorSII{
					{
						Codigo:      "001",
						Descripcion: "Error de validación",
						Detalle:     "Detalle del error",
					},
				},
			},
			wantErr: true,
		},
		{
			name:    "Respuesta nula",
			resp:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := sii.ValidarRespuestaSII(tt.resp)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidarRespuestaSII() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestProcesarRespuestaSII(t *testing.T) {
	tests := []struct {
		name    string
		resp    *models.RespuestaSII
		want    *models.EstadoSII
		wantErr bool
	}{
		{
			name: "Respuesta exitosa",
			resp: &models.RespuestaSII{
				Estado:       "OK",
				Glosa:        "Documento procesado correctamente",
				TrackID:      "123",
				FechaProceso: time.Now(),
			},
			want: &models.EstadoSII{
				Estado:  "ACEPTADO",
				Glosa:   "Documento procesado correctamente",
				TrackID: "123",
			},
			wantErr: false,
		},
		{
			name: "Respuesta con error",
			resp: &models.RespuestaSII{
				Estado:       "ERROR",
				Glosa:        "Error en el documento",
				TrackID:      "123",
				FechaProceso: time.Now(),
				Errores: []models.ErrorSII{
					{
						Codigo:      "001",
						Descripcion: "Error de validación",
						Detalle:     "El documento no cumple con los requisitos",
					},
				},
			},
			want: &models.EstadoSII{
				Estado:  "RECHAZADO",
				Glosa:   "Error en el documento",
				TrackID: "123",
				Errores: []models.ErrorSII{
					{
						Codigo:      "001",
						Descripcion: "Error de validación",
						Detalle:     "El documento no cumple con los requisitos",
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			estado, err := sii.ProcesarRespuestaSII(tt.resp)
			if (err != nil) != tt.wantErr {
				t.Errorf("ProcesarRespuestaSII() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if estado.Estado != tt.want.Estado {
					t.Errorf("ProcesarRespuestaSII() Estado = %v, want %v", estado.Estado, tt.want.Estado)
				}
				if estado.Glosa != tt.want.Glosa {
					t.Errorf("ProcesarRespuestaSII() Glosa = %v, want %v", estado.Glosa, tt.want.Glosa)
				}
				if estado.TrackID != tt.want.TrackID {
					t.Errorf("ProcesarRespuestaSII() TrackID = %v, want %v", estado.TrackID, tt.want.TrackID)
				}
				if len(estado.Errores) != len(tt.want.Errores) {
					t.Errorf("ProcesarRespuestaSII() Errores length = %v, want %v", len(estado.Errores), len(tt.want.Errores))
				}
			}
		})
	}
}
