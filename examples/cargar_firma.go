package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"FMgo/config"
	"FMgo/services"
)

func main() {
	// Cargar configuración
	supabaseConfig, err := config.NewSupabaseConfig()
	if err != nil {
		log.Fatalf("Error al cargar configuración: %v", err)
	}

	// Crear servicio de firma
	firmaService := services.NewFirmaService(supabaseConfig)

	// Leer archivo de firma
	firmaData, err := ioutil.ReadFile("firma_test/1/firma.pfx")
	if err != nil {
		log.Fatalf("Error al leer archivo de firma: %v", err)
	}

	// Extraer información de la firma
	certificado, err := firmaService.ExtraerInfoFirma(firmaData, "83559705FM")
	if err != nil {
		log.Fatalf("Error al extraer información de la firma: %v", err)
	}

	// Validar firma
	err = firmaService.ValidarFirma(certificado)
	if err != nil {
		log.Fatalf("Error al validar firma: %v", err)
	}

	// Imprimir información del certificado
	fmt.Println("Información del certificado:")
	fmt.Printf("Serial Number: %s\n", certificado.SerialNumber)
	fmt.Printf("Issuer: %s\n", certificado.Issuer)
	fmt.Printf("Subject: %s\n", certificado.Subject)
	fmt.Printf("Válido desde: %s\n", certificado.ValidFrom.Format("2006-01-02"))
	fmt.Printf("Válido hasta: %s\n", certificado.ValidTo.Format("2006-01-02"))

	// Guardar certificado en formato PEM para verificación
	err = os.WriteFile("certificado.pem", []byte(certificado.Certificate), 0644)
	if err != nil {
		log.Fatalf("Error al guardar certificado: %v", err)
	}

	fmt.Println("\nCertificado guardado en certificado.pem")
}
