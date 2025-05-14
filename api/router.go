package api

import (
	"encoding/json"
	"net/http"
)

// Router encapsula un enrutador HTTP
type Router struct {
	mux *http.ServeMux
}

// NewRouter crea un nuevo Router
func NewRouter() *Router {
	return &Router{
		mux: http.NewServeMux(),
	}
}

// Get registra un manejador para una ruta GET
func (r *Router) Get(path string, handler http.HandlerFunc) {
	r.mux.HandleFunc(path, func(w http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet {
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
			return
		}
		handler(w, req)
	})
}

// Post registra un manejador para una ruta POST
func (r *Router) Post(path string, handler http.HandlerFunc) {
	r.mux.HandleFunc(path, func(w http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
			return
		}
		handler(w, req)
	})
}

// Put registra un manejador para una ruta PUT
func (r *Router) Put(path string, handler http.HandlerFunc) {
	r.mux.HandleFunc(path, func(w http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPut {
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
			return
		}
		handler(w, req)
	})
}

// Delete registra un manejador para una ruta DELETE
func (r *Router) Delete(path string, handler http.HandlerFunc) {
	r.mux.HandleFunc(path, func(w http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodDelete {
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
			return
		}
		handler(w, req)
	})
}

// ServeHTTP implementa la interfaz http.Handler
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}

// RespondWithJSON envía una respuesta JSON al cliente
func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// RespondWithError envía un error al cliente
func RespondWithError(w http.ResponseWriter, code int, message string) {
	RespondWithJSON(w, code, map[string]string{"error": message})
}
