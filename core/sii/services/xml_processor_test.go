package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestXMLProcessor_ExtraerCertificado(t *testing.T) {
	processor := NewXMLProcessor()

	xmlData := []byte(`
		<?xml version="1.0" encoding="UTF-8"?>
		<DTE xmlns="http://www.sii.cl/SiiDte">
			<Signature xmlns="http://www.w3.org/2000/09/xmldsig#">
				<KeyInfo>
					<X509Data>
						<X509Certificate>MIIFZjCCBE6gAwIBAgIQGHf2</X509Certificate>
					</X509Data>
				</KeyInfo>
			</Signature>
		</DTE>
	`)

	cert, err := processor.extraerCertificado(xmlData)
	assert.NoError(t, err)
	assert.Equal(t, "MIIFZjCCBE6gAwIBAgIQGHf2", cert)
}

func TestXMLProcessor_ExtraerFirma(t *testing.T) {
	processor := NewXMLProcessor()

	xmlData := []byte(`
		<?xml version="1.0" encoding="UTF-8"?>
		<DTE xmlns="http://www.sii.cl/SiiDte">
			<Signature xmlns="http://www.w3.org/2000/09/xmldsig#">
				<SignatureValue>ABC123</SignatureValue>
			</Signature>
		</DTE>
	`)

	firma, err := processor.extraerFirma(xmlData)
	assert.NoError(t, err)
	assert.Equal(t, "ABC123", firma)
}

func TestXMLProcessor_ValidarEstructuraXML(t *testing.T) {
	processor := NewXMLProcessor()

	// Caso válido
	xmlValido := []byte(`
		<?xml version="1.0" encoding="UTF-8"?>
		<DTE xmlns="http://www.sii.cl/SiiDte">
			<Documento ID="DOC001">
				<Encabezado>
					<IdDoc>
						<TipoDTE>33</TipoDTE>
					</IdDoc>
					<Emisor>
						<RUTEmisor>76555555-5</RUTEmisor>
					</Emisor>
					<Receptor>
						<RUTRecep>66666666-6</RUTRecep>
					</Receptor>
				</Encabezado>
			</Documento>
		</DTE>
	`)

	err := processor.validarEstructuraXML(xmlValido)
	assert.NoError(t, err)

	// Caso inválido - falta elemento requerido
	xmlInvalido := []byte(`
		<?xml version="1.0" encoding="UTF-8"?>
		<DTE xmlns="http://www.sii.cl/SiiDte">
			<Documento ID="DOC001">
				<Encabezado>
					<IdDoc>
						<TipoDTE>33</TipoDTE>
					</IdDoc>
					<!-- Falta Emisor -->
					<Receptor>
						<RUTRecep>66666666-6</RUTRecep>
					</Receptor>
				</Encabezado>
			</Documento>
		</DTE>
	`)

	err = processor.validarEstructuraXML(xmlInvalido)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "elemento requerido no encontrado: <Emisor")
}

func TestXMLProcessor_LimpiarXML(t *testing.T) {
	processor := NewXMLProcessor()

	xmlSucio := []byte("<?xml version=\"1.0\"?>\x00<DTE>\x01<Documento>\x02</Documento>\x03</DTE>")
	xmlLimpio := processor.limpiarXML(xmlSucio)

	expected := []byte("<?xml version=\"1.0\"?><DTE><Documento></Documento></DTE>")
	assert.Equal(t, expected, xmlLimpio)
}
