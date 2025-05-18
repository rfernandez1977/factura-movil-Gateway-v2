package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"FMgo/core/firma/services"
)

func main() {
	// Configurar rutas
	certPath := os.Getenv("CERT_PATH")
	keyPath := os.Getenv("KEY_PATH")
	password := os.Getenv("CERT_PASS")
	rut := os.Getenv("RUT_FIRMANTE")

	if certPath == "" || keyPath == "" || rut == "" {
		log.Fatal("Faltan variables de entorno requeridas: CERT_PATH, KEY_PATH, RUT_FIRMANTE")
	}

	// Crear directorio de logs si no existe
	if err := os.MkdirAll("logs", 0755); err != nil {
		log.Fatalf("Error creando directorio de logs: %v", err)
	}

	// Inicializar servicios
	firmaService, err := services.NewFirmaService(certPath, keyPath, password, rut)
	if err != nil {
		log.Fatalf("Error inicializando servicio de firma: %v", err)
	}

	testService, err := services.NewTestService(firmaService, "certificacion")
	if err != nil {
		log.Fatalf("Error inicializando servicio de pruebas: %v", err)
	}

	// Ejecutar pruebas
	fmt.Println("Iniciando pruebas de integración...")

	// 1. Probar firma de XML
	fmt.Println("\n1. Probando firma de XML...")
	xmlPrueba := testService.GenerarXMLPrueba()
	if err := testService.ValidarFirma(xmlPrueba); err != nil {
		log.Printf("Error en prueba de firma: %v", err)
	} else {
		fmt.Println("✅ Prueba de firma exitosa")
	}

	// 2. Probar obtención de semilla
	fmt.Println("\n2. Probando obtención de semilla...")
	if err := testService.ProbarSemilla(); err != nil {
		log.Printf("Error en prueba de semilla: %v", err)
	} else {
		fmt.Println("✅ Prueba de semilla exitosa")
	}

	// 3. Probar obtención de token
	fmt.Println("\n3. Probando obtención de token...")
	if err := testService.ProbarToken("1234567890"); err != nil {
		log.Printf("Error en prueba de token: %v", err)
	} else {
		fmt.Println("✅ Prueba de token exitosa")
	}

	// 4. Probar flujo completo
	fmt.Println("\n4. Probando flujo completo...")
	if err := testService.ProbarFlujoCompleto(); err != nil {
		log.Printf("Error en prueba de flujo completo: %v", err)
	} else {
		fmt.Println("✅ Prueba de flujo completo exitosa")
	}

	fmt.Println("\nPruebas completadas. Revisar logs para más detalles.")
	fmt.Printf("Logs disponibles en: %s\n", filepath.Join("logs", "test"))
}
