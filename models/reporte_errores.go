package models

import (
	"time"
)

// ErrorResumen representa un resumen de error para el reporte
type ErrorResumen struct {
	Tipo         string    `json:"tipo" bson:"tipo"`
	Mensaje      string    `json:"mensaje" bson:"mensaje"`
	Cantidad     int       `json:"cantidad" bson:"cantidad"`
	PrimerError  time.Time `json:"primer_error" bson:"primer_error"`
	UltimoError  time.Time `json:"ultimo_error" bson:"ultimo_error"`
	EjemploID    string    `json:"ejemplo_id" bson:"ejemplo_id"`
	ModuloOrigen string    `json:"modulo_origen" bson:"modulo_origen"`
	Resueltos    int       `json:"resueltos" bson:"resueltos"`
	SinResolver  int       `json:"sin_resolver" bson:"sin_resolver"`
}

// EstadisticasPorModulo representa estadísticas de errores por módulo
type EstadisticasPorModulo struct {
	Modulo          string `json:"modulo" bson:"modulo"`
	TotalErrores    int    `json:"total_errores" bson:"total_errores"`
	ErroresCriticos int    `json:"errores_criticos" bson:"errores_criticos"`
	ErroresWarning  int    `json:"errores_warning" bson:"errores_warning"`
	ErroresInfo     int    `json:"errores_info" bson:"errores_info"`
	Resueltos       int    `json:"resueltos" bson:"resueltos"`
	SinResolver     int    `json:"sin_resolver" bson:"sin_resolver"`
}

// EstadisticasPorDia representa estadísticas de errores por día
type EstadisticasPorDia struct {
	Fecha           time.Time `json:"fecha" bson:"fecha"`
	TotalErrores    int       `json:"total_errores" bson:"total_errores"`
	ErroresCriticos int       `json:"errores_criticos" bson:"errores_criticos"`
	ErroresWarning  int       `json:"errores_warning" bson:"errores_warning"`
	ErroresInfo     int       `json:"errores_info" bson:"errores_info"`
	Resueltos       int       `json:"resueltos" bson:"resueltos"`
	SinResolver     int       `json:"sin_resolver" bson:"sin_resolver"`
}

// ReporteErrores representa un reporte de errores
type ReporteErrores struct {
	ID                  string                  `json:"id" bson:"_id,omitempty"`
	FechaGeneracion     time.Time               `json:"fecha_generacion" bson:"fecha_generacion"`
	FechaInicio         time.Time               `json:"fecha_inicio" bson:"fecha_inicio"`
	FechaFin            time.Time               `json:"fecha_fin" bson:"fecha_fin"`
	EmpresaID           string                  `json:"empresa_id" bson:"empresa_id"`
	UsuarioID           string                  `json:"usuario_id" bson:"usuario_id"`
	TotalErrores        int                     `json:"total_errores" bson:"total_errores"`
	ErroresCriticos     int                     `json:"errores_criticos" bson:"errores_criticos"`
	ErroresWarning      int                     `json:"errores_warning" bson:"errores_warning"`
	ErroresInfo         int                     `json:"errores_info" bson:"errores_info"`
	Resueltos           int                     `json:"resueltos" bson:"resueltos"`
	SinResolver         int                     `json:"sin_resolver" bson:"sin_resolver"`
	ErroresResumen      []ErrorResumen          `json:"errores_resumen" bson:"errores_resumen"`
	EstadisticasModulo  []EstadisticasPorModulo `json:"estadisticas_modulo" bson:"estadisticas_modulo"`
	EstadisticasDiarias []EstadisticasPorDia    `json:"estadisticas_diarias" bson:"estadisticas_diarias"`
	Filtros             map[string]interface{}  `json:"filtros" bson:"filtros"`
	FormatoReporte      string                  `json:"formato_reporte" bson:"formato_reporte"`
	ErroresDetallados   []ErrorLog              `json:"errores_detallados,omitempty" bson:"errores_detallados,omitempty"`
}

// GetField obtiene el valor de un campo
func (r *ReporteErrores) GetField(campo string) interface{} {
	switch campo {
	case "id":
		return r.ID
	case "fecha_generacion":
		return r.FechaGeneracion
	case "fecha_inicio":
		return r.FechaInicio
	case "fecha_fin":
		return r.FechaFin
	case "empresa_id":
		return r.EmpresaID
	case "usuario_id":
		return r.UsuarioID
	case "total_errores":
		return r.TotalErrores
	case "errores_criticos":
		return r.ErroresCriticos
	case "errores_warning":
		return r.ErroresWarning
	case "errores_info":
		return r.ErroresInfo
	case "resueltos":
		return r.Resueltos
	case "sin_resolver":
		return r.SinResolver
	case "formato_reporte":
		return r.FormatoReporte
	}

	// Si no se encuentra el campo específico, buscar en filtros
	if r.Filtros != nil {
		if valor, ok := r.Filtros[campo]; ok {
			return valor
		}
	}

	return nil
}

// SetField establece el valor de un campo
func (r *ReporteErrores) SetField(campo string, valor interface{}) error {
	switch campo {
	case "id":
		if id, ok := valor.(string); ok {
			r.ID = id
			return nil
		}
	case "fecha_generacion":
		if fecha, ok := valor.(time.Time); ok {
			r.FechaGeneracion = fecha
			return nil
		}
	case "fecha_inicio":
		if fecha, ok := valor.(time.Time); ok {
			r.FechaInicio = fecha
			return nil
		}
	case "fecha_fin":
		if fecha, ok := valor.(time.Time); ok {
			r.FechaFin = fecha
			return nil
		}
	case "empresa_id":
		if empresa, ok := valor.(string); ok {
			r.EmpresaID = empresa
			return nil
		}
	case "usuario_id":
		if usuario, ok := valor.(string); ok {
			r.UsuarioID = usuario
			return nil
		}
	case "total_errores":
		if total, ok := valor.(int); ok {
			r.TotalErrores = total
			return nil
		}
	case "errores_criticos":
		if criticos, ok := valor.(int); ok {
			r.ErroresCriticos = criticos
			return nil
		}
	case "errores_warning":
		if warning, ok := valor.(int); ok {
			r.ErroresWarning = warning
			return nil
		}
	case "errores_info":
		if info, ok := valor.(int); ok {
			r.ErroresInfo = info
			return nil
		}
	case "resueltos":
		if resueltos, ok := valor.(int); ok {
			r.Resueltos = resueltos
			return nil
		}
	case "sin_resolver":
		if sinResolver, ok := valor.(int); ok {
			r.SinResolver = sinResolver
			return nil
		}
	case "formato_reporte":
		if formato, ok := valor.(string); ok {
			r.FormatoReporte = formato
			return nil
		}
	default:
		// Si no es un campo específico, almacenar en filtros
		if r.Filtros == nil {
			r.Filtros = make(map[string]interface{})
		}
		r.Filtros[campo] = valor
		return nil
	}

	return NewValidationFieldError(campo, "Tipo de dato inválido para el campo", "INVALID_TYPE", valor)
}
