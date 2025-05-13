package main

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/cursor/FMgo/utils"
)

type TestInput struct {
	HasTaxes bool `json:"hasTaxes"`
	Details  []struct {
		Product struct {
			Unit struct {
				Code string `json:"code"`
			} `json:"unit"`
			Price float64 `json:"price"`
			Name  string  `json:"name"`
			Code  string  `json:"code"`
		} `json:"product"`
		Position int     `json:"position"`
		Quantity float64 `json:"quantity"`
	} `json:"details"`
	Client struct {
		Address      string `json:"address"`
		Name         string `json:"name"`
		Municipality string `json:"municipality"`
		Line         string `json:"line"`
		Code         string `json:"code"`
	} `json:"client"`
	Date     string `json:"date"`
	Currency string `json:"currency"`
}

func TestCopec(t *testing.T) {
	// Cargar configuración
	configBytes, err := os.ReadFile("config.json")
	if err != nil {
		t.Fatalf("Error leyendo archivo de configuración: %v", err)
	}

	var config struct {
		Empresa struct {
			Rut         string   `json:"rut"`
			RazonSocial string   `json:"razonSocial"`
			Giro        string   `json:"giro"`
			Correo      string   `json:"correo"`
			Actecos     []string `json:"actecos"`
			Direccion   string   `json:"direccion"`
			Comuna      string   `json:"comuna"`
			Ciudad      string   `json:"ciudad"`
		} `json:"empresa"`
		Resolucion struct {
			Numero string `json:"numero"`
			Fecha  string `json:"fecha"`
		} `json:"resolucion"`
		Firma struct {
			Rut         string `json:"rut"`
			Certificado string `json:"certificado"`
			Password    string `json:"password"`
		} `json:"firma"`
	}

	if err := json.Unmarshal(configBytes, &config); err != nil {
		t.Fatalf("Error parseando configuración: %v", err)
	}

	// Cargar input de prueba
	inputBytes, err := os.ReadFile("testdata/copec_input.json")
	if err != nil {
		t.Fatalf("Error leyendo archivo de entrada: %v", err)
	}

	var input TestInput
	if err := json.Unmarshal(inputBytes, &input); err != nil {
		t.Fatalf("Error parseando JSON de entrada: %v", err)
	}

	// Crear emisor
	emisor := utils.Emisor{
		RUT:         config.Empresa.Rut,
		RazonSocial: config.Empresa.RazonSocial,
		Giro:        config.Empresa.Giro,
		Direccion:   config.Empresa.Direccion,
		Comuna:      config.Empresa.Comuna,
		Ciudad:      config.Empresa.Ciudad,
		Correo:      config.Empresa.Correo,
		Actecos:     config.Empresa.Actecos,
	}

	// Crear receptor
	receptor := utils.Receptor{
		RUT:         input.Client.Code,
		RazonSocial: input.Client.Name,
		Giro:        input.Client.Line,
		Direccion:   input.Client.Address,
		Comuna:      input.Client.Municipality,
		Ciudad:      input.Client.Municipality,
	}

	// Cargar CAF
	cafManager := utils.NewCAFManager()
	if err := cafManager.CargarCAF("caf_test/33-cert.xml"); err != nil {
		t.Fatalf("Error cargando CAF: %v", err)
	}

	caf, err := cafManager.ObtenerCAF(33)
	if err != nil {
		t.Fatalf("Error obteniendo CAF: %v", err)
	}

	// Generar DTE
	generator := utils.NewDTEGenerator(caf)

	// Convertir detalles a formato DTE
	var detalles []utils.DetalleDTE
	for _, detalle := range input.Details {
		detalles = append(detalles, utils.DetalleDTE{
			NroLinDet:  detalle.Position,
			NombreItem: detalle.Product.Name,
			Cantidad:   detalle.Quantity,
			PrecioUnit: int(detalle.Product.Price),
			MontoItem:  int(detalle.Product.Price * detalle.Quantity),
		})
	}

	dte, err := generator.GenerarDTE(emisor, receptor, detalles)
	if err != nil {
		t.Fatalf("Error generando DTE: %v", err)
	}

	// Generar sobre
	fechaResolucion, err := time.Parse("2006-01-02", config.Resolucion.Fecha)
	if err != nil {
		t.Fatalf("Error parseando fecha de resolución: %v", err)
	}

	sobre, err := generator.GenerarSobre(emisor, receptor, []*utils.DTE{dte}, config.Resolucion.Numero, fechaResolucion)
	if err != nil {
		t.Fatalf("Error generando sobre: %v", err)
	}

	// Crear cliente SII mock
	client := utils.NewMockSIIClient()

	// Obtener semilla
	semilla, err := client.ObtenerSemilla()
	if err != nil {
		t.Fatalf("Error obteniendo semilla: %v", err)
	}

	// Obtener token
	token, err := client.ObtenerToken(semilla)
	if err != nil {
		t.Fatalf("Error obteniendo token: %v", err)
	}

	// Enviar DTE
	err = client.EnviarDTE(sobre, token)
	if err != nil {
		t.Fatalf("Error enviando DTE: %v", err)
	}

	fmt.Printf("DTE enviado exitosamente\n")
	fmt.Printf("- Folio: %d\n", dte.Documento.Encabezado.ID.Folio)
	fmt.Printf("- Fecha Emisión: %s\n", dte.Documento.Encabezado.ID.FechaEmision)
	fmt.Printf("- Monto Total: %d\n", dte.Documento.Encabezado.Totales.MontoTotal)
}
