package main

import (
	"context"
	"fmt"
	"log"

	"github.com/fmgo/supabase"
)

func main() {
	fmt.Println("=== FMgo - Ejemplo de Operaciones CRUD con Supabase ===")

	// Inicializar cliente de Supabase
	client, err := supabase.InitClientWithConfig("../config.json")
	if err != nil {
		log.Fatalf("Error inicializando cliente Supabase: %v", err)
	}

	// Crear contexto
	ctx := context.Background()

	// Listar las tablas disponibles
	tables, err := supabase.ListAllTables(ctx, client)
	if err != nil {
		log.Printf("Error listando tablas: %v", err)
	} else {
		fmt.Println("\nTablas disponibles:")
		for _, table := range tables {
			count, err := supabase.GetTableCount(ctx, client, table)
			if err != nil {
				fmt.Printf("- %s (error obteniendo conteo: %v)\n", table, err)
			} else {
				fmt.Printf("- %s (%d registros)\n", table, count)
			}
		}
	}

	// ****************************************
	// Operaciones CRUD con la tabla 'empresas'
	// ****************************************
	fmt.Println("\n=== Operaciones CRUD con tabla 'empresas' ===")

	// 1. CREAR: Insertar una nueva empresa
	fmt.Println("\n1. Creando una nueva empresa...")
	nuevaEmpresa := map[string]interface{}{
		"nombre":       "EMPRESA DE EJEMPLO CRUD",
		"rut":          "99.999.999-9",
		"direccion":    "Avenida Ejemplo 123",
		"telefono":     "+56912345678",
		"email":        "ejemplo@crud.cl",
		"rut_firma":    "99.999.999-9",
		"nombre_firma": "Ejemplo CRUD",
		"clave_firma":  "clave123",
	}

	empresaCreada, err := supabase.InsertRecord(ctx, client, "empresas", nuevaEmpresa)
	if err != nil {
		log.Printf("Error insertando empresa: %v", err)
	} else {
		fmt.Printf("Empresa creada exitosamente. ID: %s\n", empresaCreada["id"])

		// Guardar el ID para usar en otras operaciones
		empresaID := fmt.Sprintf("%v", empresaCreada["id"])

		// 2. LEER: Obtener la empresa recién creada
		fmt.Println("\n2. Leyendo la empresa creada...")
		empresa, err := supabase.GetRecordByID(ctx, client, "empresas", empresaID)
		if err != nil {
			log.Printf("Error obteniendo empresa: %v", err)
		} else {
			fmt.Printf("Empresa encontrada:\n")
			fmt.Printf("  ID: %s\n", empresa["id"])
			fmt.Printf("  Nombre: %s\n", empresa["nombre"])
			fmt.Printf("  RUT: %s\n", empresa["rut"])
			fmt.Printf("  Email: %s\n", empresa["email"])
		}

		// 3. ACTUALIZAR: Modificar la empresa
		fmt.Println("\n3. Actualizando la empresa...")
		actualizaciones := map[string]interface{}{
			"direccion": "Avenida Ejemplo Actualizada 456",
			"telefono":  "+56987654321",
			"email":     "actualizado@crud.cl",
		}

		empresaActualizada, err := supabase.UpdateRecord(ctx, client, "empresas", empresaID, actualizaciones)
		if err != nil {
			log.Printf("Error actualizando empresa: %v", err)
		} else {
			fmt.Printf("Empresa actualizada exitosamente:\n")
			fmt.Printf("  ID: %s\n", empresaActualizada["id"])
			fmt.Printf("  Nombre: %s\n", empresaActualizada["nombre"])
			fmt.Printf("  Dirección: %s\n", empresaActualizada["direccion"])
			fmt.Printf("  Teléfono: %s\n", empresaActualizada["telefono"])
			fmt.Printf("  Email: %s\n", empresaActualizada["email"])
		}

		// 4. BUSCAR: Listar empresas por filtro
		fmt.Println("\n4. Buscando empresas con filtros...")
		filtros := map[string]string{
			"rut": "99.999.999-9",
		}

		empresasFiltradas, err := supabase.QueryRecords(ctx, client, "empresas", filtros, 10)
		if err != nil {
			log.Printf("Error buscando empresas: %v", err)
		} else {
			fmt.Printf("Empresas encontradas: %d\n", len(empresasFiltradas))
			for i, e := range empresasFiltradas {
				fmt.Printf("%d. %s (%s)\n", i+1, e["nombre"], e["rut"])
			}
		}

		// 5. ELIMINAR: Eliminar la empresa creada
		fmt.Println("\n5. Eliminando la empresa creada...")
		fmt.Printf("¿Quieres eliminar la empresa con ID %s? (Esta operación se ha comentado por seguridad)\n", empresaID)

		// Comentado para evitar eliminaciones accidentales en el ejemplo
		// Descomentar para probar la eliminación
		/*
			err = supabase.DeleteRecord(ctx, client, "empresas", empresaID)
			if err != nil {
				log.Printf("Error eliminando empresa: %v", err)
			} else {
				fmt.Println("Empresa eliminada exitosamente")
			}
		*/
	}

	// Mostrar cómo listar datos de una tabla
	fmt.Println("\n=== Listado de empresas (limitado a 5) ===")
	empresas, err := supabase.ListTableData(ctx, client, "empresas", 5)
	if err != nil {
		log.Printf("Error listando empresas: %v", err)
	} else {
		for i, empresa := range empresas {
			fmt.Printf("%d. %s (%s)\n", i+1, empresa["nombre"], empresa["rut"])
		}
	}

	fmt.Println("\nEjemplo completado.")
}
