package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v4"
)

// Estructuras para parsear el input JSON

type Input struct {
	HasTaxes    bool     `json:"hasTaxes"`
	Client      Client   `json:"client"`
	Details     []Detail `json:"details"`
	NetTotal    float64  `json:"netTotal"`
	Discounts   []any    `json:"discounts"`
	Date        string   `json:"date"`
	ExemptTotal float64  `json:"exemptTotal"`
	OtherTaxes  float64  `json:"otherTaxes"`
	Taxes       float64  `json:"taxes"`
}

type Client struct {
	ID           int          `json:"id"`
	Address      string       `json:"address"`
	Email        string       `json:"email"`
	Name         string       `json:"name"`
	Municipality Municipality `json:"municipality"`
	Line         string       `json:"line"`
	Code         string       `json:"code"`
}

type Municipality struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}

type Detail struct {
	Product  Product `json:"product"`
	Quantity int     `json:"quantity"`
}

type Product struct {
	ID       int      `json:"id"`
	Unit     Unit     `json:"unit"`
	Category Category `json:"category"`
	Price    float64  `json:"price"`
	Name     string   `json:"name"`
	Code     string   `json:"code"`
}

type Unit struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}

type Category struct {
	ID       int     `json:"id"`
	OtherTax *string `json:"otherTax"`
	Name     string  `json:"name"`
	Code     string  `json:"code"`
}

type EmpresaConfig struct {
	EmpresaID        string
	Rut              string
	RutEnvia         string
	NombreFirma      string
	EmailFirma       string
	ClaveFirma       string
	NumeroResolucion string
	FechaResolucion  string
	Ambiente         string
	Certificado      []byte
	CertPassword     string
	CAF              []byte
	CAFDesde         int
	CAFHasta         int
}

func GetEmpresaConfig(conn *pgx.Conn, rut string, tipoDocumento string, folio int) (*EmpresaConfig, error) {
	var cfg EmpresaConfig
	// 1. Datos de la empresa
	err := conn.QueryRow(context.Background(), `
		SELECT id, rut, rutenvia, nombre_firma, email_firma, clave_firma, 
		       numero_resolucion, fecha_resolucion, ambiente
		FROM empresas
		WHERE rut = $1
	`, rut).Scan(
		&cfg.EmpresaID, &cfg.Rut, &cfg.RutEnvia, &cfg.NombreFirma, &cfg.EmailFirma, &cfg.ClaveFirma,
		&cfg.NumeroResolucion, &cfg.FechaResolucion, &cfg.Ambiente,
	)
	if err != nil {
		return nil, fmt.Errorf("empresa no encontrada: %w", err)
	}
	// 2. Certificado más reciente
	err = conn.QueryRow(context.Background(), `
		SELECT archivo, password
		FROM certificados
		WHERE empresa_id = $1
		ORDER BY fecha_vencimiento DESC NULLS LAST, created_at DESC
		LIMIT 1
	`, cfg.EmpresaID).Scan(&cfg.Certificado, &cfg.CertPassword)
	if err != nil {
		return nil, fmt.Errorf("certificado no encontrado: %w", err)
	}
	// 3. CAF vigente para el tipo de documento y folio
	err = conn.QueryRow(context.Background(), `
		SELECT archivo, desde, hasta
		FROM cafs
		WHERE empresa_id = $1 AND tipo_documento = $2 AND $3 BETWEEN desde AND hasta
		ORDER BY created_at DESC
		LIMIT 1
	`, cfg.EmpresaID, tipoDocumento, folio).Scan(&cfg.CAF, &cfg.CAFDesde, &cfg.CAFHasta)
	if err != nil {
		return nil, fmt.Errorf("CAF no encontrado o folio fuera de rango: %w", err)
	}
	return &cfg, nil
}

