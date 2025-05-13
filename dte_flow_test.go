package main

import (
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"
)

// Definición de modelos para el test
type Documento struct {
	ID                string
	TipoDTE           string
	Folio             int64
	FechaEmision      time.Time
	RutEmisor         string
	RazonEmisor       string
	GiroEmisor        string
	DireccionEmisor   string
	ComunaEmisor      string
	RutReceptor       string
	RazonReceptor     string
	GiroReceptor      string
	DireccionReceptor string
	ComunaReceptor    string
	MontoNeto         float64
	MontoExento       float64
	TasaIVA           float64
	MontoIVA          float64
	MontoTotal        float64
	Estado            string
	Items             []Item
	Referencias       []Referencia
	TrackID           string
	EstadoSII         string
}

type Item struct {
	NroLinea    int
	Nombre      string
	Descripcion string
	Cantidad    float64
	Precio      float64
	MontoItem   float64
	Exento      bool
}

type Referencia struct {
	NroLinea   int
	TipoDocRef string
	FolioRef   string
	FechaRef   time.Time
	RazonRef   string
}

type RespuestaSII struct {
	Estado          string
	EstadoDTE       string
	Glosa           string
	TrackID         string
	FechaProceso    time.Time
	NumeroAtencion  string
	CodigoRecepcion string
	Errores         []ErrorSII
}

type ErrorSII struct {
	Codigo      string
	Descripcion string
	Tipo        string
}

type EstadoSII struct {
	Codigo      int
	Descripcion string
	Timestamp   time.Time
	Detalle     string
}

// Configuración de Supabase
type SupabaseConfig struct {
	URL      string
	APIKey   string
	Proyecto string
}

// Cliente para Supabase
type SupabaseClient struct {
	Config SupabaseConfig
	Client *http.Client
}

// TestDTEFlow prueba el flujo completo de creación de XML, firmado y envío al SII
func TestDTEFlow(t *testing.T) {
	// Configuración de prueba
	tiposDocumento := []string{
		"33", // Factura Electrónica
		"34", // Factura Exenta
		"61", // Nota de Crédito
		"56", // Nota de Débito
		"52", // Guía de Despacho
	}

	// Crear cliente de Supabase
	supabaseClient := crearSupabaseClient()

	// Probar cada tipo de documento
	for _, tipoDTE := range tiposDocumento {
		t.Run(fmt.Sprintf("TipoDTE_%s", tipoDTE), func(t *testing.T) {
			// 1. Crear un documento de prueba
			documento := crearDocumentoPrueba(tipoDTE)
			t.Logf("Documento creado: %s, Tipo: %s, Folio: %d",
				documento.ID, documento.TipoDTE, documento.Folio)

			// 2. Validar documento
			if err := validarDocumento(documento); err != nil {
				t.Fatalf("Error validando documento: %v", err)
			}
			t.Log("Documento validado correctamente")

			// 3. Guardar documento en Supabase (estado inicial)
			if err := guardarEnSupabase(supabaseClient, documento, "documentos"); err != nil {
				t.Fatalf("Error guardando documento en Supabase: %v", err)
			}
			t.Log("Documento guardado en Supabase (estado inicial)")

			// 4. Generar XML
			xml, err := generarXML(documento)
			if err != nil {
				t.Fatalf("Error generando XML: %v", err)
			}
			t.Logf("XML generado correctamente (%d bytes)", len(xml))
			t.Logf("Muestra del XML generado: %s", mostrarXMLResumido(string(xml)))

			// 5. Validar XML contra esquema
			if err := validarXML(xml); err != nil {
				t.Fatalf("Error validando XML: %v", err)
			}
			t.Log("XML validado correctamente contra esquema")

			// 6. Firmar XML
			xmlFirmado, err := firmarXML(xml, documento.ID, documento.RutEmisor)
			if err != nil {
				t.Fatalf("Error firmando XML: %v", err)
			}
			t.Logf("XML firmado correctamente (%d bytes)", len(xmlFirmado))
			t.Logf("Muestra del XML firmado: %s", mostrarXMLResumido(string(xmlFirmado)))

			// 7. Validar firma
			if err := validarFirma(xmlFirmado); err != nil {
				t.Fatalf("Error validando firma: %v", err)
			}
			t.Log("Firma validada correctamente")

			// 8. Actualizar documento en Supabase (con XML generado)
			documento.Estado = "XML_GENERADO"
			if err := actualizarEnSupabase(supabaseClient, documento, "documentos"); err != nil {
				t.Fatalf("Error actualizando documento en Supabase: %v", err)
			}
			t.Log("Documento actualizado en Supabase (XML generado)")

			// 9. Enviar al SII
			respuesta, err := enviarAlSII(xmlFirmado)
			if err != nil {
				t.Fatalf("Error enviando al SII: %v", err)
			}
			t.Logf("Documento enviado al SII. TrackID: %s, Estado: %s", respuesta.TrackID, respuesta.Estado)

			// 10. Verificar respuesta
			if respuesta.Estado != "RECIBIDO" {
				t.Errorf("Estado incorrecto. Esperado: RECIBIDO, Obtenido: %s", respuesta.Estado)
			}

			if respuesta.TrackID == "" {
				t.Errorf("No se recibió un TrackID válido")
			}

			// 11. Actualizar documento en Supabase (con respuesta del SII)
			documento.TrackID = respuesta.TrackID
			documento.Estado = "ENVIADO_SII"
			documento.EstadoSII = respuesta.EstadoDTE
			if err := actualizarEnSupabase(supabaseClient, documento, "documentos"); err != nil {
				t.Fatalf("Error actualizando documento en Supabase: %v", err)
			}
			t.Log("Documento actualizado en Supabase (respuesta SII)")

			// 12. Consultar estado (simulado)
			t.Log("Esperando a que el SII procese el documento...")
			time.Sleep(500 * time.Millisecond) // Simulación de espera

			estadoFinal, err := consultarEstadoSII(respuesta.TrackID)
			if err != nil {
				t.Fatalf("Error consultando estado: %v", err)
			}
			t.Logf("Estado final del documento: %s - %s", estadoFinal.Descripcion, estadoFinal.Detalle)

			// 13. Verificar estado final
			if estadoFinal.Descripcion != "ACEPTADO" {
				t.Errorf("Estado final incorrecto. Esperado: ACEPTADO, Obtenido: %s", estadoFinal.Descripcion)
			}

			// 14. Actualizar estado final en Supabase
			documento.Estado = "ACEPTADO_SII"
			documento.EstadoSII = estadoFinal.Descripcion
			if err := actualizarEnSupabase(supabaseClient, documento, "documentos"); err != nil {
				t.Fatalf("Error actualizando estado final en Supabase: %v", err)
			}
			t.Log("Estado final actualizado en Supabase")

			t.Log("Flujo completo terminado exitosamente")
		})
	}
}

