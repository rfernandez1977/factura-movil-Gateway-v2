package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/fmgo/config"
	"github.com/fmgo/models"
	"github.com/fmgo/services"
)

func cargarConfiguracion() error {
	// Cargar configuración
	supabaseConfig, err := config.NewSupabaseConfig()
	if err != nil {
		return fmt.Errorf("error al cargar configuración: %v", err)
	}

	// Crear servicios
	empresaService := services.NewEmpresaService(supabaseConfig)
	certificadoService := services.NewCertificadoService(supabaseConfig)

	// Crear contexto
	ctx := context.Background()

	// Crear empresa
	empresa := &models.Empresa{
		Rut:              "76212889-6",
		RazonSocial:      "FACTURA MOVIL SPA",
		Giro:             "EMPRESAS DE SERVICIOS INTEGRALES DE INFORMÁTICA",
		Direccion:        "Vicuña Mackenna 9705",
		Comuna:           "La Florida",
		Ciudad:           "Santiago",
		Email:            "rfernandez@facturamovil.cl",
		ResolucionNumero: "12345", // Este valor debe ser actualizado con el real
		ResolucionFecha:  time.Now(),
		// Nuevos campos de firma
		FirmaRut:        "13195458-1",
		FirmaNombre:     "Rodrigo Fernandez Calderon",
		FirmaExpiracion: time.Now().AddDate(1, 0, 0), // Ajustar según la fecha real de expiración
	}

	// Guardar empresa
	err = empresaService.CrearEmpresa(ctx, empresa)
	if err != nil {
		return fmt.Errorf("error al crear empresa: %v", err)
	}

	// Leer archivo CAF
	cafXML, err := ioutil.ReadFile("caf_test/33-cert.xml")
	if err != nil {
		return fmt.Errorf("error al leer archivo CAF: %v", err)
	}

	// Crear CAF
	caf := &models.CAF{
		EmpresaID:        empresa.ID,
		TipoDocumento:    "33",
		XML:              string(cafXML),
		Estado:           "ACTIVO",
		FechaResolucion:  time.Now(), // Este valor debe ser actualizado con el real
		NumeroResolucion: "12345",    // Este valor debe ser actualizado con el real
	}

	// Guardar CAF
	err = empresaService.GuardarCAF(ctx, caf)
	if err != nil {
		return fmt.Errorf("error al guardar CAF: %v", err)
	}

	// Leer archivo de firma
	firmaData, err := ioutil.ReadFile("firma_test/1/firma.pfx")
	if err != nil {
		return fmt.Errorf("error al leer archivo de firma: %v", err)
	}

	// Extraer información del certificado
	certificado, err := certificadoService.ExtraerInfoCertificado(firmaData, "83559705FM")
	if err != nil {
		return fmt.Errorf("error al extraer información del certificado: %v", err)
	}

	// Asignar ID de empresa al certificado
	certificado.EmpresaID = empresa.ID

	// Validar certificado
	err = certificadoService.ValidarCertificado(certificado)
	if err != nil {
		return fmt.Errorf("error al validar certificado: %v", err)
	}

	// Guardar certificado
	err = empresaService.GuardarCertificado(ctx, certificado)
	if err != nil {
		return fmt.Errorf("error al guardar certificado: %v", err)
	}

	fmt.Println("Configuración cargada exitosamente:")
	fmt.Printf("Empresa ID: %s\n", empresa.ID)
	fmt.Printf("CAF ID: %s\n", caf.ID)
	fmt.Printf("Certificado ID: %s\n", certificado.ID)

	return nil
}

func main() {
	if err := cargarConfiguracion(); err != nil {
		log.Fatalf("Error: %v", err)
	}
}
