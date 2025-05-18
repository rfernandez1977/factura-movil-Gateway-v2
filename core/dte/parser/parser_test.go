package parser

import (
	"testing"
	"time"

	"FMgo/core/dte/types"
)

func TestParseXML(t *testing.T) {
	parser := NewXMLParser()

	tests := []struct {
		name    string
		xmlData []byte
		wantErr bool
	}{
		{
			name: "XML válido",
			xmlData: []byte(`<?xml version="1.0" encoding="UTF-8"?>
<DTE version="1.0">
	<Documento>
		<Encabezado>
			<IdDoc>
				<TipoDTE>33</TipoDTE>
				<Folio>1</Folio>
				<FechaEmision>2024-03-14</FechaEmision>
			</IdDoc>
			<Emisor>
				<RUTEmisor>76212889-6</RUTEmisor>
				<RznSoc>Empresa Test</RznSoc>
				<GiroEmis>Servicios</GiroEmis>
				<DirOrigen>Calle Test 123</DirOrigen>
				<CmnaOrigen>Santiago</CmnaOrigen>
				<CiudadOrigen>Santiago</CiudadOrigen>
			</Emisor>
			<Receptor>
				<RUTRecep>13195458-1</RUTRecep>
				<RznSocRecep>Cliente Test</RznSocRecep>
				<GiroRecep>Comercio</GiroRecep>
				<DirRecep>Av Test 456</DirRecep>
				<CmnaRecep>Santiago</CmnaRecep>
				<CiudadRecep>Santiago</CiudadRecep>
			</Receptor>
			<Totales>
				<MntNeto>1000</MntNeto>
				<TasaIVA>19</TasaIVA>
				<IVA>190</IVA>
				<MntTotal>1190</MntTotal>
			</Totales>
		</Encabezado>
		<Detalle>
			<NroLinDet>1</NroLinDet>
			<NmbItem>Producto Test</NmbItem>
			<QtyItem>1</QtyItem>
			<PrcItem>1000</PrcItem>
			<MontoItem>1000</MontoItem>
		</Detalle>
	</Documento>
</DTE>`),
			wantErr: false,
		},
		{
			name:    "XML inválido",
			xmlData: []byte(`<DTE><Documento>`),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dte, err := parser.ParseXML(tt.xmlData)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseXML() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && dte == nil {
				t.Error("ParseXML() returned nil DTE for valid XML")
			}
		})
	}
}

func TestGenerateXML(t *testing.T) {
	parser := NewXMLParser()

	tests := []struct {
		name    string
		dte     *types.DTE
		wantErr bool
	}{
		{
			name: "DTE válido",
			dte: &types.DTE{
				ID: "TEST-001",
				Documento: types.Documento{
					Encabezado: types.Encabezado{
						IDDocumento: types.IDDocumento{
							TipoDTE:      "33",
							Folio:        1,
							FechaEmision: time.Now(),
						},
						Emisor: types.Emisor{
							RUT:         "76212889-6",
							RazonSocial: "Empresa Test",
							Giro:        "Servicios",
							Direccion:   "Calle Test 123",
							Comuna:      "Santiago",
							Ciudad:      "Santiago",
						},
						Receptor: types.Receptor{
							RUT:         "10138666-K",
							RazonSocial: "Cliente Test",
							Giro:        "Comercio",
							Direccion:   "Av Test 456",
							Comuna:      "Santiago",
							Ciudad:      "Santiago",
						},
						Totales: types.Totales{
							MontoNeto:  1000,
							TasaIVA:    19,
							IVA:        190,
							MontoTotal: 1190,
						},
					},
					Detalles: []types.Detalle{
						{
							NumeroLinea: 1,
							Nombre:      "Producto Test",
							Cantidad:    1,
							Precio:      1000,
							MontoItem:   1000,
						},
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			xmlData, err := parser.GenerateXML(tt.dte)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateXML() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if len(xmlData) == 0 {
					t.Error("GenerateXML() returned empty XML data")
				}

				// Intentar parsear el XML generado
				parsedDTE, err := parser.ParseXML(xmlData)
				if err != nil {
					t.Errorf("Could not parse generated XML: %v", err)
				}

				if parsedDTE == nil {
					t.Error("Parsed DTE is nil")
				}
			}
		})
	}
}