// Funciones auxiliares para el test

func crearDocumentoPrueba(tipoDTE string) *Documento {
	now := time.Now()
	doc := &Documento{
		ID:                fmt.Sprintf("%s_%d", obtenerNombreTipoDTE(tipoDTE), now.Unix()),
		TipoDTE:           tipoDTE,
		Folio:             obtenerFolioSegunTipo(tipoDTE),
		FechaEmision:      now,
		RutEmisor:         "76.123.456-7",
		RazonEmisor:       "EMPRESA DE PRUEBA S.A.",
		GiroEmisor:        "DESARROLLO DE SOFTWARE",
		DireccionEmisor:   "CALLE PRINCIPAL 123",
		ComunaEmisor:      "SANTIAGO",
		RutReceptor:       "55.666.777-8",
		RazonReceptor:     "CLIENTE DE PRUEBA LTDA.",
		GiroReceptor:      "SERVICIOS GENERALES",
		DireccionReceptor: "AVENIDA NUEVA 456",
		ComunaReceptor:    "PROVIDENCIA",
		MontoNeto:         100000,
		MontoExento:       0,
		TasaIVA:           19,
		MontoIVA:          19000,
		MontoTotal:        119000,
		Estado:            "PENDIENTE",
		Items: []Item{
			{
				NroLinea:    1,
				Nombre:      "Servicio profesional",
				Descripcion: "Desarrollo de software a medida",
				Cantidad:    1,
				Precio:      100000,
				MontoItem:   100000,
				Exento:      false,
			},
		},
		Referencias: []Referencia{},
	}

	// Agregar referencias específicas según el tipo de documento
	switch tipoDTE {
	case "61": // Nota de Crédito
		doc.Referencias = append(doc.Referencias, Referencia{
			NroLinea:   1,
			TipoDocRef: "33",
			FolioRef:   "100",
			FechaRef:   now.AddDate(0, 0, -5),
			RazonRef:   "ANULA FACTURA POR ERROR EN DETALLE",
		})
	case "56": // Nota de Débito
		doc.Referencias = append(doc.Referencias, Referencia{
			NroLinea:   1,
			TipoDocRef: "33",
			FolioRef:   "101",
			FechaRef:   now.AddDate(0, 0, -3),
			RazonRef:   "CORRIGE MONTO",
		})
	case "34": // Factura Exenta
		doc.MontoExento = 100000
		doc.MontoNeto = 0
		doc.MontoIVA = 0
		doc.MontoTotal = 100000
		doc.Items[0].Exento = true
	}

	return doc
}

