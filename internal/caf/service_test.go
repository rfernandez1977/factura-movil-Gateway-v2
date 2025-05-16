package caf

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewService(t *testing.T) {
	service := NewService()
	assert.NotNil(t, service)
	assert.NotNil(t, service.validators)
}

func TestRegistrarCAF(t *testing.T) {
	service := NewService()

	tests := []struct {
		name    string
		cafXML  string
		wantErr bool
	}{
		{
			name:    "caf_valido",
			cafXML:  cafXMLValido,
			wantErr: false,
		},
		{
			name:    "caf_invalido",
			cafXML:  "<xml>malformado</xml>",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.RegistrarCAF([]byte(tt.cafXML))
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestServiceValidarFolio(t *testing.T) {
	service := NewService()
	err := service.RegistrarCAF([]byte(cafXMLValido))
	assert.NoError(t, err)

	tests := []struct {
		name    string
		rut     string
		tipoDTE int
		folio   int
		wantErr bool
	}{
		{
			name:    "folio_valido",
			rut:     "76212889-6",
			tipoDTE: 33,
			folio:   50,
			wantErr: false,
		},
		{
			name:    "folio_invalido",
			rut:     "76212889-6",
			tipoDTE: 33,
			folio:   101,
			wantErr: true,
		},
		{
			name:    "rut_no_registrado",
			rut:     "76123456-7",
			tipoDTE: 33,
			folio:   50,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.ValidarFolio(tt.rut, tt.tipoDTE, tt.folio)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestObtenerEstadoCAF(t *testing.T) {
	service := NewService()
	err := service.RegistrarCAF([]byte(cafXMLValido))
	assert.NoError(t, err)

	tests := []struct {
		name    string
		rut     string
		tipoDTE int
		wantErr bool
	}{
		{
			name:    "caf_existente",
			rut:     "76212889-6",
			tipoDTE: 33,
			wantErr: false,
		},
		{
			name:    "caf_no_existente",
			rut:     "76123456-7",
			tipoDTE: 33,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			estado, err := service.ObtenerEstadoCAF(tt.rut, tt.tipoDTE)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.NotNil(t, estado)
			assert.Equal(t, tt.rut, estado.RutEmisor)
			assert.Equal(t, tt.tipoDTE, estado.TipoDTE)
		})
	}
}
