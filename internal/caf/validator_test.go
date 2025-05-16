package caf

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const cafXMLValido = `<?xml version="1.0" encoding="UTF-8"?>
<AUTORIZACION>
	<CAF>
		<DA>
			<RE>76212889-6</RE>
			<TD>33</TD>
			<RNG>
				<D>1</D>
				<H>100</H>
			</RNG>
			<RSAPK>
				<M>2023-01-01T00:00:00Z</M>
				<E>2025-12-31T23:59:59Z</E>
			</RSAPK>
		</DA>
	</CAF>
</AUTORIZACION>`

func TestNewValidator(t *testing.T) {
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
		{
			name:    "caf_vacio",
			cafXML:  "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator, err := NewValidator([]byte(tt.cafXML))
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.NotNil(t, validator)
			assert.NotNil(t, validator.caf)
			assert.NotNil(t, validator.store)
		})
	}
}

func TestValidarFolio(t *testing.T) {
	validator, err := NewValidator([]byte(cafXMLValido))
	assert.NoError(t, err)

	tests := []struct {
		name    string
		folio   int
		wantErr bool
	}{
		{
			name:    "folio_valido",
			folio:   50,
			wantErr: false,
		},
		{
			name:    "folio_bajo_rango",
			folio:   0,
			wantErr: true,
		},
		{
			name:    "folio_sobre_rango",
			folio:   101,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidarFolio(tt.folio)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestValidarRUTEmisor(t *testing.T) {
	validator, err := NewValidator([]byte(cafXMLValido))
	assert.NoError(t, err)

	tests := []struct {
		name    string
		rut     string
		wantErr bool
	}{
		{
			name:    "rut_valido",
			rut:     "76212889-6",
			wantErr: false,
		},
		{
			name:    "rut_invalido",
			rut:     "76123456-7",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidarRUTEmisor(tt.rut)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestValidarTipoDTE(t *testing.T) {
	validator, err := NewValidator([]byte(cafXMLValido))
	assert.NoError(t, err)

	tests := []struct {
		name    string
		tipo    int
		wantErr bool
	}{
		{
			name:    "tipo_valido",
			tipo:    33,
			wantErr: false,
		},
		{
			name:    "tipo_invalido",
			tipo:    39,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidarTipoDTE(tt.tipo)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestFoliosUsados(t *testing.T) {
	validator, err := NewValidator([]byte(cafXMLValido))
	assert.NoError(t, err)

	// Validar y marcar folio como usado
	folio := 50
	err = validator.ValidarFolio(folio)
	assert.NoError(t, err)

	err = validator.MarcarFolioUsado(folio)
	assert.NoError(t, err)

	// Intentar usar el mismo folio
	err = validator.ValidarFolio(folio)
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrFolioUsado)
}