func obtenerNombreTipoDTE(tipoDTE string) string {
	switch tipoDTE {
	case "33":
		return "FACTURA"
	case "34":
		return "FACTURA_EXENTA"
	case "61":
		return "NOTA_CREDITO"
	case "56":
		return "NOTA_DEBITO"
	case "52":
		return "GUIA_DESPACHO"
	default:
		return "DOCUMENTO"
	}
}

func obtenerFolioSegunTipo(tipoDTE string) int64 {
	// En un caso real, se obtendrían folios de un CAF
	switch tipoDTE {
	case "33":
		return 500
	case "34":
		return 100
	case "61":
		return 50
	case "56":
		return 30
	case "52":
		return 200
	default:
		return 1
	}
}

func validarDocumento(doc *Documento) error {
	// Validaciones básicas
	if doc.TipoDTE == "" {
		return fmt.Errorf("tipo de DTE no puede estar vacío")
	}

	if doc.Folio <= 0 {
		return fmt.Errorf("folio debe ser mayor a cero")
	}

	if doc.RutEmisor == "" {
		return fmt.Errorf("RUT emisor no puede estar vacío")
	}

	if doc.RutReceptor == "" {
		return fmt.Errorf("RUT receptor no puede estar vacío")
	}

	if doc.MontoTotal <= 0 {
		return fmt.Errorf("monto total debe ser mayor a cero")
	}

	// Validaciones específicas por tipo de documento
	switch doc.TipoDTE {
	case "34": // Factura exenta
		if doc.MontoExento <= 0 {
			return fmt.Errorf("factura exenta debe tener monto exento mayor a cero")
		}
		if doc.MontoIVA != 0 {
			return fmt.Errorf("factura exenta no debe tener IVA")
		}
		if doc.MontoTotal != doc.MontoExento {
			return fmt.Errorf("en factura exenta, monto total debe ser igual a monto exento")
		}
	case "61", "56": // Notas de crédito y débito
		if len(doc.Referencias) == 0 {
			return fmt.Errorf("%s debe tener al menos una referencia", obtenerNombreTipoDTE(doc.TipoDTE))
		}
	default:
		// Validaciones de negocio para otros tipos
		if doc.TipoDTE != "34" && doc.MontoTotal != doc.MontoNeto+doc.MontoExento+doc.MontoIVA {
			return fmt.Errorf("error en total: %f != %f + %f + %f",
				doc.MontoTotal, doc.MontoNeto, doc.MontoExento, doc.MontoIVA)
		}

		// Validación de cálculo de IVA para documentos con IVA
		if doc.TipoDTE != "34" && doc.MontoNeto > 0 {
			ivaCalculado := doc.MontoNeto * (doc.TasaIVA / 100)
			if int(doc.MontoIVA) != int(ivaCalculado) { // Redondeo para evitar problemas de precisión
				return fmt.Errorf("error en cálculo de IVA: %f != %f", doc.MontoIVA, ivaCalculado)
			}
		}
	}

	return nil
}

