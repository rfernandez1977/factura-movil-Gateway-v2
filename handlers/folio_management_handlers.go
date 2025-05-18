package handlers

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"FMgo/models"
	"FMgo/repository" // Corregido el path de importación
)

type FolioManagementHandlers struct {
	repo  repository.DocumentRepository // Usar la interfaz
	mutex sync.Mutex
}

// Constructor para inyectar la dependencia
func NewFolioManagementHandlers(repo repository.DocumentRepository) *FolioManagementHandlers {
	return &FolioManagementHandlers{
		repo: repo,
	}
}

type ControlFolio struct {
	TipoDocumento     string    `json:"tipoDocumento"`
	RangoInicial      int       `json:"rangoInicial"`
	RangoFinal        int       `json:"rangoFinal"`
	FolioActual       int       `json:"folioActual"`
	FoliosDisponibles int       `json:"foliosDisponibles"`
	UltimoUso         time.Time `json:"ultimoUso"`
	EstadoCAF         string    `json:"estadoCAF"`
	AlertaGenerada    bool      `json:"alertaGenerada"`
}

type AsignacionFolio struct {
	TipoDocumento   string    `json:"tipoDocumento"`
	Folio           int       `json:"folio"`
	FechaAsignacion time.Time `json:"fechaAsignacion"`
	Usuario         string    `json:"usuario"`
	EstadoUso       string    `json:"estadoUso"`
}

func (h *FolioManagementHandlers) AsignarFolioHandler(c *gin.Context) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	var requestBody struct {
		TipoDocumento string `json:"tipoDocumento"`
		Usuario       string `json:"usuario"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "codigo": "REQ_001"})
		return
	}

	// Obtener control de folios desde la BD usando el repositorio
	control, err := h.repo.GetControlFolio(requestBody.TipoDocumento)
	if err != nil {
		log.Printf("Error al obtener control de folios: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "Error interno al obtener control de folios",
			"codigo": "FOLIO_500",
		})
		return
	}
	if control == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":  "Control de folios no encontrado para el tipo de documento",
			"codigo": "FOLIO_404",
		})
		return
	}

	// Validar disponibilidad de folios
	if control.FoliosDisponibles <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "No hay folios disponibles",
			"codigo":  "FOLIO_001",
			"detalle": "Solicite nuevo CAF",
		})
		return
	}

	// Verificar umbral de alerta
	if control.FoliosDisponibles < 100 && !control.AlertaGenerada {
		h.generarAlertaFolios(struct {
			TipoDocumento     string
			FoliosDisponibles int
		}{
			TipoDocumento:     control.TipoDocumento,
			FoliosDisponibles: control.FoliosDisponibles,
		})
		control.AlertaGenerada = true
	}

	// Asignar nuevo folio
	nuevoFolio := control.FolioActual + 1
	asignacion := models.AsignacionFolio{
		TipoDocumento:   control.TipoDocumento,
		Folio:           nuevoFolio,
		FechaAsignacion: time.Now(),
		Usuario:         requestBody.Usuario,
		EstadoUso:       "ASIGNADO",
	}

	// Actualizar control de folios en la BD
	control.FolioActual = nuevoFolio
	control.FoliosDisponibles--
	control.UltimoUso = time.Now()

	// Guardar la asignación y actualizar el control en una transacción si es posible
	// Idealmente, esto debería ser una transacción
	errAsignacion := h.repo.SaveAsignacionFolio(asignacion)
	errUpdate := h.repo.UpdateControlFolio(*control)

	if errAsignacion != nil || errUpdate != nil {
		log.Printf("Error al guardar asignación o actualizar control: %v, %v", errAsignacion, errUpdate)
		// Aquí se podría intentar revertir la asignación si falló la actualización del control
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "Error interno al asignar folio",
			"codigo": "FOLIO_501",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"mensaje":    "Folio asignado correctamente",
		"asignacion": asignacion,
		"control":    control,
	})
}

func (h *FolioManagementHandlers) ValidarFolioHandler(c *gin.Context) {
	var folio struct {
		TipoDocumento string `json:"tipoDocumento"`
		Folio         int    `json:"folio"`
	}

	// Validar rango de folio
	if err := h.validarRangoFolio(folio); err != nil {
		c.JSON(400, gin.H{
			"error":   "Folio fuera de rango",
			"codigo":  "FOLIO_002",
			"detalle": err.Error(),
		})
		return
	}

	// Validar uso previo
	if usado, err := h.validarUsoPrevio(folio); usado || err != nil {
		c.JSON(400, gin.H{
			"error":   "Folio ya utilizado",
			"codigo":  "FOLIO_003",
			"detalle": "El folio ya ha sido asignado previamente",
		})
		return
	}

	c.JSON(200, gin.H{
		"mensaje": "Folio válido",
		"folio":   folio,
	})
}

func (h *FolioManagementHandlers) generarAlertaFolios(control struct {
	TipoDocumento     string
	FoliosDisponibles int
}) {
	log.Printf("ALERTA: Quedan pocos folios para %s. Disponibles: %d", control.TipoDocumento, control.FoliosDisponibles)
	// Implementar envío de notificación real (email, Slack, etc.)
}

func (h *FolioManagementHandlers) validarRangoFolio(folio struct {
	TipoDocumento string `json:"tipoDocumento"`
	Folio         int    `json:"folio"`
}) error {
	// Implementar validación de rango
	return nil
}

func (h *FolioManagementHandlers) validarUsoPrevio(folio struct {
	TipoDocumento string `json:"tipoDocumento"`
	Folio         int    `json:"folio"`
}) (bool, error) {
	// Implementar validación de uso previo
	return false, nil
}
