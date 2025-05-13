package utils

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// ParseUUID parses a string into a UUID
func ParseUUID(id string) (uuid.UUID, error) {
	// Remove any surrounding whitespace
	id = strings.TrimSpace(id)

	// Parse the UUID
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid UUID format: %v", err)
	}

	return parsedUUID, nil
}

// GenerateUUID generates a new UUID v4
func GenerateUUID() string {
	return uuid.New().String()
}

// IsValidUUID checks if a string is a valid UUID
func IsValidUUID(id string) bool {
	_, err := ParseUUID(id)
	return err == nil
}

// FormatUUID formats a UUID string to ensure consistent format
func FormatUUID(id string) (string, error) {
	parsedUUID, err := ParseUUID(id)
	if err != nil {
		return "", err
	}
	return parsedUUID.String(), nil
}

// DocumentType representa el tipo de documento
type DocumentType string

const (
	DocumentTypeFactura     DocumentType = "FACTURA"
	DocumentTypeBoleta      DocumentType = "BOLETA"
	DocumentTypeNotaCredito DocumentType = "NOTA_CREDITO"
	DocumentTypeNotaDebito  DocumentType = "NOTA_DEBITO"
)

// DocumentSource representa la fuente del documento
type DocumentSource string

const (
	DocumentSourceExternal DocumentSource = "EXTERNAL" // Documento con stamp pre-generado
	DocumentSourceInternal DocumentSource = "INTERNAL" // Documento sin stamp
)

// GenerateDocumentUUIDWithSource genera un UUID específico para documentos tributarios
// considerando si el documento viene con stamp pre-generado o no
func GenerateDocumentUUIDWithSource(docType DocumentType, date time.Time, docNumber string, source DocumentSource) string {
	// Crear un UUID v5 usando el namespace de documentos tributarios
	namespace := uuid.NewSHA1(uuid.Nil, []byte("documentos_tributarios"))

	// Crear un identificador único combinando tipo, fecha, número y fuente
	uniqueID := fmt.Sprintf("%s|%s|%s|%s", docType, date.Format("2006-01-02"), docNumber, source)

	// Generar el UUID v5
	docUUID := uuid.NewSHA1(namespace, []byte(uniqueID))
	return docUUID.String()
}

// GenerateDocumentUUIDFromStamp genera un UUID basado en el stamp del documento
func GenerateDocumentUUIDFromStamp(stamp string) string {
	namespace := uuid.NewSHA1(uuid.Nil, []byte("documentos_tributarios_stamp"))
	return uuid.NewSHA1(namespace, []byte(stamp)).String()
}

// GenerateDocumentUUID es un wrapper para mantener compatibilidad
func GenerateDocumentUUID(docType string, date time.Time, docNumber string) string {
	return GenerateDocumentUUIDWithSource(DocumentType(docType), date, docNumber, DocumentSourceInternal)
}

// GenerateTransactionUUID genera un UUID específico para transacciones
// basado en el tipo de transacción y timestamp
func GenerateTransactionUUID(transactionType string) string {
	namespace := uuid.NewSHA1(uuid.Nil, []byte("transacciones"))
	timestamp := time.Now().UnixNano()
	uniqueID := fmt.Sprintf("%s|%d", transactionType, timestamp)
	return uuid.NewSHA1(namespace, []byte(uniqueID)).String()
}

// GenerateClientUUID genera un UUID específico para clientes
// basado en el RUT y nombre del cliente
func GenerateClientUUID(rut, name string) string {
	namespace := uuid.NewSHA1(uuid.Nil, []byte("clientes"))
	uniqueID := fmt.Sprintf("%s|%s", rut, name)
	return uuid.NewSHA1(namespace, []byte(uniqueID)).String()
}

// GenerateProductUUID genera un UUID específico para productos
// basado en el código y nombre del producto
func GenerateProductUUID(code, name string) string {
	namespace := uuid.NewSHA1(uuid.Nil, []byte("productos"))
	uniqueID := fmt.Sprintf("%s|%s", code, name)
	return uuid.NewSHA1(namespace, []byte(uniqueID)).String()
}
