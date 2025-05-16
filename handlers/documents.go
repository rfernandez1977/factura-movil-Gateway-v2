package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/fmgo/gateway/api"
	"github.com/fmgo/gateway/metrics"
)

// validateDocumentData validates required fields for documents
func validateDocumentData(docData map[string]interface{}, docType string) error {
	requiredFields := []string{"date", "details", "netTotal"}
	for _, field := range requiredFields {
		if _, exists := docData[field]; !exists {
			return fmt.Errorf("missing required field: %s", field)
		}
	}
	return nil
}

// CreateDocumentHandler handles creation of documents (facturas, boletas, notas, guias)
func CreateDocumentHandler(db *sql.DB, docType string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var docData map[string]interface{}
		if err := c.BindJSON(&docData); err != nil {
			log.Printf("Invalid request payload for %s: %v", docType, err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}

		// Validate document data
		if err := validateDocumentData(docData, docType); err != nil {
			log.Printf("Validation failed for %s: %v", docType, err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Determine endpoint based on document type
		var endpoint string
		switch docType {
		case "factura":
			endpoint = "/services/raw/company/29/invoice"
		case "boleta":
			endpoint = "/services/raw/company/29/ticket"
		case "nota":
			endpoint = "/services/raw/company/29/note"
		case "guia":
			endpoint = "/services/raw/company/29/waybill"
		}

		// Call Factura Móvil API
		resp, err := api.CallFacturaMovil("POST", endpoint, docData, 2)
		if err != nil {
			log.Printf("Failed to communicate with Factura Móvil for %s: %v", docType, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to communicate with Factura Móvil: " + err.Error()})
			metrics.RequestCounter.WithLabelValues(docType, "error").Inc()
			return
		}
		defer resp.Body.Close()

		// Store document in local database
		data, _ := json.Marshal(docData)
		_, err = db.Exec("INSERT INTO documents (type, data, created_at) VALUES ($1, $2, $3)",
			docType, string(data), time.Now())
		if err != nil {
			log.Printf("Failed to store %s in database: %v", docType, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store document in database"})
			return
		}

		log.Printf("%s created successfully", docType)
		c.JSON(http.StatusOK, gin.H{"status": fmt.Sprintf("%s created successfully", docType)})
		metrics.RequestCounter.WithLabelValues(docType, "success").Inc()
	}
}

// GetDocumentStatusHandler handles document status queries
func GetDocumentStatusHandler(c *gin.Context) {
	docID := c.Param("id")
	endpoint := "/services/common/company/29/document/" + docID + "/getState"

	resp, err := api.CallFacturaMovil("GET", endpoint, nil, 2)
	if err != nil {
		log.Printf("Failed to get document status for %s: %v", docID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to communicate with Factura Móvil: " + err.Error()})
		metrics.RequestCounter.WithLabelValues("document_status", "error").Inc()
		return
	}
	defer resp.Body.Close()

	var statusData map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&statusData); err != nil {
		log.Printf("Failed to parse document status response: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse response"})
		return
	}

	log.Printf("Document status retrieved for %s", docID)
	c.JSON(http.StatusOK, statusData)
	metrics.RequestCounter.WithLabelValues("document_status", "success").Inc()
}

// GetDocumentPDFHandler handles PDF downloads of documents
func GetDocumentPDFHandler(c *gin.Context) {
	docID := c.Param("id")
	endpoint := "/document/toPdf/" + docID

	resp, err := api.CallFacturaMovil("GET", endpoint, nil, 2)
	if err != nil {
		log.Printf("Failed to get document PDF for %s: %v", docID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to communicate with Factura Móvil: " + err.Error()})
		metrics.RequestCounter.WithLabelValues("document_pdf", "error").Inc()
		return
	}
	defer resp.Body.Close()

	c.Header("Content-Type", "application/pdf")
	c.Writer.WriteHeader(http.StatusOK)
	if _, err := c.Writer.WriteFrom(resp.Body); err != nil {
		log.Printf("Failed to serve PDF for %s: %v", docID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to serve PDF"})
		return
	}

	log.Printf("Document PDF retrieved for %s", docID)
	metrics.RequestCounter.WithLabelValues("document_pdf", "success").Inc()
}
