package verify

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4"
)

func VerifySupabase() error {
	// Cadena de conexión directa (reemplaza TU_PASSWORD por la real)
	connStr := "postgresql://postgres.hptxgcuajsdupooptsax:rTFJE9FdUiJg2emR@aws-0-us-west-1.pooler.supabase.com:6543/postgres"

	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		return fmt.Errorf("No se pudo conectar a la base de datos: %v", err)
	}
	defer conn.Close(context.Background())

	fmt.Println("\n✅ Conexión directa a PostgreSQL establecida correctamente")

	tableOrder := []string{"empresas", "documentos", "certificados", "sesiones", "cafs", "xml_files", "sobres_xml"}
	tables := map[string]string{
		"empresas": `
			CREATE TABLE IF NOT EXISTS empresas (
				id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
				nombre VARCHAR(255) NOT NULL,
				rut VARCHAR(20) NOT NULL UNIQUE,
				direccion TEXT,
				telefono VARCHAR(50),
				email VARCHAR(255),
				rut_firma VARCHAR(20) NOT NULL,
				nombre_firma VARCHAR(255) NOT NULL,
				clave_firma VARCHAR(255) NOT NULL,
				created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
				updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
			);`,
		"documentos": `
			CREATE TABLE IF NOT EXISTS documentos (
				id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
				empresa_id UUID REFERENCES empresas(id),
				tipo_documento VARCHAR(50) NOT NULL,
				numero_documento VARCHAR(50) NOT NULL,
				fecha_emision TIMESTAMP WITH TIME ZONE NOT NULL,
				monto DECIMAL(15,2),
				estado VARCHAR(50) NOT NULL,
				created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
				updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
			);`,
		"certificados": `
			CREATE TABLE IF NOT EXISTS certificados (
				id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
				empresa_id UUID REFERENCES empresas(id),
				nombre VARCHAR(255) NOT NULL,
				archivo BYTEA NOT NULL,
				password VARCHAR(255),
				fecha_vencimiento TIMESTAMP WITH TIME ZONE,
				created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
				updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
			);`,
		"sesiones": `
			CREATE TABLE IF NOT EXISTS sesiones (
				id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
				empresa_id UUID REFERENCES empresas(id),
				token VARCHAR(255) NOT NULL,
				fecha_inicio TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
				fecha_fin TIMESTAMP WITH TIME ZONE,
				estado VARCHAR(50) NOT NULL,
				created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
				updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
			);`,
		"cafs": `
			CREATE TABLE IF NOT EXISTS cafs (
				id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
				empresa_id UUID REFERENCES empresas(id),
				tipo_documento VARCHAR(50) NOT NULL,
				desde INT NOT NULL,
				hasta INT NOT NULL,
				archivo BYTEA NOT NULL,
				fecha_vencimiento TIMESTAMP WITH TIME ZONE,
				created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
				updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
			);`,
		"xml_files": `
			CREATE TABLE IF NOT EXISTS xml_files (
				id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
				documento_id UUID REFERENCES documentos(id),
				nombre_archivo VARCHAR(255) NOT NULL,
				contenido XML NOT NULL,
				created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
				updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
			);`,
		"sobres_xml": `
			CREATE TABLE IF NOT EXISTS sobres_xml (
				id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
				empresa_id UUID REFERENCES empresas(id),
				nombre_archivo VARCHAR(100) NOT NULL,
				contenido BYTEA NOT NULL,
				fecha_creacion TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
			);
		`,
	}

	fmt.Println("\nVerificando y creando tablas si es necesario...")

	// Eliminar tabla empresas si existe para recrearla con la nueva estructura
	_, err = conn.Exec(context.Background(), "DROP TABLE IF EXISTS empresas CASCADE;")
	if err != nil {
		fmt.Printf("❌ Error al eliminar tabla empresas: %v\n", err)
	}

	for _, tableName := range tableOrder {
		createSQL := tables[tableName]
		_, err := conn.Exec(context.Background(), createSQL)
		if err != nil {
			fmt.Printf("❌ Error al crear/verificar tabla %s: %v\n", tableName, err)
			continue
		}
		fmt.Printf("✅ Tabla %s verificada/creada correctamente\n", tableName)

		// Contar registros
		var count int
		err = conn.QueryRow(context.Background(), fmt.Sprintf("SELECT COUNT(*) FROM %s", tableName)).Scan(&count)
		if err != nil {
			fmt.Printf("   ⚠️  No se pudo contar registros en %s: %v\n", tableName, err)
		} else {
			fmt.Printf("   Registros en %s: %d\n", tableName, count)
		}
	}

	// --- Inserción de la empresa ---
	fmt.Println("\nInsertando empresa inicial...")
	_, err = conn.Exec(context.Background(), `
		INSERT INTO empresas (nombre, rut, direccion, telefono, email, rut_firma, nombre_firma, clave_firma)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (rut) DO NOTHING
	`,
		"FACTURA MOVIL SPA",
		"76212889-6",
		"Vicuña Mackenna 9705",
		"",
		"rfernandez@facturamovil.cl",
		"76212889-6",                 // RUT de la firma electrónica
		"Rodrigo Fernandez Calderon", // Nombre del titular de la firma
		"",                           // Clave de la firma (debe ser proporcionada)
	)
	if err != nil {
		fmt.Printf("❌ Error insertando empresa: %v\n", err)
	} else {
		fmt.Println("✅ Empresa insertada correctamente")
	}

	// --- Inserción del CAF ---
	fmt.Println("Insertando CAF...")
	cafData, err := os.ReadFile("../../caf_test/33-cert.xml")
	if err != nil {
		fmt.Printf("❌ Error leyendo archivo CAF: %v\n", err)
	} else {
		_, err = conn.Exec(context.Background(), `
			INSERT INTO cafs (empresa_id, tipo_documento, desde, hasta, archivo)
			VALUES ((SELECT id FROM empresas WHERE rut = $1), $2, $3, $4, $5)
		`,
			"76212889-6",
			"33",
			1,   // desde (ajusta si es necesario)
			100, // hasta (ajusta si es necesario)
			cafData,
		)
		if err != nil {
			fmt.Printf("❌ Error insertando CAF: %v\n", err)
		} else {
			fmt.Println("✅ CAF insertado correctamente")
		}
	}

	// --- Inserción del certificado ---
	fmt.Println("Insertando certificado...")
	certData, err := os.ReadFile("../../firma_test/1/firma.pfx")
	if err != nil {
		fmt.Printf("❌ Error leyendo archivo de firma: %v\n", err)
	} else {
		_, err = conn.Exec(context.Background(), `
			INSERT INTO certificados (empresa_id, nombre, archivo, password)
			VALUES ((SELECT id FROM empresas WHERE rut = $1), $2, $3, $4)
		`,
			"76212889-6",
			"Rodrigo Fernandez Calderon",
			certData,
			"", // password si aplica
		)
		if err != nil {
			fmt.Printf("❌ Error insertando certificado: %v\n", err)
		} else {
			fmt.Println("✅ Certificado insertado correctamente")
		}
	}

	return nil
}
