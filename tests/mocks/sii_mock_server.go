package mocks

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// EstadoDTE representa el estado de un documento tributario en el SII
type EstadoDTE struct {
	TrackID string    `json:"trackId"`
	Estado  string    `json:"estado"`
	Glosa   string    `json:"glosa"`
	Fecha   time.Time `json:"fecha,omitempty"`
}

// StartMockSIIServer inicia el servidor mock del SII
func StartMockSIIServer(port int) {
	// Configurar rutas
	http.HandleFunc("/api/v1/dte/emitidos", handleEmitirDocumento)
	http.HandleFunc("/api/v1/dte/estado", handleConsultarEstado)

	// Iniciar servidor
	addr := fmt.Sprintf(":%d", port)
	fmt.Printf("Iniciando servidor mock del SII en puerto %d...\n", port)
	go func() {
		log.Fatal(http.ListenAndServe(addr, nil))
	}()
}

// handleEmitirDocumento simula la emisión de un documento al SII
func handleEmitirDocumento(w http.ResponseWriter, r *http.Request) {
	// Verificar método
	if r.Method != "POST" {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Leer XML del documento
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, "Error al leer el cuerpo de la petición", http.StatusBadRequest)
		return
	}

	// Obtener headers
	rutEmisor := r.Header.Get("X-RUT-Emisor")
	tipoDTE := r.Header.Get("X-Tipo-DTE")
	folio := r.Header.Get("X-Folio")

	log.Printf("Recibida solicitud para emitir documento. RUT: %s, Tipo: %s, Folio: %s\n",
		rutEmisor, tipoDTE, folio)
	log.Printf("Contenido XML recibido (%d bytes): %s\n", len(body), truncateString(string(body), 100))

	// Simular respuesta exitosa
	response := map[string]string{
		"trackId": "123456789",
		"estado":  "RECIBIDO",
	}

	// Responder con JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// handleConsultarEstado simula la consulta de estado de un documento en el SII
func handleConsultarEstado(w http.ResponseWriter, r *http.Request) {
	// Verificar método
	if r.Method != "GET" {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Obtener trackID de la consulta
	trackID := r.URL.Query().Get("trackId")
	if trackID == "" {
		http.Error(w, "TrackID requerido", http.StatusBadRequest)
		return
	}

	log.Printf("Recibida consulta para documento con trackID: %s\n", trackID)

	// Simular respuesta exitosa para trackID 123456789
	if trackID == "123456789" {
		estado := EstadoDTE{
			TrackID: trackID,
			Estado:  "ACEPTADO",
			Glosa:   "Documento aceptado correctamente",
			Fecha:   time.Now(),
		}

		// Responder con JSON
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(estado)
		return
	}

	// Para cualquier otro trackID, simular documento no encontrado
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]string{
		"error": "Documento no encontrado",
	})
}

// truncateString trunca strings largos en logs
func truncateString(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max] + "..."
}
