package models

import (
	"time"
)

// ErrorLog representa un registro de error en el sistema
type ErrorLog struct {
	ID               string                 `json:"id" bson:"_id,omitempty"`
	Tipo             string                 `json:"tipo" bson:"tipo"` // ERROR, WARNING, INFO
	Mensaje          string                 `json:"mensaje" bson:"mensaje"`
	Detalles         string                 `json:"detalles" bson:"detalles"`
	FechaHora        time.Time              `json:"fecha_hora" bson:"fecha_hora"`
	ModuloOrigen     string                 `json:"modulo_origen" bson:"modulo_origen"`
	Funcion          string                 `json:"funcion" bson:"funcion"`
	Linea            int                    `json:"linea" bson:"linea"`
	UsuarioID        string                 `json:"usuario_id" bson:"usuario_id"`
	IP               string                 `json:"ip" bson:"ip"`
	Contexto         map[string]interface{} `json:"contexto" bson:"contexto,omitempty"`
	DocumentoID      string                 `json:"documento_id" bson:"documento_id,omitempty"`
	TipoDocumento    string                 `json:"tipo_documento" bson:"tipo_documento,omitempty"`
	EmpresaID        string                 `json:"empresa_id" bson:"empresa_id,omitempty"`
	Resuelto         bool                   `json:"resuelto" bson:"resuelto"`
	UsuarioResolutor string                 `json:"usuario_resolutor" bson:"usuario_resolutor,omitempty"`
	FechaResolucion  time.Time              `json:"fecha_resolucion" bson:"fecha_resolucion,omitempty"`
	NotaResolucion   string                 `json:"nota_resolucion" bson:"nota_resolucion,omitempty"`
}

// NewErrorLog crea una nueva instancia de ErrorLog
func NewErrorLog(tipo, mensaje, detalles, usuario string) *ErrorLog {
	return &ErrorLog{
		Tipo:      tipo,
		Mensaje:   mensaje,
		Detalles:  detalles,
		UsuarioID: usuario,
		FechaHora: time.Now(),
		Resuelto:  false,
	}
}

// GetField obtiene el valor de un campo
func (e *ErrorLog) GetField(campo string) interface{} {
	switch campo {
	case "id":
		return e.ID
	case "tipo":
		return e.Tipo
	case "mensaje":
		return e.Mensaje
	case "detalles":
		return e.Detalles
	case "fecha_hora":
		return e.FechaHora
	case "modulo_origen":
		return e.ModuloOrigen
	case "funcion":
		return e.Funcion
	case "linea":
		return e.Linea
	case "usuario_id":
		return e.UsuarioID
	case "ip":
		return e.IP
	case "documento_id":
		return e.DocumentoID
	case "tipo_documento":
		return e.TipoDocumento
	case "empresa_id":
		return e.EmpresaID
	case "resuelto":
		return e.Resuelto
	}

	// Si no se encuentra el campo específico, buscar en contexto
	if e.Contexto != nil {
		if valor, ok := e.Contexto[campo]; ok {
			return valor
		}
	}

	return nil
}

// SetField establece el valor de un campo
func (e *ErrorLog) SetField(campo string, valor interface{}) error {
	switch campo {
	case "id":
		if id, ok := valor.(string); ok {
			e.ID = id
			return nil
		}
	case "tipo":
		if tipo, ok := valor.(string); ok {
			e.Tipo = tipo
			return nil
		}
	case "mensaje":
		if mensaje, ok := valor.(string); ok {
			e.Mensaje = mensaje
			return nil
		}
	case "detalles":
		if detalles, ok := valor.(string); ok {
			e.Detalles = detalles
			return nil
		}
	case "fecha_hora":
		if fecha, ok := valor.(time.Time); ok {
			e.FechaHora = fecha
			return nil
		}
	case "modulo_origen":
		if modulo, ok := valor.(string); ok {
			e.ModuloOrigen = modulo
			return nil
		}
	case "funcion":
		if funcion, ok := valor.(string); ok {
			e.Funcion = funcion
			return nil
		}
	case "linea":
		if linea, ok := valor.(int); ok {
			e.Linea = linea
			return nil
		}
	case "usuario_id":
		if usuario, ok := valor.(string); ok {
			e.UsuarioID = usuario
			return nil
		}
	case "ip":
		if ip, ok := valor.(string); ok {
			e.IP = ip
			return nil
		}
	case "documento_id":
		if doc, ok := valor.(string); ok {
			e.DocumentoID = doc
			return nil
		}
	case "tipo_documento":
		if tipo, ok := valor.(string); ok {
			e.TipoDocumento = tipo
			return nil
		}
	case "empresa_id":
		if empresa, ok := valor.(string); ok {
			e.EmpresaID = empresa
			return nil
		}
	case "resuelto":
		if resuelto, ok := valor.(bool); ok {
			e.Resuelto = resuelto
			return nil
		}
	default:
		// Si no es un campo específico, almacenar en contexto
		if e.Contexto == nil {
			e.Contexto = make(map[string]interface{})
		}
		e.Contexto[campo] = valor
		return nil
	}

	return NewValidationFieldError(campo, "Tipo de dato inválido para el campo", "INVALID_TYPE", valor)
}