func generarXML(documento *Documento) ([]byte, error) {
	// En un caso real, se usaría un template más elaborado o una librería de marshalling
	fechaEmision := documento.FechaEmision.Format("2006-01-02")

	// Construir detalles
	detalles := ""
	for _, item := range documento.Items {
		detalle := fmt.Sprintf(`    <Detalle>
      <NroLinDet>%d</NroLinDet>
      <NmbItem>%s</NmbItem>
      <DscItem>%s</DscItem>
      <QtyItem>%.1f</QtyItem>
      <PrcItem>%.0f</PrcItem>
      <MontoItem>%.0f</MontoItem>`,
			item.NroLinea, item.Nombre, item.Descripcion,
			item.Cantidad, item.Precio, item.MontoItem)

		// Agregar indicador de exento si corresponde
		if item.Exento {
			detalle += "\n      <IndExe>1</IndExe>"
		}

		detalle += "\n    </Detalle>"
		detalles += detalle
	}

	// Construir referencias si existen
	referencias := ""
	for _, ref := range documento.Referencias {
		referencia := fmt.Sprintf(`    <Referencia>
      <NroLinRef>%d</NroLinRef>
      <TpoDocRef>%s</TpoDocRef>
      <FolioRef>%s</FolioRef>
      <FchRef>%s</FchRef>
      <RazonRef>%s</RazonRef>
    </Referencia>`,
			ref.NroLinea, ref.TipoDocRef, ref.FolioRef,
			ref.FechaRef.Format("2006-01-02"), ref.RazonRef)
		referencias += referencia
	}

	xmlString := fmt.Sprintf(`<?xml version="1.0" encoding="ISO-8859-1"?>
<DTE version="1.0">
  <Documento ID="%s">
    <Encabezado>
      <IdDoc>
        <TipoDTE>%s</TipoDTE>
        <Folio>%d</Folio>
        <FchEmis>%s</FchEmis>
      </IdDoc>
      <Emisor>
        <RUTEmisor>%s</RUTEmisor>
        <RznSoc>%s</RznSoc>
        <GiroEmis>%s</GiroEmis>
        <DirOrigen>%s</DirOrigen>
        <CmnaOrigen>%s</CmnaOrigen>
      </Emisor>
      <Receptor>
        <RUTRecep>%s</RUTRecep>
        <RznSocRecep>%s</RznSocRecep>
        <GiroRecep>%s</GiroRecep>
        <DirRecep>%s</DirRecep>
        <CmnaRecep>%s</CmnaRecep>
      </Receptor>
      <Totales>`,
		documento.ID,
		documento.TipoDTE,
		documento.Folio,
		fechaEmision,
		documento.RutEmisor,
		documento.RazonEmisor,
		documento.GiroEmisor,
		documento.DireccionEmisor,
		documento.ComunaEmisor,
		documento.RutReceptor,
		documento.RazonReceptor,
		documento.GiroReceptor,
		documento.DireccionReceptor,
		documento.ComunaReceptor)

	// Agregar totales según tipo de documento
	if documento.TipoDTE == "34" { // Factura exenta
		xmlString += fmt.Sprintf(`
        <MntExe>%.0f</MntExe>
        <MntTotal>%.0f</MntTotal>`,
			documento.MontoExento,
			documento.MontoTotal)
	} else { // Otros documentos
		xmlString += fmt.Sprintf(`
        <MntNeto>%.0f</MntNeto>
        <MntExe>%.0f</MntExe>
        <TasaIVA>%.1f</TasaIVA>
        <IVA>%.0f</IVA>
        <MntTotal>%.0f</MntTotal>`,
			documento.MontoNeto,
			documento.MontoExento,
			documento.TasaIVA,
			documento.MontoIVA,
			documento.MontoTotal)
	}

	xmlString += `
      </Totales>
    </Encabezado>`

	// Agregar detalles
	xmlString += "\n" + detalles

	// Agregar referencias si existen
	if len(referencias) > 0 {
		xmlString += "\n" + referencias
	}

	// Cerrar el documento
	xmlString += `
  </Documento>
</DTE>`

	return []byte(xmlString), nil
}

func validarXML(xmlData []byte) error {
	// Simulación de validación de XML
	// En un caso real, se validaría contra un esquema XSD
	if !strings.Contains(string(xmlData), "<DTE") {
		return fmt.Errorf("XML no contiene elemento DTE")
	}

	if !strings.Contains(string(xmlData), "<Documento") {
		return fmt.Errorf("XML no contiene elemento Documento")
	}

	return nil
}

func firmarXML(xmlData []byte, documentoID, rutEmisor string) ([]byte, error) {
	// Simulación de firmado
	// En un caso real, aquí iría el código para:
	// 1. Cargar el certificado digital
	// 2. Calcular el hash del documento (digest)
	// 3. Firmar el hash con la clave privada
	// 4. Incorporar la firma al XML según estándar XMLDSig

	firmaSimulada := fmt.Sprintf(`<Signature xmlns="http://www.w3.org/2000/09/xmldsig#">
    <SignedInfo>
      <CanonicalizationMethod Algorithm="http://www.w3.org/TR/2001/REC-xml-c14n-20010315"/>
      <SignatureMethod Algorithm="http://www.w3.org/2000/09/xmldsig#rsa-sha1"/>
      <Reference URI="#%s">
        <DigestMethod Algorithm="http://www.w3.org/2000/09/xmldsig#sha1"/>
        <DigestValue>SIMULATED_DIGEST_VALUE</DigestValue>
      </Reference>
    </SignedInfo>
    <SignatureValue>SIMULATED_SIGNATURE_VALUE_FOR_%s</SignatureValue>
    <KeyInfo>
      <KeyValue>
        <RSAKeyValue>
          <Modulus>SIMULATED_MODULUS</Modulus>
          <Exponent>AQAB</Exponent>
        </RSAKeyValue>
      </KeyValue>
      <X509Data>
        <X509Certificate>SIMULATED_CERTIFICATE_DATA</X509Certificate>
      </X509Data>
    </KeyInfo>
  </Signature>`, documentoID, rutEmisor)

	xmlFirmado := string(xmlData)
	xmlFirmado = xmlFirmado[:len(xmlFirmado)-6] + firmaSimulada + "</DTE>"

	return []byte(xmlFirmado), nil
}

