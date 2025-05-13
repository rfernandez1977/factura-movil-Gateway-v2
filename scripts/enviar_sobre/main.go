package main

import (
	"bytes"
	"context"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/jackc/pgx/v4"
	"software.sslmate.com/src/go-pkcs12"
)

func main() {
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

	// 1. Obtener el último sobre
	var sobreID, empresaID, nombreArchivo string
	var contenido []byte
	err = conn.QueryRow(context.Background(), `
		SELECT id, empresa_id, nombre_archivo, contenido
		FROM sobres_xml
		ORDER BY fecha_creacion DESC
		LIMIT 1
	`).Scan(&sobreID, &empresaID, &nombreArchivo, &contenido)
	if err != nil {
		fmt.Println("No se pudo obtener el último sobre:", err)
		os.Exit(1)
	}
	fmt.Println("Sobre a enviar:", nombreArchivo, "(ID:", sobreID, ")")

	// 2. Obtener el certificado y clave de la empresa emisora (más reciente)
	var certData []byte
	var password string
	err = conn.QueryRow(context.Background(), `
		SELECT archivo, password
		FROM certificados
		WHERE empresa_id = $1
		ORDER BY fecha_vencimiento DESC NULLS LAST, created_at DESC
		LIMIT 1
	`, empresaID).Scan(&certData, &password)
	if err != nil {
		fmt.Println("No se pudo obtener el certificado de la empresa:", err)
		os.Exit(1)
	}

	// 3. Obtener el RUT de la empresa emisora
	var rutEmpresa string
	err = conn.QueryRow(context.Background(), `SELECT rut FROM empresas WHERE id = $1`, empresaID).Scan(&rutEmpresa)
	if err != nil {
		fmt.Println("No se pudo obtener el RUT de la empresa:", err)
		os.Exit(1)
	}

	// 4. Decodificar el PFX y extraer la clave privada y el certificado
	_, _, err = decodePFX(certData, password)
	if err != nil {
		fmt.Println("Error al decodificar el PFX:", err)
		os.Exit(1)
	}
	fmt.Println("Certificado y clave privada extraídos correctamente.")

	// 5. Firmar el sobre XML (aquí solo mostramos el paso, la firma real requiere goxmldsig o similar)
	// --- Aquí deberías usar una librería como github.com/russellhaering/goxmldsig para firmar el XML ---
	// Por ahora, enviamos el XML sin firmar para mostrar el flujo.
	// xmlFirmado := firmarXML(contenido, privateKey, cert)
	xmlFirmado := contenido // TODO: reemplazar por el XML firmado

	// 6. Enviar el sobre firmado al SII (certificación)
	endpoint := "https://maullin.sii.cl/cgi_dte/UPL/DTEUpload"
	resp, err := http.Post(endpoint, "application/xml", bytes.NewReader(xmlFirmado))
	if err != nil {
		fmt.Println("Error al enviar el sobre al SII:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	fmt.Println("Respuesta del SII:")
	fmt.Println(string(body))
}

// decodePFX decodifica un archivo PFX y retorna la clave privada y el certificado
func decodePFX(pfxData []byte, password string) (interface{}, *x509.Certificate, error) {
	privateKey, cert, err := pkcs12.Decode(pfxData, password)
	if err != nil {
		return nil, nil, err
	}
	return privateKey, cert, nil
}
