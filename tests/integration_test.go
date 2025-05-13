package tests

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/assert"
)

// TestData contiene los datos de prueba
type TestData struct {
	Empresa struct {
		Nombre      string
		Rut         string
		Direccion   string
		Telefono    string
		Email       string
		RutFirma    string
		NombreFirma string
		ClaveFirma  string
	}
	CAF struct {
		TipoDocumento string
		Desde         int
		Hasta         int
	}
}

func setupTestData(testName string) TestData {
	// Generar un sufijo único corto basado en el nombre de la prueba y los últimos 6 dígitos del timestamp
	timestamp := time.Now().UnixNano() % 1000000
	suffix := fmt.Sprintf("%s%06d", testName, timestamp)
	// El RUT y el email no deben exceder 20 caracteres
	rut := fmt.Sprintf("76212889-%s", suffix[:min(8, len(suffix))])      // máximo 20 caracteres
	email := fmt.Sprintf("prueba%s@em.cl", suffix[:min(8, len(suffix))]) // email corto

	return TestData{
		Empresa: struct {
			Nombre      string
			Rut         string
			Direccion   string
			Telefono    string
			Email       string
			RutFirma    string
			NombreFirma string
			ClaveFirma  string
		}{
			Nombre:      fmt.Sprintf("EMPRESA DE PRUEBA %s", suffix),
			Rut:         rut,
			Direccion:   "Calle de Prueba 123",
			Telefono:    "+56912345678",
			Email:       email,
			RutFirma:    rut,
			NombreFirma: "Rodrigo Fernandez Calderon",
			ClaveFirma:  "123456",
		},
		CAF: struct {
			TipoDocumento string
			Desde         int
			Hasta         int
		}{
			TipoDocumento: "33",
			Desde:         1,
			Hasta:         100,
		},
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func getDBConnection() (*pgx.Conn, error) {
	connStr := "postgresql://postgres.hptxgcuajsdupooptsax:rTFJE9FdUiJg2emR@aws-0-us-west-1.pooler.supabase.com:6543/postgres"
	config, err := pgx.ParseConfig(connStr)
	if err != nil {
		return nil, fmt.Errorf("error parsing connection string: %v", err)
	}

	// Deshabilitar declaraciones preparadas
	config.PreferSimpleProtocol = true

	return pgx.ConnectConfig(context.Background(), config)
}

func cleanupTestData(conn *pgx.Conn, testData TestData) error {
	// Eliminar datos en orden inverso a las dependencias
	_, err := conn.Exec(context.Background(), fmt.Sprintf(`
		DELETE FROM xml_files 
		WHERE documento_id IN (
			SELECT id FROM documentos 
			WHERE empresa_id IN (
				SELECT id FROM empresas 
				WHERE rut = '%s'
			)
		)`, testData.Empresa.Rut))
	if err != nil {
		return fmt.Errorf("error limpiando xml_files: %v", err)
	}

	_, err = conn.Exec(context.Background(), fmt.Sprintf(`
		DELETE FROM documentos 
		WHERE empresa_id IN (
			SELECT id FROM empresas 
			WHERE rut = '%s'
		)`, testData.Empresa.Rut))
	if err != nil {
		return fmt.Errorf("error limpiando documentos: %v", err)
	}

	_, err = conn.Exec(context.Background(), fmt.Sprintf(`
		DELETE FROM cafs 
		WHERE empresa_id IN (
			SELECT id FROM empresas 
			WHERE rut = '%s'
		)`, testData.Empresa.Rut))
	if err != nil {
		return fmt.Errorf("error limpiando cafs: %v", err)
	}

	_, err = conn.Exec(context.Background(), fmt.Sprintf(`
		DELETE FROM certificados 
		WHERE empresa_id IN (
			SELECT id FROM empresas 
			WHERE rut = '%s'
		)`, testData.Empresa.Rut))
	if err != nil {
		return fmt.Errorf("error limpiando certificados: %v", err)
	}

	_, err = conn.Exec(context.Background(), fmt.Sprintf(`
		DELETE FROM sesiones 
		WHERE empresa_id IN (
			SELECT id FROM empresas 
			WHERE rut = '%s'
		)`, testData.Empresa.Rut))
	if err != nil {
		return fmt.Errorf("error limpiando sesiones: %v", err)
	}

	_, err = conn.Exec(context.Background(), fmt.Sprintf(`
		DELETE FROM empresas 
		WHERE rut = '%s'`, testData.Empresa.Rut))
	if err != nil {
		return fmt.Errorf("error limpiando empresas: %v", err)
	}

	return nil
}

func TestSupabaseIntegration(t *testing.T) {
	// 1. Prueba de conexión
	t.Run("Conexión a Supabase", func(t *testing.T) {
		conn, err := getDBConnection()
		assert.NoError(t, err, "Debería poder conectarse a Supabase")
		defer conn.Close(context.Background())

		// Verificar que la conexión está activa
		err = conn.Ping(context.Background())
		assert.NoError(t, err, "La conexión debería estar activa")
	})

	// 2. Prueba de CRUD de Empresa
	t.Run("CRUD de Empresa", func(t *testing.T) {
		testData := setupTestData("CRUD-Empresa")
		conn, err := getDBConnection()
		assert.NoError(t, err)
		defer conn.Close(context.Background())
		defer cleanupTestData(conn, testData)

		// Insertar empresa
		var empresaID string
		err = conn.QueryRow(context.Background(), fmt.Sprintf(`
			INSERT INTO empresas (
				nombre, rut, direccion, telefono, email, 
				rut_firma, nombre_firma, clave_firma
			) VALUES (
				'%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s'
			)
			RETURNING id
		`,
			testData.Empresa.Nombre,
			testData.Empresa.Rut,
			testData.Empresa.Direccion,
			testData.Empresa.Telefono,
			testData.Empresa.Email,
			testData.Empresa.RutFirma,
			testData.Empresa.NombreFirma,
			testData.Empresa.ClaveFirma,
		)).Scan(&empresaID)
		assert.NoError(t, err, "Debería poder insertar la empresa")
		assert.NotEmpty(t, empresaID, "Debería obtener un ID de empresa")

		// Leer empresa
		var nombre, rut string
		err = conn.QueryRow(context.Background(), fmt.Sprintf(`
			SELECT nombre, rut FROM empresas WHERE id = '%s'
		`, empresaID)).Scan(&nombre, &rut)
		assert.NoError(t, err, "Debería poder leer la empresa")
		assert.Equal(t, testData.Empresa.Nombre, nombre, "El nombre debería coincidir")
		assert.Equal(t, testData.Empresa.Rut, rut, "El RUT debería coincidir")

		// Actualizar empresa
		_, err = conn.Exec(context.Background(), fmt.Sprintf(`
			UPDATE empresas 
			SET nombre = 'EMPRESA ACTUALIZADA', telefono = '+56987654321'
			WHERE id = '%s'
		`, empresaID))
		assert.NoError(t, err, "Debería poder actualizar la empresa")

		// Verificar actualización
		err = conn.QueryRow(context.Background(), fmt.Sprintf(`
			SELECT nombre FROM empresas WHERE id = '%s'
		`, empresaID)).Scan(&nombre)
		assert.NoError(t, err, "Debería poder leer la empresa actualizada")
		assert.Equal(t, "EMPRESA ACTUALIZADA", nombre, "El nombre debería estar actualizado")

		// Eliminar empresa
		_, err = conn.Exec(context.Background(), fmt.Sprintf(`
			DELETE FROM empresas WHERE id = '%s'
		`, empresaID))
		assert.NoError(t, err, "Debería poder eliminar la empresa")

		// Verificar eliminación
		err = conn.QueryRow(context.Background(), fmt.Sprintf(`
			SELECT id FROM empresas WHERE id = '%s'
		`, empresaID)).Scan(&empresaID)
		assert.Error(t, err, "No debería encontrar la empresa eliminada")
	})

	// 3. Prueba de CAF
	t.Run("Manejo de CAF", func(t *testing.T) {
		testData := setupTestData("CAF")
		conn, err := getDBConnection()
		assert.NoError(t, err)
		defer conn.Close(context.Background())
		defer cleanupTestData(conn, testData)

		// Insertar empresa primero (necesario para la FK)
		var empresaID string
		err = conn.QueryRow(context.Background(), fmt.Sprintf(`
			INSERT INTO empresas (
				nombre, rut, direccion, telefono, email, 
				rut_firma, nombre_firma, clave_firma
			) VALUES (
				'%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s'
			)
			RETURNING id
		`,
			testData.Empresa.Nombre,
			testData.Empresa.Rut,
			testData.Empresa.Direccion,
			testData.Empresa.Telefono,
			testData.Empresa.Email,
			testData.Empresa.RutFirma,
			testData.Empresa.NombreFirma,
			testData.Empresa.ClaveFirma,
		)).Scan(&empresaID)
		assert.NoError(t, err)

		// Insertar CAF
		var cafID string
		err = conn.QueryRow(context.Background(), fmt.Sprintf(`
			INSERT INTO cafs (
				empresa_id, tipo_documento, desde, hasta, archivo
			) VALUES (
				'%s', '%s', %d, %d, 'contenido de prueba del CAF'
			)
			RETURNING id
		`,
			empresaID,
			testData.CAF.TipoDocumento,
			testData.CAF.Desde,
			testData.CAF.Hasta,
		)).Scan(&cafID)
		assert.NoError(t, err, "Debería poder insertar el CAF")
		assert.NotEmpty(t, cafID, "Debería obtener un ID de CAF")

		// Verificar CAF
		var tipoDoc string
		var desde, hasta int
		err = conn.QueryRow(context.Background(), fmt.Sprintf(`
			SELECT tipo_documento, desde, hasta 
			FROM cafs WHERE id = '%s'
		`, cafID)).Scan(&tipoDoc, &desde, &hasta)
		assert.NoError(t, err, "Debería poder leer el CAF")
		assert.Equal(t, testData.CAF.TipoDocumento, tipoDoc)
		assert.Equal(t, testData.CAF.Desde, desde)
		assert.Equal(t, testData.CAF.Hasta, hasta)
	})

	// 4. Prueba de Certificados
	t.Run("Manejo de Certificados", func(t *testing.T) {
		testData := setupTestData("Certificados")
		conn, err := getDBConnection()
		assert.NoError(t, err)
		defer conn.Close(context.Background())
		defer cleanupTestData(conn, testData)

		// Insertar empresa primero
		var empresaID string
		err = conn.QueryRow(context.Background(), fmt.Sprintf(`
			INSERT INTO empresas (
				nombre, rut, direccion, telefono, email, 
				rut_firma, nombre_firma, clave_firma
			) VALUES (
				'%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s'
			)
			RETURNING id
		`,
			testData.Empresa.Nombre,
			testData.Empresa.Rut,
			testData.Empresa.Direccion,
			testData.Empresa.Telefono,
			testData.Empresa.Email,
			testData.Empresa.RutFirma,
			testData.Empresa.NombreFirma,
			testData.Empresa.ClaveFirma,
		)).Scan(&empresaID)
		assert.NoError(t, err)

		// Insertar certificado
		var certID string
		err = conn.QueryRow(context.Background(), fmt.Sprintf(`
			INSERT INTO certificados (
				empresa_id, nombre, archivo, password
			) VALUES (
				'%s', 'Certificado de Prueba', 'contenido de prueba del certificado', 'password123'
			)
			RETURNING id
		`, empresaID)).Scan(&certID)
		assert.NoError(t, err, "Debería poder insertar el certificado")
		assert.NotEmpty(t, certID, "Debería obtener un ID de certificado")

		// Verificar certificado
		var nombre string
		err = conn.QueryRow(context.Background(), fmt.Sprintf(`
			SELECT nombre FROM certificados WHERE id = '%s'
		`, certID)).Scan(&nombre)
		assert.NoError(t, err, "Debería poder leer el certificado")
		assert.Equal(t, "Certificado de Prueba", nombre)
	})

	// 5. Prueba de Documentos
	t.Run("Manejo de Documentos", func(t *testing.T) {
		testData := setupTestData("Documentos")
		conn, err := getDBConnection()
		assert.NoError(t, err)
		defer conn.Close(context.Background())
		defer cleanupTestData(conn, testData)

		// Insertar empresa primero
		var empresaID string
		err = conn.QueryRow(context.Background(), fmt.Sprintf(`
			INSERT INTO empresas (
				nombre, rut, direccion, telefono, email, 
				rut_firma, nombre_firma, clave_firma
			) VALUES (
				'%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s'
			)
			RETURNING id
		`,
			testData.Empresa.Nombre,
			testData.Empresa.Rut,
			testData.Empresa.Direccion,
			testData.Empresa.Telefono,
			testData.Empresa.Email,
			testData.Empresa.RutFirma,
			testData.Empresa.NombreFirma,
			testData.Empresa.ClaveFirma,
		)).Scan(&empresaID)
		assert.NoError(t, err)

		// Insertar documento
		var docID string
		err = conn.QueryRow(context.Background(), fmt.Sprintf(`
			INSERT INTO documentos (
				empresa_id, tipo_documento, numero_documento,
				fecha_emision, monto, estado
			) VALUES (
				'%s', '33', '1', NOW(), 1000.00, 'PENDIENTE'
			)
			RETURNING id
		`, empresaID)).Scan(&docID)
		assert.NoError(t, err, "Debería poder insertar el documento")
		assert.NotEmpty(t, docID, "Debería obtener un ID de documento")

		// Insertar XML asociado
		var xmlID string
		err = conn.QueryRow(context.Background(), fmt.Sprintf(`
			INSERT INTO xml_files (
				documento_id, nombre_archivo, contenido
			) VALUES (
				'%s', 'documento_1.xml', '<xml>contenido de prueba</xml>'
			)
			RETURNING id
		`, docID)).Scan(&xmlID)
		assert.NoError(t, err, "Debería poder insertar el XML")
		assert.NotEmpty(t, xmlID, "Debería obtener un ID de XML")

		// Verificar documento y XML
		var tipoDoc, estado string
		var monto float64
		err = conn.QueryRow(context.Background(), fmt.Sprintf(`
			SELECT d.tipo_documento, d.estado, d.monto
			FROM documentos d
			WHERE d.id = '%s'
		`, docID)).Scan(&tipoDoc, &estado, &monto)
		assert.NoError(t, err, "Debería poder leer el documento")
		assert.Equal(t, "33", tipoDoc)
		assert.Equal(t, "PENDIENTE", estado)
		assert.Equal(t, 1000.00, monto)
	})

	// 6. Prueba de creación y guardado de XML DTE y sobres
	t.Run("Creación y guardado de XML DTE y sobres", func(t *testing.T) {
		testData := setupTestData("XMLDTE")
		conn, err := getDBConnection()
		assert.NoError(t, err)
		defer conn.Close(context.Background())
		defer cleanupTestData(conn, testData)

		// Insertar empresa primero
		var empresaID string
		err = conn.QueryRow(context.Background(), fmt.Sprintf(`
			INSERT INTO empresas (
				nombre, rut, direccion, telefono, email, 
				rut_firma, nombre_firma, clave_firma
			) VALUES (
				'%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s'
			)
			RETURNING id
		`,
			testData.Empresa.Nombre,
			testData.Empresa.Rut,
			testData.Empresa.Direccion,
			testData.Empresa.Telefono,
			testData.Empresa.Email,
			testData.Empresa.RutFirma,
			testData.Empresa.NombreFirma,
			testData.Empresa.ClaveFirma,
		)).Scan(&empresaID)
		assert.NoError(t, err)

		// Insertar documento DTE
		var docID string
		err = conn.QueryRow(context.Background(), fmt.Sprintf(`
			INSERT INTO documentos (
				empresa_id, tipo_documento, numero_documento,
				fecha_emision, monto, estado
			) VALUES (
				'%s', '33', '1001', NOW(), 38859.07, 'PENDIENTE'
			)
			RETURNING id
		`, empresaID)).Scan(&docID)
		assert.NoError(t, err)
		assert.NotEmpty(t, docID)

		// XML DTE de ejemplo
		xmlDTE := `<DTE><Encabezado><IdDoc><TipoDTE>33</TipoDTE><Folio>1001</Folio></IdDoc></Encabezado><Receptor><RUTRecep>76071974-9</RUTRecep><RznSocRecep>AGRICOLA LOS DOS LIMITADA</RznSocRecep><DirRecep>PARCELA 16 s/n, Villa/Pob. P.P.LOS CRISTALES</DirRecep><CmnaRecep>Curico</CmnaRecep></Receptor><Detalle><NmbItem>Servicio Mensual Plan Copihue</NmbItem><QtyItem>1</QtyItem><PrcItem>38859.07</PrcItem></Detalle></DTE>`

		// Guardar XML DTE en xml_files
		var xmlID string
		err = conn.QueryRow(context.Background(), fmt.Sprintf(`
			INSERT INTO xml_files (
				documento_id, nombre_archivo, contenido
			) VALUES (
				'%s', 'dte_1001.xml', '%s'
			)
			RETURNING id
		`, docID, xmlDTE)).Scan(&xmlID)
		assert.NoError(t, err, "Debería poder guardar el XML DTE")
		assert.NotEmpty(t, xmlID)

		// XML Sobre de ejemplo (agrupa 1 DTE)
		xmlSobre := `<EnvioDTE><SetDTE><Caratula><RutEmisor>76212889-6</RutEmisor><RutEnvia>76212889-6</RutEnvia><RutReceptor>60803000-K</RutReceptor></Caratula>` + xmlDTE + `</SetDTE></EnvioDTE>`

		// Guardar sobre en sobres_xml
		var sobreID string
		err = conn.QueryRow(context.Background(), fmt.Sprintf(`
			INSERT INTO sobres_xml (
				empresa_id, nombre_archivo, contenido
			) VALUES (
				'%s', 'sobre_envio_1001.xml', $1
			)
			RETURNING id
		`, empresaID), []byte(xmlSobre)).Scan(&sobreID)
		assert.NoError(t, err, "Debería poder guardar el sobre XML")
		assert.NotEmpty(t, sobreID)

		// Verificar que ambos existen
		var count int
		err = conn.QueryRow(context.Background(), fmt.Sprintf(`SELECT COUNT(*) FROM xml_files WHERE id = '%s'`, xmlID)).Scan(&count)
		assert.NoError(t, err)
		assert.Equal(t, 1, count, "El XML DTE debe existir")

		err = conn.QueryRow(context.Background(), fmt.Sprintf(`SELECT COUNT(*) FROM sobres_xml WHERE id = '%s'`, sobreID)).Scan(&count)
		assert.NoError(t, err)
		assert.Equal(t, 1, count, "El sobre XML debe existir")
	})
}