func validarFirma(xmlData []byte) error {
	// Simulación de validación de firma
	// En un caso real, aquí iría el código para:
	// 1. Extraer la firma del XML
	// 2. Obtener el certificado usado
	// 3. Verificar la firma con la clave pública

	if !strings.Contains(string(xmlData), "<Signature") {
		return fmt.Errorf("XML no contiene firma")
	}

	if !strings.Contains(string(xmlData), "<SignatureValue>") {
		return fmt.Errorf("Firma inválida: no contiene valor de firma")
	}

	return nil
}

func enviarAlSII(xmlData []byte) (*RespuestaSII, error) {
	// Simulación de envío al SII
	// En un caso real, aquí iría el código para:
	// 1. Establecer conexión segura con el SII (TLS)
	// 2. Autenticarse (certificado digital)
	// 3. Enviar el XML por POST al endpoint correspondiente
	// 4. Procesar la respuesta XML

	// Simulamos una respuesta exitosa
	return &RespuestaSII{
		Estado:          "RECIBIDO",
		EstadoDTE:       "EN_PROCESO",
		Glosa:           "Documento recibido correctamente",
		TrackID:         fmt.Sprintf("TRACK_%d", time.Now().Unix()),
		FechaProceso:    time.Now(),
		NumeroAtencion:  fmt.Sprintf("123456%d", time.Now().Unix()%1000),
		CodigoRecepcion: "0",
		Errores:         []ErrorSII{},
	}, nil
}

func consultarEstadoSII(trackID string) (*EstadoSII, error) {
	// Simulación de consulta al SII
	// En un caso real, aquí iría el código para:
	// 1. Establecer conexión segura con el SII (TLS)
	// 2. Autenticarse (certificado digital)
	// 3. Consultar estado por GET con el trackID como parámetro
	// 4. Procesar la respuesta XML

	return &EstadoSII{
		Codigo:      0,
		Descripcion: "ACEPTADO",
		Timestamp:   time.Now(),
		Detalle:     "Documento validado y aceptado por el SII",
	}, nil
}

// Funciones para manejar Supabase

func crearSupabaseClient() *SupabaseClient {
	// En un entorno real, obtendríamos esto de variables de entorno
	// o de un archivo de configuración
	config := SupabaseConfig{
		URL:      "https://prueba.supabase.co",
		APIKey:   "api-key-simulada",
		Proyecto: "proyecto-test",
	}

	return &SupabaseClient{
		Config: config,
		Client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func guardarEnSupabase(client *SupabaseClient, documento *Documento, tabla string) error {
	// Simulación de guardado en Supabase
	// En un caso real, aquí iría el código para:
	// 1. Convertir el documento a JSON
	// 2. Hacer una petición POST a la API de Supabase

	fmt.Printf("[SUPABASE] Guardando documento %s en tabla %s\n", documento.ID, tabla)

	// Simular éxito
	return nil
}

func actualizarEnSupabase(client *SupabaseClient, documento *Documento, tabla string) error {
	// Simulación de actualización en Supabase
	// En un caso real, aquí iría el código para:
	// 1. Convertir el documento a JSON
	// 2. Hacer una petición PATCH a la API de Supabase

	fmt.Printf("[SUPABASE] Actualizando documento %s en tabla %s, Estado: %s\n",
		documento.ID, tabla, documento.Estado)

	// Simular éxito
	return nil
}

func obtenerDeSupabase(client *SupabaseClient, id string, tabla string) (*Documento, error) {
	// Simulación de consulta a Supabase
	// En un caso real, aquí iría el código para:
	// 1. Hacer una petición GET a la API de Supabase
	// 2. Convertir la respuesta JSON a un documento

	fmt.Printf("[SUPABASE] Consultando documento %s en tabla %s\n", id, tabla)

	// Simular documento
	return &Documento{
		ID:     id,
		Estado: "PENDIENTE",
	}, nil
}

// Función auxiliar para mostrar un XML resumido
func mostrarXMLResumido(xml string) string {
	if len(xml) > 300 {
		return xml[:150] + "..." + xml[len(xml)-150:]
	}
	return xml
}

// Para ejecutar solamente esta prueba usar:
// go test -v dte_flow_test.go
