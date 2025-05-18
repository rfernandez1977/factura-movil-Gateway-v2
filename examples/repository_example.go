package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"FMgo/repository"
)

func main() {
	fmt.Println("=== FMgo - Ejemplo de Patrón Repositorio con Supabase ===")

	// Inicializar el repositorio
	repo, err := repository.InitializeRepository("../config.json")
	if err != nil {
		log.Fatalf("Error inicializando repositorio: %v", err)
	}

	// Crear contexto
	ctx := context.Background()

	// Listar empresas existentes
	fmt.Println("\n=== Listando empresas existentes ===")
	empresas, err := repo.ListEmpresas(ctx, 5)
	if err != nil {
		log.Printf("Error listando empresas: %v", err)
	} else {
		fmt.Printf("Total de empresas: %d\n", len(empresas))
		for i, empresa := range empresas {
			fmt.Printf("%d. %s (%s)\n", i+1, empresa.Nombre, empresa.RUT)
		}
	}

	// Crear una nueva empresa
	fmt.Println("\n=== Creando nueva empresa ===")
	nuevaEmpresa := &repository.Empresa{
		Nombre:      "EMPRESA REPOSITORIO EJEMPLO",
		RUT:         "88.888.888-8",
		Direccion:   "Calle Repositorio 123",
		Telefono:    "+56987654321",
		Email:       "repositorio@ejemplo.cl",
		RUTFirma:    "88.888.888-8",
		NombreFirma: "Ejemplo Repositorio",
		ClaveFirma:  "clave123",
	}

	empresaCreada, err := repo.CreateEmpresa(ctx, nuevaEmpresa)
	if err != nil {
		log.Printf("Error creando empresa: %v", err)
	} else {
		fmt.Printf("Empresa creada exitosamente. ID: %s\n", empresaCreada.ID)

		// Obtener empresa por ID
		fmt.Println("\n=== Obteniendo empresa por ID ===")
		empresa, err := repo.GetEmpresaByID(ctx, empresaCreada.ID)
		if err != nil {
			log.Printf("Error obteniendo empresa: %v", err)
		} else {
			fmt.Printf("Empresa encontrada:\n")
			fmt.Printf("  ID: %s\n", empresa.ID)
			fmt.Printf("  Nombre: %s\n", empresa.Nombre)
			fmt.Printf("  RUT: %s\n", empresa.RUT)
			fmt.Printf("  Email: %s\n", empresa.Email)
		}

		// Actualizar empresa
		fmt.Println("\n=== Actualizando empresa ===")
		empresaCreada.Direccion = "Calle Repositorio Actualizada 456"
		empresaCreada.Telefono = "+56912345678"
		empresaCreada.Email = "actualizado@repositorio.cl"

		empresaActualizada, err := repo.UpdateEmpresa(ctx, empresaCreada)
		if err != nil {
			log.Printf("Error actualizando empresa: %v", err)
		} else {
			fmt.Printf("Empresa actualizada exitosamente:\n")
			fmt.Printf("  ID: %s\n", empresaActualizada.ID)
			fmt.Printf("  Nombre: %s\n", empresaActualizada.Nombre)
			fmt.Printf("  Dirección: %s\n", empresaActualizada.Direccion)
			fmt.Printf("  Teléfono: %s\n", empresaActualizada.Telefono)
			fmt.Printf("  Email: %s\n", empresaActualizada.Email)
		}

		// Crear documento para la empresa
		fmt.Println("\n=== Creando documento para la empresa ===")
		now := time.Now().Format(time.RFC3339)
		nuevoDocumento := &repository.Documento{
			EmpresaID:       empresaCreada.ID,
			TipoDocumento:   "33",
			NumeroDocumento: "1001",
			FechaEmision:    now,
			Monto:           1000.0,
			Estado:          "PENDIENTE",
		}

		documentoCreado, err := repo.CreateDocumento(ctx, nuevoDocumento)
		if err != nil {
			log.Printf("Error creando documento: %v", err)
		} else {
			fmt.Printf("Documento creado exitosamente. ID: %s\n", documentoCreado.ID)

			// Obtener documento por ID
			fmt.Println("\n=== Obteniendo documento por ID ===")
			documento, err := repo.GetDocumentoByID(ctx, documentoCreado.ID)
			if err != nil {
				log.Printf("Error obteniendo documento: %v", err)
			} else {
				fmt.Printf("Documento encontrado:\n")
				fmt.Printf("  ID: %s\n", documento.ID)
				fmt.Printf("  Tipo: %s\n", documento.TipoDocumento)
				fmt.Printf("  Número: %s\n", documento.NumeroDocumento)
				fmt.Printf("  Empresa ID: %s\n", documento.EmpresaID)
				fmt.Printf("  Estado: %s\n", documento.Estado)
				fmt.Printf("  Monto: %.2f\n", documento.Monto)
			}

			// Actualizar estado del documento
			fmt.Println("\n=== Actualizando estado del documento ===")
			err = repo.UpdateDocumentoEstado(ctx, documentoCreado.ID, "ACEPTADO")
			if err != nil {
				log.Printf("Error actualizando estado del documento: %v", err)
			} else {
				fmt.Println("Estado del documento actualizado correctamente")

				// Verificar cambio de estado
				docActualizado, _ := repo.GetDocumentoByID(ctx, documentoCreado.ID)
				fmt.Printf("Estado actualizado: %s\n", docActualizado.Estado)
			}

			// Listar documentos de la empresa
			fmt.Println("\n=== Listando documentos de la empresa ===")
			documentos, err := repo.ListDocumentosByEmpresa(ctx, empresaCreada.ID, 10)
			if err != nil {
				log.Printf("Error listando documentos: %v", err)
			} else {
				fmt.Printf("Total de documentos: %d\n", len(documentos))
				for i, doc := range documentos {
					fmt.Printf("%d. Tipo: %s, Número: %s, Estado: %s\n",
						i+1, doc.TipoDocumento, doc.NumeroDocumento, doc.Estado)
				}
			}
		}

		// Eliminación de empresa (comentado por seguridad)
		fmt.Println("\n=== Eliminación de empresa ===")
		fmt.Printf("¿Quieres eliminar la empresa con ID %s? (Esta operación se ha comentado por seguridad)\n", empresaCreada.ID)

		// Comentado para evitar eliminaciones accidentales en el ejemplo
		// Descomentar para probar la eliminación
		/*
			err = repo.DeleteEmpresa(ctx, empresaCreada.ID)
			if err != nil {
				log.Printf("Error eliminando empresa: %v", err)
			} else {
				fmt.Println("Empresa eliminada exitosamente")
			}
		*/
	}

	fmt.Println("\nEjemplo completado.")
}