func main() {
	rutEmpresa := "76212889-6"
	tipoDocumento := "33"

	connStr := "postgresql://postgres.hptxgcuajsdupooptsax:rTFJE9FdUiJg2emR@aws-0-us-west-1.pooler.supabase.com:6543/postgres"
	config, err := pgx.ParseConfig(connStr)
	if err != nil {
		fmt.Println("Error en la configuración de conexión:", err)
		os.Exit(1)
	}
	config.PreferSimpleProtocol = true
	conn, err := pgx.ConnectConfig(context.Background(), config)
	if err != nil {
		fmt.Println("Error al conectar a la base de datos:", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	// 1. Determinar el siguiente folio disponible dentro del rango CAF
	var folio int
	err = conn.QueryRow(context.Background(), `
		SELECT COALESCE(MAX(CAST(numero_documento AS INTEGER)), 0) + 1
		FROM documentos
		WHERE empresa_id = (SELECT id FROM empresas WHERE rut = $1) AND tipo_documento = $2
	`, rutEmpresa, tipoDocumento).Scan(&folio)
	if err != nil {
		fmt.Println("Error al obtener el siguiente folio:", err)
		os.Exit(1)
	}

	// 2. Obtener la configuración completa de la empresa
	cfg, err := GetEmpresaConfig(conn, rutEmpresa, tipoDocumento, folio)
	if err != nil {
		fmt.Println("Error al obtener la configuración:", err)
		os.Exit(1)
	}

	// 3. Validar que el folio está dentro del rango CAF
	if folio < cfg.CAFDesde || folio > cfg.CAFHasta {
		fmt.Printf("Folio %d fuera del rango CAF (%d - %d)\n", folio, cfg.CAFDesde, cfg.CAFHasta)
		os.Exit(1)
	}

	// 4. Generar cliente y producto de ejemplo
	clienteRUT := "12345678-9"
	clienteNombre := "Cliente de Prueba S.A."
	clienteDireccion := "Av. Ejemplo 123, Santiago"
	clienteComuna := "Santiago"
	productoNombre := "Producto de Prueba"
	productoCantidad := 2
	productoPrecio := 5000.0
	neto := float64(productoCantidad) * productoPrecio
	iva := neto * 0.19
	total := neto + iva

	// 5. Insertar documento DTE
	var docID string
	err = conn.QueryRow(context.Background(), `
		INSERT INTO documentos (
			empresa_id, tipo_documento, numero_documento, fecha_emision, monto, estado
		) VALUES (
			$1, $2, $3, $4, $5, $6
		) RETURNING id
	`, cfg.EmpresaID, tipoDocumento, folio, time.Now(), total, "PENDIENTE").Scan(&docID)
	if err != nil {
		fmt.Println("Error al insertar el documento:", err)
		os.Exit(1)
	}

	// 6. Generar XML de ejemplo
	xmlDTE := fmt.Sprintf(`<DTE><Encabezado><IdDoc><TipoDTE>%s</TipoDTE><Folio>%d</Folio></IdDoc></Encabezado><Emisor><RUTEmisor>%s</RUTEmisor></Emisor><Receptor><RUTRecep>%s</RUTRecep><RznSocRecep>%s</RznSocRecep><DirRecep>%s</DirRecep><CmnaRecep>%s</CmnaRecep></Receptor><Detalle><NmbItem>%s</NmbItem><QtyItem>%d</QtyItem><PrcItem>%.2f</PrcItem></Detalle><Totales><MntNeto>%.2f</MntNeto><IVA>%.2f</IVA><MntTotal>%.2f</MntTotal></Totales></DTE>`,
		tipoDocumento, folio, cfg.Rut, clienteRUT, clienteNombre, clienteDireccion, clienteComuna, productoNombre, productoCantidad, productoPrecio, neto, iva, total)

	// 7. Guardar XML en xml_files
	var xmlID string
	err = conn.QueryRow(context.Background(), `
		INSERT INTO xml_files (
			documento_id, nombre_archivo, contenido
		) VALUES (
			$1, $2, $3
		) RETURNING id
	`, docID, fmt.Sprintf("dte_%d.xml", folio), xmlDTE).Scan(&xmlID)
	if err != nil {
		fmt.Println("Error al guardar el XML:", err)
		os.Exit(1)
	}

	fmt.Println("Factura emitida para empresa:", rutEmpresa)
	fmt.Println("Folio asignado:", folio)
	fmt.Println("Documento DTE insertado con ID:", docID)
	fmt.Println("XML guardado con ID:", xmlID, "y nombre:", fmt.Sprintf("dte_%d.xml", folio))
}
