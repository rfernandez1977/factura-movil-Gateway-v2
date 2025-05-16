package caf

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	// Ejemplo de CAF con firma (valores de prueba)
	cafConFirma = `<?xml version="1.0" encoding="UTF-8"?>
<DA>
	<RE>76212889-6</RE>
	<TD>33</TD>
	<RNG>
		<D>1</D>
		<H>100</H>
	</RNG>
	<RSAPK>
		<M>4hzcHj89dX1P5K7+zx5jntS+G9NLQ8L7BgNT28JwYvHHJMJaNJrJpW5EE/Zqn1rG</M>
		<E>Aw==</E>
	</RSAPK>
</DA>`

	// Módulo y exponente de prueba (valores de ejemplo)
	moduloBase64    = "4hzcHj89dX1P5K7+zx5jntS+G9NLQ8L7BgNT28JwYvHHJMJaNJrJpW5EE/Zqn1rG"
	exponenteBase64 = "Aw=="
)

func TestParsePublicKey(t *testing.T) {
	key := RSAKey{
		Modulus:  moduloBase64,
		Exponent: exponenteBase64,
	}

	publicKey, err := ParsePublicKey(key)
	assert.NoError(t, err)
	assert.NotNil(t, publicKey)
	assert.Equal(t, 3, publicKey.E) // El exponente debe ser 3 (Aw== en base64)
}

func TestParsePublicKeyInvalid(t *testing.T) {
	tests := []struct {
		name    string
		key     RSAKey
		wantErr bool
		errType error
	}{
		{
			name: "modulo_invalido",
			key: RSAKey{
				Modulus:  "invalid base64",
				Exponent: exponenteBase64,
			},
			wantErr: true,
			errType: ErrClavePublicaInvalida,
		},
		{
			name: "exponente_invalido",
			key: RSAKey{
				Modulus:  moduloBase64,
				Exponent: "invalid base64",
			},
			wantErr: true,
			errType: ErrClavePublicaInvalida,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParsePublicKey(tt.key)
			assert.Error(t, err)
			assert.ErrorIs(t, err, tt.errType)
		})
	}
}

func TestExtractSignedInfo(t *testing.T) {
	info, err := ExtractSignedInfo([]byte(cafConFirma))
	assert.NoError(t, err)
	assert.NotNil(t, info)
	assert.Equal(t, "76212889-6", info.RE)
	assert.Equal(t, 33, info.TD)
	assert.Equal(t, 1, info.RNG.D)
	assert.Equal(t, 100, info.RNG.H)
	assert.Equal(t, moduloBase64, info.RSAPK.Modulus)
	assert.Equal(t, exponenteBase64, info.RSAPK.Exponent)
}

func TestVerifySignature(t *testing.T) {
	// Extraer información firmada
	info, err := ExtractSignedInfo([]byte(cafConFirma))
	assert.NoError(t, err)

	// Parsear clave pública
	publicKey, err := ParsePublicKey(info.RSAPK)
	assert.NoError(t, err)

	// Crear verificador
	verifier := NewSignatureVerifier(publicKey)

	// Canonicalizar datos
	canonicalData, err := CanonicalizeXML([]byte(cafConFirma))
	assert.NoError(t, err)

	// Verificar firma
	err = verifier.VerifySignature(canonicalData, "MCwCFQDjX8CfAUVDJWALC/Z4T8gL3My7EwIVALgTNfwJhh9L2vTzXxvTKjOQzQ==")
	// Nota: Esta prueba fallará porque usamos datos de ejemplo
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrVerificacionFirma)
}

func TestVerifySignatureInvalid(t *testing.T) {
	// Crear clave pública de prueba
	key := RSAKey{
		Modulus:  moduloBase64,
		Exponent: exponenteBase64,
	}
	publicKey, err := ParsePublicKey(key)
	assert.NoError(t, err)

	verifier := NewSignatureVerifier(publicKey)

	tests := []struct {
		name      string
		data      []byte
		signature string
		wantErr   bool
		errType   error
	}{
		{
			name:      "firma_invalida",
			data:      []byte("test data"),
			signature: "invalid base64",
			wantErr:   true,
			errType:   ErrVerificacionFirma,
		},
		{
			name:      "firma_incorrecta",
			data:      []byte("test data"),
			signature: base64.StdEncoding.EncodeToString([]byte("wrong signature")),
			wantErr:   true,
			errType:   ErrVerificacionFirma,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := verifier.VerifySignature(tt.data, tt.signature)
			assert.Error(t, err)
			assert.ErrorIs(t, err, tt.errType)
		})
	}
}
